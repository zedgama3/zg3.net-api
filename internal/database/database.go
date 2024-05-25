package database

import (
	"database/sql"
	"fmt"
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
func New(cfg DatabaseConfig) (*sql.DB, error) {

	// Connect to database
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Database)

	var db *sql.DB
	if newDb, err := sql.Open("postgres", connStr); err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	} else {
		db = newDb
	}
	//defer db.Close()

	// Check the database connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	return db, nil
}
