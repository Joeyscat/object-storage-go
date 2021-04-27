package temp

import (
    "fmt"
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/joeyscat/object-storage-go/pkg/rs"
    "net/http"
    "strings"
)

func head(w http.ResponseWriter, r *http.Request) {
    token := strings.Split(r.URL.EscapedPath(), "/")[2]
    stream, err := rs.NewRSResumablePutStreamFromToken(token)
    if err != nil {
        log.Warn(err.Error())
        w.WriteHeader(http.StatusForbidden)
        return
    }
    current := stream.CurrentSize()
    if current == -1 {
        w.WriteHeader(http.StatusNotFound)
        return
    }
    w.Header().Set("content-length", fmt.Sprintf("%d", current))
}
