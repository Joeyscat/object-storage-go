package store

import (
	"context"

	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
)

type BucketStore interface {
	Create(ctx context.Context, bucket *v1.Bucket) error
	Delete(ctx context.Context, bucketName string, userID v1.UserID) error
	List(ctx context.Context, userID v1.UserID) ([]*v1.Bucket, error)
}
