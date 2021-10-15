package temp

import (
    "go.uber.org/zap"
    "net/http"
    "os"
    "strings"

    "github.com/joeyscat/object-storage-go/pkg/log"
)

func put(w http.ResponseWriter, r *http.Request) {
    uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
    tempInfo_, err := readFromFile(uuid)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusNotFound)
        return
    }
    infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
    datFile := infoFile + ".dat"
    f, err := os.Open(datFile)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    defer f.Close()
    info, err := f.Stat()
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    actual := info.Size()
    os.Remove(infoFile)
    if actual != tempInfo_.Size {
        os.Remove(datFile)
        log.Warn("actual size mismatch", zap.Int64("expect", tempInfo_.Size), zap.Int64("actual", actual))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    CommitTempObject(datFile, tempInfo_)
}
