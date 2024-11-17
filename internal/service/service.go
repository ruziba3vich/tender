/*
 * @Author: javohir-a abdusamatovjavohir@gmail.com
 * @Date: 2024-11-17 04:57:41
 * @LastEditors: javohir-a abdusamatovjavohir@gmail.com
 * @LastEditTime: 2024-11-17 12:40:05
 * @FilePath: /tender/internal/service/service.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package service

import (
	"log/slog"

	"github.com/zohirovs/internal/storage"
	"github.com/zohirovs/internal/storage/redis"
)

type (
	Service struct {
		User         *UserService
		Notification *NotificationService
		Tender       *TenderService
		Bid          *BidService
	}
)

func NewService(cache *redis.RedisService, logger *slog.Logger, repo storage.StorageI) *Service {
	return &Service{
		User:         NewUserService(repo.UserRepo(), logger),
		Notification: NewNotificationService(repo.NotificationRepo(), cache.Notification, logger),
		Tender:       NewTenderService(repo.TenderRepo(), cache.Tender, logger),
		Bid:          NewBidService(repo.BidRepo(), repo.TenderRepo(), logger),
	}
}
