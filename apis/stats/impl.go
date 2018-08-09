package stats

import (
	"github.com/gin-gonic/gin"
	"github.com/kjj6198/drink-bot/app"
)

type rank struct {
	UserName string `json:"username"`
	MyRank   int    `json:"myrank"`
	Sum      int    `json:"sum"`
}

const (
	rankSQL = `SELECT
	username,
	SUM(price),
	rank() OVER (ORDER BY sum(price) DESC) AS myrank
	FROM orders
	INNER JOIN users ON orders.user_id = users.id
	GROUP BY user_id, users.username
	`
	drinkShopSumSQL = `SELECT
	drink_shops.name AS name,
	sum(price) AS total
	FROM orders
	INNER JOIN menus ON menus.id = orders.menu_id
	INNER JOIN drink_shops ON menus.drink_shop_id = drink_shops.id
	GROUP BY drink_shops.name
	ORDER BY total DESC`
)

type drinkShopStat struct {
	Name  string `json:"name"`
	Total int    `json:"total"`
}

func drinkShopSum(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	rows, err := appContext.DB.Raw(drinkShopSumSQL).Rows()
	defer rows.Close()
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	var results []*drinkShopStat
	for rows.Next() {
		var result drinkShopStat
		rows.Scan(&result.Name, &result.Total)
		results = append(results, &result)
	}

	c.JSON(200, results)
}

func ranks(c *gin.Context) {
	appContext := c.MustGet("app").(app.AppContext)
	rows, err := appContext.DB.Raw(rankSQL).Rows()
	defer rows.Close()
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	var results []*rank

	for rows.Next() {
		var result rank
		rows.Scan(&result.UserName, &result.Sum, &result.MyRank)
		results = append(results, &result)
	}

	c.JSON(200, results)
}

func RegisterStatsHandler(router *gin.RouterGroup) {
	router.GET("/all_ranks", ranks)
	router.GET("/drink_shops", drinkShopSum)
}
