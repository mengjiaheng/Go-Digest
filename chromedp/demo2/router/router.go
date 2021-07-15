package router

import (
	"demo2/api"
	"demo2/core"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.POST("/index", func(c *gin.Context) {
		url := "http://127.0.0.1:8090/index"
		if auth := api.ClientRequest(url); auth != "" {
			core.NewChrome(url, auth)
			fmt.Println(auth)
		}

		c.String(200, "OK")
	})
	return router
}
