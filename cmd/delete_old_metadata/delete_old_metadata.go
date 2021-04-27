package main

import (
    "github.com/joeyscat/object-storage-go/pkg/log"
    "github.com/joeyscat/object-storage-go/pkg/mongo"
)

const MinVersionCount = 5

func main() {
	buckets, err := mongo.SearchVersionStatus(MinVersionCount + 1)
	if err != nil {
		log.Warn(err.Error())
		return
	}

	for _, bucket := range buckets {
		for v := 0; v < bucket.DocCount-MinVersionCount; v++ {
			mongo.DelMetadata(bucket.Key, v+int(bucket.MinVersion.Value))
		}
	}
}
