package app

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

type AppContext struct {
	DB     *gorm.DB
	Client *redis.Client
}
