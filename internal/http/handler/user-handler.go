package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/models"
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

func (h *UserHandler) Register(c *gin.Context) {
	var user models.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
}
