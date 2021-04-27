package rs

import (
    "encoding/base64"
    "encoding/json"
    "fmt"
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/joeyscat/object-storage-go/pkg/objectstream"
    "github.com/joeyscat/object-storage-go/pkg/utils"
    "io"
    "net/http"
)

type resumableToken struct {
    Name    string
    Size    int64
    Hash    string
    Servers []string
    Uuids   []string
}

type RSResumablePutStream struct {
    *RSPutStream
    *resumableToken
}

func NewRsResumablePutStream(dataServers []string, name, hash string, size int64) (*RSResumablePutStream, error) {
    putStream, err := NewRSPutStream(dataServers, hash, size)
    if err != nil {
        return nil, err
    }
    uuids := make([]string, AllShards)
    for i := range uuids {
        uuids[i] = putStream.writers[i].(*objectstream.TempPutStream).Uuid
    }
    token := &resumableToken{
        Name:    name,
        Size:    size,
        Hash:    hash,
        Servers: dataServers,
        Uuids:   uuids,
    }
    return &RSResumablePutStream{
        RSPutStream:    putStream,
        resumableToken: token,
    }, nil
}

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {
    b, err := base64.StdEncoding.DecodeString(token)
    if err != nil {
        return nil, err
    }
    var t resumableToken
    err = json.Unmarshal(b, &t)
    if err != nil {
        return nil, err
    }

    writers := make([]io.Writer, AllShards)
    for i := range writers {
        writers[i] = &objectstream.TempPutStream{
            Server: t.Servers[i],
            Uuid:   t.Uuids[i],
        }
    }
    enc := NewEncoder(writers)
    return &RSResumablePutStream{
        RSPutStream:    &RSPutStream{enc},
        resumableToken: &t,
    }, nil
}

func (s *RSResumablePutStream) ToToken() string {
    b, _ := json.Marshal(s)
    return base64.StdEncoding.EncodeToString(b)
}

func (s *RSResumablePutStream) CurrentSize() int64 {
    url := fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0])
    resp, err := http.Head(url)
    if err != nil {
        log.Warn(err.Error())
        return -1
    }
    if resp.StatusCode != http.StatusOK {
        log.Warn(fmt.Sprintf("%d - %s", resp.StatusCode, url))
        return -1
    }
    size := utils.GetSizeFromHeader(resp.Header) * DataShards
    if size > s.Size {
        size = s.Size
    }
    return size
}
