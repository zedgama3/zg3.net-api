package main

import (
	"encoding/json"
	"fmt"
	"log"

	"zg3.net-api/internal/app/auth"
	"zg3.net-api/internal/config"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	database "zg3.net-api/internal/service"
)

type Config = config.Config

// main function to load configuration and print it
func main() {

	// Load configuration from file
	var cfg Config
	if newCfg, err := config.New("./config/config.json"); err != nil {
		log.Fatal("error reading config,", err)
	} else {
		cfg = *newCfg
	}

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
