package handlers

import (
	"database/sql"
	"net/http"
)

func GetRecipes(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		w.Write([]byte("Recipes fetched!")) // Заменим на реальные данные позже
	}
}
