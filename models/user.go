package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

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

func (u *User) Create(db *gorm.DB) {
	db.Create(u)
}
