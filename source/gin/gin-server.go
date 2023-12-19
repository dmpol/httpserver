package main

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const lifeTimeJWT = 15

var jwtSecretKey = []byte("my-super-secret-key-12345")

type Person struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users = map[string]string{
	"aaa": "12345",
	"bbb": "23456",
	"ccc": "34567",
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Login(c *gin.Context) {
	var person Person

	if err := c.BindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userPassword, ok := users[person.Username]

	if !ok || userPassword != person.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	tokenActionTime := time.Now().Add(lifeTimeJWT * time.Minute)

	claims := &Claims{
		Username: person.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenActionTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "internal server error"})
		return
	}

	c.SetCookie("token", tokenString, lifeTimeJWT*60,
		"/", "localhost", false, true)
}

func Welcome(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := &Claims{}

	tokn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !tokn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, "Welcome, "+claims.Username)
}

func Refresh(c *gin.Context) {
	tokenStr, err := c.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims := &Claims{}

	tokn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !tokn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
		return
	}

	if time.Until(claims.ExpiresAt.Time) > lifeTimeJWT*60*time.Second {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenActionTime := time.Now().Add(lifeTimeJWT * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(tokenActionTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "internal server error"})
		return
	}

	c.SetCookie("token", tokenString, lifeTimeJWT*60,
		"/", "localhost", false, true)
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", 0,
		"/", "localhost", false, true)
}

func main() {
	router := gin.Default()

	router.GET("/login", Login)
	router.GET("/welcome", Welcome)
	router.GET("/refresh", Refresh)
	router.GET("/logout", Logout)

	router.Run("localhost:8080")

}
