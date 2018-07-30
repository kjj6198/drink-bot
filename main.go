package main

import (
	"log"
	"os"

	"github.com/apex/gateway"

	apis "github.com/kjj6198/drink-bot/apis/command"
	"github.com/kjj6198/drink-bot/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.Load()
	router := gin.Default()

	command := router.Group("/")
	apis.RegisterCommandHandler(command)

	if os.Getenv("ENV") == "development" {
		router.Run()
	} else {
		log.Fatal(gateway.ListenAndServe(":3000", router))
	}
}
