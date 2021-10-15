package v1

import "github.com/joeyscat/object-storage-go/internal/api_server/store"

type Service interface {
	Buckets() BucketSrv
	Objects() ObjectSrv
}

type service struct {
	store store.Factory
}

func NewService(store store.Factory) Service {
	return &service{
		store: store,
	}
}

func (s *service) Buckets() BucketSrv {
	return newBuckets(s)
}

func (s *service) Objects() ObjectSrv {
	return newObjects(s)
}
