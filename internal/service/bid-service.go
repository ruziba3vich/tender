package service

import (
	"log/slog"

	"github.com/zohirovs/internal/repos"
	"github.com/zohirovs/internal/storage/redis"
)

type BidService struct {
	bidRepo  repos.BidRepo
	bidCache *redis.BidCaching
	logger   *slog.Logger
}

func NewBidService(bidRepo repos.BidRepo, cache *redis.BidCaching, logger *slog.Logger) *BidService {
	return &BidService{
		bidRepo:  bidRepo,
		bidCache: cache,
		logger:   logger,
	}
}
