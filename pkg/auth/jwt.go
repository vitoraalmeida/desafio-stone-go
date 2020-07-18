package auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

func CreateToken(user_id int) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}
