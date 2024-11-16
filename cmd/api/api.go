package api

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/casbin/casbin"
	"github.com/zohirovs/internal/config"
	"github.com/zohirovs/internal/http/app"
	"github.com/zohirovs/internal/http/handler"
	"github.com/zohirovs/internal/service"
	"github.com/zohirovs/internal/storage"
	mongo "github.com/zohirovs/internal/storage/mongoDB"
	"github.com/zohirovs/internal/storage/redis"
)

func Run() error {
	// Load configuration and handle errors
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
		return err // Add return statement after log.Fatal for better error handling
	}

	// Set up logger with file output
	logFile, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer logFile.Close()

	// Initialize structured logger with JSON format
	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// Initialize MongoDB connection
	db, err := mongo.ConnectDB(cfg)
	if err != nil {
		logger.Error("Error while connecting to MongoDB", slog.String("err", err.Error()))
		return err
	}

	// Initialize Redis client and service
	redisService := redis.New(redis.NewRedisClient(cfg), logger)

	// Initialize storage layer with MongoDB and Redis
	storage := storage.New(db, logger, redisService)

	// Initialize service layer
	service := service.NewService(redisService, logger, storage)

	// Initialize HTTP handler
	handler := handler.NewHandler(logger, service)

	// Set up Casbin enforcer for authorization
	modelPath := filepath.Join("internal", "casbin", "model.conf")
	policyPath := filepath.Join("internal", "casbin", "policy.csv")

	enforcer, err := casbin.NewEnforcerSafe(modelPath, policyPath)
	if err != nil {
		logger.Error("Error initializing Casbin enforcer", slog.String("err", err.Error()))
		return err
	}

	// Start the HTTP server
	return app.Run(handler, logger, cfg, enforcer)
}
