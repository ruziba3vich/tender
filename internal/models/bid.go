package models

import "time"

type (
	Bid struct {
		BidId        string
		TenderId     string
		ContractorId string
		Price        float32
		DeliveryTime time.Time
		Comment      string
		Status       Status
	}
)

// Bid (id, tender_id, contractor_id, price, delivery_time, comments, status)
