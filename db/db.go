package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	connStr := "user=user password=password dbname=recipe_db sslmode=disable host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
