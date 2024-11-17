package models

type CreateBidRequest struct {
	TenderID    string  `json:"tender_id" binding:"required"`
	Amount      float64 `json:"amount" binding:"required"`
	Description string  `json:"description"`
}
