package storage

import (
	"log/slog"

	"github.com/zohirovs/internal/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type TenderStorage struct {
	db          *mongo.Collection
	logger      *slog.Logger
	tenderCache *redis.TenderCaching
}

func NewTenderStorage(db *mongo.Database, logger *slog.Logger, cache *redis.TenderCaching) *TenderStorage {
	return &TenderStorage{
		db:          db.Collection("tender"),
		logger:      logger,
		tenderCache: cache,
	}
}
