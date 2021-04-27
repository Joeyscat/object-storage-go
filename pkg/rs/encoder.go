package rs

import (
    "io"

    "github.com/klauspost/reedsolomon"
)

type encoder struct {
    writers []io.Writer
    enc     reedsolomon.Encoder
    cache   []byte
}

func NewEncoder(writers []io.Writer) *encoder {
    enc, _ := reedsolomon.New(DataShards, ParityShards)
    return &encoder{writers, enc, nil}
}

func (e *encoder) Write(data []byte) (n int, err error) {
    length := len(data)
    current := 0
    for length != 0 {
        next := BlockSize - len(e.cache)
        if next > length {
            next = length
        }
        e.cache = append(e.cache, data[current:current+next]...)
        if len(e.cache) == BlockSize {
            e.Flush()
        }
        current += next
        length -= next
    }
    return len(data), nil
}

func (e *encoder) Flush() {
    if len(e.cache) == 0 {
        return
    }
    shards, _ := e.enc.Split(e.cache)
    e.enc.Encode(shards)
    for i := range shards {
        e.writers[i].Write(shards[i])
    }
    e.cache = []byte{}
}
