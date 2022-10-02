package jwt

import (
	"nglapi/global"
	"os"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/golang-jwt/jwt/v4"
)

func ParseJWTToken(token string) *jwt.Token {
	parsed, err := jwt.ParseWithClaims(token, &NglClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil
	}

	return parsed
}

func VerifyJWTToken(token string) bool {
	parsed := ParseJWTToken(token)
	if parsed == nil || !parsed.Valid || parsed.Claims.(NglClaims).ExpiresAt.Unix() <= time.Now().Unix() {
		return false
	} else {
		_, err := global.RedisClient.Get(global.ContextConsume, parsed.Claims.(NglClaims).UserId).Result()
		if err != nil || err == redis.Nil {
			return false
		}

		return true
	}
}
