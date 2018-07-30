package apis

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/services/drink"
	"github.com/kjj6198/drink-bot/services/slack"
	"github.com/kjj6198/drink-bot/utils"
)

func createMenu(c *gin.Context) {
	payload := []byte(c.PostForm("payload"))
	input := new(slack.SlackDialogParams)
	json.Unmarshal(payload, input)
	log.Println(input.User)
	err := drink.CreateMenu(&drink.MenuParams{
		Email:       input.User.Name,
		Name:        input.Submission["name"].(string),
		Duration:    utils.ParseInt(input.Submission["duration"].(string)),
		DrinkShopID: utils.ParseInt(input.Submission["drink_shop"].(string)),
		Channel:     "#" + input.SlackChannel.Name,
	})

	if err != nil {
		log.Println(err)
		c.AbortWithStatus(400)
		return
	}

	c.Status(200)
}

func openDialog(c *gin.Context) {
	input := new(slack.SlackMessageInput)
	c.Bind(input)
	_, err := slack.OpenDialog(slack.CreateMenuDialog(input.Text, input.TriggerID), map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("SLACK_ACCESS_TOKEN")),
	})

	if err == nil {
		c.Status(200)
		return
	}

	c.AbortWithStatus(400)
}

func getDrinkShops(c *gin.Context) {
	result, _ := drink.GetDrinkShops()
	c.JSON(200, result)
}

func RegisterCommandHandler(router *gin.RouterGroup) {
	router.POST("/drink_shops", getDrinkShops)
	router.POST("/", openDialog)
	router.POST("/create_menu", createMenu)
}
