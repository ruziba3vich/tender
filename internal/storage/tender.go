package storage

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

type TenderStorage struct {
	db     *mongo.Collection
	logger *slog.Logger
}

func NewTenderStorage(db *mongo.Database, logger *slog.Logger) *TenderStorage {
	return &TenderStorage{
		db:     db.Collection("tenders"),
		logger: logger,
	}
}

func (s *TenderStorage) CreateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error) {
	if tender.TenderId == "" {
		tender.TenderId = primitive.NewObjectID().Hex()
	}

	deadline, err := time.Parse(time.RFC3339, tender.Deadline)
	if err != nil {
		return nil, fmt.Errorf("invalid deadline format: %w", err)
	}

	if deadline.Before(time.Now()) {
		return nil, errors.New("deadline must be in the future")
	}

	_, err = s.db.InsertOne(ctx, tender)
	if err != nil {
		s.logger.Error("failed to create tender",
			"error", err,
			"tender_id", tender.TenderId)
		return nil, fmt.Errorf("failed to create tender: %w", err)
	}

	return tender, nil
}

func (s *TenderStorage) GetTender(ctx context.Context, id string) (*models.Tender, error) {
	var tender models.Tender

	err := s.db.FindOne(ctx, bson.M{"tender_id": id}).Decode(&tender)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("tender not found: %s", id)
		}
		s.logger.Error("failed to get tender",
			"error", err,
			"tender_id", id)
		return nil, fmt.Errorf("failed to get tender: %w", err)
	}

	return &tender, nil
}

func (s *TenderStorage) UpdateTender(ctx context.Context, updatedTender *models.Tender) (*models.Tender, error) {
	if updatedTender.Deadline != "" {
		deadline, err := time.Parse(time.RFC3339, updatedTender.Deadline)
		if err != nil {
			return nil, fmt.Errorf("invalid deadline format: %w", err)
		}

		if updatedTender.Status == "open" && deadline.Before(time.Now()) {
			return nil, errors.New("deadline must be in the future for open tenders")
		}
	}

	update := bson.M{
		"$set": bson.M{
			"client_id":      updatedTender.ClientId,
			"title":          updatedTender.Title,
			"description":    updatedTender.Description,
			"budget":         updatedTender.Budget,
			"status":         updatedTender.Status,
			"deadline":       updatedTender.Deadline,
			"attachment_url": updatedTender.AttachmentUrl,
		},
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var result models.Tender
	err := s.db.FindOneAndUpdate(
		ctx,
		bson.M{"tender_id": updatedTender.TenderId},
		update,
		opts,
	).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("tender not found: %s", updatedTender.TenderId)
		}
		s.logger.Error("failed to update tender",
			"error", err,
			"tender_id", updatedTender.TenderId)
		return nil, fmt.Errorf("failed to update tender: %w", err)
	}

	return &result, nil
}

func (s *TenderStorage) DeleteTender(ctx context.Context, id string) error {
	result, err := s.db.DeleteOne(ctx, bson.M{"tender_id": id})
	if err != nil {
		s.logger.Error("failed to delete tender",
			"error", err,
			"tender_id", id)
		return fmt.Errorf("failed to delete tender: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("tender not found: %s", id)
	}

	return nil
}

func (s *TenderStorage) ListTenders(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.Tender, error) {
	cursor, err := s.db.Find(ctx, filter, opts)
	if err != nil {
		s.logger.Error("failed to list tenders",
			"error", err)
		return nil, fmt.Errorf("failed to list tenders: %w", err)
	}
	defer cursor.Close(ctx)

	var tenders []*models.Tender
	if err = cursor.All(ctx, &tenders); err != nil {
		return nil, fmt.Errorf("failed to decode tenders: %w", err)
	}

	return tenders, nil
}

func (s *TenderStorage) ListTendersByClient(ctx context.Context, clientID string) ([]*models.Tender, error) {
	filter := bson.M{"client_id": clientID}
	return s.ListTenders(ctx, filter, nil)
}

func (s *TenderStorage) ListOpenTenders(ctx context.Context) ([]*models.Tender, error) {
	now := time.Now().Format(time.RFC3339)
	filter := bson.M{
		"status":   "open",
		"deadline": bson.M{"$gt": now},
	}
	return s.ListTenders(ctx, filter, nil)
}

func (s *TenderStorage) UpdateStatus(ctx context.Context, id string, status models.Status) error {
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	result, err := s.db.UpdateOne(ctx, bson.M{"tender_id": id}, update)
	if err != nil {
		s.logger.Error("failed to update tender status",
			"error", err,
			"tender_id", id,
			"status", status)
		return fmt.Errorf("failed to update tender status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("tender not found: %s", id)
	}

	return nil
}

func (s *TenderStorage) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tender_id", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "client_id", Value: 1},
			},
		},
		{
			Keys: bson.D{
				{Key: "status", Value: 1},
				{Key: "deadline", Value: 1},
			},
		},
	}

	_, err := s.db.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		s.logger.Error("failed to create indexes",
			"error", err)
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	return nil
}
