package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"utilserver/pkg/clients"
	"utilserver/pkg/endpoint"
	"utilserver/pkg/spotify"
	"utilserver/pkg/storage"

	"github.com/gorilla/handlers"
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
	cache, err := storage.NewCache(os.Getenv("REDIS_CONNECTION_STRING"))
	if err != nil {
		panic(err)
	}
	storage, err := storage.NewStorage(os.Getenv("MONGODB_CONNECTION_STRING"), os.Getenv("MONGODB_DATABASE"))
	if err != nil {
		panic(err)
	}
	httpClient := clients.New(5)
	Services := spotify.NewServices(storage, httpClient, cache)

	router := endpoint.NewHandler(cache, Services)

	fmt.Printf("Starting server at port 8090\n")
	// allow CORS and start listening
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	log.Fatal(http.ListenAndServe(":8090", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
