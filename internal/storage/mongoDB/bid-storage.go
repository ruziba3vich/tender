package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/zohirovs/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BidStorage struct {
	db     *mongo.Collection
	logger *slog.Logger
}

func NewBidStorage(db *mongo.Database, logger *slog.Logger) *BidStorage {
	return &BidStorage{
		db:     db.Collection("Bids"),
		logger: logger,
	}
}

func (s *BidStorage) CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error) {
	if bid.BidId == "" {
		bid.BidId = primitive.NewObjectID().Hex()
	}

	bid.CreatedAt = time.Now()
	bid.Status = "pending"

	_, err := s.db.InsertOne(ctx, bid)
	if err != nil {
		s.logger.Error("failed to create bid",
			"error", err,
			"bid_id", bid.BidId)
		return nil, fmt.Errorf("failed to create bid: %w", err)
	}

	return bid, nil
}

func (s *BidStorage) GetBid(ctx context.Context, id string) (*models.Bid, error) {
	var bid models.Bid

	err := s.db.FindOne(ctx, bson.M{"bid_id": id}).Decode(&bid)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("bid not found: %s", id)
		}
		s.logger.Error("failed to get bid",
			"error", err,
			"bid_id", id)
		return nil, fmt.Errorf("failed to get bid: %w", err)
	}

	return &bid, nil
}

func (s *BidStorage) ListBidsForTender(ctx context.Context, tenderId string, filter map[string]interface{}) ([]*models.Bid, error) {
	baseFilter := bson.M{"tender_id": tenderId}
	if filter != nil {
		for k, v := range filter {
			baseFilter[k] = v
		}
	}
	cursor, err := s.db.Find(ctx, baseFilter)
	if err != nil {
		s.logger.Error("failed to list bids",
			"error", err,
			"tender_id", tenderId)
		return nil, fmt.Errorf("failed to list bids: %w", err)
	}
	defer cursor.Close(ctx)

	var bids []*models.Bid
	if err = cursor.All(ctx, &bids); err != nil {
		return nil, fmt.Errorf("failed to decode bids: %w", err)
	}

	return bids, nil
}

func (s *BidStorage) UpdateBidStatus(ctx context.Context, bidId string, status string) error {
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	result, err := s.db.UpdateOne(ctx, bson.M{"bid_id": bidId}, update)
	if err != nil {
		s.logger.Error("failed to update bid status",
			"error", err,
			"bid_id", bidId)
		return fmt.Errorf("failed to update bid status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("bid not found: %s", bidId)
	}

	return nil
}

func (s *BidStorage) ListBidsByContractor(ctx context.Context, contractorId string) ([]*models.Bid, error) {
	cursor, err := s.db.Find(ctx, bson.M{"contractor_id": contractorId})
	if err != nil {
		s.logger.Error("failed to list contractor bids",
			"error", err,
			"contractor_id", contractorId)
		return nil, fmt.Errorf("failed to list contractor bids: %w", err)
	}
	defer cursor.Close(ctx)

	var bids []*models.Bid
	if err = cursor.All(ctx, &bids); err != nil {
		return nil, fmt.Errorf("failed to decode bids: %w", err)
	}

	return bids, nil
}

func (s *BidStorage) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "bid_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "tender_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "contractor_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "price", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "delivery_time", Value: 1},
			},
		},
	}

	_, err := s.db.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		s.logger.Error("failed to create bid indexes",
			"error", err)
		return fmt.Errorf("failed to create bid indexes: %w", err)
	}

	return nil
}
