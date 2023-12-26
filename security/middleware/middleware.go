package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jwt2 "myhttpserver/security/jwt"
	"net/http"
)

func MiddlewareSecurity() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("Mid1")

		tokenStr, err := c.Cookie("token")
		if err != nil {
			fmt.Println("Mid2")

			//if err == http.ErrNoCookie {
			//	fmt.Println("Mid3")
			//	c.Redirect(http.StatusSeeOther, "/login")
			//	//c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			//	return
			//}
			//c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		fmt.Println("Mid4")

		username, tokn, err := jwt2.DemarshJWT(tokenStr)
		if err != nil {

			fmt.Println("Mid5")

			if username == "" {

				fmt.Println("Mid6")

				c.Redirect(http.StatusSeeOther, "/login")
				return
			}
			//if err == jwt.ErrSignatureInvalid {
			//	fmt.Println("Mid7")
			//	c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			//	return
			//}
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		fmt.Println("Mid8")

		c.Set("username", username)

		if !tokn.Valid {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		c.Next()
	}
}
