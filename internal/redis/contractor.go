package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type ContractorCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}
