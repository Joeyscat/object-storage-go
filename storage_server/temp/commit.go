package temp

import (
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/utils"

	"github.com/joeyscat/object-storage-go/storage_server/locate"
)

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo) id() (int, error) {
	s := strings.Split(t.Name, ".")
	id, err := strconv.Atoi(s[1])
	if err != nil {
		return 0, err
	}
	return id, nil
}

func CommitTempObject(datFile string, t *tempInfo) {
	f, err := os.Open(datFile)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer f.Close()
	d := url.PathEscape(utils.CalculateHash(f))
	f.Seek(0, io.SeekStart)

	// 压缩到正式对象文件，然后删除临时对象文件
	w, err := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + t.Name + "." + d)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer w.Close()
	w2 := gzip.NewWriter(w)
	defer w2.Close()
	io.Copy(w2, f)

	os.Remove(datFile)

	id, err := t.id()
	if err != nil {
		log.Warn(err.Error())
		return
	}
	locate.Add(t.hash(), id)
}
