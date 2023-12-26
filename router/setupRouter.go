package router

import (
	"github.com/gin-gonic/gin"
	security_handlers "myhttpserver/handlers"
	"myhttpserver/security/middleware"
)

func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/welcome", security_handlers.Welcome)
	router.GET("/login", security_handlers.Login)

	router.Use(middleware.MiddlewareSecurity())

	user := router.Group("/user/")
	user.GET("/welcome_user", security_handlers.WelcomeUser)
	user.GET("/refresh", security_handlers.Refresh)
	user.GET("/logout", security_handlers.Logout)

	return router
}
