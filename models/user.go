package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

// User is a normal user
type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Picture     string    `json:"picture"`
	SignInCount int       `json:"sign_in_count"`
	IsAdmin     bool      `json:"is_admin"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SeralizedUser is seralized user type for response
// don't put important data in it.
type SeralizedUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Picture  string `json:"picture"`
	IsAdmin  bool   `json:"is_admin"`
}

func (u *User) Create(db *gorm.DB) *User {
	return db.Create(u).Value.(*User)
}
