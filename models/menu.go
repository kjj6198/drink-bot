package models

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/kjj6198/drink-bot/services/slack"
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
	Channel     string     `json:"channel"`
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
		User      *User      `json:"user"`
		Sum       int        `json:"sum"`
	}{
		ID:        m.ID,
		Name:      m.Name,
		EndTime:   m.EndTime,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		DrinkShop: m.DrinkShop,
		Orders:    m.Orders,
		User:      m.User,
		Sum:       m.GetSum(),
	})
}

func (m *Menu) GetMenu(db *gorm.DB) {
	m = db.
		Preload("User").
		Preload("DrinkShop").
		Preload("Orders").
		Preload("Orders.User").
		Model(m).
		First(m).
		Value.(*Menu)
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
		Preload("Orders.User").
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

	m.SendMenuInfoToChannel()
}

func (m *Menu) NotifyCountdown(db *gorm.DB) {
	// TODO: refactor to worker mode.
	countdown := m.EndTime.Sub(time.Now().UTC())
	endHint := fmt.Sprintf("*%s* 已經結束囉，沒有訂到的哭哭哦 :jack-see-you: 查詢訂單：\n`@yuile :order_id`", m.Name)

	time.AfterFunc(countdown, func() {
		go db.Model(m).Update("is_active", false)
		slack.SendMessage(endHint, []slack.SlackAttachment{
			slack.SlackAttachment{
				AuthorName: "17 Drink",
				AuthorLink: os.Getenv("HOST_URL"),
				Color:      "#fe6565",
				TitleLink:  fmt.Sprintf("%s/menus/%s", os.Getenv("HOST_URL"), fmt.Sprint(m.ID)),
				Text:       m.Name,
				Fields: []map[string]string{
					map[string]string{"title": "訂單編號", "value": strconv.FormatUint(uint64(m.ID), 10)},
					map[string]string{"title": "訂單金額", "value": strconv.Itoa(m.GetSum())},
					map[string]string{"title": "杯數", "value": strconv.Itoa(len(m.Orders))},
					map[string]string{"title": "店家電話", "value": m.DrinkShop.Phone},
				},
			},
		}, m.Channel)
	})
}

func (m *Menu) SendMenuInfoToChannel() {
	hint := fmt.Sprintf(
		"*%s* 發起了訂飲料活動",
		m.User.Username,
	)

	loc, _ := time.LoadLocation("Asia/Taipei")
	// TODO: error handling
	go slack.SendMessage(hint, []slack.SlackAttachment{
		slack.SlackAttachment{
			AuthorName: "17 Drink",
			AuthorLink: os.Getenv("HOST_URL"),
			Title:      "訂單連結",
			Color:      "#27cc85",
			TitleLink:  fmt.Sprintf("%s/menus/%s", os.Getenv("HOST_URL"), fmt.Sprint(m.ID)),
			Text:       m.Name,
			Fields: []map[string]string{
				map[string]string{"title": "店家名稱", "value": m.DrinkShop.Name},
				map[string]string{"title": "開始時間", "value": m.CreatedAt.Format("15:04:05")},
				map[string]string{"title": "結束時間", "value": m.EndTime.In(loc).Format("15:04:05")},
				map[string]string{"title": "剩餘時間", "value": fmt.Sprintf("%d 分", int(m.EndTime.Sub(time.Now()).Minutes()))},
			},
			ImageURL:  m.DrinkShop.ImageURL,
			Timestamp: time.Now().Unix(),
		},
	}, m.Channel)
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
	m.EndTime = endTime.UTC()
	m.Name = name
	m.IsActive = true

	menu := db.Model(&Menu{}).Create(m).Value.(*Menu)

	if menu != nil {
		go menu.NotifyCountdown(db)
	}

	return menu
}
