package temp

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/joeyscat/object-storage-go/pkg/log"
)

func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempInfo_, err := readFromFile(uuid)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := os.Getenv("STORAGE_ROOT") + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, err := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := f.Stat()
	if err != nil {
		log.Warn(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	if actual > tempInfo_.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		log.Warn("actual size mismatch", zap.Int64("expect", tempInfo_.Size), zap.Int64("actual", actual))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/temp/" + uuid)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bs, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var info tempInfo
	json.Unmarshal(bs, &info)
	return &info, nil
}
