package db

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	DatabaseType = "postgres"
)

type password string

func (p password) String() string {
	if p != "" {
		return fmt.Sprintf("password=%s", string(p))
	}

	return ""
}

// Connect connect postgre db
// TODO: error handling and better logger.
func Connect() *gorm.DB {
	var shouldEnableSSL = "disable"
	var password = password(os.Getenv("DB_PASSWORD"))
	if os.Getenv("ENV") != "development" {
		shouldEnableSSL = "require"
	}

	dbConnStr := fmt.Sprintf(
		"host=%s port=%s user=%s %s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		password,
		os.Getenv("DB_NAME"),
		shouldEnableSSL,
	)

	db, err := gorm.Open(DatabaseType, dbConnStr)

	if os.Getenv("ENV") == "development" {
		db.LogMode(true)
	}

	if err != nil {
		panic(err)
	}

	return db
}
