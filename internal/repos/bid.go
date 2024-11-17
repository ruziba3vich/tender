package repos

import (
	"context"

	"github.com/zohirovs/internal/models"
)

type BidRepo interface {
	CreateBid(ctx context.Context, bid *models.Bid) (*models.Bid, error)
	GetBid(ctx context.Context, id string) (*models.Bid, error)
	ListBidsForTender(ctx context.Context, tenderId string, filter map[string]interface{}) ([]*models.Bid, error)
	UpdateBidStatus(ctx context.Context, bidId string, status string) error
	ListBidsByContractor(ctx context.Context, contractorId string) ([]*models.Bid, error)
}
