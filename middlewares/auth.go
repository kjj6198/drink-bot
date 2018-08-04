package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
)

func AllowOrigin() func(c *gin.Context) {
	return func(c *gin.Context) {
		if os.Getenv("ENV") == "development" {
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		} else {
			c.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, OPTIONS, POST, DELETE")
			c.Header("Access-Control-Allow-Origin", "https://drink-17.heroku.com")
			c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept")
		}

		if c.Request.Method == "OPTIONS" {
			c.Status(200)
			return
		}

		c.Next()
	}
}
