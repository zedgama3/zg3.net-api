package handler

import (
	"myapi/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProduct handles the GET request to retrieve a product
func GetProduct(c *gin.Context) {
	// For demonstration, creating a product instance statically
	product := model.NewProduct(1, "Example Product", 19.99)

	c.JSON(http.StatusOK, gin.H{
		"id":    product.ID,
		"name":  product.Name,
		"price": product.Price,
	})
}
