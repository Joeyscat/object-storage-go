package temp

import (
    "github.com/joeyscat/object-storage-go/pkg/log"
    "io"
    "net/http"
    "os"
    "strings"
)

func get(w http.ResponseWriter, r *http.Request) {
    uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
    f, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid + ".dat")
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusNotFound)
        return
    }
    defer f.Close()
    io.Copy(w, f)
}
