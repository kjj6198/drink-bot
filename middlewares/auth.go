package middlewares

import (
	"github.com/kjj6198/drink-bot/app"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/models"
	"github.com/kjj6198/drink-bot/services/token"
)

// AllowOrigin enables common CORS
func AllowOrigin() func(c *gin.Context) {
	// TODO: in production, we only want to whitelist hosts to access API.
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, OPTIONS, POST, DELETE")
		c.Header("Access-Control-Allow-Headers", "Origin, Accept, Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.Status(200)
			return
		}

		c.Next()
	}
}

// Auth identify user by cookie auth_token
// TODO: support header authorization
func Auth() func(c *gin.Context) {
	return func(c *gin.Context) {
		appContext := c.MustGet("app").(app.AppContext)
		tokenVal, err := c.Cookie("auth_token")
		if err != nil {
			// TODO: prettify error
			c.AbortWithStatusJSON(401, err)
			return
		}

		user, err := token.Parse(tokenVal)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{
				"message": ErrNoPermission,
			})
			return
		}

		dbUser := appContext.
			DB.
			Where("email = ?", user["email"]).
			First(new(models.User)).
			Value.(*models.User)

		if dbUser.ID == 0 {
			c.AbortWithStatusJSON(401, gin.H{
				"message": "can not find user",
			})
			return
		}

		c.Set("current_user", dbUser)

		c.Next()
	}
}
