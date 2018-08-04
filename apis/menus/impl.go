package menus

import (
	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/middlewares"
	"github.com/kjj6198/drink-bot/models"
)

func getMenus(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	menu := new(models.Menu)

	c.JSON(200, menu.GetMenus(appContext.DB, 100))
}

func RegisterMenusHandler(router *gin.RouterGroup) {
	router.OPTIONS("", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Next()
	})
	router.GET("", middlewares.AllowOrigin(), getMenus)
}
