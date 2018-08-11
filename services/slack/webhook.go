package slack

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/kjj6198/requests"
)

type SlackAttachment struct {
	Fallback   string              `json:"fallback,omitempty"`
	Title      string              `json:"title,omitempty"`
	TitleLink  string              `json:"title_link,omitempty"`
	Text       string              `json:"text,omitempty"`
	Color      string              `json:"color,omitempty"`
	Pretext    string              `json:"pretest,omitempty"`
	AuthorName string              `json:"author_name,omitempty"`
	AuthorLink string              `json:"author_link,omitempty"`
	ImageURL   string              `json:"image_url,omitempty"`
	Timestamp  int64               `json:"ts"`
	Fields     []map[string]string `json:"fields"`
}

func getChannel(channel string) string {
	if channel == "" {
		return "#frontend-underground"
	} else if strings.HasPrefix("#", channel) {
		return channel
	}

	return fmt.Sprintf("#%s", channel)
}
func SendMessage(message string, attachments []SlackAttachment, channel string) error {
	slackURL := os.Getenv("SLACK_WEBHOOK_URL")

	if os.Getenv("ENV") == "development" {
		channel = "#bon-appetit"
	}

	payload := map[string]interface{}{
		"text":        message,
		"username":    "yuile",
		"channel":     getChannel(channel),
		"attachments": attachments,
	}

	config := requests.Config{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Method: "POST",
		URL:    slackURL,
		Body:   payload,
	}

	fmt.Println(requests.Request(context.Background(), config))

	return nil
}
