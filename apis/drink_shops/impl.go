package drink_shops

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/models"
	"github.com/kjj6198/drink-bot/services/uploader"
	"github.com/kjj6198/drink-bot/utils"
)

const (
	baseURL = "http://17drink.s3-ap-northeast-1.amazonaws.com"
)

func create(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	multipart, err := c.MultipartForm()
	uploader := uploader.NewUploader(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_ACCESS_SECRET_KEY"),
	)

	if err != nil {
		c.JSON(400, gin.H{
			"message": "can not read multipart data.",
		})
		return
	}

	files := multipart.File["file"]
	file := files[0]
	f, _ := file.Open()
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Println("error during uploading file, skip.")
	}

	filename, err := uploader.Upload(
		"uploads",
		utils.GenerateUUID(),
		"image/jpeg",
		data,
	)

	drinkShop := &models.DrinkShop{
		Name:     multipart.Value["name"][0],
		Phone:    multipart.Value["photo"][0],
		Address:  multipart.Value["address"][0],
		ImageURL: fmt.Sprintf("%s/%s", baseURL, filename),
		Comment:  multipart.Value["comment"][0],
	}

	_, drinkShop = drinkShop.CreateDrinkShop(appContext.DB)
	c.JSON(200, drinkShop)
}

func drinkShop(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)

	id, _ := strconv.ParseUint(c.Param("dirnk_shop_id"), 10, 64)
	drinkShop := &models.DrinkShop{
		ID: uint(id),
	}

	result := drinkShop.Find(appContext.DB)
	c.JSON(200, result)
}

func delete(c *gin.Context) {
	if c.Request.Method != "DELETE" {
		log.Println("method must be DELETE")
		c.AbortWithStatus(400)
		return
	}

	appContext := c.MustGet("app").(app.AppContext)
	id, _ := strconv.ParseUint(c.Param("drink_shop_id"), 10, 64)

	drinkShop := &models.DrinkShop{
		ID: uint(id),
	}

	if drinkShop.DeleteDrinkShop(appContext.DB) {
		c.JSON(200, gin.H{"message": "ok"})
		return
	}

	c.AbortWithStatus(400)
}

func RegisterDrinkShopsHandler(router *gin.RouterGroup) {
	router.POST("", middlewares.Auth(), create)
	router.GET("/:drink_shop_id", drinkShop)
	router.DELETE("/:drink_shop_id", middlewares.Auth(), delete)
}
