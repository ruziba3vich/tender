package handler

import (
	"log/slog"

	"github.com/zohirovs/internal/service"
)

type BidHandler struct {
	logger     *slog.Logger
	bidService *service.BidService
}

func NewBidHandler(logger *slog.Logger, bid *service.BidService) *BidHandler {
	return &BidHandler{
		logger:     logger,
		bidService: bid,
	}
}
