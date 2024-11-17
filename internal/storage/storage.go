/*
 * @Author: javohir-a abdusamatovjavohir@gmail.com
 * @Date: 2024-11-17 06:27:56
 * @LastEditors: javohir-a abdusamatovjavohir@gmail.com
 * @LastEditTime: 2024-11-17 12:58:05
 * @FilePath: /tender/internal/storage/storage.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
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
		bidRepo:          mongodb.NewBidStorage(db, logger),
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
