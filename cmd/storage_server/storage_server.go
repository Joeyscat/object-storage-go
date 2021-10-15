package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"

	"github.com/joeyscat/object-storage-go/internal/storage_server/heartbeat"
	"github.com/joeyscat/object-storage-go/internal/storage_server/locate"
	"github.com/joeyscat/object-storage-go/internal/storage_server/objects"
	"github.com/joeyscat/object-storage-go/internal/storage_server/temp"
)

func main() {
	defer natsmq.CloseSingletonNats()

	locate.CollectObjects()

	go heartbeat.StartHeartbeat()
	go locate.StartLocate()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	addr := os.Getenv("LISTEN_ADDRESS")
	log.Info(fmt.Sprintf("Listening on %s\n", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err.Error())
	}
}
