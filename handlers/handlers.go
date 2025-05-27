package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
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
		log.Println("Received request for /recipes endpoint")

		log.Println("Querying database for recipes")
		rows, err := db.Query("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes")
		if err != nil {
			log.Printf("Failed to query database: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		log.Println("Database query successful")

		log.Println("Processing query results")
		var recipes []Recipe
		for rows.Next() {
			var r Recipe
			if err := rows.Scan(&r.ID, &r.Name, &r.Protein, &r.Fat, &r.Carbs, &r.Calories, &r.Instructions); err != nil {
				log.Printf("Failed to scan row: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			recipes = append(recipes, r)
		}
		log.Printf("Processed %d recipes", len(recipes))

		log.Println("Setting response headers")
		w.Header().Set("Content-Type", "application/json")
		log.Println("Encoding response as JSON")
		if err := json.NewEncoder(w).Encode(recipes); err != nil {
			log.Printf("Failed to encode response: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		log.Println("Response sent successfully")
	}
}
