package mongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/joeyscat/object-storage-go/internal/api_server/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type datastore struct {
	db *mongo.Database
}

func (ds *datastore) Buckets() store.BucketStore {
	return newBuckets(ds)
}

func (ds *datastore) Objects() store.ObjectStore {
	panic("")
}

func (ds *datastore) Close() error {
	return ds.db.Client().Disconnect(context.Background())
}

var (
	mongodbFactory store.Factory
	once           sync.Once
)

func GetMongoFactoryOr(ctx context.Context, opts ...*options.ClientOptions) (store.Factory, error) {

	var err error
	var dbInst *mongo.Database
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var client *mongo.Client
		client, err = mongo.Connect(ctx, opts...)
		if err != nil {
			return
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			return
		}
		dbInst = client.Database("object_storage")

		mongodbFactory = &datastore{dbInst}

		err = createIndexes(ctx, dbInst)
	})

	if mongodbFactory == nil || err != nil {
		return nil, fmt.Errorf("failed to get mongo store factory, error: %w", err)
	}
	return mongodbFactory, nil
}

func createIndexes(ctx context.Context, db *mongo.Database) error {
	trueConst := true
	s, err := db.Collection("buckets").Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.M{"name": 1},
		Options: &options.IndexOptions{Unique: &trueConst}},
	)
	if err != nil {
		return err
	}
	println(s)

	return nil
}
