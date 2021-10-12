package rs

import (
	"fmt"
	"io"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/objectstream"
	"go.uber.org/zap"
)

type RSPutStream struct {
	*encoder
}

func NewRSPutStream(dataServers []string, hash string, size int64) (*RSPutStream, error) {
	if len(dataServers) != AllShards {
		return nil, fmt.Errorf("data-servers number mismatch")
	}

	perShard := (size + DataShards - 1) / DataShards
	writers := make([]io.Writer, AllShards)
	var err error
	for i := range writers {
		writers[i], err = objectstream.NewTempPutStream(dataServers[i],
			fmt.Sprintf("%s.%d", hash, i), uint64(perShard))
		if err != nil {
			log.Error("NewTempPutStream error", zap.String("err", err.Error()))
			return nil, err
		}
	}
	enc, err := NewEncoder(writers)
	if err != nil {
		return nil, err
	}

	return &RSPutStream{enc}, nil
}

// Commit 提交上传结果
// 当 success=true 时,将临时对象转正,否则删除临时对象
func (s *RSPutStream) Commit(success bool) {
	s.Flush()
	for i := range s.writers {
		s.writers[i].(*objectstream.TempPutStream).Commit(success)
	}
}
