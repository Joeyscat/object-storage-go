package main

import (
	"github.com/joeyscat/object-storage-go/pkg/natsmq"

	"github.com/joeyscat/object-storage-go/internal/api_server"
	"github.com/joeyscat/object-storage-go/internal/api_server/heartbeat"
)

func main() {
	defer natsmq.CloseSingletonNats()

	go heartbeat.ListenHeartbeat()

	api_server.InitRouter()
}
