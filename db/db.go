package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	log.Println("Attempting to connect to PostgreSQL database")
	connStr := "user=user password=password dbname=recipe_db sslmode=disable host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, err
	}

	log.Println("Pinging database to verify connection")
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v", err)
		db.Close()
		return nil, err
	}

	log.Println("Database connection established successfully")
	return db, nil
}
