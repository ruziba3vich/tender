package repos

import (
	"context"

	"github.com/zohirovs/internal/models"
)

type (
	TenderRepo interface {
		CreateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error)
		GetTender(ctx context.Context, id string) (*models.Tender, error)
		UpdateTender(ctx context.Context, updatedTender *models.Tender) (*models.Tender, error)
		DeleteTender(ctx context.Context, id string) (*models.Tender, error)
	}
)
