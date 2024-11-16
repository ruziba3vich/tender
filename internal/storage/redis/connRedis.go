/*
 * @Author: javohir-a abdusamatovjavohir@gmail.com
 * @Date: 2024-11-16 23:47:59
 * @LastEditors: javohir-a abdusamatovjavohir@gmail.com
 * @LastEditTime: 2024-11-17 01:39:16
 * @FilePath: /tender/internal/storage/redis/connRedis.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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
