package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"zg3.net-api/internal/app/auth"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	database "zg3.net-api/internal/service"
)

// Config struct to hold all configuration
type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Username string `json:"username"`
		Password string `json:"password"`
		Schema   string `json:"schema"`
		Database string `json:"database"`
	} `json:"database"`
}

// main function to load configuration and print it
func main() {
	// Load configuration from file
	configFileName := "./config/config.json"

	if !checkFileExists(configFileName) {
		log.Fatalf("Configuration not found: %s.\n", configFileName)
	}
	cfg, err := loadConfig(configFileName)
	if err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	// Use the loaded configuration
	fmt.Printf("Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("Database Port: %d\n", cfg.Database.Port)
	fmt.Printf("Database User: %s\n", cfg.Database.Username)
	fmt.Printf("Database Name: %s\n", cfg.Database.Database)
	fmt.Printf("Database Schema: %s\n", cfg.Database.Schema)

	jsonData, err := json.MarshalIndent(cfg, "", "    ")
	if err != nil {
		log.Fatal("Error formating configuration data:", err)
	}
	fmt.Println(string(jsonData))

	db, err := database.Connect(cfg.Database)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Setup API Routes
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"}) // Trust only localhost

	// Middleware to inject db into context
	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.POST("/login", auth.Login)
	//router.GET("/files", handler.AuthenticateUser)

	router.Run(":8080")

}

func checkFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}

// loadConfig reads a config file and decodes it to Config struct
func loadConfig(path string) (Config, error) {
	var config Config
	file, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}
