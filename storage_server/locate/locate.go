package locate

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"

    "github.com/joeyscat/object-storage-go/pkg/log"

    "github.com/joeyscat/object-storage-go/pkg/rabbitmq"
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
    q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
    defer q.Close()
    q.Bind("data-server")
    c := q.Consume()

    for msg := range c {
        hash, err := strconv.Unquote(string(msg.Body))
        if err != nil {
            panic(err)
        }
        id := Locate(hash)
        if id != -1 {
            q.Send(msg.ReplyTo, Message{
                Addr: os.Getenv("LISTEN_ADDRESS"),
                Id:   id,
            })
        }
    }
}

func CollectObjects() {
    files, _ := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
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
