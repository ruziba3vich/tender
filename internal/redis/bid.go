package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type BidCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}
