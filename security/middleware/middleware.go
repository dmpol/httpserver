package middleware

import (
	"github.com/gin-gonic/gin"
	jwt2 "myhttpserver/security/jwt"
	"net/http"
)

func MiddlewareSecurity() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		username, tokn, err := jwt2.DemarshJWT(tokenStr)
		if err != nil {
			if username == "" {
				c.Redirect(http.StatusSeeOther, "/login")
				return
			}
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}

		c.Set("username", username)

		if !tokn.Valid {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		c.Next()
	}
}
