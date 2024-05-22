package main

import (
	"myapi/internal/app/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{"127.0.0.1"}) // Trust only localhost

	// Setup route
	router.GET("/user", handler.GetUser)
	router.GET("/product", handler.GetProduct)

	// Start server
	router.Run(":8080")
}
