package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"zg3.net-api/internal/user"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

func main() {

	// Load configuration from file
	cfg, err := newConfig("./config/config.json")
	if err != nil {
		log.Fatal("error reading config,", err)
	}

	// Setup database connection
	db, err := newDatabase(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Setup API Routes
	router := mux.NewRouter()

	// Middleware to inject db into context
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set db in context
			ctx := context.WithValue(r.Context(), "db", db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	})

	// User endpoints

	user.SetConfig(cfg.User)
	router.HandleFunc("/login", user.Login).Methods("POST")

	err = http.ListenAndServe(":8012", router)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

}

// Connects to a PostgreSQL database and returns a handle to that database.
func newDatabase(cfg DatabaseConfig) (*sql.DB, error) {

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

// Read the file listed as "f" and return a config object
func newConfig(f string) (*Config, error) {
	var c Config

	if _, err := os.Stat(f); err != nil {
		return nil, fmt.Errorf("file does not exist: %v", f)
	}

	if file, err := os.ReadFile(f); err != nil {
		return nil, fmt.Errorf("unable to read config file: %v", err)
	} else {
		if err := json.Unmarshal(file, &c); err != nil {
			return nil, fmt.Errorf("error parsing json: %v", err)
		} else {
			return &c, nil
		}
	}
}
