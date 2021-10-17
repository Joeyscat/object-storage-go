package v1

import (
	"net/http"
	"strings"

	v1 "github.com/joeyscat/object-storage-go/internal/api_server/model/v1"
	srvv1 "github.com/joeyscat/object-storage-go/internal/api_server/service/v1"
	"github.com/joeyscat/object-storage-go/internal/api_server/store"
	"github.com/joeyscat/object-storage-go/internal/pkg/auth"
	"github.com/labstack/echo/v4"
)

type BocketController struct {
	srv srvv1.Service
}

func NewBucketController(store store.Factory) *BocketController {
	return &BocketController{
		srv: srvv1.NewService(store),
	}
}

func (b *BocketController) GetBucketList(c echo.Context) error {
	u, err := auth.GetUser(c)
	if err != nil {
		return auth.UserInfoNotFoundInRequest(c)
	}

	bs, err := b.srv.Buckets().List(c.Request().Context(), v1.UserID(u.UserID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, bs)
}

func (b *BocketController) CreateBucket(c echo.Context) error {
	u, err := auth.GetUser(c)
	if err != nil {
		return auth.UserInfoNotFoundInRequest(c)
	}
	bucketName := c.Param("bucketName")
	if strings.TrimSpace(bucketName) == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	var bucket = v1.Bucket{
		UserID: v1.UserID(u.UserID),
		Name:   bucketName,
	}
	err = b.srv.Buckets().Create(c.Request().Context(), &bucket)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusCreated, nil)
}

func (b *BocketController) DeleteBucket(c echo.Context) error {
	u, err := auth.GetUser(c)
	if err != nil {
		return auth.UserInfoNotFoundInRequest(c)
	}
	bucketName := c.Param("bucketName")
	if strings.TrimSpace(bucketName) == "" {
		return c.JSON(http.StatusBadRequest, nil)
	}

	err = b.srv.Buckets().Delete(c.Request().Context(), bucketName, v1.UserID(u.UserID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, nil)
}
