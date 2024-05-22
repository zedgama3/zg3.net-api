package handler

import (
	"myapi/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser handles the GET request to retrieve a user
func GetUser(c *gin.Context) {
	// For demonstration, creating a user instance statically
	user := model.NewUser(1, "John Doe", "johndoe@example.com")

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
