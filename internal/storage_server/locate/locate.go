package locate

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Message struct {
	Addr string
	Id   int
}

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	log.Info(fmt.Sprintf("Locate: %s", hash))
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	log.Info(fmt.Sprintf("Add: %s %d", hash, id))
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

func StartLocate() {
	nc, err := natsmq.GetSingletonNats(os.Getenv("NATS_URL"), nats.Name("object_locate_sub"))
	if err != nil {
		log.Error("GetSingletonNats", zap.Any("error", err))
		return
	}

	natsmq.SubscribeWithReply(nc, os.Getenv("NATS_SUBJECT_OBJ_LOCATE"), func(msg *nats.Msg) ([]byte, error) {
		hash := string(msg.Data)
		id := Locate(hash)
		if id != -1 {
			return json.Marshal(Message{
				Addr: os.Getenv("LISTEN_ADDRESS"),
				Id:   id,
			})
		}
		return nil, fmt.Errorf("hash[%s] not found", hash)
	})
}

func CollectObjects() {
	files, err := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	if err != nil {
		log.Warn(err.Error())
		return
	}
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, err := strconv.Atoi(file[1])
		if err != nil {
			panic(err)
		}
		objects[hash] = id
	}
}
