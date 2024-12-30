package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtsecretkey = []byte("my-secret-key")

func CreateJWT(email string, password string) (string, error) {

	claims := jwt.MapClaims{
		"email":    email,
		"password": password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtsecretkey)

}
