package main

import (
	"context"
	"home-test-app/handlers"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var (
	dbPool *pgxpool.Pool
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Initialize PostgreSQL connection
	dbPool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to PostgreSQL: %v", err)
	}
	defer dbPool.Close()

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc("/items", handlers.GetItems(dbPool)).Methods("GET")
	router.HandleFunc("/item", handlers.CreateItem(dbPool)).Methods("POST")
	router.HandleFunc("/item/{id}", handlers.UpdateItem(dbPool)).Methods("PUT")
	router.HandleFunc("/item/{id}", handlers.DeleteItem(dbPool)).Methods("DELETE")
	router.HandleFunc("/ping", handlers.Ping()).Methods("GET")

	log.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
