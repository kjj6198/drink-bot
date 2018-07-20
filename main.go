package main

import (
	"fmt"
	"os"

	"github.com/kjj6198/drink-bot/config"
	"github.com/kjj6198/drink-bot/services/drink"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/services/slack"
)

func main() {
	config.Load()
	router := gin.Default()
	router.POST("/drink_shops", func(c *gin.Context) {
		result, _ := drink.GetDrinkShops()
		c.JSON(200, result)
	})

	router.POST("/create_menu", func(c *gin.Context) {

		c.JSON(200, &gin.H{
			"message": "trigger",
		})
	})

	router.POST("/", func(c *gin.Context) {
		input := new(slack.SlackMessageInput)
		c.Bind(input)
		fmt.Println(input)
		result, err := slack.OpenDialog(slack.DialogOptions{
			Dialog: slack.Dialog{
				CallbackID:  "submit-menu",
				Title:       "建立飲料訂單",
				SubmitLabel: "建立飲料訂單",
				Elements: []slack.Element{
					slack.Element{
						Label:    "訂單名稱",
						Name:     "name",
						Value:    input.Text,
						Type:     "text",
						Hint:     "please select drink_shop",
						Optional: false,
					},
					slack.Element{
						Label:      "店家名稱",
						Name:       "drink_shop",
						Type:       "select",
						Hint:       "please select drink_shop",
						DataSource: "external",
						Optional:   false,
					},
					slack.Element{
						Label:       "時間（分）",
						Name:        "duration",
						Type:        "text",
						SubType:     "number",
						Placeholder: "e.g: 900",
						Hint:        "單位為(分)",
						Optional:    false,
					},
					slack.Element{
						Label:      "notify to channel",
						Name:       "channel",
						Type:       "select",
						DataSource: "channels",
						Optional:   false,
					},
				},
			},
			TriggerID: input.TriggerID,
		}, map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("SLACK_ACCESS_TOKEN")),
		})

		fmt.Println(result, err)

		c.JSON(200, &gin.H{
			"message": "ok",
		})
	})

	router.Run()
}
