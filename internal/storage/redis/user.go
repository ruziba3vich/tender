package redis

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zohirovs/internal/models"
)

type UserCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewUserCaching(client *redis.Client, logger *slog.Logger) *UserCaching {
	return &UserCaching{
		redisClient: client,
		logger:      logger,
	}
}

func (u *UserCaching) SetUser(ctx context.Context, user *models.User) error {
	err := u.redisClient.HSet(ctx, "user:"+user.ID, map[string]interface{}{
		"id":        user.ID,
		"username":  user.Username,
		"email":     user.Email,
		"password":  user.Password,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}).Err()

	if err != nil {
		u.logger.Error("failed to set user in redis",
			"error", err,
			"user_id", user.ID,
		)
		return err
	}

	return nil
}

func (u *UserCaching) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	user := models.User{}

	values, err := u.redisClient.HGetAll(ctx, "user:"+userID).Result()
	if err != nil {
		u.logger.Error("failed to get user from redis",
			"error", err,
			"user_id", userID,
		)
		return &user, err
	}

	if len(values) == 0 {
		return &user, redis.Nil
	}
	createdAt, _ := time.Parse(time.RFC3339, values["createdAt"])
	updatedAt, _ := time.Parse(time.RFC3339, values["updatedAt"])
	deletedAt, _ := time.Parse(time.RFC3339, values["deletedAt"])

	user = models.User{
		ID:        values["id"],
		Username:  values["username"],
		Email:     values["email"],
		Password:  values["password"],
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		DeletedAt: deletedAt,
	}
	return &user, nil
}

func (u *UserCaching) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user := models.User{}

	// Get all users from Redis
	keys, err := u.redisClient.Keys(ctx, "user:*").Result()
	if err != nil {
		u.logger.Error("failed to get user keys from redis",
			"error", err,
			"email", email,
		)
		return &user, err
	}

	// Iterate through users to find matching email
	for _, key := range keys {
		values, err := u.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			u.logger.Error("failed to get user from redis",
				"error", err,
				"key", key,
			)
			continue
		}

		if values["email"] == email {
			createdAt, _ := time.Parse(time.RFC3339, values["createdAt"])
			updatedAt, _ := time.Parse(time.RFC3339, values["updatedAt"])
			deletedAt, _ := time.Parse(time.RFC3339, values["deletedAt"])

			user = models.User{
				ID:        values["id"],
				Username:  values["username"],
				Email:     values["email"],
				Password:  values["password"],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				DeletedAt: deletedAt,
			}
			return &user, nil
		}
	}

	return &user, redis.Nil
}

func (u *UserCaching) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := models.User{}

	// Get all users from Redis
	keys, err := u.redisClient.Keys(ctx, "user:*").Result()
	if err != nil {
		u.logger.Error("failed to get user keys from redis",
			"error", err,
			"username", username,
		)
		return &user, err
	}

	// Iterate through users to find matching username
	for _, key := range keys {
		values, err := u.redisClient.HGetAll(ctx, key).Result()
		if err != nil {
			u.logger.Error("failed to get user from redis",
				"error", err,
				"key", key,
			)
			continue
		}

		if values["username"] == username {
			createdAt, _ := time.Parse(time.RFC3339, values["createdAt"])
			updatedAt, _ := time.Parse(time.RFC3339, values["updatedAt"])
			deletedAt, _ := time.Parse(time.RFC3339, values["deletedAt"])

			user = models.User{
				ID:        values["id"],
				Username:  values["username"],
				Email:     values["email"],
				Password:  values["password"],
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
				DeletedAt: deletedAt,
			}
			return &user, nil
		}
	}

	return &user, redis.Nil
}

func (u *UserCaching) StoreEmailAndCode(ctx context.Context, email string, code int) error {
	codeKey := "verification_code:" + email
	err := u.redisClient.Set(ctx, codeKey, code, time.Minute*1).Err()
	if err != nil {
		u.logger.Error("Error while storing verification code", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (u *UserCaching) GetCodeByEmail(ctx context.Context, email string) (int, error) {
	codeKey := "verification_code:" + email
	codeStr, err := u.redisClient.Get(ctx, codeKey).Result()
	if err == redis.Nil {
		return 0, nil
	} else if err != nil {
		u.logger.Error("Error while getting verification code", slog.String("error", err.Error()))
		return 0, err
	}

	var code int
	_, err = fmt.Sscanf(codeStr, "%d", &code)
	if err != nil {
		u.logger.Error("Error while parsing verification code", slog.String("error", err.Error()))
		return 0, err
	}

	return code, nil
}
