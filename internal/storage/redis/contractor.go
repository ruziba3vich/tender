package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
)

type ContractorCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewContractorCaching(client *redis.Client, logger *slog.Logger) *ContractorCaching {
	return &ContractorCaching{
		redisClient: client,
		logger:      logger,
	}
}
