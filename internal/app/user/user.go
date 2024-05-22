package user

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "MHcCAQEEIOqsEsDcG1fKjvhyQxZvBsUCORF+IFBSCANguJ5h4iWHoAoGCCqGSM49AwEHoUQDQgAE8P45+bQ8zWLyStx9oe1/VKnNaJ4XiC0eLwrEKWocwS1XZnb3hE4zZHK5GMsI6HARamQtyZ+cChJSLnHohzmiVA=="

// Authenticate a user based on their username and password and return a User
func (u *User) PasswordLogin(username string, password string) *User {
	return u
}

// Generate a new JWT access token
func (u *User) NewToken() (*Token, error) {
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
	claims["admin"] = u.Username == "zedgama3"
	claims["authorized"] = true
	claims["username"] = u.Username
	claims["permissions"] = u.Permissions
	claims["user_id"] = u.UserId

	key, _ := base64.StdEncoding.DecodeString(secretKey)

	// Sign the token with the secret key
	var myToken Token
	if tokenString, err := token.SignedString(key); err != nil {
		return nil, &HttpError{http.StatusInternalServerError, gin.H{"error": "Unable to sign JWT."}, err}
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
