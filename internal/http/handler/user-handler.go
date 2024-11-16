package handler

import (
	"log/slog"

	"github.com/zohirovs/internal/service"
)

type UserHandler struct {
	logger      *slog.Logger
	userService *service.UserService
}

func NewUserHandler(logger *slog.Logger, user *service.UserService) *UserHandler {
	return &UserHandler{
		logger:      logger,
		userService: user,
	}
}
