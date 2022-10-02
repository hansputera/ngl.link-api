package database

import (
	"log"
	"os"

	"github.com/go-redis/redis/v9"
	"nglapi/global"
)

func InitDatabase() {
	opt, err := redis.ParseURL(os.Getenv("REDIS_URI"))
	if err != nil {
		log.Fatal(err)
	}
	rdb := redis.NewClient(opt)

	global.RedisClient = rdb
}
