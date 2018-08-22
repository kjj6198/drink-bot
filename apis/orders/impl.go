package orders

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/middlewares"
	"github.com/kjj6198/drink-bot/models"
)

type orderParams struct {
	Name   string `json:"name"`
	MenuID uint   `json:"menu_id"`
	Price  int    `json:"price"`
	Note   string `json:"note"`
}

func createOrder(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	currentUser, ok := c.Get("current_user")

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	params := new(orderParams)
	c.BindJSON(params)

	if params.Price <= 10 {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "price can not be less than 10",
		})
		return
	}

	if params.Name == "" {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "name is required",
		})
		return
	}

	newOrder := &models.Order{}
	c.JSON(200, newOrder.CreateOrder(
		appContext.DB,
		currentUser.(*models.User).ID,
		params.MenuID,
		params.Name,
		params.Price,
		params.Note,
	))
}

func updateOrder(c *gin.Context) {
	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	appContext := c.MustGet("app").(app.AppContext)
	currentUser, _ := c.Get("current_user")
	orderID, _ := strconv.ParseUint(c.Param("order_id"), 10, 64)
	params := new(orderParams)
	c.BindJSON(params)
	order := appContext.
		DB.
		Preload("User").
		Preload("Menu").
		First(&models.Order{ID: uint(orderID)}).Value.(*models.Order)

	if order.UserID != currentUser.(*models.User).ID {
		c.AbortWithStatusJSON(401, gin.H{
			"message": "不可以修改其他人的訂單！",
		})
		return
	}

	if !order.Menu.IsActive {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "訂單已結束，不可修改訂單",
		})
		return
	}

	c.JSON(200, order.UpdateOrder(
		appContext.DB,
		params.Name,
		params.Price,
		params.Note,
	))
}

func RegisterOrdersHandler(router *gin.RouterGroup) {
	router.OPTIONS("", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		c.Next()
	})
	// TODO 可以用 before action 的方式來註冊 middleware 嗎？
	router.POST("", middlewares.AllowOrigin(), middlewares.Auth(), createOrder)
	router.PUT("/:order_id", middlewares.AllowOrigin(), middlewares.Auth(), updateOrder)
}
