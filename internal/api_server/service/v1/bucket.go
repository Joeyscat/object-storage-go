package v1

import (
	"context"
	"errors"
	"strings"

	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
	"github.com/joeyscat/object-storage-go/internal/api_server/store"
)

type BucketSrv interface {
	Create(ctx context.Context, bucket *v1.Bucket) error
	Delete(ctx context.Context, bucketName string, userID v1.UserID) error
	List(ctx context.Context, userID v1.UserID) ([]*v1.Bucket, error)
}

type bucketService struct {
	store store.Factory
}

var _ BucketSrv = (*bucketService)(nil)

func newBuckets(srv *service) *bucketService {
	return &bucketService{store: srv.store}
}

func (s *bucketService) Create(ctx context.Context, bucket *v1.Bucket) error {
	// check bucket
	if bucket.UserID == "" {
		return errors.New("bucket.UserID is 0")
	}
	if strings.TrimSpace(bucket.Name) == "" {
		return errors.New("bucket name if empty")
	}

	return s.store.Buckets().Create(ctx, bucket)
}

func (s *bucketService) Delete(ctx context.Context, bucketName string, userID v1.UserID) error {
	if userID == "" {
		return errors.New("bucket.UserID is 0")
	}
	if strings.TrimSpace(bucketName) == "" {
		return errors.New("bucket name if empty")
	}

	return s.store.Buckets().Delete(ctx, bucketName, userID)
}

func (s *bucketService) List(ctx context.Context, userID v1.UserID) ([]*v1.Bucket, error) {
	return s.store.Buckets().List(ctx, userID)
}
