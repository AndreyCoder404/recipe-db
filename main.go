package main

import (
	"github.com/AndreyCoder404/recipe-db/db"
	"github.com/AndreyCoder404/recipe-db/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	log.Println("Starting application...")

	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Successfully connected to database")
	defer func() {
		log.Println("Closing database connection")
		db.Close()
	}()

	log.Println("Setting up router")
	r := mux.NewRouter()
	r.HandleFunc("/recipes", handlers.GetRecipes(db))
	log.Println("Router setup complete, registered /recipes endpoint")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
