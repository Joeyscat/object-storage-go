package mongo

import (
	"context"
	"testing"

	"github.com/joeyscat/object-storage-go/internal/api_server/store"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	TEST_MONGODB_URI = "mongodb://object_storage_rw:5QXj_hVQ7_5r5oOr1KVXjGam00qVgCZ35d5BmTxTYDpemN4d7o7SxCp1euiGtCR3@127.0.0.1:27017/object_storage"
)

func TestGetMongoFactoryOr(t *testing.T) {
	type args struct {
		ctx  context.Context
		opts []*options.ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		want    store.Factory
		wantErr bool
	}{
		{
			"OK",
			args{context.Background(), []*options.ClientOptions{options.Client().ApplyURI(TEST_MONGODB_URI)}},
			&datastore{nil},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetMongoFactoryOr(tt.args.ctx, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMongoFactoryOr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.want == nil {
				assert.Nil(t, got)
			}
		})
	}
}
