package core

import (
	"net/http"

	"github.com/joeyscat/object-storage-go/internal/pkg/errcode"
	"github.com/labstack/echo/v4"
)

func ToResponse(ctx echo.Context, data interface{}) error {
	if data == nil {
		data = map[string]interface{}{}
	}
	return ctx.JSON(http.StatusOK, data)
}

func ToResponseList(ctx echo.Context, list interface{}, totalRows int) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
		"pager": Pager{
			Page:      GetPage(ctx),
			PageSize:  GetPageSize(ctx),
			TotalRows: totalRows,
		},
	})
}

func ToErrorResponse(ctx echo.Context, err *errcode.Error) error {
	response := map[string]interface{}{"code": err.Code, "msg": err.Msg}
	details := err.Details
	if len(details) > 0 {
		response["details"] = details
	}
	return ctx.JSON(err.StatusCode(), response)
}
