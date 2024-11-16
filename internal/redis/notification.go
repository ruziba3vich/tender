package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type NotificationCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewNotificationCaching(client *redis.Client, logger *slog.Logger) *NotificationCaching {
	return &NotificationCaching{
		redisClient: client,
		logger:      logger,
	}
}
