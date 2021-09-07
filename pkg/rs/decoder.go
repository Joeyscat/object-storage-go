package rs

import (
	"io"

	"github.com/klauspost/reedsolomon"
)

type decoder struct {
	readers   []io.Reader
	writers   []io.Writer
	enc       reedsolomon.Encoder
	size      uint64
	cache     []byte
	cacheSize int
	total     uint64
}

func NewDecoder(readers []io.Reader, writers []io.Writer, size uint64) (*decoder, error) {
	enc, err := reedsolomon.New(DataShards, ParityShards)
	if err != nil {
		return nil, err
	}
	return &decoder{readers, writers, enc, size, nil, 0, 0}, nil
}

func (d *decoder) Read(buf []byte) (n int, err error) {
	if d.cacheSize == 0 {
		err := d.getData()
		if err != nil {
			return 0, err
		}
	}
	length := len(buf)
	if d.cacheSize < length {
		length = d.cacheSize
	}
	d.cacheSize -= length
	copy(buf, d.cache[:length])
	d.cache = d.cache[length:]
	return length, nil
}

func (d *decoder) getData() error {
	if d.total == d.size {
		return io.EOF
	}
	shards := make([][]byte, AllShards)
	repairIds := make([]int, 0)
	for i := range shards {
		if d.readers[i] == nil {
			repairIds = append(repairIds, i)
		} else {
			shards[i] = make([]byte, BlockPerShard)
			n, err := io.ReadFull(d.readers[i], shards[i])
			if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
				shards[i] = nil
			} else if n != BlockPerShard {
				shards[i] = shards[i][:n]
			}
		}
	}
	err := d.enc.Reconstruct(shards)
	if err != nil {
		return err
	}
	for i := range repairIds {
		id := repairIds[i]
		d.writers[id].Write(shards[id])
	}
	for i := 0; i < DataShards; i++ {
		shardSize := uint64(len(shards[i]))
		if d.total+shardSize > d.size {
			shardSize -= d.total + shardSize - d.size
		}
		d.cache = append(d.cache, shards[i][:shardSize]...)
		d.cacheSize += int(shardSize)
		d.total += shardSize
	}
	return nil
}
