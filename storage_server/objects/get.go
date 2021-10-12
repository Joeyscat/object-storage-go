package objects

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/log"

	"github.com/joeyscat/object-storage-go/storage_server/locate"
)

func get(w http.ResponseWriter, r *http.Request) {
	f := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if f == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	writeFileWithGzip(w, f)
}

func getFile(name string) string {
	files, err := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/" + name + ".*")
	if err != nil {
		panic(err)
	}
	file := files[0]
	h := sha256.New()
	writeFileWithGzip(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		log.Warn(fmt.Sprintf("object hash mismatch, remove %s", file))
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	return file
}

func writeFileWithGzip(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer f.Close()

	gzipStream, err := gzip.NewReader(f)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer gzipStream.Close()

	io.Copy(w, gzipStream)
}
