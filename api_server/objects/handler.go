package objects

import (
    "fmt"
    "net/http"

    "github.com/joeyscat/object-storage-go/pkg/log"
)

func Handler(w http.ResponseWriter, r *http.Request) {
    log.Info(fmt.Sprintf("%s %s\n", r.Method, r.URL))
    m := r.Method
    if m == http.MethodPut {
        put(w, r)
        return
    }
    if m == http.MethodGet {
        get(w, r)
        return
    }
    if m == http.MethodPost {
        post(w, r)
        return
    }
    if m == http.MethodDelete {
        del(w, r)
        return
    }
    w.WriteHeader(http.StatusMethodNotAllowed)
}
