package objectstream

import (
    "fmt"
    "io"
    "net/http"

    "github.com/joeyscat/object-storage-go/pkg/log"
)

type GetStream struct {
    reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {
    log.Info(fmt.Sprintf("newGetStream: %s", url))
    resp, err := http.Get(url)
    if err != nil {
        return nil, err
    }
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("data-server return http code: %d", resp.StatusCode)
    }
    return &GetStream{resp.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
    if server == "" || object == "" {
        return nil, fmt.Errorf("invalid server %s, object %s", server, object)
    }
    return newGetStream("http://" + server + "/objects/" + object)
}

func (s *GetStream) Read(buf []byte) (n int, err error) {
    return s.reader.Read(buf)
}
