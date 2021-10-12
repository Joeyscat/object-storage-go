package locate

import (
	"encoding/json"
	"os"
	"time"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"
	"github.com/joeyscat/object-storage-go/pkg/rs"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Message struct {
	Addr string
	Id   int
}

func Locate(name string) map[int]string {
	log.Debug("Locate", zap.String("name", name))
	nc, err := natsmq.GetSingletonNats(os.Getenv("NATS_URL"), nats.Name("object_locate_pub"))
	if err != nil {
		log.Error("GetSingletonNats", zap.Any("error", err))
		return nil
	}

	rs, err := natsmq.PublichAndWaitForReply(nc, os.Getenv("NATS_SUBJECT_OBJ_LOCATE"), []byte(name), time.Second, rs.AllShards)
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
	}
	return locateInfo
}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DataShards
}
