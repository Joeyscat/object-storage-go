package objects

import (
    "fmt"
    "github.com/joeyscat/object-storage-go/pkg/mongo"
    "net/http"
    "strings"

    "github.com/joeyscat/object-storage-go/pkg/log"
)

func del(w http.ResponseWriter, r *http.Request) {
    name := strings.Split(r.URL.EscapedPath(), "/")[2]
    version, err := mongo.SearchLatestVersion(name)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    if version.Hash == "" && version.Size == 0 {
        log.Info(fmt.Sprintf("this object [%s] had deleted, could not delete again", name))
        w.WriteHeader(http.StatusNotFound)
        return
    }
    err = mongo.AddVersion(name, "", version.Version+1, 0)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
}
