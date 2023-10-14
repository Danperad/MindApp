package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	connection, err := connectMongoDb()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.cancel()
	defer func() {
		if err := connection.client.Disconnect(connection.ctx); err != nil {
			log.Fatal(err)
		}
	}()
	port := os.Getenv("SERVER_PORT")
	if len(port) == 0 {
		port = "8080"
	}
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := connection.CreateServer()
	router.HandleFunc("/diary/", server.getAllUserNotes).Methods("GET")
	http.Handle("/", router)
}
