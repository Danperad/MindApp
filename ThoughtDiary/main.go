package main

import (
	"fmt"
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
	if port == "" {
		port = "8080"
	}
	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/diary/", connection.getAllUserNotes).Methods("GET")
	log.Printf("Listen on :%v\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), router))
}
