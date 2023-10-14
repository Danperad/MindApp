package main

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type Connection struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
	cancel     context.CancelFunc
}

func connectMongoDb() (*Connection, error) {
	serverUrl := os.Getenv("MONGO_URL")
	if serverUrl == "" {
		return nil, errors.New("MONGO_URL environment does not set")
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
	collection := client.Database("thought_diary").Collection("notes")
	connect := Connection{cancel: cancel, ctx: ctx, client: client, collection: collection}
	return &connect, nil
}
