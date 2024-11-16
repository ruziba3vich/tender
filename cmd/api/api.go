package api

import (
	"log"
	"log/slog"
	"os"

	"github.com/zohirovs/internal/config"
	mongo "github.com/zohirovs/internal/storage/mongoDB"
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

}
