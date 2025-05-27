package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
)

func TestGetRecipes(t *testing.T) {
	log.Println("Starting TestGetRecipes")

	// Создаём мок базы данных
	log.Println("Creating sqlmock database")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	log.Println("sqlmock database created successfully")

	// Тест успешного случая
	t.Run("successful case", func(t *testing.T) {
		log.Println("Running successful case test")
		rows := sqlmock.NewRows([]string{"id", "name", "protein", "fat", "carbs", "calories", "instructions"}).
			AddRow(1, "Test Recipe", 10.0, 5.0, 20.0, 200, "Cook for 10 minutes")
		mock.ExpectQuery("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes").
			WillReturnRows(rows)
		log.Println("Set up mock expectations for successful case")

		log.Println("Creating HTTP request")
		req, err := http.NewRequest("GET", "/recipes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := GetRecipes(db)
		log.Println("Executing handler")
		handler(rr, req)

		log.Printf("Checking status code: got %v, expected %v", rr.Code, http.StatusOK)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		log.Println("Decoding response body")
		var recipes []Recipe
		if err := json.NewDecoder(rr.Body).Decode(&recipes); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		log.Printf("Decoded %d recipes", len(recipes))
		if len(recipes) != 1 {
			t.Errorf("expected 1 recipe, got %d", len(recipes))
		}
		if recipes[0].Name != "Test Recipe" {
			t.Errorf("expected name 'Test Recipe', got '%s'", recipes[0].Name)
		}
		log.Println("Successful case test passed")
	})

	// Тест ошибки базы данных
	t.Run("database error", func(t *testing.T) {
		log.Println("Running database error test")
		mock.ExpectQuery("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes").
			WillReturnError(sql.ErrNoRows)
		log.Println("Set up mock expectations for database error")

		log.Println("Creating HTTP request")
		req, err := http.NewRequest("GET", "/recipes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := GetRecipes(db)
		log.Println("Executing handler")
		handler(rr, req)

		log.Printf("Checking status code: got %v, expected %v", rr.Code, http.StatusInternalServerError)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		log.Println("Database error test passed")
	})

	log.Println("Verifying mock expectations")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
	log.Println("TestGetRecipes completed")
}

func TestGetRecipesIntegration(t *testing.T) {
	log.Println("Starting TestGetRecipesIntegration")

	// Пропускаем тест, если база недоступна
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}
	log.Println("Integration test mode enabled")

	log.Println("Attempting to connect to real database")
	db, err := sql.Open("postgres", "user=user password=password dbname=recipe_db sslmode=disable host=localhost")
	if err != nil {
		t.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("Successfully connected to real database")

	log.Println("Creating HTTP request")
	req, err := http.NewRequest("GET", "/recipes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := GetRecipes(db)
	log.Println("Executing handler with real database")
	handler(rr, req)

	log.Printf("Checking status code: got %v, expected %v", rr.Code, http.StatusOK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	log.Println("Decoding response body")
	var recipes []Recipe
	if err := json.NewDecoder(rr.Body).Decode(&recipes); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	log.Printf("Decoded %d recipes", len(recipes))
	if len(recipes) != 1 {
		t.Errorf("expected 1 recipe, got %d", len(recipes))
	}
	log.Println("TestGetRecipesIntegration completed")
}
