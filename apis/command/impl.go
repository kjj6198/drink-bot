package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/models"
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
	app := c.MustGet("app").(app.AppContext)
	result := new(models.DrinkShop).GetDrinkShops(app.DB, app.Client)
	res := make([]map[string]string, len(result))

	for i, val := range result {
		res[i] = make(map[string]string)
		res[i]["label"] = val.Name
		res[i]["value"] = strconv.FormatUint(uint64(val.ID), 10)
	}

	if c.Request.Method == "GET" {
		c.JSON(200, result)
		return
	}

	c.JSON(200, gin.H{
		"options": res,
	})
}

func RegisterCommandHandler(router *gin.RouterGroup) {
	router.POST("/drink_shops", getDrinkShops)
	router.GET("/drink_shops", getDrinkShops)
	router.POST("/", openDialog)
	router.POST("/create_menu", createMenu)
}
