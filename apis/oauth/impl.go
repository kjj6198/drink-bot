package oauth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

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
		Secure:   os.Getenv("ENV") != "development",
		HttpOnly: true,
		MaxAge:   3600,
	})
}

func googleoauth(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	input := new(oauthInput)
	c.ShouldBindJSON(input)

	if input.IDToken == "" {
		c.AbortWithStatusJSON(400, makeError(ErrNoIdTokenField))
		return
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

		setTokenCookie(c, "auth_token", jwtToken)

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
	setTokenCookie(c, "auth_token", jwtToken)

	c.JSON(200, gin.H{
		"id":       dbUser.ID,
		"email":    dbUser.Email,
		"username": dbUser.Username,
		"is_admin": dbUser.IsAdmin,
		"picture":  dbUser.Picture,
	})
}

func logout(c *gin.Context) {
	if str, err := c.Cookie("auth_token"); str != "" && err == nil {
		c.SetCookie("auth_token", "", 0, "/", "/", false, true)
		c.Status(200)
		return
	}

	c.JSON(400, makeError(ErrNotLogin))
}

func auth(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	tokenVal, err := c.Cookie("token")

	if err != nil {
		log.Println("request doesn't have token in cookie, fallback to Authorization header...")
	}

	// if don't have token in cookie
	// read Authorization Header as fallback
	if tokenVal == "" && err != nil {
		authorization := c.Request.Header.Get("Authorization")
		values := strings.Split(authorization, " ")

		if values[0] == "Bearer" && len(values) == 2 {
			tokenVal = values[1]
		}
	}

	if tokenVal == "" {
		c.Status(200)
		return
	}

	user, err := token.Parse(tokenVal)
	if err != nil {
		c.AbortWithStatusJSON(400, makeError(ErrNoIdTokenField))
		return
	}

	dbUser := appContext.DB.Where("email = ?", user["email"]).First(new(models.User)).Value.(*models.User)

	if dbUser.ID == 0 {
		c.AbortWithStatusJSON(401, makeError(ErrPermissionDenied))
		return
	}

	c.JSON(200, &models.SeralizedUser{
		ID:       dbUser.ID,
		Email:    dbUser.Email,
		Picture:  dbUser.Picture,
		Username: dbUser.Username,
		IsAdmin:  dbUser.IsAdmin,
	})
}

func allowCookie() func(c *gin.Context) {
	return func(c *gin.Context) {
		if os.Getenv("ENV") == "development" {
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
			c.Header("Access-Control-Allow-Credentials", "true")
		} else {
			c.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, OPTIONS, POST, DELETE")
			c.Header("Access-Control-Allow-Origin", "https://drink-17.heroku.com")
			c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			return
		}
		c.Next()
	}
}

// TODO: Move to middlewares
func allowOrigin() func(c *gin.Context) {
	return func(c *gin.Context) {
		if os.Getenv("ENV") == "development" {
			c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
		} else {
			c.Header("Access-Control-Allow-Methods", "GET, PUT, PATCH, OPTIONS, POST, DELETE")
			c.Header("Access-Control-Allow-Origin", "https://drink-17.heroku.com")
			c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers: Origin, X-Requested-With, Content-Type, Accept")
		}

		if c.Request.Method == "OPTIONS" {
			return
		}
		c.Next()
	}
}

// RegisterOAuthHandler register oauth api route
func RegisterOAuthHandler(router *gin.RouterGroup) {
	// if user already has account, just normally sign in
	// if user don't have account, create one and sign in for him.
	router.OPTIONS("/sign_in", allowCookie())
	router.OPTIONS("/auth", allowCookie())
	router.OPTIONS("/logout", allowCookie())
	router.POST("/sign_in", allowCookie(), googleoauth)
	router.POST("/auth", allowCookie(), auth)
	router.POST("/logout", allowCookie(), logout)
}
