package versions

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/log"

	"github.com/joeyscat/object-storage-go/pkg/mongo"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	from := 0
	size := 1000
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	for {
		metas, err := mongo.SearchAllVersions(name, int64(from), int64(size))
		if err != nil {
			log.Warn(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		for i := range metas {
			b, err := json.Marshal(metas[i])
			if err != nil {
				log.Warn(err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(b)
			w.Write([]byte("\n"))
		}
		if len(metas) != size {
			return
		}
		from += size
	}
}
