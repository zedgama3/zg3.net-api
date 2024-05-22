package model

// User represents a user in the system
type User struct {
    ID    int
    Name  string
    Email string
}

// NewUser creates and returns a new User instance
func NewUser(id int, name, email string) *User {
    return &User{ID: id, Name: name, Email: email}
}
