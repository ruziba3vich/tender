package service

import (
	"log/slog"

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

func NewService(cache *redis.TenderCaching, logger *slog.Logger) *Service {
	return &Service{}
}
