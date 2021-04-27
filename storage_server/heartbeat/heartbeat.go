package heartbeat

import (
    "os"
    "time"

    "github.com/joeyscat/object-storage-go/pkg/rabbitmq"
)

// StartHeartbeat 每5秒向指定exchange发送一条消息, 用于暴露本服务节点的监听地址
func StartHeartbeat() {
    q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
    defer q.Close()
    for {
        q.Publish("api-server", os.Getenv("LISTEN_ADDRESS"))
        time.Sleep(5 * time.Second)
    }
}
