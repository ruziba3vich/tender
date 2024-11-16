package handler

import (
	"log/slog"

	"github.com/zohirovs/internal/service"
)

type Handler struct {
	UserHandler         *UserHandler
	BidHandler          *BidHandler
	NotificationHandler *NotificationHandler
	TenderHandler       *TenderHandler
}

func NewHandler(logger *slog.Logger, service *service.Service) *Handler {
	return &Handler{
		UserHandler:         NewUserHandler(logger, service.User),
		BidHandler:          NewBidHandler(logger, service.Bid),
		NotificationHandler: NewNotificationHandler(logger, service.Notification),
		TenderHandler:       NewTenderHandler(logger, service.Tender),
	}
}
