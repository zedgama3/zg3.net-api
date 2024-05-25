/*
	Using this package:
	1. Set the secret key using `SetSecret(secret string)`
	  - These can be generated using `openssl rand -base64 128`
	2. Call a method to return a user object
	  - PasswordLogin
*/

package user

import (
	"crypto/subtle"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

// Package level unexported JWT secret.

var config Config

// Exported function to allow setting the JWT secret.  This must be at least 32 characters.
func SetConfig(newConfig Config) {
	config = newConfig
}

// HTTP endpoint handler
func Login(c *gin.Context) {

	// Check that needed components are present
	db, exists := c.MustGet("db").(*sql.DB)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Handle not found"})
	}
	if len(config.Secret) < 32 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Missing secret"})
	}

	// Handle request based on requested grant type
	switch c.PostForm("grant_type") {
	case "password":

		// Get user data and confirm password
		var user User
		if newUser, err := PasswordLogin(db, c.PostForm("username"), c.PostForm("password")); err != nil {
			c.JSON(err.Code, gin.H{"error": err.Message})
		} else {
			user = *newUser
		}

		// Create and return token
		if token, err := user.NewToken(); err != nil {
			fmt.Printf("Error creating token: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create token"})
		} else {
			c.JSON(http.StatusOK, gin.H{"token": token})
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
	}
}

// Authenticate a user based on their username and password and return a User
func PasswordLogin(db *sql.DB, username string, password string) (*User, *HttpError) {
	/*
		TODO:
		- Validate password
		- Populate User struct
	*/
	u := User{
		Username: username,
	}

	if err := dbGetUser(db, &u); err != nil {
		return nil, err
		//TODO: Process error
	}

	if isMatch, err := verifyPassword(password, u.passwordHash); err != nil {
		return nil, err
	} else if !isMatch {
		return nil, &HttpError{http.StatusUnauthorized, gin.H{"error": "User is not authorized"}, errors.New("Password verification failed.")}
	}

	return &u, nil
}

// Generate a new JWT access token
func (u *User) NewToken() (*Token, error) {

	// Create the token
	token := jwt.New(jwt.SigningMethodES256)

	// Set claims - Notes starting with * are required.
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "https://zg3.net"                    // *Issuer
	claims["sub"] = "home-1"                             // *Subject
	claims["aud"] = ""                                   // *Audience
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // *Expiration Time
	claims["iat"] = time.Now().Unix()                    // *Issued At
	claims["auth_time"] = time.Now().Unix()              // *Authentication Time
	claims["amr"] = []string{"pwd"}                      // Authentication methods reference
	claims["name"] = ""
	claims["admin"] = u.Username == "zedgama3"
	claims["authorized"] = true
	claims["username"] = u.Username
	claims["permissions"] = u.Permissions
	claims["user_id"] = u.UserId

	// Decode the base64 encoded private key
	keyData, err := base64.StdEncoding.DecodeString(config.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %v", err)
	}

	block, _ := pem.Decode(keyData)
	if block == nil || block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing the private key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse EC private key: %v", err)
	}

	// Sign the token with the ECDSA private key
	var myToken Token
	if tokenString, err := token.SignedString(privateKey); err != nil {
		return nil, fmt.Errorf("unable to sign JWT: %v", err)
	} else {
		myToken = Token(tokenString)
	}

	return &myToken, nil
}

// Create a whole user from token
// Is this really useful???  Should it be here?
func (u *User) TokenLogin(t *Token) (*User, error) {
	return nil, nil
}

func dbGetUser(db *sql.DB, user *User) *HttpError {

	// Querying the database
	query := `SELECT user_id, username, password, status, permissions FROM api.users WHERE username = $1 LIMIT 1;`
	rows, err := db.Query(query, user.Username)
	if err != nil {
		return &HttpError{http.StatusInternalServerError, gin.H{"error": "Unable to connect to database"}, err}
	}
	defer rows.Close()

	// Iterate through the result set
	if rows.Next() {
		err := rows.Scan(&user.UserId, &user.Username, &user.passwordHash, &user.Status, &user.Permissions)
		if err != nil {
			return &HttpError{http.StatusInternalServerError, gin.H{"error": "Database error."}, err}
		}
		//c.JSON(http.StatusOK, user)
	} else {
		return &HttpError{http.StatusUnauthorized, gin.H{"error": "User not authorized."}, errors.New("User not found")}
	}

	return nil
}

// Check if plaintext password matches stored hash.
func verifyPassword(password, encodedHash string) (isMatch bool, err *HttpError) {

	// Parse the hash string
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, &HttpError{http.StatusInternalServerError, gin.H{"error": "invalid hash format"}, errors.New("invalid hash format")}
	}
	var memory, time, threads int
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return false, &HttpError{http.StatusInternalServerError, gin.H{"error": "invalid hash format"}, err}
	}

	// Generate a new hash for testing.
	var salt []byte
	if newSalt, err := base64.RawStdEncoding.DecodeString(parts[4]); err != nil {
		return false, &HttpError{}
	} else {
		salt = newSalt
	}

	var hash []byte
	if newHash, err := base64.RawStdEncoding.DecodeString(parts[5]); err != nil {
		return false, &HttpError{http.StatusInternalServerError, gin.H{"error": "Unable to build hash."}, err}
	} else {
		hash = newHash
	}

	// Generate a hash with the same parameters and salt
	otherHash := argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(threads), uint32(32)) // Assuming keyLen is constant

	// Compare the hashes
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}
