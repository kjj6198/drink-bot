package oauth

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/models"
	"github.com/kjj6198/drink-bot/services/token"
	"github.com/kjj6198/drink-bot/utils"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/requests"
)

const (
	GoogleTokenInfoURL = "https://www.googleapis.com/oauth2/v3/tokeninfo"
)

type oauthInput struct {
	IDToken string `json:"id_token"`
}

type googleUser struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	HD      string `json:"hd"`
	Picture string `json:"picture"`
}

func (u *googleUser) Is17User() bool {
	return u.HD == "17.media"
}

func googleoauth(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	input := new(oauthInput)
	c.ShouldBindJSON(input)

	_, data, _ := requests.Request(context.Background(), requests.Config{
		Method: "GET",
		URL:    GoogleTokenInfoURL,
		Params: url.Values{
			"id_token": []string{input.IDToken},
		},
	})

	user := new(googleUser)
	json.Unmarshal([]byte(data), user)
	dbUser := appContext.DB.Where("email = ?", user.Email).First(new(models.User)).Value.(*models.User)

	if dbUser == nil && user.Is17User() {
		newUser := &models.User{
			Email:       user.Email,
			Username:    user.Name,
			Picture:     user.Picture,
			SignInCount: 1,
			IsAdmin:     false,
		}

		newUser.Create(appContext.DB)
	}

	if !user.Is17User() {
		c.AbortWithStatus(401)
		return
	}

	signedStr, err := token.Sign(dbUser)
	if hasErr := utils.ErrorHandler(err, c); hasErr {
		return
	}

	dbUser.Picture = user.Picture
	dbUser.Email = user.Email
	dbUser.SignInCount++
}

func RegisterOAuthHandler(router *gin.RouterGroup) {
	router.POST("/sign_in", googleoauth)
}
