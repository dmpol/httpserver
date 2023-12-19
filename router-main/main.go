package main

import (
	"github.com/gin-gonic/gin"
	security_handlers "myhttpserver/handlers/security-handlers"
)

func main() {
	router := gin.Default()

	router.GET("/login", security_handlers.Login)

	router.Use(security_handlers.MiddlewareSecurity())

	router.GET("/welcome", security_handlers.Welcome)
	router.GET("/refresh", security_handlers.Refresh)
	router.GET("/logout", security_handlers.Logout)

	router.Run("localhost:8080")

}
