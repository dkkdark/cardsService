package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type ServiceImpl struct {
	client *mongo.Client
	c      context.Context
}

type InitParams struct {
	Host string
	Port string
}

func New(params *InitParams, context context.Context) (*ServiceImpl, error) {
	serviceImpl := &ServiceImpl{}
	log.Println("Starting connect to MongoDB")

	connMongoStr := fmt.Sprintf("mongodb://%s:%s", params.Host, params.Port)

	clientOptions := options.Client().ApplyURI(connMongoStr)
	client, err := mongo.Connect(context, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// check connection here
	err = client.Ping(context, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")

	serviceImpl.client = client
	serviceImpl.c = context

	return serviceImpl, nil
}
