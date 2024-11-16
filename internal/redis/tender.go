package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/zohirovs/internal/models"
)

const (
	tenderKeyPrefix = "tender:"
	defaultTTL      = 24 * time.Hour
)

type TenderCaching struct {
	redisClient *redis.Client
	logger      *slog.Logger
}

func NewTenderCaching(client *redis.Client, logger *slog.Logger) *TenderCaching {
	return &TenderCaching{
		redisClient: client,
		logger:      logger,
	}
}

func (tc *TenderCaching) Set(ctx context.Context, tender *models.Tender) error {
	key := tc.generateKey(tender.TenderId)

	data, err := json.Marshal(tender)
	if err != nil {
		return fmt.Errorf("failed to marshal tender: %w", err)
	}

	err = tc.redisClient.Set(ctx, key, data, defaultTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set tender in cache: %w", err)
	}

	return nil
}

func (tc *TenderCaching) Get(ctx context.Context, id string) (*models.Tender, error) {
	key := tc.generateKey(id)

	data, err := tc.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("tender not found in cache")
		}
		return nil, fmt.Errorf("failed to get tender from cache: %w", err)
	}

	var tender models.Tender
	if err := json.Unmarshal(data, &tender); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tender: %w", err)
	}

	return &tender, nil
}

func (tc *TenderCaching) Delete(ctx context.Context, id string) error {
	key := tc.generateKey(id)

	err := tc.redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete tender from cache: %w", err)
	}

	return nil
}

func (tc *TenderCaching) generateKey(id string) string {
	return fmt.Sprintf("%s%s", tenderKeyPrefix, id)
}
