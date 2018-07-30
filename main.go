package main

import (
	"log"
	"os"

	"github.com/kjj6198/drink-bot/app"

	"github.com/apex/gateway"

	apis "github.com/kjj6198/drink-bot/apis/command"
	"github.com/kjj6198/drink-bot/config"
	"github.com/kjj6198/drink-bot/db"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	router := gin.Default()

	command := router.Group("/")
	pg := db.Connect()
	client := db.NewClient()
	command.Use(func(c *gin.Context) {
		appContext := app.AppContext{
			DB:     pg,
			Client: client,
		}

		c.Set("app", appContext)
	})

	apis.RegisterCommandHandler(command)

	if os.Getenv("ENV") != "development" {
		router.Run()
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", router))
	}
}
