package main

import (
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/joeyscat/object-storage-go/pkg/mongo"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

func main() {
    files, err := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
    if err != nil {
        log.Warn(err.Error())
        return
    }

    for _, file := range files {
        hash := strings.Split(filepath.Base(file), ".")[0]
        hashInMetadata, err := mongo.HasHash(hash)
        if err != nil {
            log.Warn(err.Error())
            return
        }
        if !hashInMetadata {
            del(hash)
        }
    }
}

func del(hash string) {
    log.Info("DELETE " + hash)
    url := "http://" + os.Getenv("LISTEN_ADDRESS") + "/objects/" + hash
    request, err := http.NewRequest("DELETE", url, nil)
    if err != nil {
        log.Warn(err.Error())
        return
    }
    client := http.Client{}
    if _, err = client.Do(request); err != nil {
        log.Warn(err.Error())
    }

}
