package main

import (
	"github.com/AndreyCoder404/recipe-db/db"
	"github.com/AndreyCoder404/recipe-db/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel) // Устанавливаем уровень логирования (Debug для детализации)
	logrus.Debug("Starting application...")

	db, err := db.Connect()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database")
	}
	logrus.Info("Successfully connected to database")
	defer func() {
		logrus.Info("Closing database connection")
		db.Close()
	}()

	logrus.Info("Setting up router")
	r := mux.NewRouter()
	r.HandleFunc("/recipes", handlers.GetRecipes(db))
	logrus.Info("Router setup complete, registered /recipes endpoint")

	logrus.Info("Server running on :8080")
	logrus.Fatal(http.ListenAndServe(":8080", r))
}
