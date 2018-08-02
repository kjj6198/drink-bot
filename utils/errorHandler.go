package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func ErrorHandler(err error, c *gin.Context) bool {
	if err != nil {
		log.Println(errors.Wrap(err, "[ErrorHandler]"))
		c.AbortWithError(400, err)
		return true
	}

	return false
}
