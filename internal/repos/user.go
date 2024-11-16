package repos

import (
	"context"

	"github.com/zohirovs/internal/models"
)

type UserRepo interface {
	// 1
	RegisterUser(ctx context.Context, user *models.User) (string, error)
	// 2
	GetUserByUserID(ctx context.Context, userID string) (*models.User, error)
	// 3
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	// 4
	ChangeUserRole(ctx context.Context, userID string, role string) error
	// 5
	ChangeUserPassword(ctx context.Context, resetPassword *models.ResetPassword) error
	// 6
	SendVerificationCode(ctx context.Context, email string) error
	// 7
	Login(ctx context.Context, login *models.LoginRequest) (string, error)
}
