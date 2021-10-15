package v1

import "github.com/joeyscat/object-storage-go/internal/api_server/store"

type BucketSrv interface {
}

type bucketService struct {
	store store.Factory
}

var _ BucketSrv = (*bucketService)(nil)

func newBuckets(srv *service) *bucketService {
	return &bucketService{store: srv.store}
}
