package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct {
	writer *io.PipeWriter
	c      chan error
}

func NewPutStream(server, object string) *PutStream {
	reader, writer := io.Pipe()
	c := make(chan error)
	go func() {
		request, err := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		if err != nil {
			c <- err
			return
		}
		client := http.Client{}
		rsp, err := client.Do(request)
		if err == nil && rsp.StatusCode != http.StatusOK {
			err = fmt.Errorf("data-server return http code %d", rsp.StatusCode)
		}
		c <- err
	}()
	return &PutStream{writer, c}
}

func (p *PutStream) Write(data []byte) (n int, err error) {
	return p.writer.Write(data)
}

func (p *PutStream) Close() error {
	p.writer.Close()
	return <-p.c
}
