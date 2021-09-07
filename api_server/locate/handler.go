package locate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/log"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(info)
	if err != nil {
		log.Warn(fmt.Sprintf("parse location info error: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
