package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Token string

type Config struct {
	Secret     string `json:"secret"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}

type User struct {
	UserId       string `json:"user_id"`
	Username     string `json:"username"`
	Status       string `json:"status"`
	Permissions  []byte `json:"permissions"`
	passwordHash string
}

// This should be moved to a common package.
type HttpError struct {
	Code    int   // HTTP response code
	Message gin.H // Message body to be returned
	Err     error // Internal error to be logged - may contain sensitive data
}

func (he *HttpError) Error() string {
	return fmt.Sprintf("status %d: %v - %v", he.Code, he.Message, he.Err)
}
