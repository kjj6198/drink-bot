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

	if dbUser.ID == 0 && !user.Is17User() {
		newUser := &models.User{
			Email:       user.Email,
			Username:    user.Name,
			Picture:     user.Picture,
			SignInCount: 1,
			IsAdmin:     false,
		}

		u := newUser.Create(appContext.DB)
		jwtToken, err := token.Sign(u)

		if utils.ErrorHandler(err, c) {
			return
		}

		c.SetCookie("token", jwtToken, 108000, "/", "", false, true)

		// TODO: find a better serailizer to do this work.
		c.JSON(200, gin.H{
			"id":       u.ID,
			"email":    u.Email,
			"username": u.Username,
			"is_admin": u.IsAdmin,
			"picture":  u.Picture,
		})
	} else {
		c.AbortWithStatus(401)
		return
	}

	// TODO: move this operation to model layer.
	dbUser.Picture = user.Picture
	dbUser.Email = user.Email
	dbUser.SignInCount++
	appContext.DB.Model(&models.User{}).Update(dbUser)

	jwtToken, _ := token.Sign(dbUser)
	c.SetCookie("token", jwtToken, 108000, "/", "", false, true)

	c.JSON(200, gin.H{
		"id":       dbUser.ID,
		"email":    dbUser.Email,
		"username": dbUser.Username,
		"is_admin": dbUser.IsAdmin,
		"picture":  dbUser.Picture,
	})
}

func RegisterOAuthHandler(router *gin.RouterGroup) {
	router.POST("/sign_in", googleoauth)
}
