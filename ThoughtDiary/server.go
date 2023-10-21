package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"net/http"
	"strings"
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

func (s Connection) GetAllNotes(userId uuid.UUID) (*[]Note, error) {
	filter := bson.D{{"user_id", userId}}
	var notes []Note
	collection := s.client.Database("thought_diaries").Collection("diaries")
	cursor, err := collection.Find(s.ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(s.ctx, &notes); err != nil {
		return nil, err
	}
	return &notes, nil
}

func (s Connection) checkAndParseToken(token string) (*uuid.UUID, error) {
	jwtToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.signKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("error on parse authorization token")
	}
	res, err := uuid.Parse(fmt.Sprintf("%v", claims["user_id"]))
	return &res, err
}

func (s Connection) getAllUserNotes(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	authorization := req.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	authorization = strings.TrimPrefix(authorization, "Bearer ")
	token, err := s.checkAndParseToken(authorization)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	notes, err := s.GetAllNotes(*token)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(notes)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
