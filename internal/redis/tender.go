package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type TenderCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}
