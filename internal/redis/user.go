package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type UserCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewUserCaching(client *redis.Client, logger *slog.Logger) *UserCaching {
	return &UserCaching{
		redisClient: client,
		logger:      logger,
	}
}
