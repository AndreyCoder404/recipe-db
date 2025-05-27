package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel) // Устанавливаем уровень логирования для тестов
}

func TestGetRecipes(t *testing.T) {
	logrus.Debug("Starting TestGetRecipes")

	// Создаём мок базы данных
	logrus.Debug("Creating sqlmock database")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()
	logrus.Debug("sqlmock database created successfully")

	// Тест успешного случая
	t.Run("successful case", func(t *testing.T) {
		logrus.Debug("Running successful case test")
		rows := sqlmock.NewRows([]string{"id", "name", "protein", "fat", "carbs", "calories", "instructions"}).
			AddRow(1, "Test Recipe", 10.0, 5.0, 20.0, 200, "Cook for 10 minutes")
		mock.ExpectQuery("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes").
			WillReturnRows(rows)
		logrus.Debug("Set up mock expectations for successful case")

		logrus.Debug("Creating HTTP request")
		req, err := http.NewRequest("GET", "/recipes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := GetRecipes(db)
		logrus.Debug("Executing handler")
		handler(rr, req)

		logrus.Debugf("Checking status code: got %v, expected %v", rr.Code, http.StatusOK)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		logrus.Debug("Decoding response body")
		var recipes []Recipe
		if err := json.NewDecoder(rr.Body).Decode(&recipes); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		logrus.Debugf("Decoded %d recipes", len(recipes))
		if len(recipes) != 1 {
			t.Errorf("expected 1 recipe, got %d", len(recipes))
		}
		if recipes[0].Name != "Test Recipe" {
			t.Errorf("expected name 'Test Recipe', got '%s'", recipes[0].Name)
		}
		logrus.Info("Successful case test passed")
	})

	// Тест ошибки базы данных
	t.Run("database error", func(t *testing.T) {
		logrus.Debug("Running database error test")
		mock.ExpectQuery("SELECT id, name, protein, fat, carbs, calories, instructions FROM recipes").
			WillReturnError(sql.ErrNoRows)
		logrus.Debug("Set up mock expectations for database error")

		logrus.Debug("Creating HTTP request")
		req, err := http.NewRequest("GET", "/recipes", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler := GetRecipes(db)
		logrus.Debug("Executing handler")
		handler(rr, req)

		logrus.Debugf("Checking status code: got %v, expected %v", rr.Code, http.StatusInternalServerError)
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
		}
		logrus.Info("Database error test passed")
	})

	logrus.Debug("Verifying mock expectations")
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
	logrus.Info("TestGetRecipes completed")
}

func TestGetRecipesIntegration(t *testing.T) {
	logrus.Debug("Starting TestGetRecipesIntegration")

	// Пропускаем тест, если база недоступна
	if testing.Short() {
		logrus.Warn("skipping integration test in short mode")
		t.Skip("skipping integration test in short mode")
	}
	logrus.Info("Integration test mode enabled")

	logrus.Debug("Attempting to connect to real database")
	db, err := sql.Open("postgres", "user=user password=password dbname=recipe_db sslmode=disable host=localhost")
	if err != nil {
		logrus.WithError(err).Fatal("failed to connect to database")
	}
	defer db.Close()
	logrus.Info("Successfully connected to real database")

	logrus.Debug("Creating HTTP request")
	req, err := http.NewRequest("GET", "/recipes", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := GetRecipes(db)
	logrus.Debug("Executing handler with real database")
	handler(rr, req)

	logrus.Debugf("Checking status code: got %v, expected %v", rr.Code, http.StatusOK)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	logrus.Debug("Decoding response body")
	var recipes []Recipe
	if err := json.NewDecoder(rr.Body).Decode(&recipes); err != nil {
		logrus.WithError(err).Fatal("failed to decode response")
	}
	logrus.Debugf("Decoded %d recipes", len(recipes))
	if len(recipes) != 1 {
		t.Errorf("expected 1 recipe, got %d", len(recipes))
	}
	logrus.Info("TestGetRecipesIntegration completed")
}
