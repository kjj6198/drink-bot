package db

import (
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	DatabaseType = "postgres"
)

type DrinkShop struct {
	ID        int
	Name      string
	Phone     string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (d *DrinkShop) AfterCreate() {
	fmt.Println("after save")
}

func Connect() *gorm.DB {
	dbConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(DatabaseType, dbConnStr)
	db.LogMode(true)

	if err != nil {
		panic(err)
	}

	return db
}
