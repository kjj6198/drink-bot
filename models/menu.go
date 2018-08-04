package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

type Menu struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DrinkShopID uint      `json:"-"`
	UserID      uint      `json:"-"`
	User        User      `json:"user"`
	DrinkShop   DrinkShop `json:"drink_shop"`
	Orders      []Order   `json:"orders"`
}

// GetSum calculate menu sum
// TODO: calculate from sql
// SELECT SUM(price) from orders where orders.id = menu.id
func (m *Menu) GetSum() int {
	sum := 0
	for _, v := range m.Orders {
		sum += v.Price
	}

	return sum
}

func (m *Menu) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        uint      `json:"id"`
		Name      string    `json:"name"`
		EndTime   time.Time `json:"end_time"`
		IsActive  bool      `json:"is_active"`
		CreatedAt time.Time `json:"created_at"`
		DrinkShop DrinkShop `json:"drink_shop"`
		Orders    []Order   `json:"orders"`
		Sum       int       `json:"sum"`
	}{
		ID:        m.ID,
		Name:      m.Name,
		EndTime:   m.EndTime,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		DrinkShop: m.DrinkShop,
		Orders:    m.Orders,
		Sum:       m.GetSum(),
	})
}

func (m *Menu) GetMenus(db *gorm.DB, limit int) *[]*Menu {
	if limit <= 0 {
		limit = 50
	}

	var menus []*Menu

	return db.
		Limit(limit).
		Order("menus.created_at DESC").
		Preload("User").
		Preload("DrinkShop").
		Preload("Orders").
		Find(&menus).
		Value.(*[]*Menu)
}

func (m *Menu) CreateMenu(
	name string,
	endTime time.Time,
	drinkShopID uint,
	userID uint,
) {

}
