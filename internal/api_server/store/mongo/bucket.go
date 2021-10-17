package mongo

import (
	"context"
	"errors"

	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type buckets struct {
	db *mongo.Database
}

func newBuckets(ds *datastore) *buckets {
	return &buckets{db: ds.db}
}

type Bucket struct {
	ID     primitive.ObjectID `bson:"_id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Name   string
}

func (b *buckets) Create(ctx context.Context, bucket *v1.Bucket) error {
	userID, err := primitive.ObjectIDFromHex(string(bucket.UserID))
	if err != nil {
		return err
	}

	bk := Bucket{
		ID:     primitive.NewObjectID(),
		UserID: userID,
		Name:   bucket.Name,
	}
	_, err = b.db.Collection("buckets").InsertOne(ctx, bk)
	return err
}

func (b *buckets) Delete(ctx context.Context, bucketName string, userID v1.UserID) error {
	userID1, err := primitive.ObjectIDFromHex(string(userID))
	if err != nil {
		return err
	}

	r, err := b.db.Collection("buckets").DeleteOne(ctx, bson.M{
		"name":    bucketName,
		"user_id": userID1,
	})
	if r.DeletedCount < 1 {
		return errors.New("deleted no buckets")
	}

	return err
}

func (b *buckets) List(ctx context.Context, userID v1.UserID) ([]*v1.Bucket, error) {
	userID1, err := primitive.ObjectIDFromHex(string(userID))
	if err != nil {
		return nil, err
	}

	r, err := b.db.Collection("buckets").Find(ctx, bson.M{
		"user_id": userID1,
	})
	if err != nil {
		return nil, err
	}

	var buckets []*v1.Bucket
	err = r.All(ctx, &buckets)

	return buckets, err
}
