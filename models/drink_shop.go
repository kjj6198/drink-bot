package models

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
)

const (
	DrinkShopKey = "drink_shops"
)

type DrinkShop struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	ImageURL  string    `json:"image_url"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (d *DrinkShop) GetDrinkShops(db *gorm.DB, client *redis.Client) []DrinkShop {
	var result []DrinkShop

	if client.LLen(DrinkShopKey).Val() != 0 {
		for _, val := range client.LRange(DrinkShopKey, 0, -1).Val() {
			in := new(DrinkShop)
			json.Unmarshal([]byte(val), in)
			result = append(result, *in)
		}

		return result
	}

	result = *db.Limit(150).Find(&[]DrinkShop{}).Value.(*[]DrinkShop)
	go func() {
		for _, val := range result {
			data, _ := json.Marshal(&val)
			client.RPush(DrinkShopKey, data)
		}
	}()

	return result
}
