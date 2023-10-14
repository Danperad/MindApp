package main

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"time"
)

type Note struct {
	ID         primitive.ObjectID `bson:"_id"`
	UserId     uuid.UUID          `bson:"user_id"`
	Title      string             `bson:"title"`
	Date       time.Time          `bson:"date"`
	LastEdit   time.Time          `bson:"last_edit"`
	Text       string             `bson:"text"`
	IsEditable bool               `bson:"is_editable"`
}

type DiaryServer struct {
	collection *mongo.Collection
	ctx        context.Context
}

func (c Connection) CreateServer() *DiaryServer {
	return &DiaryServer{collection: c.collection}
}

func (s DiaryServer) GetAllNotes(userId uuid.UUID) (*[]Note, error) {
	filter := bson.D{{"user_id", userId}}
	var notes []Note
	cursor, err := s.collection.Find(s.ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(s.ctx, &notes); err != nil {
		return nil, err
	}
	return &notes, nil
}

func (s DiaryServer) getAllUserNotes(w http.ResponseWriter, req *http.Request) {

}
