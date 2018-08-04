package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/services/token"
)

func myRank(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	token.Parse()
	return
}

func RegisterProfileHandler(router *gin.RouterGroup) {

}
