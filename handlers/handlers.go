package handlers

import (
	"database/sql"
	"encoding/json"
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
		rows, err := db.Query("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var recipes []Recipe
		for rows.Next() {
			var r Recipe
			if err := rows.Scan(&r.ID, &r.Name, &r.Protein, &r.Fat, &r.Carbs, &r.Calories, &r.Instructions); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			recipes = append(recipes, r)
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(recipes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
