package service

import (
	"log/slog"

	"github.com/zohirovs/internal/repos"
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

func NewService(cache *redis.RedisService, logger *slog.Logger, repo repos.Repos) *Service {
	return &Service{
		User:         NewUserService(repo.UserRepo, cache.User, logger),
		Notification: NewNotificationService(repo.NotificationRepo, cache.Notification, logger),
		Tender:       NewTenderService(repo.TenderRepo, cache.Tender, logger),
		Bid:          NewBidService(repo.BidRepo, cache.Bid, logger),
	}
}
