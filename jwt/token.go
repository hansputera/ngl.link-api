package jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func GetToken(userid string, igid string) (*string, error) {
	claims := GetClaims(&userid, &igid)

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, *claims)
	signedStr, err := token.SignedString(os.Getenv("JWT_SECRET"))
	if err != nil {
		return nil, err
	}

	return &signedStr, nil
}
