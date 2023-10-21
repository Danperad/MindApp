package main

import (
	"context"
	"crypto/rsa"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"os"
	"strings"
)

var db *pgx.Conn
var signKey *rsa.PrivateKey

type (
	SignInModel struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	User struct {
		UserId   uint
		Login    string
		Password string
	}
	AuthTokenModel struct {
		AuthToken string `json:"auth_token"`
	}
)

func signIn(w http.ResponseWriter, r *http.Request) {
	var signin SignInModel
	err := json.NewDecoder(r.Body).Decode(&signin)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	var user User
	err = db.QueryRow(context.Background(), "SELECT * FROM users WHERE login=$1 AND password=$2", signin.Login, signin.Password).Scan(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"user_id": user.UserId,
	})
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	authModel := AuthTokenModel{AuthToken: tokenString}
	err = json.NewEncoder(w).Encode(authModel)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getPublicKey(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	header = strings.TrimPrefix(header, "Bearer ")
}

func main() {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		log.Fatal("Environment DATABASE_URL does not set")
	}
	db, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer func(db *pgx.Conn, ctx context.Context) {
		err := db.Close(ctx)
		if err != nil {
			log.Fatalf("Unable to disconnect from database: %v\n", err)
		}
	}(db, context.Background())

	router := mux.NewRouter()
	router.HandleFunc("auth/signin", signIn).Methods("POST")
	router.HandleFunc("data/publickey", getPublicKey).Methods("GET").Headers("Authorization")
	err = http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
