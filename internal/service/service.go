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
		Bid:          NewBidService(repo.BidRepo(), cache.Bid, logger),
	}
}
