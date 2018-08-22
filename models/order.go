package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	Price     int       `json:"price"`
	UserID    uint      `json:"user_id"`
	User      *User     `json:"user"`
	MenuID    uint      `json:"menu_id"`
	Menu      *Menu     `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SerializedOrder struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	Price     int       `json:"price"`
	UserID    uint      `json:"-"`
	User      User      `json:"-"`
}

func (o *Order) GetOrderMenu(db *gorm.DB) *Menu {
	return db.Where("menu_id = ?", o.MenuID).Model(o.Menu).Value.(*Menu)
}

func (o *Order) GetUserOrders(db *gorm.DB) *[]*Menu {
	var menus []*Menu
	return db.
		Limit(200).
		Order("created_at DESC").
		Where("user_id = ?", o.UserID).
		Find(&menus).Value.(*[]*Menu)
}

func (o *Order) AfterCreate(tx *gorm.DB) {
	o.User = tx.Find(&User{ID: o.UserID}).Value.(*User)
	o.Menu = tx.Find(&Menu{ID: o.MenuID}).Value.(*Menu)
}

func (o *Order) CreateOrder(
	db *gorm.DB,
	userID uint,
	menuID uint,
	name string,
	price int,
	note string,
) *Order {
	o.UserID = userID
	o.MenuID = menuID
	o.Name = name
	o.Price = price
	o.Note = note

	return db.Create(o).Value.(*Order)
}

func (o *Order) UpdateOrder(
	db *gorm.DB,
	name string,
	price int,
	note string,
) *Order {
	if price <= 10 {
		price = o.Price
	}

	return db.
		Preload("User").
		Preload("Menu").
		Model(o).
		Update(&Order{
			ID:    o.ID,
			Name:  name,
			Price: price,
			Note:  note,
		}).Value.(*Order)
}
