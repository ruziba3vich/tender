package service

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/repos"
)

type UserService struct {
	userRepo repos.UserRepo
	logger   *slog.Logger
}

func NewUserService(userRepo repos.UserRepo, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

// 1
func (s *UserService) RegisterUser(ctx context.Context, user *models.RegisterUser) (string, error) {
	_, err := s.userRepo.GetUserByEmail(ctx, user.Email)
	if err == nil {
		return "", fmt.Errorf("user with this email already exists")
	}

	isValid := s.isEmailExists(user.Email)
	if !isValid {
		return "", fmt.Errorf("invalid email format")
	}

	isValid = s.isValidPassword(user.Password)
	if !isValid {
		return "", fmt.Errorf("invalid password format")
	}

	new_user := models.User{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
		Role:     user.Role,
	}

	token, err := s.userRepo.RegisterUser(ctx, &new_user)
	if err != nil {
		return "", fmt.Errorf("failed to register user: %w", err)
	}

	return token, nil
}

// 2
func (s *UserService) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.GetUserByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by userID: %w", err)
	}
	return user, nil
}

// 3
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return user, nil
}

// 4
func (s *UserService) ChangeUserRole(ctx context.Context, userID string, role string) error {
	err := s.userRepo.ChangeUserRole(ctx, userID, role)
	if err != nil {
		return fmt.Errorf("failed to change user role: %w", err)
	}
	return nil
}

// 5
func (s *UserService) ChangeUserPassword(ctx context.Context, resetPassword *models.ResetPassword) error {
	if !s.isValidPassword(resetPassword.NewPassword) {
		return fmt.Errorf("invalid password format")
	}

	err := s.userRepo.ChangeUserPassword(ctx, resetPassword)
	if err != nil {
		return fmt.Errorf("failed to change user password: %w", err)
	}
	return nil
}

// 6
func (s *UserService) SendVerificationCode(ctx context.Context, email string) error {
	_, err := s.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("user does not exists: %w", err)
	}

	err = s.userRepo.SendVerificationCode(ctx, email)
	if err != nil {
		return fmt.Errorf("failed to send verification code: %w", err)
	}

	return nil
}

// 7
func (s *UserService) Login(ctx context.Context, login *models.LoginRequest) (string, error) {
	token, err := s.userRepo.Login(ctx, login)
	if err != nil {
		return "", fmt.Errorf("failed to login: %w", err)
	}

	return token, nil
}

func (s *UserService) isEmailExists(email string) bool {
	// Check email format
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		s.logger.Info("Invalid email format")
		return false
	}

	// Get email domain
	domain := parts[1]

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		s.logger.Info("MX record not found", "domain", domain)
		return false
	}

	// Connect to SMTP server
	for _, mx := range mxRecords {
		// Try to connect to SMTP server
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:25", mx.Host), 5*time.Second)
		if err != nil {
			continue
		}
		conn.Close()

		// If connection is successful, email address is valid
		s.logger.Info("Successfully connected to SMTP server", "host", mx.Host)
		return true
	}

	return false
}

func (s *UserService) isValidPassword(password string) bool {
	// Regex to validate password (uppercase-lowercase letters, numbers and special characters)
	var validPasswordRegex = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]{8,}$`
	re := regexp.MustCompile(validPasswordRegex)
	return re.MatchString(password)
}
