package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/repos"
)

type BidService struct {
	bidRepo    repos.BidRepo
	tenderRepo repos.TenderRepo
	logger     *slog.Logger
}

func NewBidService(bidRepo repos.BidRepo, tenderRepo repos.TenderRepo, logger *slog.Logger) *BidService {
	return &BidService{
		bidRepo:    bidRepo,
		tenderRepo: tenderRepo,
		logger:     logger,
	}
}

func (s *BidService) CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error) {
	// Validate tender exists and is open
	tender, err := s.tenderRepo.GetTender(ctx, bid.TenderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get tender: %w", err)
	}

	if tender.Status != "OPEN" {
		return nil, fmt.Errorf("tender is not open for bidding")
	}

	deadline, err := time.Parse(time.RFC3339, tender.Deadline)
	if err != nil {
		return nil, fmt.Errorf("invalid tender deadline format: %w", err)
	}

	if time.Now().After(deadline) {
		return nil, fmt.Errorf("tender deadline has passed")
	}

	// Create the bid
	createdBid, err := s.bidRepo.CreateBid(ctx, bid)
	if err != nil {
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return createdBid, nil
}

func (s *BidService) GetBid(ctx context.Context, id string) (*models.Bid, error) {
	bid, err := s.bidRepo.GetBid(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get bid: %w", err)
	}
	return bid, nil
}

func (s *BidService) ListBidsForTender(ctx context.Context, tenderId string, filter map[string]interface{}) ([]*models.Bid, error) {
	bids, err := s.bidRepo.ListBidsForTender(ctx, tenderId, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list bids: %w", err)
	}
	return bids, nil
}

func (s *BidService) UpdateBidStatus(ctx context.Context, bidId string, status string) error {
	err := s.bidRepo.UpdateBidStatus(ctx, bidId, status)
	if err != nil {
		return fmt.Errorf("failed to update bid status: %w", err)
	}
	return nil
}

func (s *BidService) ListBidsByContractor(ctx context.Context, contractorId string) ([]*models.Bid, error) {
	bids, err := s.bidRepo.ListBidsByContractor(ctx, contractorId)
	if err != nil {
		return nil, fmt.Errorf("failed to list contractor bids: %w", err)
	}
	return bids, nil
}
