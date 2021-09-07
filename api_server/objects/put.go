package objects

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/mongo"

	"github.com/joeyscat/object-storage-go/pkg/log"

	"github.com/joeyscat/object-storage-go/api_server/heartbeat"
	"github.com/joeyscat/object-storage-go/api_server/locate"
	"github.com/joeyscat/object-storage-go/pkg/rs"
	"github.com/joeyscat/object-storage-go/pkg/utils"
)

func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Warn("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	size, err := utils.GetSizeFromHeader(r.Header)
	if err != nil {
		log.Warn(fmt.Sprintf("get size from header err: %v", err))
		return
	}
	c, err := storeObject(r.Body, hash, size)
	if err != nil {
		log.Warn(fmt.Sprintf("storeObject err: %v", err))
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	latestVersion, err := mongo.SearchLatestVersion(name)
	if err != nil {
		log.Warn(fmt.Sprintf("SearchLatestVersion err: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	if latestVersion.Version == 0 {
		err = mongo.PutMetadata(name, hash, uint64(size))
	} else {
		err = mongo.AddVersion(name, hash, latestVersion.Version+1, uint64(size))
	}
	if err != nil {
		log.Warn(fmt.Sprintf("AddVersion err: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
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
