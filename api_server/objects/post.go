package objects

import (
    "github.com/joeyscat/object-storage-go/api_server/heartbeat"
    "github.com/joeyscat/object-storage-go/api_server/locate"
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/joeyscat/object-storage-go/pkg/mongo"
    "github.com/joeyscat/object-storage-go/pkg/rs"
    "github.com/joeyscat/object-storage-go/pkg/utils"
    "net/http"
    "net/url"
    "strconv"
    "strings"
)

func post(w http.ResponseWriter, r *http.Request) {
    name := strings.Split(r.URL.EscapedPath(), "/")[2]
    size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusForbidden)
        return
    }
    hash := utils.GetHashFromHeader(r.Header)
    if hash == "" {
        log.Warn("missing object hash in digest header")
        w.WriteHeader(http.StatusBadRequest)
        return
    }
    if locate.Exist(url.PathEscape(hash)) {
        v, err := mongo.SearchLatestVersion(name)
        if err != nil {
            log.Warn(err.Error())
            w.WriteHeader(http.StatusInternalServerError)
            return
        }
        err = mongo.AddVersion(name, hash, v.Version+1, uint64(size))
        if err != nil {
            log.Warn(err.Error())
            w.WriteHeader(http.StatusInternalServerError)
        } else {
            w.WriteHeader(http.StatusOK)
        }
        return
    }
    ds := heartbeat.ChooseRandomDataServer(rs.AllShards, nil)
    if len(ds) != rs.AllShards {
        log.Warn("cannot find enough data-server")
        w.WriteHeader(http.StatusServiceUnavailable)
        return
    }
    stream, err := rs.NewRsResumablePutStream(ds, name, url.PathEscape(hash), size)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("location", "/temp/"+url.PathEscape(stream.ToToken()))
    w.WriteHeader(http.StatusCreated)
}
