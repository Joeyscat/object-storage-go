package objects

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/utils"

	"github.com/joeyscat/object-storage-go/pkg/log"

	"github.com/joeyscat/object-storage-go/api_server/heartbeat"
	"github.com/joeyscat/object-storage-go/api_server/locate"
	"github.com/joeyscat/object-storage-go/pkg/mongo"
	"github.com/joeyscat/object-storage-go/pkg/rs"
)

func get(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	versionId := r.URL.Query()["version"]
	version := 0
	var err error
	if len(versionId) != 0 {
		version, err = strconv.Atoi(versionId[0])
		if err != nil {
			log.Warn(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	meta, err := mongo.GetMetadata(name, version)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	hash := url.PathEscape(meta.Hash)
	stream, err := GetStream(hash, meta.Size)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer stream.Close()
	// 分段下载
	offset, err := utils.GetOffsetFromHeader(r.Header)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
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
		w2 := gzip.NewWriter(w)
		defer w2.Close()
		_, err = io.Copy(w2, stream)
	} else {
		_, err = io.Copy(w, stream)
	}

	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func GetStream(hash string, size uint64) (*rs.RSGetStream, error) {
	locateInfo := locate.Locate(hash)
	if len(locateInfo) < rs.DataShards {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.AllShards {
		dataServers = heartbeat.ChooseRandomDataServer(rs.AllShards-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}
