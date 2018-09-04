package main

import (
	"log"
	"os"

	"github.com/kjj6198/drink-bot/apis"
	"github.com/kjj6198/drink-bot/apis/drink_shops"
	"github.com/kjj6198/drink-bot/middlewares"

	"github.com/kjj6198/drink-bot/apis/stats"

	"github.com/kjj6198/drink-bot/apis/orders"

	"github.com/kjj6198/drink-bot/apis/menus"

	"github.com/kjj6198/drink-bot/apis/command"
	"github.com/kjj6198/drink-bot/apis/oauth"
	"github.com/kjj6198/drink-bot/app"

	"github.com/apex/gateway"

	"github.com/kjj6198/configo"
	"github.com/kjj6198/drink-bot/db"

	"github.com/gin-gonic/gin"
)

func main() {
	configo.Load("./config/env.yml")
	router := gin.Default()

	router.NoRoute(middlewares.AllowOrigin())

	api := router.Group("/")
	pg := db.Connect()
	client := db.NewClient()

	api.Use(func(c *gin.Context) {
		appContext := app.AppContext{
			DB:     pg,
			Client: client,
		}

		c.Set("app", appContext)
	})

	command.RegisterCommandHandler(api)
	api.Use(middlewares.AllowOrigin())
	routerMap := apis.RouterMap{
		"/user":       oauth.RegisterOAuthHandler,
		"/menus":      menus.RegisterMenusHandler,
		"/orders":     orders.RegisterOrdersHandler,
		"/stats":      stats.RegisterStatsHandler,
		"/drink_shop": drink_shops.RegisterDrinkShopsHandler,
	}

	for path, handler := range routerMap {
		// TODO: add recover, panic, logger as middlewares
		handler(api.Group(path))
	}

	if os.Getenv("ENV") == "development" {
		router.Run()
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", router))
	}
}
