package main

import (
	"cardsService/internal/config"
	"cardsService/internal/repository"
	srv "cardsService/internal/server"
	"cardsService/internal/token"
	"fmt"
	"log"
)

func main() {
	cfg := config.GetConfigs()
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.FirstDatabaseConfig.User, cfg.FirstDatabaseConfig.Password, cfg.FirstDatabaseConfig.Host, cfg.FirstDatabaseConfig.DbName)

	serviceRepository, err := repository.New(&repository.InitParams{
		DBConnection: connStr,
		DBType:       cfg.FirstDatabaseConfig.DBConnectionType,
	})
	if err != nil {
		log.Fatalf("error init repository, err: %+v", err)
	}

	tokenService := token.New(config.PublicKey, config.PrivateKey)

	server := srv.New(cfg.Connection.HTTPPort, serviceRepository, tokenService, cfg.Connection.MasterPassword)
	server.StartServer()
}
