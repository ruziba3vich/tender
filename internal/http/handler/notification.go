package handler

import (
	"log/slog"

	"github.com/zohirovs/internal/service"
)

type NotificationHandler struct {
	logger              *slog.Logger
	notificationService *service.NotificationService
}

func NewNotificationHandler(logger *slog.Logger, notification *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		logger:              logger,
		notificationService: notification,
	}
}
