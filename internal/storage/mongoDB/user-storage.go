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

func NewUserStorage(db *mongo.Database, cfg *config.Config, logger *slog.Logger, cache *redis.UserCaching) repos.UserRepo {
	return &UserStorage{
		db:        db.Collection("Users"),
		logger:    logger,
		userCache: cache,
		cfg:       cfg,
	}
}

// 1
func (u *UserStorage) RegisterUser(ctx context.Context, user *models.User) (string, error) {
	u.logger.Info("starting user registration", "email", user.Email, "username", user.Username)

	// Foydalanuvchini yaratish vaqtini belgilash
	user.CreatedAt = time.Now()

	// Parolni hash qilish
	hashedPassword, err := u.hashPassword(user.Password)
	if err != nil {
		u.logger.Error("failed to hash password", "error", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hashedPassword)

	// Foydalanuvchini MongoDB'ga yozish
	result, err := u.db.InsertOne(ctx, user)
	if err != nil {
		u.logger.Error("failed to insert user", "error", err)
		return "", err
	}

	// Yozilgan ID'ni olish
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		u.logger.Error("failed to convert inserted ID to ObjectID")
		return "", fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	// Foydalanuvchini cache'ga saqlash
	if err := u.userCache.SetUser(ctx, user); err != nil {
		u.logger.Warn("failed to cache user data", "error", err)
		// Bu yerda xatoni qaytarmaymiz, chunki foydalanuvchi bazaga muvaffaqiyatli saqlandi
	}

	// JWT token generatsiya qilish uchun token ma'lumotlarini tayyorlash
	tokenClaims := models.TokenClaims{
		UserID:   insertedID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}

	token, err := jwttokens.GenerateAccessToken(u.cfg.JWT.SecretKey, &tokenClaims)
	if err != nil {
		u.logger.Error("failed to generate access token", "error", err)
		return "", err
	}

	u.logger.Info("user registration completed successfully", "userID", insertedID.Hex())
	return token, nil
}

// 2
func (u *UserStorage) GetUserByUserID(ctx context.Context, userID string) (*models.User, error) {
	u.logger.Info("fetching user by ID", "userID", userID)

	// First try to get user from cache
	if cachedUser, err := u.userCache.GetUserByUserID(ctx, userID); err == nil {
		u.logger.Debug("user found in cache", "userID", userID)
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
			u.logger.Warn("user not found", "userID", userID)
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

	u.logger.Info("successfully fetched user", "userID", userID)
	return &user, nil
}

// 3
func (u *UserStorage) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	u.logger.Info("fetching user by email", "email", email)

	// First try to get user from cache
	if cachedUser, err := u.userCache.GetUserByEmail(ctx, email); err == nil {
		u.logger.Debug("user found in cache", "email", email)
		return cachedUser, nil
	}

	// Find user in MongoDB
	var user models.User
	err := u.db.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			u.logger.Warn("user not found", "email", email)
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

	u.logger.Info("successfully fetched user", "email", email)
	return &user, nil
}

// 4
func (u *UserStorage) ChangeUserRole(ctx context.Context, userID string, role string) error {
	u.logger.Info("changing user role", "userID", userID, "newRole", role)

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		u.logger.Error("failed to convert userID to ObjectID", "error", err)
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Start a MongoDB transaction
	session, err := u.db.Database().Client().StartSession()
	if err != nil {
		u.logger.Error("failed to start session", "error", err)
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		u.logger.Debug("starting MongoDB transaction")
		if err := session.StartTransaction(); err != nil {
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Update user role in MongoDB
		filter := bson.M{"_id": objectID}
		update := bson.M{"$set": bson.M{"role": role}}
		result, err := u.db.UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			u.logger.Error("failed to update user role", "error", err)
			return fmt.Errorf("failed to update user role: %w", err)
		}

		if result.MatchedCount == 0 {
			session.AbortTransaction(sc)
			u.logger.Warn("user not found", "userID", userID)
			return fmt.Errorf("user not found")
		}

		// Get updated user data
		var updatedUser models.User
		err = u.db.FindOne(sc, filter).Decode(&updatedUser)
		if err != nil {
			session.AbortTransaction(sc)
			u.logger.Error("failed to fetch updated user", "error", err)
			return fmt.Errorf("failed to fetch updated user: %w", err)
		}

		// Commit the transaction
		if err = session.CommitTransaction(sc); err != nil {
			u.logger.Error("failed to commit transaction", "error", err)
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		u.logger.Debug("successfully committed transaction")

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

	u.logger.Info("successfully changed user role", "userID", userID, "newRole", role)
	return nil
}

// 5
func (u *UserStorage) ChangeUserPassword(ctx context.Context, resetPassword *models.ResetPassword) error {
	u.logger.Info("starting password change process", "email", resetPassword.Email)

	// Verify email and code first
	err := u.verifyEmail(ctx, resetPassword.Email, resetPassword.Code)
	if err != nil {
		u.logger.Error("email verification failed", "error", err)
		return fmt.Errorf("invalid verification code: %w", err)
	}

	// Start a session for transaction
	session, err := u.db.Database().Client().StartSession()
	if err != nil {
		u.logger.Error("failed to start session", "error", err)
		return fmt.Errorf("failed to start session: %w", err)
	}
	defer session.EndSession(ctx)

	// Execute transaction
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err := session.StartTransaction(); err != nil {
			u.logger.Error("failed to start transaction", "error", err)
			return fmt.Errorf("failed to start transaction: %w", err)
		}

		// Get current user data
		var user models.User
		filter := bson.M{"email": resetPassword.Email}
		err = u.db.FindOne(sc, filter).Decode(&user)
		if err != nil {
			session.AbortTransaction(sc)
			u.logger.Error("failed to fetch user", "error", err)
			return fmt.Errorf("failed to fetch user: %w", err)
		}

		// Hash new password
		hashedPassword, err := u.hashPassword(resetPassword.NewPassword)
		if err != nil {
			session.AbortTransaction(sc)
			u.logger.Error("failed to hash password", "error", err)
			return fmt.Errorf("failed to hash new password: %w", err)
		}

		// Update password in MongoDB
		update := bson.M{"$set": bson.M{"password": hashedPassword}}
		result, err := u.db.UpdateOne(sc, filter, update)
		if err != nil {
			session.AbortTransaction(sc)
			u.logger.Error("failed to update password in database", "error", err)
			return fmt.Errorf("failed to update password: %w", err)
		}

		if result.MatchedCount == 0 {
			session.AbortTransaction(sc)
			u.logger.Warn("user not found during password update", "email", resetPassword.Email)
			return fmt.Errorf("user not found")
		}

		// Get updated user data
		var updatedUser models.User
		err = u.db.FindOne(sc, filter).Decode(&updatedUser)
		if err != nil {
			session.AbortTransaction(sc)
			u.logger.Error("failed to fetch updated user", "error", err)
			return fmt.Errorf("failed to fetch updated user: %w", err)
		}

		// Commit the transaction
		if err = session.CommitTransaction(sc); err != nil {
			u.logger.Error("failed to commit transaction", "error", err)
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Update cache only after successful MongoDB transaction
		if err := u.userCache.SetUser(ctx, &updatedUser); err != nil {
			u.logger.Warn("failed to update user cache", "error", err)
			// Don't return error here as the DB update was successful
		}

		u.logger.Info("password successfully changed", "email", resetPassword.Email)
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

// 6
func (u *UserStorage) SendVerificationCode(ctx context.Context, email string) error {
	u.logger.Info("starting verification code sending process", "email", email)

	code := u.generateVerificationCode()
	u.logger.Debug("verification code generated", "email", email)

	m := gomail.NewMessage()
	m.SetHeader("From", u.cfg.Email.SmtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %d", code))

	d := gomail.NewDialer(u.cfg.Email.SmtpHost, u.cfg.Email.SmtpPort, u.cfg.Email.SmtpUser, u.cfg.Email.SmtpPass)

	if err := d.DialAndSend(m); err != nil {
		u.logger.Error("failed to send verification email", "error", err, "email", email)
		return err
	}

	if err := u.userCache.StoreEmailAndCode(ctx, email, code); err != nil {
		u.logger.Error("failed to store verification code in cache", "error", err, "email", email)
		return err
	}

	u.logger.Info("verification code sent successfully", "email", email)
	return nil
}

// 7
func (u *UserStorage) Login(ctx context.Context, login *models.LoginRequest) (string, error) {
	u.logger.Info("starting login process", "email", login.Email)

	// Find user by email
	user, err := u.GetUserByEmail(ctx, login.Email)
	if err != nil {
		u.logger.Error("failed to get user during login", "error", err, "email", login.Email)
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user exists
	if user == nil {
		u.logger.Warn("login attempt for non-existent user", "email", login.Email)
		return "", errors.New("user not found")
	}

	// Verify password
	isValid, err := u.checkPassword(user.Password, login.Password)
	if err != nil {
		u.logger.Error("failed to verify password", "error", err, "email", login.Email)
		return "", fmt.Errorf("failed to verify password: %w", err)
	}

	if !isValid {
		u.logger.Warn("invalid login credentials", "email", login.Email)
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
		u.logger.Error("failed to generate access token", "error", err, "email", login.Email)
		return "", err
	}

	u.logger.Info("login successful", "email", login.Email)
	return token, nil
}

func (u *UserStorage) hashPassword(password string) (string, error) {
	u.logger.Debug("hashing password")
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		u.logger.Error("failed to hash password", "error", err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	u.logger.Debug("password hashed successfully")
	return string(hashedPassword), nil
}

func (u *UserStorage) checkPassword(hashedPassword, password string) (bool, error) {
	u.logger.Debug("checking password")
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		u.logger.Debug("password check failed")
		return false, fmt.Errorf("failed to compare password: %w", err)
	}
	u.logger.Debug("password check successful")
	return true, nil
}

func (u *UserStorage) generateVerificationCode() int {
	u.logger.Debug("generating verification code")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := r.Intn(90000) + 10000
	u.logger.Debug("verification code generated")
	return code
}

func (u *UserStorage) verifyEmail(ctx context.Context, email string, code int) error {
	u.logger.Debug("verifying email code", "email", email)
	c, err := u.userCache.GetCodeByEmail(ctx, email)
	if err != nil {
		u.logger.Error("failed to get code from cache", "error", err, "email", email)
		return err
	}
	if c != code {
		u.logger.Warn("invalid verification code provided", "email", email)
		return errors.New("invalide code")
	}
	u.logger.Debug("email code verified successfully", "email", email)
	return nil
}
