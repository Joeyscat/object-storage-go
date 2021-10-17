package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
	srvv1 "github.com/joeyscat/object-storage-go/internal/api_server/service/v1"
	"github.com/joeyscat/object-storage-go/internal/pkg/auth"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	expectBucketList = []*v1.Bucket{}
	userForAuth      = auth.User{UserID: 1, UserName: "JoJo"}
)

func TestBocketController_GetBucketList(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	s, err := json.Marshal(userForAuth)
	assert.Nil(t, err)
	req.Header.Set(auth.HEADER_USER_KEY, string(s))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	bucketService := srvv1.NewMockBucketSrv(ctrl)
	mockService.EXPECT().Buckets().Return(bucketService)
	bucketService.EXPECT().List(gomock.Any(), gomock.Any()).Return(expectBucketList, nil)

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"OK",
			fields{mockService},
			args{c},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BocketController{
				srv: tt.fields.srv,
			}
			if err := b.GetBucketList(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("BocketController.GetBucketList() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, http.StatusOK, c.Response().Status)
		})
	}
}

func TestBocketController_CreateBucket(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/buckets/:bucketName")
	c.SetParamNames("bucketName")
	c.SetParamValues("bucket01")
	s, err := json.Marshal(userForAuth)
	assert.Nil(t, err)
	req.Header.Set(auth.HEADER_USER_KEY, string(s))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	bucketService := srvv1.NewMockBucketSrv(ctrl)
	mockService.EXPECT().Buckets().Return(bucketService)
	bucketService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"OK",
			fields{mockService},
			args{c},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BocketController{
				srv: tt.fields.srv,
			}
			if err := b.CreateBucket(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("BocketController.CreateBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBocketController_DeleteBucket(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/buckets/:bucketName")
	c.SetParamNames("bucketName")
	c.SetParamValues("bucket01")
	s, err := json.Marshal(userForAuth)
	assert.Nil(t, err)
	req.Header.Set(auth.HEADER_USER_KEY, string(s))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := srvv1.NewMockService(ctrl)
	bucketService := srvv1.NewMockBucketSrv(ctrl)
	mockService.EXPECT().Buckets().Return(bucketService)
	bucketService.EXPECT().Delete(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	type fields struct {
		srv srvv1.Service
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"OK",
			fields{mockService},
			args{c},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BocketController{
				srv: tt.fields.srv,
			}
			if err := b.DeleteBucket(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("BocketController.DeleteBucket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
