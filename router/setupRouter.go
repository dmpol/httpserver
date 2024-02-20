package router

import (
	"github.com/gin-gonic/gin"
	security_handlers "myhttpserver/handlers"
	"myhttpserver/security/middleware"
)

func Setup() *gin.Engine {
	router := gin.Default()

	router.GET("/welcome", security_handlers.Home)
	router.GET("/login", security_handlers.Login)

	rest := router.Group("/api/v1/")
	rest.GET("/persons", security_handlers.GetPersons)
	rest.GET("/persons/:id", security_handlers.GetPerson)
	rest.POST("/persons", security_handlers.CreatePerson)
	rest.PUT("/persons/:id", security_handlers.UpdatePerson)
	rest.DELETE("/persons/:id", security_handlers.DeletePerson)

	user := router.Group("/user/")
	user.Use(middleware.MiddlewareSecurity())
	user.GET("/welcome_user", security_handlers.WelcomeUser)
	user.GET("/refresh", security_handlers.Refresh)
	user.GET("/logout", security_handlers.Logout)

	return router
}
