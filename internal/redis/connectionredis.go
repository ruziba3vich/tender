package redis

import (
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

func NewRedisClient(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURI,
		Password: "",
		DB:       0,
	})

	return rdb
}
