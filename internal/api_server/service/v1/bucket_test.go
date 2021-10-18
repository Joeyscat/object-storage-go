package v1

import (
	"context"
	"reflect"
	"testing"

	gomock "github.com/golang/mock/gomock"
	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
	"github.com/joeyscat/object-storage-go/internal/api_server/store"
)

var (
	expectBucketList = []*v1.Bucket{}
)

func Test_bucketService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockFactory(ctrl)
	bucketStore := store.NewMockBucketStore(ctrl)
	mockStore.EXPECT().Buckets().Return(bucketStore)
	bucketStore.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
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
			fields{mockStore},
			args{context.Background(), &v1.Bucket{UserID: "1", Name: "bucket01"}},
			false,
		},
		{
			"bucket filed empty error",
			fields{mockStore},
			args{context.Background(), &v1.Bucket{Name: "bucket01"}},
			true,
		},
		{
			"bucket filed empty error",
			fields{mockStore},
			args{context.Background(), &v1.Bucket{UserID: "1"}},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &bucketService{
				store: tt.fields.store,
			}
			if err := s.Create(tt.args.ctx, tt.args.bucket); (err != nil) != tt.wantErr {
				t.Errorf("bucketService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bucketService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockFactory(ctrl)
	bucketStore := store.NewMockBucketStore(ctrl)
	mockStore.EXPECT().Buckets().Return(bucketStore)
	bucketStore.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		store store.Factory
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
			fields{mockStore},
			args{context.Background(), "bucket01", "xx"},
			false,
		},
		{
			"bucket filed empty error",
			fields{mockStore},
			args{context.Background(), "", "xx"},
			true,
		},
		{
			"bucket filed empty error",
			fields{mockStore},
			args{context.Background(), "bucket01", ""},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &bucketService{
				store: tt.fields.store,
			}
			if err := s.Delete(tt.args.ctx, tt.args.bucketName, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("bucketService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_bucketService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := store.NewMockFactory(ctrl)
	bucketStore := store.NewMockBucketStore(ctrl)
	mockStore.EXPECT().Buckets().Return(bucketStore)
	bucketStore.EXPECT().List(gomock.Any(), gomock.Any()).Return(expectBucketList, nil)

	type fields struct {
		store store.Factory
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
			fields{mockStore},
			args{context.Background(), "1"},
			expectBucketList,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &bucketService{
				store: tt.fields.store,
			}
			got, err := s.List(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("bucketService.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bucketService.List() = %v, want %v", got, tt.want)
			}
		})
	}
}
