package db

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf(
			"%s:%s",
			os.Getenv("REDIS_HOST"),
			os.Getenv("REDIS_PORT"),
		),
		DB: 0,
	})

	return client
}
