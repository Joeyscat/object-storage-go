package main

import (
	"net/http"
	"os"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/natsmq"

	"github.com/joeyscat/object-storage-go/api_server/heartbeat"
	"github.com/joeyscat/object-storage-go/api_server/locate"
	"github.com/joeyscat/object-storage-go/api_server/objects"
	"github.com/joeyscat/object-storage-go/api_server/temp"
	"github.com/joeyscat/object-storage-go/api_server/versions"
)

func main() {
	defer natsmq.CloseSingletonNats()

	go heartbeat.ListenHeartbeat()

	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/locate/", locate.Handler)
	http.HandleFunc("/versions/", versions.Handler)

	addr := os.Getenv("LISTEN_ADDRESS")
	log.Info(addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err.Error())
	}
}
