package mongo

import (
	"context"
	"testing"

	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
	"github.com/joeyscat/object-storage-go/internal/api_server/store"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Test_buckets_Create(t *testing.T) {
	ctx := context.Background()
	factory, err := GetMongoFactoryOr(ctx, options.Client().ApplyURI(TEST_MONGODB_URI))
	assert.Nil(t, err)

	type fields struct {
		store store.BucketStore
	}
	type args struct {
		ctx    context.Context
		bucket *v1.Bucket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"OK",
			fields{factory.Buckets()},
			args{ctx, &v1.Bucket{UserID: v1.UserID("616be1fb033267765dca1ef6"), Name: "bucket01"}},
			false,
		},
		{
			"OK",
			fields{factory.Buckets()},
			args{ctx, &v1.Bucket{UserID: v1.UserID("616be1fb033267765dca1ef6"), Name: "bucket02"}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.fields.store
			if err := b.Create(tt.args.ctx, tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("buckets.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buckets_Delete(t *testing.T) {
	ctx := context.Background()
	factory, err := GetMongoFactoryOr(ctx, options.Client().ApplyURI(TEST_MONGODB_URI))
	assert.Nil(t, err)

	type fields struct {
		store store.BucketStore
	}
	type args struct {
		ctx        context.Context
		bucketName string
		userID     v1.UserID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"OK",
			fields{factory.Buckets()},
			args{ctx, "bucket01", v1.UserID("616be1fb033267765dca1ef6")},
			false,
		},
		{
			"NOT OK",
			fields{factory.Buckets()},
			args{ctx, "not_exists_bucket", v1.UserID("616be1fb033267765dca1ef6")},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.fields.store
			if err := b.Delete(tt.args.ctx, tt.args.bucketName, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("buckets.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_buckets_List(t *testing.T) {
	ctx := context.Background()
	factory, err := GetMongoFactoryOr(ctx, options.Client().ApplyURI(TEST_MONGODB_URI))
	assert.Nil(t, err)

	type fields struct {
		store store.BucketStore
	}
	type args struct {
		ctx    context.Context
		userID v1.UserID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*v1.Bucket
		wantErr bool
	}{
		{
			"OK",
			fields{factory.Buckets()},
			args{ctx, v1.UserID("616be1fb033267765dca1ef6")},
			[]*v1.Bucket{{}},
			false,
		},
		{
			"No Result",
			fields{factory.Buckets()},
			args{ctx, v1.UserID("xx")},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.fields.store
			got, err := b.List(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("buckets.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(tt.want) > 0 {
				assert.NotEmpty(t, got)

				for _, v := range got {
					t.Logf("%v", v)
				}
			}
		})
	}
}
