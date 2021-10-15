package temp

import (
    "fmt"
    "github.com/joeyscat/object-storage-go/pkg/log"
    "net/http"
    "os"
    "strings"
)

func head(w http.ResponseWriter, r *http.Request) {
    uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
    f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusNotFound)
        return
    }
    defer f.Close()
    stat, err := f.Stat()
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    w.Header().Set("content-length", fmt.Sprintf("%d", stat.Size()))
}
