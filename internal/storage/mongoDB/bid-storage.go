package mongodb

import (
	"log/slog"

	"github.com/zohirovs/internal/storage/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidStorage struct {
	db       *mongo.Collection
	logger   *slog.Logger
	bidCache *redis.BidCaching
}

func NewBidStorage(db *mongo.Database, logger *slog.Logger, cache *redis.BidCaching) *BidStorage {
	return &BidStorage{
		db:       db.Collection("Bids"),
		logger:   logger,
		bidCache: cache,
	}
}
