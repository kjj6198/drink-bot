package middlewares

import (
	"os"

	"github.com/kjj6198/drink-bot/app"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/models"
	"github.com/kjj6198/drink-bot/services/token"
)

func AllowOrigin() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, OPTIONS, POST, DELETE")
		if os.Getenv("ENV") == "development" {
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "*")
		} else {
			c.Header("Access-Control-Allow-Origin", "https://drink-17.heroku.com")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		}

		if c.Request.Method == "OPTIONS" {
			c.Status(200)
			return
		}

		c.Next()
	}
}

func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		appContext := c.MustGet("app").(app.AppContext)
		tokenVal, err := c.Cookie("token")
		if err != nil {
			// TODO: prettify error
			c.AbortWithStatus(400)
			return
		}

		user, err := token.Parse(tokenVal)
		if err != nil {
			c.AbortWithStatus(400)
			return
		}

		dbUser := appContext.DB.Where("email = ?", user["email"]).First(new(models.User)).Value.(*models.User)

		if dbUser.ID == 0 {
			c.AbortWithStatus(401)
			return
		}

		c.Set("current_user", dbUser)

		c.Next()
	}
}
