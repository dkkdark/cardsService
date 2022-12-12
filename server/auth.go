package server

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func (s *ServiceImpl) getAuthToken(c echo.Context) string {
	auth := c.Request().Header.Get("authorization")
	return strings.Replace(auth, "Bearer ", "", -1)
}
