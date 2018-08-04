package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Order struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	Price     int       `json:"price"`
	UserID    uint      `json:"-"`
	User      User      `json:"-"`
	MenuID    uint      `json:"-"`
	Menu      Menu      `json:"-"`
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
