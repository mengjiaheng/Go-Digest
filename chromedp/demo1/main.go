package main

import (
	"demo1/router"
)

func main() {
	r := router.InitRouter()
	r.Run(":9000")
}
