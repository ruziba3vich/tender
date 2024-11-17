package models

import "time"

type Bid struct {
	BidId        string    `json:"bid_id,omitempty" bson:"bid_id"`
	TenderId     string    `json:"tender_id" bson:"tender_id" binding:"required"`
	ContractorId string    `json:"contractor_id" bson:"contractor_id"`
	Price        float64   `json:"price" bson:"price" binding:"required"`
	DeliveryTime int       `json:"delivery_time" bson:"delivery_time" binding:"required"` // in days
	Comments     string    `json:"comments" bson:"comments"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
	Status       string    `json:"status" bson:"status"` // pending, accepted, rejected
}

type CreateBid struct {
	TenderId     string  `json:"tender_id" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	DeliveryTime int     `json:"delivery_time" binding:"required"`
	Comments     string  `json:"comments"`
}
