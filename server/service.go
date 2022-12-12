package server

import (
	"cardsService/repository"
	"cardsService/token"
	"github.com/labstack/echo/v4"
)

type ServiceImpl struct {
	server         *echo.Echo
	httpPort       string
	repository     *repository.ServiceImpl
	tokenService   *token.ServiceImpl
	masterPassword string
}

func New(httpPort string, repository *repository.ServiceImpl,
	tokenService *token.ServiceImpl,
	masterPassword string) *ServiceImpl {
	service := &ServiceImpl{
		httpPort:       httpPort,
		repository:     repository,
		tokenService:   tokenService,
		masterPassword: masterPassword,
	}

	service.server = echo.New()
	service.setupRoutes()
	return service
}

func (s *ServiceImpl) StartServer() {
	s.server.Logger.Fatal(s.server.Start(":" + s.httpPort))
}
