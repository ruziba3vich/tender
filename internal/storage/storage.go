package storage

import (
	"log/slog"

	"github.com/zohirovs/internal/config"
	"github.com/zohirovs/internal/repos"
	mongodb "github.com/zohirovs/internal/storage/mongoDB"
	"github.com/zohirovs/internal/storage/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type StorageI interface {
	UserRepo() repos.UserRepo
	TenderRepo() repos.TenderRepo
	BidRepo() repos.BidRepo
	NotificationRepo() repos.NotificationRepo
}

type Storage struct {
	userRepo         repos.UserRepo
	tenderRepo       repos.TenderRepo
	bidRepo          repos.BidRepo
	notificationRepo repos.NotificationRepo
}

func New(db *mongo.Database, cfg *config.Config, logger *slog.Logger, cache *redis.RedisService) StorageI {
	return &Storage{
		userRepo:         mongodb.NewUserStorage(db, cfg, logger, cache.User),
		tenderRepo:       mongodb.NewTenderStorage(db, logger, cache.Tender),
		bidRepo:          mongodb.NewBidStorage(db, logger, cache.Bid),
		notificationRepo: mongodb.NewNotificationStorage(db, logger, cache.Notification),
	}
}

func (s *Storage) UserRepo() repos.UserRepo {
	return s.userRepo
}

func (s *Storage) TenderRepo() repos.TenderRepo {
	return s.tenderRepo
}

func (s *Storage) BidRepo() repos.BidRepo {
	return s.bidRepo
}

func (s *Storage) NotificationRepo() repos.NotificationRepo {
	return s.notificationRepo
}
