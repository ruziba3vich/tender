// repos/tender.go
package repos

import (
	"context"

	"github.com/zohirovs/internal/models"
)

type TenderRepo interface {
	CreateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error)
	GetTender(ctx context.Context, id string) (*models.Tender, error)
	UpdateTender(ctx context.Context, tender *models.Tender) (*models.Tender, error)
	DeleteTender(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status models.Status) error
}
