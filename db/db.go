package db

import (
	"database/sql"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	logrus.Debug("Attempting to connect to PostgreSQL database")
	connStr := "user=user password=password dbname=recipe_db sslmode=disable host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithError(err).Error("Failed to open database connection")
		return nil, err
	}

	logrus.Debug("Pinging database to verify connection")
	if err := db.Ping(); err != nil {
		logrus.WithError(err).Error("Failed to ping database")
		db.Close()
		return nil, err
	}

	logrus.Info("Database connection established successfully")
	return db, nil
}
