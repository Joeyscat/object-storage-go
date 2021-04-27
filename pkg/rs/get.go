package rs

import (
    "fmt"
    "io"

    "github.com/joeyscat/object-storage-go/pkg/objectstream"
)

type RSGetStream struct {
    *decoder
}

func NewRSGetStream(locateInfo map[int]string, dataServers []string, hash string, size uint64) (*RSGetStream, error) {
    if len(locateInfo)+len(dataServers) != AllShards {
        return nil, fmt.Errorf("data-servers number mismatch")
    }

    readers := make([]io.Reader, AllShards)
    for i := 0; i < AllShards; i++ {
        server := locateInfo[i]
        if server == "" {
            locateInfo[i] = dataServers[0]
            dataServers = dataServers[1:]
            continue
        }
        reader, err := objectstream.NewGetStream(server, fmt.Sprintf("%s.%d", hash, i))
        if err == nil {
            readers[i] = reader
        }
    }

    writers := make([]io.Writer, AllShards)
    perShard := (size + DataShards - 1) / DataShards
    var err error
    for i := range readers {
        if readers[i] == nil {
            writers[i], err = objectstream.NewTempPutStream(locateInfo[i], fmt.Sprintf("%s.%d", hash, i), perShard)
            if err != nil {
                return nil, err
            }
        }
    }

    dec := NewDecoder(readers, writers, size)
    return &RSGetStream{dec}, nil
}

func (s *RSGetStream) Close() {
    for i := range s.writers {
        if s.writers[i] != nil {
            s.writers[i].(*objectstream.TempPutStream).Commit(true)
        }
    }
}

func (s *RSGetStream) Seek(offset int64, whence int) (int64, error) {
    if whence != io.SeekCurrent {
        panic("only support SeekCurrent")
    }
    if offset < 0 {
        panic("only support SeekCurrent")
    }
    for offset != 0 {
        length := int64(BlockSize)
        if offset < length {
            length = offset
        }
        buf := make([]byte, length)
        io.ReadFull(s, buf)
        offset -= length
    }
    return offset, nil
}
