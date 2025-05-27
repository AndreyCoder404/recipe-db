package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetRecipes(t *testing.T) {
	// Создаём мок базы данных
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	// Тест успешного случая
	t.Run("successful case", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "protein", "fat", "carbs", "calories", "instructions"}).
			AddRow(1, "Test Recipe", 10.0, 5.0, 20.0, 200, "Cook for 10 minutes")
		mock.ExpectQuery("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes").
			WillReturnRows(rows)

		req, err := http.NewRequest("GET", "/recipes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := GetRecipes(db)

		handler(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		var recipes []Recipe
		if err := json.NewDecoder(rr.Body).Decode(&recipes); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(recipes) != 1 {
			t.Errorf("expected 1 recipe, got %d", len(recipes))
		}
		if recipes[0].Name != "Test Recipe" {
			t.Errorf("expected name 'Test Recipe', got '%s'", recipes[0].Name)
		}
	})

	// Тест ошибки базы данных
	t.Run("database error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes").
			WillReturnError(sql.ErrNoRows)

		req, err := http.NewRequest("GET", "/recipes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := GetRecipes(db)

		handler(rr, req)

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
