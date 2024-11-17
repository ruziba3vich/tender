package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/service"
	"golang.org/x/crypto/bcrypt"
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

// @Summary Register a new user
// @Description Register a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.RegisterUser true "User registration details"
// @Success 201 {object} models.User "Successfully registered user"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	h.logger.Info("Register user")

	var user models.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(401, gin.H{"error": "Invalid request"})
		return
	}

	if user.Role != "client" && user.Role != "contractor" {
		h.logger.Error("invalid role")
		c.JSON(400, gin.H{"message": "invalid role"})
		return
	}

	if user.Email == "" || user.Username == "" {
		h.logger.Error("username or email cannot be empty")
		c.JSON(400, gin.H{"message": "username or email cannot be empty"})
		return
	}

	isValid := h.isValidEmail(user.Email)
	if isValid == false {
		h.logger.Error("invalid email format")
		c.JSON(400, gin.H{"message": "invalid email format"})
		return
	}

	token, err := h.userService.RegisterUser(c.Request.Context(), &user)
	if token == "Duplicate" {
		c.JSON(400, gin.H{"message": "Email already exists"})
		return
	}

	if err != nil {
		h.logger.Error("failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"token": token})
	h.logger.Info("User registered successfully", "email", user.Email)
}

// @Summary Login user
// @Description Authenticate user and return JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.RegisterUser true "User login credentials"
// @Success 200 {object} gin.H "Successfully logged in with token"
// @Failure 400 {object} gin.H "Invalid request"
// @Failure 500 {object} gin.H "Internal server error"
// @Router /login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	h.logger.Info("Login user")

	var user models.LoginRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if user.Username == "" || user.Password == "" {
		h.logger.Error("username or password cannot be empty")
		c.JSON(400, gin.H{"message": "Username and password are required"})
		return
	}

	token, err := h.userService.Login(c.Request.Context(), &user)

	if token == "not found" {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	if err != nil {
		h.logger.Error("failed to login user", "error", err)
		c.JSON(401, gin.H{"message": "Invalid username or password"})
		return
	}

	c.JSON(200, gin.H{"token": token})
	h.logger.Info("User logged in successfully", "username", user.Username)
}

func (h *UserHandler) isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if re == nil {
		return false
	}
	return re.MatchString(email)
}

func (h *UserHandler) hashPassword(password string) (string, error) {
	h.logger.Debug("hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		h.logger.Error("failed to hash password", "error", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	h.logger.Debug("password hashed successfully")
	return string(hashedPassword), nil
}
