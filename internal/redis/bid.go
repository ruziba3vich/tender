package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type BidCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewBidCaching(client *redis.Client, logger *slog.Logger) *BidCaching {
	return &BidCaching{
		redisClient: client,
		logger:      logger,
	}
}
