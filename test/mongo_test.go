package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	uri = "mongodb://object_storage_rw:5QXj_hVQ7_5r5oOr1KVXjGam00qVgCZ35d5BmTxTYDpemN4d7o7SxCp1euiGtCR3@127.0.0.1:27017/object_storage"
)

func TestConnectingMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	assert.Nil(t, err)

	err = client.Ping(ctx, readpref.Primary())
	assert.Nil(t, err)

	r, err := client.Database("object_storage").ListCollections(context.Background(),
		bson.D{})
	assert.Nil(t, err)

	t.Logf("%v", r)

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
