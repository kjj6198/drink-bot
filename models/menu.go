package models

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
)

type Menu struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	EndTime     time.Time  `json:"end_time"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DrinkShopID uint       `json:"-"`
	UserID      uint       `json:"-"`
	User        *User      `json:"user"`
	DrinkShop   *DrinkShop `json:"drink_shop"`
	Orders      []Order    `json:"orders"`
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
		ID        uint       `json:"id"`
		Name      string     `json:"name"`
		EndTime   time.Time  `json:"end_time"`
		IsActive  bool       `json:"is_active"`
		CreatedAt time.Time  `json:"created_at"`
		DrinkShop *DrinkShop `json:"drink_shop"`
		Orders    []Order    `json:"orders"`
		Sum       int        `json:"sum"`
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

func (m *Menu) AfterCreate(tx *gorm.DB) {
	drinkShop := tx.Find(&DrinkShop{
		ID: m.DrinkShopID,
	}).Value.(*DrinkShop)

	user := tx.Find(&User{
		ID: m.UserID,
	}).Value.(*User)

	m.DrinkShop = drinkShop
	m.User = user
}

func (m *Menu) AfterSave(tx *gorm.DB) {
	if os.Getenv("ENV") == "development" {
		fmt.Println(m.DrinkShop, m.User)
	}
}

func (m *Menu) CreateMenu(
	db *gorm.DB,
	name string,
	endTime time.Time,
	drinkShopID uint,
	userID uint,
) *Menu {
	m.UserID = userID
	m.DrinkShopID = drinkShopID
	m.EndTime = endTime
	m.Name = name
	m.IsActive = true

	menu := db.Model(&Menu{}).Create(m).Value.(*Menu)
	return menu
}
