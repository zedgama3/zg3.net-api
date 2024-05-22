package database

import (
	"database/sql"
	"fmt"
	"log"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
	Database string `json:"database"`
}

// Connects to a PostgreSQL database and returns a handle to that database.
func Connect(cfg DatabaseConfig) (*sql.DB, error) {
	fmt.Println(cfg)

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	//defer db.Close()

	// Check the database connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	fmt.Println("Successfully connected to the database!")

	//TODO: Error handling.  Pass error details to calling function.
	return db, nil
}
