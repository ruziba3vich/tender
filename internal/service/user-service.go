package service

import (
	"log/slog"

	"github.com/zohirovs/internal/repos"
	"github.com/zohirovs/internal/storage/redis"
)

type UserService struct {
	userRepo  repos.UserRepo
	userCache *redis.UserCaching
	logger    *slog.Logger
}

func NewUserService(userRepo repos.UserRepo, cache *redis.UserCaching, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo:  userRepo,
		userCache: cache,
		logger:    logger,
	}
}
