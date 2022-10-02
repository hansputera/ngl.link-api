package global

import (
	"context"

	"github.com/go-redis/redis/v9"
)

var ContextConsume = context.Background()
var RedisClient *redis.Client