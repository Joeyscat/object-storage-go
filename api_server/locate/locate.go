package locate

import (
    "encoding/json"
    "log"
    "os"
    "time"

    "github.com/joeyscat/object-storage-go/pkg/rabbitmq"
    "github.com/joeyscat/object-storage-go/pkg/rs"
)

type Message struct {
    Addr string
    Id   int
}

func Locate(name string) (locateInfo map[int]string) {
    q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
    q.Publish("data-server", name)
    c := q.Consume()
    go func() {
        time.Sleep(time.Second)
        q.Close()
    }()

    locateInfo = make(map[int]string)
    for i := 0; i < rs.AllShards; i++ {
        msg := <-c
        if len(msg.Body) == 0 {
            return
        }
        var info Message
        json.Unmarshal(msg.Body, &info)
        log.Printf("Locate: %+v", info)
        locateInfo[info.Id] = info.Addr
    }
    return
}

func Exist(name string) bool {
    return len(Locate(name)) >= rs.DataShards
}
