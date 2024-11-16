package redis

import (
	"log/slog"

	"github.com/go-redis/redis/v8"
	"github.com/zohirovs/internal/config"
)

type RedisService struct {
	notification *NotificationCaching
	user         *UserCaching
	tender       *TenderCaching
	bid          *BidCaching
	contractor   *ContractorCaching
}

func New(redisDb *redis.Client, logger *slog.Logger) *RedisService {
	return &RedisService{
		logger:  logger,
		redisDb: redisDb,
	}
}

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURI,
		Password: "",
		DB:       0,
	})

	return rdb
}
