package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/zohirovs/internal/config"
	jwttokens "github.com/zohirovs/internal/jwt"
	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/repos"
	"github.com/zohirovs/internal/storage/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserStorage struct {
	db        *mongo.Collection
	logger    *slog.Logger
	userCache *redis.UserCaching
	cfg       *config.Config
}

func NewUserStorage(db *mongo.Database, logger *slog.Logger, cache *redis.UserCaching) repos.UserRepo {
	return &UserStorage{
		db:        db.Collection("Users"),
		logger:    logger,
		userCache: cache,
	}
}

// 1
func (u *UserStorage) RegisterUser(ctx context.Context, user *models.User) (string, error) {
	// Set creation timestamp
	user.CreatedAt = time.Now()

	// Hash the password
	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		u.logger.Error("failed to hash password", "error", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Start a MongoDB session
	session, err := u.db.Database().Client().StartSession()
	if err != nil {
		u.logger.Error("failed to start session", "error", err)
		return "", err
	}
	defer session.EndSession(ctx)

	// Start transaction
	var insertedID primitive.ObjectID
	err = mongo.WithSession(ctx, session, func(sessionContext mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return err
		}

		// Insert user into MongoDB
		result, err := u.db.InsertOne(sessionContext, user)
		if err != nil {
			u.logger.Error("failed to insert user", "error", err)
			return err
		}

		// Get the inserted ID
		var ok bool
		insertedID, ok = result.InsertedID.(primitive.ObjectID)
		if !ok {
			u.logger.Error("failed to convert inserted ID to ObjectID")
			return err
		}

		// Commit the transaction
		if err = session.CommitTransaction(sessionContext); err != nil {
			u.logger.Error("failed to commit transaction", "error", err)
			return err
		}

		// Cache the user data only if MongoDB transaction was successful
		if err := u.userCache.SetUser(ctx, user); err != nil {
			u.logger.Warn("failed to cache user data", "error", err)
			// Don't return error here as the user is already saved in DB
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	tokenClaims := models.TokenClaims{
		UserID:   insertedID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}

	token, err := jwttokens.GenerateAccessToken(u.cfg.JWT.SecretKey, &tokenClaims)
	if err != nil {
		return "", err
	}

	return token, nil
}

// 2
func (u *UserStorage) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	// First try to get user from cache
	if cachedUser, err := u.userCache.GetUserByUserID(ctx, userID); err == nil {
		return cachedUser, nil
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		u.logger.Error("failed to convert userID to ObjectID", "error", err)
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	// Find user in MongoDB
	var user models.User
	err = u.db.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		u.logger.Error("failed to fetch user from database", "error", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Cache the user data for future requests
	if err := u.userCache.SetUser(ctx, &user); err != nil {
		u.logger.Warn("failed to cache user data", "error", err)
		// Don't return error here as we still have the user data
	}

	return &user, nil
}

// 3
func (u *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// First try to get user from cache
	if cachedUser, err := u.userCache.GetUserByEmail(ctx, email); err == nil {
		return cachedUser, nil
	}

	// Find user in MongoDB
	var user models.User
	err := u.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		u.logger.Error("failed to fetch user from database", "error", err)
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	// Cache the user data for future requests
	if err := u.userCache.SetUser(ctx, &user); err != nil {
		u.logger.Warn("failed to cache user data", "error", err)
		// Don't return error here as we still have the user data
	}

	return &user, nil
}

// 4
func (u *UserStorage) ChangeUserRole(ctx context.Context, userID string, role string) error {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		u.logger.Error("failed to convert userID to ObjectID", "error", err)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Start a MongoDB transaction
	session, err := u.db.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Update user role in MongoDB
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": bson.M{"role": role}}
		result, err := u.db.UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update user role: %w", err)
		}

		if result.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return fmt.Errorf("user not found")
		}

		// Get updated user data
		var updatedUser models.User
		err = u.db.FindOne(sc, filter).Decode(&updatedUser)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to fetch updated user: %w", err)
		}

		// Commit the transaction
		if err = session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Update cache only after successful MongoDB transaction
		if err := u.userCache.SetUser(ctx, &updatedUser); err != nil {
			u.logger.Warn("failed to update user cache", "error", err)
			// Don't return error here as the DB update was successful
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 5
func (u *UserStorage) ChangeUserPassword(ctx context.Context, resetPassword *models.ResetPassword) error {
	// Verify email and code first
	err := u.verifyEmail(ctx, resetPassword.Email, resetPassword.Code)
	if err != nil {
		return fmt.Errorf("invalid verification code: %w", err)
	}

	// Start a session for transaction
	session, err := u.db.Database().Client().StartSession()
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Get current user data
		var user models.User
		filter := bson.M{"email": resetPassword.Email}
		err = u.db.FindOne(sc, filter).Decode(&user)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to fetch user: %w", err)
		}

		// Hash new password
		hashedPassword, err := u.hashPassword(resetPassword.NewPassword)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to hash new password: %w", err)
		}

		// Update password in MongoDB
		update := bson.M{"$set": bson.M{"password": hashedPassword}}
		result, err := u.db.UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to update password: %w", err)
		}

		if result.MatchedCount == 0 {
			session.AbortTransaction(sc)
			return fmt.Errorf("user not found")
		}

		// Get updated user data
		var updatedUser models.User
		err = u.db.FindOne(sc, filter).Decode(&updatedUser)
		if err != nil {
			session.AbortTransaction(sc)
			return fmt.Errorf("failed to fetch updated user: %w", err)
		}

		// Commit the transaction
		if err = session.CommitTransaction(sc); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Update cache only after successful MongoDB transaction
		if err := u.userCache.SetUser(ctx, &updatedUser); err != nil {
			u.logger.Warn("failed to update user cache", "error", err)
			// Don't return error here as the DB update was successful
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 6
func (u *UserStorage) SendVerificationCode(ctx context.Context, email string) error {
	code := u.generateVerificationCode()

	m := gomail.NewMessage()
	m.SetHeader("From", u.cfg.Email.SmtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %d", code))

	d := gomail.NewDialer(u.cfg.Email.SmtpHost, u.cfg.Email.SmtpPort, u.cfg.Email.SmtpUser, u.cfg.Email.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return u.userCache.StoreEmailAndCode(ctx, email, code)
}

// 7
func (u *UserStorage) Login(ctx context.Context, login *models.LoginRequest) (string, error) {
	// Find user by email
	user, err := u.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user exists
	if user == nil {
		return "", errors.New("user not found")
	}

	// Verify password
	isValid, err := u.checkPassword(user.Password, login.Password)
	if err != nil {
		return "", fmt.Errorf("failed to verify password: %w", err)
	}

	if !isValid {
		return "", errors.New("invalid credentials")
	}

	tokenClaims := models.TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}

	token, err := jwttokens.GenerateAccessToken(u.cfg.JWT.SecretKey, &tokenClaims)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *UserStorage) hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

func (u *UserStorage) checkPassword(hashedPassword, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return false, fmt.Errorf("failed to compare password: %w", err)
	}
	return true, nil
}

func (u *UserStorage) generateVerificationCode() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(90000) + 10000
}

func (u *UserStorage) verifyEmail(ctx context.Context, email string, code int) error {
	c, err := u.userCache.GetCodeByEmail(ctx, email)
	if err != nil {
		return err
	}
	if c != code {
		return errors.New("invalide code")
	}
	return nil
}
