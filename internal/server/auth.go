package server

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func (s *ServiceImpl) GetAuthToken(c echo.Context) string {
	auth := c.Request().Header.Get("authorization")
	return strings.Replace(auth, "Bearer ", "", -1)
}
