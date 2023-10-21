package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type Connection struct {
	client  *mongo.Client
	ctx     context.Context
	cancel  context.CancelFunc
	signKey []byte
}

func connectMongoDb() (*Connection, error) {
	serverUrl := os.Getenv("MONGO_URL")
	if serverUrl == "" {
		return nil, fmt.Errorf("MONGO_URL environment does not set")
	}
	signKey := os.Getenv("SIGN_KEY")
	if signKey == "" {
		return nil, fmt.Errorf("SIGN_KEY environment does not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("Try connect to mongodb database")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverUrl))
	if err != nil {
		cancel()
		return nil, err
	}
	log.Println("Check Connection mongodb database")
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			cancel()
			return nil, err
		}
		cancel()
		return nil, err
	}
	log.Println("Connection successful")
	connect := Connection{cancel: cancel, ctx: ctx, client: client, signKey: []byte(signKey)}
	return &connect, nil
}
