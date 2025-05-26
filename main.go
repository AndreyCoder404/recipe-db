package main

import (
	"github.com/AndreyCoder404/recipe-db/db"
	"github.com/AndreyCoder404/recipe-db/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/recipes", handlers.GetRecipes(db))

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
