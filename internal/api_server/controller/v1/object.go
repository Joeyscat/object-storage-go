package v1

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/joeyscat/object-storage-go/internal/api_server/heartbeat"
	srvv1 "github.com/joeyscat/object-storage-go/internal/api_server/service/v1"
	"github.com/joeyscat/object-storage-go/internal/api_server/store"
	"github.com/joeyscat/object-storage-go/internal/pkg/core"
	"github.com/joeyscat/object-storage-go/internal/pkg/object"
	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/mongo"
	"github.com/joeyscat/object-storage-go/pkg/rs"
	"github.com/joeyscat/object-storage-go/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ObjectController struct {
	srv srvv1.Service
}

func NewObjectController(store store.Factory) *ObjectController {
	return &ObjectController{
		srv: srvv1.NewService(store),
	}
}

func (o *ObjectController) PutObject(c echo.Context) error {
	r := c.Request()

	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Warn("missing object hash in digest header")
		return c.JSON(http.StatusBadRequest, nil)
	}

	size, err := utils.GetSizeFromHeader(r.Header)
	if err != nil {
		log.Warn(fmt.Sprintf("get size from header err: %v", err))
		return c.JSON(http.StatusBadRequest, nil)
	}

	code, err := storeObject(r.Body, hash, size)
	if err != nil {
		log.Warn(fmt.Sprintf("storeObject err: %v", err))
		return c.JSON(code, nil)
	}
	if code != http.StatusOK {
		return c.JSON(code, nil)
	}

	name := c.Param("name")
	latestVersion, err := mongo.SearchLatestVersion(name)
	if err != nil {
		log.Warn(fmt.Sprintf("SearchLatestVersion err: %v", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}

	if latestVersion.Version == 0 {
		err = mongo.PutMetadata(name, hash, uint64(size))
	} else {
		err = mongo.AddVersion(name, hash, latestVersion.Version+1, uint64(size))
	}
	if err != nil {
		log.Warn(fmt.Sprintf("AddVersion err: %v", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	return core.ToResponse(c, nil)
}

func (o *ObjectController) CreateObject(c echo.Context) error {
	r := c.Request()
	name := c.Param("name")

	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		log.Warn("parse size error", zap.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, nil)
	}

	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Warn("missing object hash in digest header")
		return c.JSON(http.StatusBadRequest, nil)
	}

	if object.Exist(url.PathEscape(hash)) {
		var v *mongo.Metadata
		v, err = mongo.SearchLatestVersion(name)
		if err != nil {
			log.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, nil)
		}
		err = mongo.AddVersion(name, hash, v.Version+1, uint64(size))
		if err != nil {
			log.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, nil)
		}

		return core.ToResponse(c, nil)
	}

	ds := heartbeat.ChooseRandomDataServer(rs.AllShards, nil)
	if len(ds) != rs.AllShards {
		log.Warn("cannot find enough data-server")
		return c.JSON(http.StatusServiceUnavailable, nil)
	}
	stream, err := rs.NewRsResumablePutStream(ds, name, url.PathEscape(hash), size)
	if err != nil {
		log.Warn(err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}
	token, err := stream.ToToken()
	if err != nil {
		log.Warn(err.Error())
		return c.JSON(http.StatusInternalServerError, nil)
	}
	c.Response().Header().Set("location", "/temp/"+url.PathEscape(token))
	c.Response().WriteHeader(http.StatusCreated)
	return nil
}

func (o *ObjectController) GetObject(c echo.Context) error {
	r := c.Request()
	w := c.Response()

	name := c.Param("name")
	versionId := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			log.Warn("parse version error", zap.String("error", err.Error()))
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("version invalid: %v", err))
		}
	}

	meta, err := mongo.GetMetadata(name, version)
	if err != nil {
		log.Warn("GetMetadata error", zap.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("GetMetadata error: %v", err))
	}
	if meta.Hash == "" {
		log.Warn("object not found: meta.Hash empty", zap.String("name", name))
		return c.JSON(http.StatusNotFound, "object not found")
	}

	hash := url.PathEscape(meta.Hash)
	stream, err := object.GetStream(hash, meta.Size)
	if err != nil {
		log.Warn("GetStream error", zap.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("GetStream error: %v", err))
	}
	defer stream.Close()
	// 分段下载
	offset, err := utils.GetOffsetFromHeader(r.Header)
	if err != nil {
		log.Warn("GetOffsetFromHeader error", zap.String("error", err.Error()))
		return c.JSON(http.StatusBadRequest, fmt.Sprintf("GetOffsetFromHeader error: %v", err))
	}
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes %d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	// 数据压缩
	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	if acceptGzip {
		w.Header().Set("content-encoding", "gzip")
		zw := gzip.NewWriter(w)
		defer zw.Close()
		_, err = io.Copy(zw, stream)
	} else {
		_, err = io.Copy(w, stream)
	}

	if err != nil {
		log.Warn("write data error", zap.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("write data error: %v", err))
	}
	return nil
}

func (o *ObjectController) DeleteObject(c echo.Context) error {
	name := c.Param("name")
	version, err := mongo.SearchLatestVersion(name)
	if err != nil {
		log.Warn("SearchLatestVersion error", zap.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("SearchLatestVersion error: %v", err))
	}
	if version.Hash == "" && version.Size == 0 {
		log.Debug(fmt.Sprintf("this object [%s] had deleted, could not delete again", name))
		return c.JSON(http.StatusNotFound, nil)
	}
	err = mongo.AddVersion(name, "", version.Version+1, 0)
	if err != nil {
		log.Warn("AddVersion error", zap.String("error", err.Error()))
		return c.JSON(http.StatusInternalServerError, fmt.Sprintf("AddVersion error: %v", err))
	}

	return core.ToResponse(c, nil)
}

func (o *ObjectController) GetObjectLocate(c echo.Context) error {
	w := c.Response()
	hash := c.Param("hash")

	info := object.Locate(hash)
	if len(info) == 0 {
		return c.JSON(http.StatusNotFound, nil)
	}
	b, err := json.Marshal(info)
	if err != nil {
		log.Warn(fmt.Sprintf("parse location info error: %v", err))
		return c.JSON(http.StatusInternalServerError, nil)
	}
	w.Write(b)
	return nil
}

func (o *ObjectController) HeadObjectVersion(c echo.Context) error {
	w := c.Response()
	name := c.Param("name")

	from := 0
	size := 1000
	for {
		metas, err := mongo.SearchAllVersions(name, int64(from), int64(size))
		if err != nil {
			log.Warn(err.Error())
			return c.JSON(http.StatusInternalServerError, nil)
		}
		for i := range metas {
			b, err := json.Marshal(metas[i])
			if err != nil {
				log.Warn(err.Error())
				return c.JSON(http.StatusInternalServerError, nil)
			}
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return nil
		}
		from += size
	}
}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if object.Exist(url.PathEscape(hash)) {
		// TODO 未校验数据与hash
		return http.StatusOK, nil
	}

	stream, err := putStream(url.PathEscape(hash), size)
	if err != nil {
		return http.StatusServiceUnavailable, err
	}

	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch, calculated=%s, requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServer(rs.AllShards, nil)
	if len(servers) != rs.AllShards {
		return nil, fmt.Errorf("cannot find enough data-server")
	}

	return rs.NewRSPutStream(servers, hash, size)
}
