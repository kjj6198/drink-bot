package menus

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/middlewares"
	"github.com/kjj6198/drink-bot/models"
)

type menuParams struct {
	Name         string `json:"name"`
	EndTime      int64  `json:"end_time"`
	DrinkShopID  uint   `json:"drink_shop_id"`
	ShouldNotify bool   `json:"should_notify,emitempty"`
}

func getMenus(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	menu := new(models.Menu)

	c.JSON(200, menu.GetMenus(appContext.DB, 100))
}

func createMenu(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	currentUser, ok := c.Get("current_user")
	if !ok {
		c.AbortWithStatus(401)
		return
	}

	params := new(menuParams)
	c.BindJSON(params)

	menu := new(models.Menu)

	timestamp := time.Unix(params.EndTime, 0)

	if timestamp.Before(time.Now()) {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "結束時間不得小於現在時間",
		})
		return
	}

	if timestamp.IsZero() {
		c.AbortWithStatusJSON(400, gin.H{
			"message": "end_time is required",
		})
		return
	}

	menu.CreateMenu(
		appContext.DB,
		params.Name,
		timestamp,
		params.DrinkShopID,
		currentUser.(*models.User).ID,
	)
}

func RegisterMenusHandler(router *gin.RouterGroup) {
	router.OPTIONS("", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	router.GET("", middlewares.AllowOrigin(), getMenus)
	router.POST("", middlewares.AllowOrigin(), middlewares.Auth(), createMenu)
}
