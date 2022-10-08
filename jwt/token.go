package jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func GetToken(igid string) (*string, error) {
	claims := GetClaims(&igid)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, *claims)
	signedStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	return &signedStr, nil
}
