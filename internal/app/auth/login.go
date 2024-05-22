package auth

/*
TODO:
- [ ] Read secretKey from config and use a helper function to set it so that it's not publicly readable
- [ ] Understand the difference between an ID token and an access token and if I need this much complication
  - I probably will because of the ability for third-party account impersonation - e.g., Coaches being able to access data
- [ ] Move all the helper functions somewhere else and clean up this mess of code
- [ ] Update the OpenAPI document to match implementation
*/

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

var secretKey = "MHcCAQEEIOqsEsDcG1fKjvhyQxZvBsUCORF+IFBSCANguJ5h4iWHoAoGCCqGSM49AwEHoUQDQgAE8P45+bQ8zWLyStx9oe1/VKnNaJ4XiC0eLwrEKWocwS1XZnb3hE4zZHK5GMsI6HARamQtyZ+cChJSLnHohzmiVA=="

func Login(c *gin.Context) {

	// Determine grant type
	switch c.PostForm("grant_type") {
	case "password":
		passwordLogin(c)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
	}
}

func verifyPasswordCredentials(db *sql.DB, username string, password string) (user User, httpErr *HttpError) {

	// Querying the database
	query := `SELECT user_id, username, password, status, permissions FROM api.users WHERE username = $1 LIMIT 1;`
	rows, err := db.Query(query, username)
	if err != nil {
		return User{}, &HttpError{http.StatusInternalServerError, gin.H{"error": "Unable to connect to database"}, err}
	}
	defer rows.Close()

	// Iterate through the result set
	if rows.Next() {
		err := rows.Scan(&user.UserId, &user.Username, &user.Password, &user.Status, &user.Permissions)
		if err != nil {
			return User{}, &HttpError{http.StatusInternalServerError, gin.H{"error": "Database error."}, err}
		}
		//c.JSON(http.StatusOK, user)
	} else {
		return User{}, &HttpError{http.StatusUnauthorized, gin.H{"error": "User not authorized."}, errors.New("User not found")}
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return User{}, &HttpError{http.StatusInternalServerError, gin.H{"error": "There was an error when trying to access the database.  System logs contain additional details."}, err}
	}

	if user.Password == nil {
		hash, _ := generateHash(password)
		return User{}, &HttpError{http.StatusUnauthorized, gin.H{"hash": hash}, errors.New("password not set")}
	}

	// Verify password
	match, err := verifyPassword(password, *user.Password)
	if err != nil {
		return User{}, &HttpError{http.StatusInternalServerError, gin.H{"error": "User not authorized."}, err}
	}
	if match {
		// s, err := createJWT()
		// if err != nil {
		// 	return User{}, HttpError{http.StatusInternalServerError, gin.H{"error": err}, err}
		// }
		// c.JSON(http.StatusOK, gin.H{"token": s})

		//TODO: Return User struct and generate the JWT elsewhere.
		return user, nil
	} else {
		return User{}, &HttpError{http.StatusUnauthorized, gin.H{"error": "Invalid Password"}, errors.New("invalid password")}
	}

}

func passwordLogin(c *gin.Context) {

	// Get database handle
	db, exists := c.MustGet("db").(*sql.DB)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB Handle not found"})
	}

	// Verify credentials
	user, err := verifyPasswordCredentials(db, c.PostForm("username"), c.PostForm("password"))
	if err != nil {
		c.JSON(err.Code, err.Message)
	}

	// Create JWT
	s, err := createJWT(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err})
	}
	c.JSON(http.StatusOK, gin.H{"token": s})

}

func generateHash(password string) (string, error) {
	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Argon2id parameters
	time := 1
	memory := 65536 // 64 * 1024, memory in KiB
	threads := 4
	keyLen := 32

	// Generate the hash
	hash := argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(threads), uint32(keyLen))

	// Encode salt and hash
	saltEncoded := base64.RawStdEncoding.EncodeToString(salt)
	hashEncoded := base64.RawStdEncoding.EncodeToString(hash)

	// Return formatted string with all parameters
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s", memory, time, threads, saltEncoded, hashEncoded), nil
}

func verifyPassword(password, encodedHash string) (bool, error) {
	// Parse the hash string
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}
	var memory, time, threads int
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads)
	if err != nil {
		return false, err
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}

	// Generate a hash with the same parameters and salt
	otherHash := argon2.IDKey([]byte(password), salt, uint32(time), uint32(memory), uint8(threads), uint32(32)) // Assuming keyLen is constant

	// Compare the hashes
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func createJWT(user User) (string, *HttpError) {
	//TODO:
	//- [X] Modify function to accept claim information
	//- [ ] Modify function to accept nonce and only return nonce if it was given.

	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims - Notes starting with * are required.
	claims := token.Claims.(jwt.MapClaims)
	claims["iss"] = "https://zg3.net"                    // *Issuer:              Identifier for the issuer of the token. This should be a URL.
	claims["sub"] = "home-1"                             // *Subject:             Identifier for the end-user. It must be unique within the issuer's domain.
	claims["aud"] = ""                                   // *Audience:            Identifies the recipients that the token is intended for. It must contain the client ID of the relying party.
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // *Expiration Time:     The expiration time on or after which the token must not be accepted for processing. It is represented in Unix time (seconds since the epoch).
	claims["iat"] = time.Now().Unix()                    // *Issued At:           The time at which the token was issued. It is also represented in Unix time.
	claims["auth_time"] = time.Now().Unix()              // *Authentication Time: Time when the end-user authentication occurred. This is required if the max_age parameter was used during the authentication request.
	claims["nonce"] = ""                                 // Unguessable, case-sensitive string value passed in authentication request from the relaying party
	claims["amr"] = []string{"pwd"}                      // Authentication methods reference. e.g., pwd: password, mfa, otp, sms, etc...
	claims["name"] = ""
	claims["admin"] = user.Username == "zedgama3"
	claims["authorized"] = true
	claims["username"] = user.Username
	claims["permissions"] = user.Permissions
	claims["user_id"] = user.UserId

	key, _ := base64.StdEncoding.DecodeString(secretKey)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", &HttpError{http.StatusInternalServerError, gin.H{"error": "Unable to sign JWT."}, err}
	}

	return tokenString, nil
}
