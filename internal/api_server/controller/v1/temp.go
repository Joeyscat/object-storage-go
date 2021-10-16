package v1

import "github.com/labstack/echo/v4"

type TempController struct {
}

func NewTempController() *TempController {
	return &TempController{}
}

func (t *TempController) HeadTempObject(c echo.Context) error {
	panic("xx")
}

func (t *TempController) PutTempObject(c echo.Context) error {
	panic("xx")
}
