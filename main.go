package main

import (
	router_main "myhttpserver/router"
)

func main() {
	router := router_main.Setup()
	router.Run("localhost:8080")
}
