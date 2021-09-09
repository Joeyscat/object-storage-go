package heartbeat

import (
	"os"
	"sync"
	"time"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	nc, err := natsmq.GetSingletonNats(os.Getenv("NATS_URL"), nats.Name("storage_heartbeat_sub"))
	if err != nil {
		log.Error("GetSingletonNats", zap.Any("error", err))
		return
	}

	go removeExpiredDataServer()

	nc.Subscribe(os.Getenv("NATS_SUBJECT_STORAG_HEARTBEAT"), func(msg *nats.Msg) {
		ss := string(msg.Data)
		mutex.Lock()
		dataServers[ss] = time.Now()
		mutex.Unlock()
	})
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}

func getDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s := range dataServers {
		ds = append(ds, s)
	}
	return ds
}
