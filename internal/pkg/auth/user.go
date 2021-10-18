package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/joeyscat/object-storage-go/pkg/log"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

const (
	HEADER_USER_KEY = "OS-USER"
)

type User struct {
	UserID   string
	UserName string
}

func GetUser(c echo.Context) (*User, error) {
	userStr := c.Request().Header.Get(HEADER_USER_KEY)
	if strings.TrimSpace(userStr) == "" {
		return nil, fmt.Errorf("request header [%s] not found", HEADER_USER_KEY)
	}

	var u User
	err := json.Unmarshal([]byte(userStr), &u)
	if err != nil {
		return nil, err
	}

	// check user fields
	log.Debug("GetUser", zap.Any("user", u))
	if u.UserID == "" || u.UserName == "" {
		return nil, errors.New("parse user from header failed")
	}

	return &u, nil
}

func UserInfoNotFoundInRequest(c echo.Context) error {
	return c.JSON(http.StatusUnauthorized, map[string]string{"msg": "user info not found in request"})
}
