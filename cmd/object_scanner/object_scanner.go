package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joeyscat/object-storage-go/internal/pkg/object"
	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/joeyscat/object-storage-go/pkg/mongo"
	"github.com/joeyscat/object-storage-go/pkg/utils"
)

func main() {
	files, err := filepath.Glob(os.Getenv("STORAGE_ROOT") + "/objects/*")
	if err != nil {
		log.Warn(err.Error())
		return
	}

	for _, file := range files {
		hash := strings.Split(filepath.Base(file), ".")[0]
		verify(hash)
	}
}

func verify(hash string) {
	log.Info("verify " + hash)
	size, err := mongo.SearchHashSize(hash)
	if err != nil {
		log.Warn(err.Error())
		return
	}
	stream, err := object.GetStream(hash, uint64(size))
	if err != nil {
		log.Warn(err.Error())
		return
	}
	defer stream.Close()
	d := utils.CalculateHash(stream)
	if d != hash {
		log.Warn(fmt.Sprintf("object hash mismatch, calculated=%s, requested=%s", d, hash))
	}
}
