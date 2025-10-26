package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func SetDB() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	userDB := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, userDB, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("error al conectar con la base de datos: %w", err)
	}

	return db, nil
}
