package v1

import "github.com/joeyscat/object-storage-go/internal/api_server/store"

type ObjectSrv interface {
}

type objectService struct {
	store store.Factory
}

var _ ObjectSrv = (*objectService)(nil)

func newObjects(srv *service) *objectService {
	return &objectService{store: srv.store}
}
