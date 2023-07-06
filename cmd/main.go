package main

import (
	"cardsService/internal/config"
	"cardsService/internal/repository/mongo"
	"cardsService/internal/repository/postgres"
	srv "cardsService/internal/server"
	"cardsService/internal/token"
	"context"
	"fmt"
	"log"
)

func main() {
	cfg := config.GetConfigs()
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.FirstDatabaseConfig.User, cfg.FirstDatabaseConfig.Password, cfg.FirstDatabaseConfig.Host, cfg.FirstDatabaseConfig.DbName)

	serviceRepository, err := postgres.New(&postgres.InitParams{
		DBConnection: connStr,
		DBType:       cfg.FirstDatabaseConfig.DBConnectionType,
	})
	if err != nil {
		log.Fatalf("error init postgres repository, err: %+v", err)
	}

	mongoClient, err := mongo.New(&mongo.InitParams{
		Host: cfg.SecondDatabaseConfig.Host,
		Port: cfg.SecondDatabaseConfig.Port,
	}, context.Background())
	if err != nil {
		log.Fatalf("error init mongo client, err: %+v", err)
	}

	tokenService := token.New(config.PublicKey, config.PrivateKey)

	server := srv.New(cfg.Connection.HTTPPort, serviceRepository, mongoClient, tokenService, cfg.Connection.MasterPassword)
	server.StartServer()
}
