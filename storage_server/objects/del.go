package objects

import (
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/joeyscat/object-storage-go/storage_server/locate"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func del(w http.ResponseWriter, r *http.Request) {
    hash := strings.Split(r.URL.EscapedPath(), "/")[2]
    files, err := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + hash + ".*")
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    locate.Del(hash)
    // TODO files[0] ????
    err = os.Rename(files[0], os.Getenv("STORAGE_ROOT")+"/garbage/"+filepath.Base(files[0]))
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
    }
}
