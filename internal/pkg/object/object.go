package object

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/joeyscat/object-storage-go/internal/api_server/heartbeat"
	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"
	"github.com/joeyscat/object-storage-go/pkg/rs"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func GetStream(hash string, size uint64) (*rs.RSGetStream, error) {
	locateInfo := Locate(hash)
	if len(locateInfo) < rs.DataShards {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.AllShards {
		dataServers = heartbeat.ChooseRandomDataServer(rs.AllShards-len(locateInfo), locateInfo)
	}
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)
}

func Locate(hash string) map[int]string {
	log.Debug("Locate", zap.String("hash", hash))
	nc, err := natsmq.GetSingletonNats(os.Getenv("NATS_URL"), nats.Name("object_locate_pub"))
	if err != nil {
		log.Error("GetSingletonNats", zap.Any("error", err))
		return nil
	}

	rs, err := natsmq.PublichAndWaitForReply(nc, os.Getenv("NATS_SUBJECT_OBJ_LOCATE"), []byte(hash), time.Second, rs.AllShards)
	if err != nil {
		log.Error("PublichAndWaitForReply", zap.Any("error", err))
		return nil
	}

	locateInfo := make(map[int]string)
	for _, r := range rs {
		var info Message
		err := json.Unmarshal(r.Data, &info)
		if err != nil {
			log.Error("Unmarshal msg", zap.Any("error", err))
			continue
		}
		locateInfo[info.Id] = info.Addr
		log.Debug("Locator receive", zap.String("msg", info.Addr))
	}
	log.Warn("Locator receive no msg")
	return locateInfo
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DataShards
}

type Message struct {
	Addr string
	Id   int
}
