package auth

import "github.com/gin-gonic/gin"

type User struct {
	UserId      string  `json:"user_id"`
	Username    string  `json:"username"`
	Password    *string `json:"password,omitempty"`
	Status      string  `json:"status"`
	Permissions []byte  `json:"permissions"`
}

type HttpError struct {
	Code          int   // HTTP response code
	Message       gin.H // Message body to be returned
	InternalError error // Internal error to be logged - may contain sensitive data
}
