package temp

import (
    "fmt"
    "github.com/joeyscat/object-storage-go/pkg/log"
    "net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    log.Info(fmt.Sprintf("%s %s\n", r.Method, r.URL))
    m := r.Method
    if m == http.MethodPut {
        put(w, r)
        return
    }
    if m == http.MethodHead {
        head(w, r)
        return
    }

    w.WriteHeader(http.StatusMethodNotAllowed)
}
