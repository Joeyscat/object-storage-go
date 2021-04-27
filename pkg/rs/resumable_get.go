package rs

import (
    "github.com/joeyscat/object-storage-go/pkg/objectstream"
    "io"
)

type RSResumableGetStream struct {
    *decoder
}

func NewRsResumableGetStream(dataServers []string, uuids []string, size int64) (*RSResumableGetStream, error) {
    readers := make([]io.Reader, AllShards)
    var err error
    for i := 0; i < AllShards; i++ {
        readers[i], err = objectstream.NewTempGetStream(dataServers[i], uuids[i])
        if err != nil {
            return nil, err
        }
    }
    writers := make([]io.Writer, AllShards)
    dec := NewDecoder(readers, writers, uint64(size))
    return &RSResumableGetStream{dec}, nil
}
