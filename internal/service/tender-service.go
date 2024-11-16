package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/repos"
	"github.com/zohirovs/internal/storage/redis"
)

type TenderService struct {
	tenderRepo  repos.TenderRepo
	tenderCache *redis.TenderCaching
	logger      *slog.Logger
}

func NewTenderService(tenderRepo repos.TenderRepo, cache *redis.TenderCaching, logger *slog.Logger) *TenderService {
	return &TenderService{
		tenderRepo:  tenderRepo,
		tenderCache: cache,
		logger:      logger,
	}
}

func (s *TenderService) CreateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error) {
	createdTender, err := s.tenderRepo.CreateTender(ctx, tender)
	if err != nil {
		return nil, fmt.Errorf("failed to create tender: %w", err)
	}

	if s.tenderCache != nil {
		if err := s.tenderCache.Set(ctx, createdTender); err != nil {
			s.logger.Warn("failed to cache tender after creation",
				"error", err,
				"tender_id", createdTender.TenderId)
		}
	}

	return createdTender, nil
}

func (s *TenderService) GetTender(ctx context.Context, id string) (*models.Tender, error) {
	if s.tenderCache != nil {
		tender, err := s.tenderCache.Get(ctx, id)
		if err == nil {
			return tender, nil
		}
		s.logger.Debug("cache miss for tender",
			"tender_id", id,
			"error", err)
	}

	tender, err := s.tenderRepo.GetTender(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get tender: %w", err)
	}

	if s.tenderCache != nil {
		if err := s.tenderCache.Set(ctx, tender); err != nil {
			s.logger.Warn("failed to update cache after getting tender",
				"error", err,
				"tender_id", id)
		}
	}

	return tender, nil
}

func (s *TenderService) UpdateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error) {
	updatedTender, err := s.tenderRepo.UpdateTender(ctx, tender)
	if err != nil {
		return nil, fmt.Errorf("failed to update tender: %w", err)
	}

	if s.tenderCache != nil {
		if err := s.tenderCache.Set(ctx, updatedTender); err != nil {
			s.logger.Warn("failed to update cache after tender update",
				"error", err,
				"tender_id", tender.TenderId)
		}
	}

	return updatedTender, nil
}

func (s *TenderService) DeleteTender(ctx context.Context, id string) error {
	if err := s.tenderRepo.DeleteTender(ctx, id); err != nil {
		return fmt.Errorf("failed to delete tender: %w", err)
	}

	if s.tenderCache != nil {
		if err := s.tenderCache.Delete(ctx, id); err != nil {
			s.logger.Warn("failed to delete tender from cache",
				"error", err,
				"tender_id", id)
		}
	}

	return nil
}

func (s *TenderService) UpdateTenderStatus(ctx context.Context, id string, status models.Status) error {
	if err := s.tenderRepo.UpdateStatus(ctx, id, status); err != nil {
		return fmt.Errorf("failed to update tender status: %w", err)
	}

	if s.tenderCache != nil {
		if err := s.tenderCache.Delete(ctx, id); err != nil {
			s.logger.Warn("failed to invalidate cache after status update",
				"error", err,
				"tender_id", id)
		}
	}

	return nil
}
