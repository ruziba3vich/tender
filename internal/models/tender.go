package models

type (
	Tender struct {
		TenderId      string `json:"tender_id,omitempty"`
		ClientId      string `json:"client_id,omitempty"`
		Title         string `json:"title" binding:"required"`
		Description   string `json:"description" binding:"required"` // Fix typo here
		Budget        int    `json:"budget" binding:"required"`
		Status        string `json:"status,omitempty"` // Optional
		Deadline      string `json:"deadline" binding:"required"`
		AttachmentUrl string `json:"attachment_url" binding:"required"`
	}
	CreateTender struct {
		Title         string `json:"title" binding:"required"`
		Description   string `json:"description" binding:"required"`
		Budget        int    `json:"budget" binding:"required"`
		Deadline      string `json:"deadline" binding:"required"`
		AttachmentUrl string `json:"attachment_url" binding:"required"`
	}

	UpdateTender struct {
		TenderId      string `json:"tender_id"`
		Title         string `json:"title"`
		Description   string `json:"desription"`
		Budget        int    `json:"budget"`
		Deadline      string `json:"deadline"`
		AttachmentUrl string `json:"attachment_url"`
	}
)

// Tender (id, client_id, title, description, deadline, budget, status)
