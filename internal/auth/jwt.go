package auth

import (
	"github.com/golang-jwt/jwt"
	"time"
)

func NewJWT(login string, exp time.Time, secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   exp.Unix(),
	})

	return token.SignedString([]byte(secret))
}
