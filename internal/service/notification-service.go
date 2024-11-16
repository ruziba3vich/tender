package service

import (
	"log/slog"

	"github.com/zohirovs/internal/repos"
	"github.com/zohirovs/internal/storage/redis"
)

type NotificationService struct {
	notificationRepo  repos.NotificationRepo
	notificationCache *redis.NotificationCaching
	logger            *slog.Logger
}

func NewNotificationService(notificationRepo repos.NotificationRepo, cache *redis.NotificationCaching, logger *slog.Logger) *NotificationService {
	return &NotificationService{
		notificationRepo:  notificationRepo,
		notificationCache: cache,
		logger:            logger,
	}
}
