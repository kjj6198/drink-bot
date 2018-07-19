package slack

import (
	"fmt"
	"net/http"
)

func OpenDialog() {
	resp, err := http.PostForm(fmt.Sprintf("%s/%s", SlackAPIURL, ActionOpenDialog))
}
