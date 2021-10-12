package objectstream

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/log"
)

type TempPutStream struct {
	Server string
	Uuid   string
}

func NewTempPutStream(server, object string, size uint64) (*TempPutStream, error) {
	req, err := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("post object to data server not OK")
	}
	uuid, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{server, string(uuid)}, nil
}

func (w *TempPutStream) Write(bs []byte) (n int, err error) {
	req, err := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid,
		strings.NewReader(string(bs)))
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("data-server return http code %d", resp.StatusCode)
	}
	return len(bs), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	req, err := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	client := http.Client{}
	client.Do(req)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
