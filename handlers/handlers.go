package handlers

import (
	"github.com/gin-gonic/gin"
	jwt2 "myhttpserver/security/jwt"
	"net/http"
)

type Person struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]string{
	"Dmitry":  "12345",
	"Valerya": "23456",
	"ccc":     "34567",
}

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, "Welcome!")
}

func Login(c *gin.Context) {
	var person Person

	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error3": err.Error()})
		return
	}

	userPassword, ok := users[person.Username]

	if !ok || userPassword != person.Password {
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
