package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type NglClaims struct {
	IgId   string `json:"gid"`
	jwt.RegisteredClaims
}

func GetClaims(igid *string) *NglClaims {
	return &NglClaims{
		*igid,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add((24 * 7) * time.Hour),
			),
			Issuer: "Ngl-Clone Project",
		},
	}
}
