package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Recipe представляет структуру рецепта
type Recipe struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Protein      float64 `json:"protein"`
	Fat          float64 `json:"fat"`
	Carbs        float64 `json:"carbs"`
	Calories     int     `json:"calories"`
	Instructions string  `json:"instructions"`
}

// GetRecipes возвращает список рецептов
func GetRecipes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logrus.Info("Received request for /recipes endpoint")

		logrus.Debug("Querying database for recipes")
		rows, err := db.Query("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes")
		if err != nil {
			logrus.WithError(err).Error("Failed to query database")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		logrus.Debug("Database query successful")

		logrus.Debug("Processing query results")
		var recipes []Recipe
		for rows.Next() {
			var r Recipe
			if err := rows.Scan(&r.ID, &r.Name, &r.Protein, &r.Fat, &r.Carbs, &r.Calories, &r.Instructions); err != nil {
				logrus.WithError(err).Error("Failed to scan row")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			recipes = append(recipes, r)
		}
		logrus.Info("Processed ", len(recipes), " recipes")

		logrus.Debug("Setting response headers")
		w.Header().Set("Content-Type", "application/json")
		logrus.Debug("Encoding response as JSON")
		if err := json.NewEncoder(w).Encode(recipes); err != nil {
			logrus.WithError(err).Error("Failed to encode response")
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		logrus.Info("Response sent successfully")
	}
}
