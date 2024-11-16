package mongodb

import (
	"log/slog"

	"github.com/zohirovs/internal/storage/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type NotificationStorage struct {
	db                *mongo.Collection
	logger            *slog.Logger
	notificationCache *redis.NotificationCaching
}

func NewNotificationStorage(db *mongo.Database, logger *slog.Logger, cache *redis.NotificationCaching) *NotificationStorage {
	return &NotificationStorage{
		db:                db.Collection("Notifications"),
		logger:            logger,
		notificationCache: cache,
	}
}
