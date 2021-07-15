package main

import "demo2/router"

func main() {
	r := router.InitRouter()
	r.Run(":9000")
}
