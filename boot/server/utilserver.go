package main

import (
	"fmt"
	"log"
	"net/http"
	"utilserver/pkg/clients"
	"utilserver/pkg/endpoint"
	"utilserver/pkg/spotify"
	"utilserver/pkg/storage"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
		log.Fatal("Error loading .env file")
	}
}

func main() {
	repository, err := storage.New()
	if err != nil {
		log.Fatal(err)
	}
	httpClient := clients.New(5)
	spotifyAuthService := spotify.New(repository, httpClient)

	router := endpoint.Handler(spotifyAuthService)

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8080", router))
}
