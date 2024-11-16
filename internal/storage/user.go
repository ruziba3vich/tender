package storage

import (
	"log/slog"

	"github.com/zohirovs/internal/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStorage struct {
	db        *mongo.Collection
	logger    *slog.Logger
	userCache *redis.UserCaching
}

func NewUserStorage(db *mongo.Database, logger *slog.Logger, cache *redis.UserCaching) *UserStorage {
	return &UserStorage{
		db:        db.Collection("User"),
		logger:    logger,
		userCache: cache,
	}
}
