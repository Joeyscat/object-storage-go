package temp

import (
    "encoding/json"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"

    "github.com/joeyscat/object-storage-go/pkg/log"
)

type tempInfo struct {
    Uuid string
    Name string
    Size int64
}

func (t *tempInfo) writeToFile() error {
    f, err := os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid)
    if err != nil {
        return err
    }
    defer f.Close()
    bs, _ := json.Marshal(t)
    f.Write(bs)
    return nil
}

func post(w http.ResponseWriter, r *http.Request) {
    output, err := exec.Command("uuidgen").Output()
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    uuid := strings.TrimSuffix(string(output), "\n")
    name := strings.Split(r.URL.EscapedPath(), "/")[2]
    size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    t := tempInfo{uuid, name, size}
    err = t.writeToFile()
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusInternalServerError)
        return
    }
    os.Create(os.Getenv("STORAGE_ROOT") + "/temp/" + t.Uuid + ".dat")
    w.Write([]byte(uuid))
}
