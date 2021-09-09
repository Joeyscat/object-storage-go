package heartbeat

import (
	"os"
	"time"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func StartHeartbeat() {
	nc, err := natsmq.GetSingletonNats(os.Getenv("NATS_URL"), nats.Name("storage_heartbeat_pub"))
	if err != nil {
		log.Error("GetSingletonNats", zap.Any("error", err))
		return
	}

	for {
		nc.Publish(os.Getenv("NATS_SUBJECT_STORAG_HEARTBEAT"), []byte(os.Getenv("LISTEN_ADDRESS")))
		time.Sleep(5 * time.Second)
	}
}
