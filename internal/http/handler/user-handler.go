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

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Register a new user with the provided details
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body     models.RegisterUser true "User registration details"
// @Success      201 {object} gin.H               "Successfully registered user"
// @Failure      400 {object} gin.H               "Invalid request"
// @Failure      500 {object} gin.H               "Internal server error"
// @Router       /register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	h.logger.Info("Register user")

	var user models.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userService.RegisterUser(c.Request.Context(), &user)
	if err != nil {
		h.logger.Error("failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
	h.logger.Info("User registered successfully", "email", user.Email)
}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body     models.LoginRequest true "User login credentials"
// @Success      200 {object} gin.H               "Successfully logged in with token"
// @Failure      400 {object} gin.H               "Invalid request"
// @Failure      500 {object} gin.H               "Internal server error"
// @Router       /login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	h.logger.Info("Login user")

	var user models.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), &user)
	if err != nil {
		h.logger.Error("failed to login user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
	h.logger.Info("User logged in successfully", "email", user.Email)
}
