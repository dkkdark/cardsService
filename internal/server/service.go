package server

import (
	"cardsService/internal/repository/mongo"
	"cardsService/internal/repository/postgres"
	"cardsService/internal/token"
	"github.com/labstack/echo/v4"
)

type ServiceImpl struct {
	server         *echo.Echo
	httpPort       string
	Repository     *postgres.ServiceImpl
	MongoClient    *mongo.ServiceImpl
	TokenService   *token.ServiceImpl
	masterPassword string
}

func New(httpPort string, repository *postgres.ServiceImpl, client *mongo.ServiceImpl,
	tokenService *token.ServiceImpl,
	masterPassword string) *ServiceImpl {
	service := &ServiceImpl{
		httpPort:       httpPort,
		Repository:     repository,
		MongoClient:    client,
		TokenService:   tokenService,
		masterPassword: masterPassword,
	}

	service.server = echo.New()
	service.setupRoutes()
	return service
}

func (s *ServiceImpl) StartServer() {
	s.server.Logger.Fatal(s.server.Start(":" + s.httpPort))
}
