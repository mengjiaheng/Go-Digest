package router

import (
	"demo1/core"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/", func(c *gin.Context) {
		core.NewCrawler()
		c.String(200, "OK")
	})
	return router
}
