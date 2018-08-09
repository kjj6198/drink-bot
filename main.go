package main

import (
	"log"
	"os"

	"github.com/kjj6198/drink-bot/apis/drink_shops"

	"github.com/kjj6198/drink-bot/apis/profile"
	"github.com/kjj6198/drink-bot/apis/stats"

	"github.com/kjj6198/drink-bot/apis/orders"

	"github.com/kjj6198/drink-bot/apis/menus"

	"github.com/kjj6198/drink-bot/apis/command"
	"github.com/kjj6198/drink-bot/apis/oauth"
	"github.com/kjj6198/drink-bot/app"

	"github.com/apex/gateway"

	"github.com/kjj6198/drink-bot/config"
	"github.com/kjj6198/drink-bot/db"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	router := gin.Default()

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
	userGroup := api.Group("/user")
	menuGroup := api.Group("/menus")
	orderGroup := api.Group("/orders")
	statsGroup := api.Group("/stats")
	drinkShopGroup := api.Group("/drink_shop")

	oauth.RegisterOAuthHandler(userGroup)
	profile.RegisterProfileHandler(userGroup)
	menus.RegisterMenusHandler(menuGroup)
	orders.RegisterOrdersHandler(orderGroup)
	stats.RegisterStatsHandler(statsGroup)
	drink_shops.RegisterDrinkShopsHandler(drinkShopGroup)

	if os.Getenv("ENV") == "development" {
		router.Run()
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", router))
	}
}
