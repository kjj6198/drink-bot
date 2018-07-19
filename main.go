package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		fmt.Println(c.Request.Header)
		fmt.Println(c.Request.ParseForm())
		fmt.Println(c.Request.PostForm)
		input := new(SlackMessageInput)
		c.Bind(input)
		fmt.Println(input)

		c.JSON(200, &gin.H{
			"message": "ok",
		})
	})

	router.Run()
}
