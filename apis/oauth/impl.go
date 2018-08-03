package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/models"
	"github.com/kjj6198/drink-bot/services/token"
	"github.com/kjj6198/drink-bot/utils"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/requests"
)

const (
	// GoogleTokenInfoURL is used for sign in token exchange.
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

func setTokenCookie(c *gin.Context, name string, value string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
		Expires:  time.Now().Add(24 * 30 * time.Hour),
	})
}

func googleoauth(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	input := new(oauthInput)
	c.ShouldBindJSON(input)

	if input.IDToken == "" {
		c.AbortWithStatusJSON(400, makeError(ErrNoIdTokenField))
	}

	// if already has token cookie, forward it to auth func.
	if val, err := c.Cookie("token"); err != nil && val != "" {
		auth(c)
		return
	}

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

	if dbUser.ID == 0 && user.Is17User() {
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

		setTokenCookie(c, "token", jwtToken)

		// TODO: find a better serailizer to do this work.
		c.JSON(200, gin.H{
			"id":       u.ID,
			"email":    u.Email,
			"username": u.Username,
			"is_admin": u.IsAdmin,
			"picture":  u.Picture,
		})
	} else if !user.Is17User() {
		c.AbortWithStatusJSON(401, makeError(ErrNot17User))
		return
	}

	// TODO: move this operation to model layer.
	dbUser.Picture = user.Picture
	dbUser.Email = user.Email
	dbUser.SignInCount++
	appContext.DB.Model(&models.User{}).Update(dbUser)

	jwtToken, _ := token.Sign(dbUser)
	setTokenCookie(c, "token", jwtToken)

	c.JSON(200, gin.H{
		"id":       dbUser.ID,
		"email":    dbUser.Email,
		"username": dbUser.Username,
		"is_admin": dbUser.IsAdmin,
		"picture":  dbUser.Picture,
	})
}

func auth(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	tokenVal, err := c.Cookie("token")

	if !utils.ErrorHandler(err, c) {
		return
	}

	// if don't have token in cookie
	// read Authorization Header as fallback
	if tokenVal == "" {
		authorization := c.Request.Header.Get("Authorization")
		values := strings.Split(authorization, " ")

		if values[0] == "Bearer" && len(values) == 2 {
			tokenVal = values[1]
		}
	}

	user, err := token.Parse(tokenVal)
	if !utils.ErrorHandler(err, c) {
		return
	}

	dbUser := appContext.DB.Where("email = ?", user["email"]).First(new(models.User)).Value.(*models.SeralizedUser)

	if dbUser.ID == 0 {
		c.AbortWithStatusJSON(401, makeError(ErrPermissionDenied))
	}
	fmt.Println(dbUser)
	c.JSON(200, dbUser)
}

// RegisterOAuthHandler register oauth api route
func RegisterOAuthHandler(router *gin.RouterGroup) {
	// if user already has account, just normally sign in
	// if user don't have account, create one and sign in for him.
	router.POST("/sign_in", googleoauth)
	router.POST("/auth", auth)
}
