package api

import (
	"log/slog"

	"github.com/zohirovs/internal/config"
	"github.com/zohirovs/internal/storage"
)

func Run(config *config.Config, logger *slog.Logger) error {
	// Connect to the database

	db, err := storage.ConnectDB(config)
	if err != nil {
		logger.Error("Error while connecting to MongoDB", slog.String("err", err.Error()))
		return err
	}

	
}
