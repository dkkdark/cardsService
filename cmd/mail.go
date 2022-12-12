package main

import (
	"cardsService/config"
	"cardsService/repository"
	srv "cardsService/server"
	"cardsService/token"
	"log"
)

func main() {
	serviceRepository, err := repository.New(&repository.InitParams{
		DBConnection: config.DBConnection,
		DBType:       config.DBConnectionType,
	})
	if err != nil {
		log.Fatalf("error init repository, err: %+v", err)
	}

	tokenService := token.New(config.PublicKey, config.PrivateKey)

	server := srv.New(config.HTTPPort, serviceRepository, tokenService, config.MasterPassword)
	server.StartServer()
}
