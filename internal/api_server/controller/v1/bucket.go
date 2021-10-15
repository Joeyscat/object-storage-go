package v1

import (
	"net/http"

	srvv1 "github.com/joeyscat/object-storage-go/internal/api_server/service/v1"
	"github.com/joeyscat/object-storage-go/internal/api_server/store"
	"github.com/labstack/echo/v4"
)

type BocketController struct {
	srv srvv1.Service
}

func NewBocketController(store store.Factory) *BocketController {
	return &BocketController{
		srv: srvv1.NewService(store),
	}
}

func (b *BocketController) GetBucket(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}
