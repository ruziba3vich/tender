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
	config, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	// Set up logger
	logFile, err := os.OpenFile("application.log", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(logFile, nil))

	// Connect to the database
	db, err := mongo.ConnectDB(config)
	if err != nil {
		logger.Error("Error while connecting to MongoDB", slog.String("err", err.Error()))
		return err
	}

	redisService := redis.New(redis.NewRedisClient(config), logger)

	storage := storage.New(db, logger, redisService)

	service := service.NewService(redisService, logger, storage)

	handler := handler.NewHandler(logger, service)

	// Set up Casbin enforcer
	modelPath := filepath.Join("internal", "pkg", "casbin", "model.conf")
	policyPath := filepath.Join("internal", "pkg", "casbin", "policy.csv")

	enforcer, err := casbin.NewEnforcerSafe(modelPath, policyPath)
	if err != nil {
		log.Fatal(err)
	}

	return app.Run(handler, logger, config, enforcer)
}
