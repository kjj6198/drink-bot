package apis

import "github.com/gin-gonic/gin"

//type Handler interface {
//	RegisterHandler(router *gin.RouterGroup)
//}

type Handler func(router *gin.RouterGroup)

type RouterMap map[string]Handler
