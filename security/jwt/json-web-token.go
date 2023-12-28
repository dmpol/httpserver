package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const lifeTimeJWT int = 15

var jwtSecretKey = []byte("my-super-secret-key-12345")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func CreatingJWT(userName string) (string, error) {

	tokenActionTime := time.Now().Add(time.Duration(lifeTimeJWT) * time.Minute)
	claims := &Claims{
		Username: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenActionTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	return tokenString, err
}

func DemarshJWT(tokenStr string) (string, *jwt.Token, error) {

	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})
	return claims.Username, tkn, err
}

func TimeRefresh(claims *Claims) bool {
	return time.Until(claims.ExpiresAt.Time) > time.Duration(lifeTimeJWT*60)*time.Second
}

func GetLifeTimeJWT() int {
	return lifeTimeJWT
}
