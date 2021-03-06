package menus

import (
	"strconv"
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
	ShouldNotify bool   `json:"should_notify,omitempty"`
	Channel      string `json:"channel"`
}

func getMenus(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	menu := new(models.Menu)

	c.JSON(200, menu.GetMenus(appContext.DB, 20))
}

func getMenu(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	id, _ := strconv.ParseUint(c.Param("menu_id"), 10, 64)
	menu := &models.Menu{ID: uint(id)}
	menu.GetMenu(appContext.DB)

	c.JSON(200, menu)
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

	m := menu.CreateMenu(
		appContext.DB,
		params.Name,
		timestamp,
		params.Channel,
		params.DrinkShopID,
		currentUser.(*models.User).ID,
	)

	if m != nil {
		c.JSON(200, m)
		return
	}

	c.JSON(400, gin.H{
		"message": "cannot create menu",
	})
}

func RegisterMenusHandler(router *gin.RouterGroup) {
	router.OPTIONS("", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Next()
	})

	router.GET("", middlewares.AllowOrigin(), getMenus)
	router.GET("/:menu_id", middlewares.AllowOrigin(), getMenu)
	router.POST("", middlewares.AllowOrigin(), middlewares.Auth(), createMenu)
}
