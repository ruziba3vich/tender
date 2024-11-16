package repos

import "github.com/zohirovs/internal/models"

type (
	TenderRepo interface {
		CreateTender(*models.Tender) (*models.Tender, error)
	}
)
