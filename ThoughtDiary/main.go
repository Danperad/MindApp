package main

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

type connection struct {
	client     *mongo.Client
	collection *mongo.Collection
	ctx        context.Context
	cancel     context.CancelFunc
}

type Note struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserId     uuid.UUID          `bson:"user_id"`
	Title      string             `bson:"title"`
	Date       time.Time          `bson:"date"`
	LastEdit   time.Time          `bson:"last_edit"`
	Text       string             `bson:"text"`
	IsEditable bool               `bson:"is_editable"`
}

var connect connection

func connectMongoDb() error {
	serverUrl := os.Getenv("MONGO_URL")
	if len(serverUrl) == 0 {
		return errors.New("MONGO_URL environment does not set")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	log.Println("Try connect to mongodb database")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(serverUrl))
	if err != nil {
		cancel()
		return err
	}
	log.Println("Check connection mongodb database")
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		err := client.Disconnect(ctx)
		if err != nil {
			cancel()
			return err
		}
		cancel()
		return err
	}
	log.Println("Connection successful")
	databaseName := os.Getenv("MONGO_DATABASE")
	if len(databaseName) == 0 {
		cancel()
		return errors.New("MONGO_DATABASE environment does not set")
	}
	collection := client.Database(databaseName).Collection("diaries")
	connect = connection{cancel: cancel, ctx: ctx, client: client, collection: collection}
	return nil
}

func main() {
	if err := connectMongoDb(); err != nil {
		log.Fatal(err)
	}
	defer connect.cancel()
	defer func() {
		if err := connect.client.Disconnect(connect.ctx); err != nil {
			log.Fatal(err)
		}
	}()

}
