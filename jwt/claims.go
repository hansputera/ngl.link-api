package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type NglClaims struct {
	UserId string `json:"id"`
	IgId   string `json:"gid"`
	jwt.RegisteredClaims
}

func GetClaims(userid *string, igid *string) *NglClaims {
	return &NglClaims{
		*userid,
		*igid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(3 * time.Hour),
			),
			Issuer: "Ngl-Clone Project",
		},
	}
}
