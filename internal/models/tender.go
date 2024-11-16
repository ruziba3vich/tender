package models

type (
	Tender struct {
		TenderId    string  `json:"tender_id"`
		ClientId    string  `json:"client_id"`
		Title       string  `json:"title"`
		Description string  `json:"desription"`
		Budget      float32 `json:"budget"`
		Status      Status  `json:"status"`
		Deadline    string  `json:"deadline"`
	}
)

// Tender (id, client_id, title, description, deadline, budget, status)