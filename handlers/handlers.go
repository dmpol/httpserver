package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
	"myhttpserver/db/connection"
	jwt2 "myhttpserver/security/jwt"
	"net/http"
)

func Home(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome!")
}

func Login(c *gin.Context) {
	var personAuth PersonAuth

	if err := c.BindJSON(&personAuth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error Login": err.Error()})
		return
	}

	var person PersonWithId

	err := connection.GetConnect().QueryRow("SELECT user_id, user_name, password_hash, email FROM users WHERE user_name = $1", personAuth.Username).
		Scan(&person.Id, &person.Username, &person.Password, &person.Email)
	if err != nil {
		log.Printf("Ошибка запроса к DB: %s", err)
	}

	err = comparePasswords(person.Password, personAuth.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	tknStr, err := jwt2.CreatingJWT(person.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "internal server error"})
		return
	}

	c.SetCookie("token", tknStr, jwt2.GetLifeTimeJWT()*60,
		"/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/user/welcome_user")

}

func WelcomeUser(c *gin.Context) {

	username := c.GetString("username")
	if username == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	c.JSON(http.StatusOK, "Welcome, "+username)
}

func Refresh(c *gin.Context) {

	username := c.GetString("username")
	if username == "" {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	tknStr, err := jwt2.CreatingJWT(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "internal server error"})
		return
	}

	c.SetCookie("token", tknStr, jwt2.GetLifeTimeJWT()*60,
		"/", "localhost", false, true)
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", 0,
		"/", "localhost", false, true)
}
