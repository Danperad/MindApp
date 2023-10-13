package main

import (
	"github.com/elastic/go-elasticsearch/v8"
	"log"
	"os"
)

func getElasticClient() (*elasticsearch.Client, error) {
	server := os.Getenv("ELASTIC_URL")
	username := os.Getenv("ELASTIC_USER")
	password := os.Getenv("ELASTIC_PASSWORD")
	cfg := elasticsearch.Config{
		Addresses: []string{server},
		Username:  username,
		Password:  password,
	}
	return elasticsearch.NewClient(cfg)
}

func main() {
	client, err := getElasticClient()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// TODO: Close connection
	}()
}
