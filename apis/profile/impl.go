package profile

import (
	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
	"github.com/kjj6198/drink-bot/middlewares"
	"github.com/kjj6198/drink-bot/models"
)

type rank struct {
	UserName string `json:"username"`
	MyRank   int    `json:"myrank"`
	Sum      int    `json:"sum"`
}

func myRank(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	currentUser, ok := c.Get("current_user")

	if !ok {
		c.AbortWithStatus(401)
		return
	}

	// TODO: extract to model layer
	row := appContext.DB.Raw(`SELECT
	username,
	SUM(price),
	rank() OVER (ORDER BY sum(price) DESC) AS myrank
	FROM orders
	INNER JOIN users ON orders.user_id = users.id
	WHERE user_id = ?
	GROUP BY user_id, users.username`, currentUser.(*models.User).ID).Row()

	var orders []*models.Order
	results := appContext.DB.
		Limit(200).
		Where("user_id = ?", currentUser.(*models.User).ID).
		Find(&orders).Value.(*[]*models.Order)

	var result rank
	row.Scan(&result.UserName, &result.Sum, &result.MyRank)
	c.JSON(200, gin.H{
		"rank":   result,
		"orders": results,
	})
}

func RegisterProfileHandler(router *gin.RouterGroup) {
	router.GET("/myrank", middlewares.Auth(), myRank)
}
