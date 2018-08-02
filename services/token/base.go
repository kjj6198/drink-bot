package token

import (
	"fmt"
	"log"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kjj6198/drink-bot/models"
)

func Sign(user *models.User) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email":   user.Email,
		"picture": user.Picture,
	})

	return token.SignedString([]byte("64characterslongstring"))
}

func Parse(jwtStr string) (result map[string]interface{}, err error) {
	token, _ := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method %+v", token.Header["alg"])
		}

		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})

	if token == nil {
		log.Fatal("can not parse jwt string")
		return nil, err
	}

	return token.Claims.(jwt.MapClaims), nil
}
