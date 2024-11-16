package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-redis/redis/v8"
	"github.com/zohirovs/internal/config"
)

type RedisService struct {
	Notification *NotificationCaching
	User         *UserCaching
	Tender       *TenderCaching
	Bid          *BidCaching
	Contractor   *ContractorCaching
}

func New(redisDb *redis.Client, logger *slog.Logger) *RedisService {
	return &RedisService{
		Notification: NewNotificationCaching(redisDb, logger),
		User:         NewUserCaching(redisDb, logger),
		Tender:       NewTenderCaching(redisDb, logger),
		Bid:          NewBidCaching(redisDb, logger),
		Contractor:   NewContractorCaching(redisDb, logger),
	}
}

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURI,
		Password: "", // Change this if password is required
		DB:       0,
	})

	// Check Redis connection
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return rdb, nil
}
