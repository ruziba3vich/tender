package handler

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/service"
)

// TenderHandler handles tender-related HTTP requests
type TenderHandler struct {
	ser    *service.TenderService
	logger *slog.Logger
}

// NewTenderHandler creates a new TenderHandler
func NewTenderHandler(logger *slog.Logger, ser *service.TenderService) *TenderHandler {
	return &TenderHandler{
		ser:    ser,
		logger: logger,
	}
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// CreateTender godoc
// @Summary      Create a new tender
// @Description  Create a new tender and store it in the database
// @Tags         tenders
// @Accept       json
// @Produce      json
// @Param        tender body     models.Tender true "Tender object"
// @Success      201 {object} models.Tender
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /tenders [post]
func (h *TenderHandler) CreateTender(c *gin.Context) {
	var tender models.Tender
	if err := c.ShouldBindJSON(&tender); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	createdTender, err := h.ser.CreateTender(c.Request.Context(), &tender)
	if err != nil {
		h.logger.Error("failed to create tender", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create tender"})
		return
	}

	c.JSON(http.StatusCreated, createdTender)
}

// GetTender godoc
// @Summary      Get tender by ID
// @Description  Retrieve a tender from the database by its ID
// @Tags         tenders
// @Accept       json
// @Produce      json
// @Param        id path     string true "Tender ID"
// @Success      200 {object} models.Tender
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /tenders/{id} [get]
func (h *TenderHandler) GetTender(c *gin.Context) {
	id := c.Param("id")

	tender, err := h.ser.GetTender(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get tender", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get tender"})
		return
	}

	if tender == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "Tender not found"})
		return
	}

	c.JSON(http.StatusOK, tender)
}

// UpdateTender godoc
// @Summary      Update a tender
// @Description  Update the details of an existing tender
// @Tags         tenders
// @Accept       json
// @Produce      json
// @Param        tender body     models.Tender true "Tender object"
// @Success      200 {object} models.Tender
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /tenders [put]
func (h *TenderHandler) UpdateTender(c *gin.Context) {
	var tender models.Tender
	if err := c.ShouldBindJSON(&tender); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	updatedTender, err := h.ser.UpdateTender(c.Request.Context(), &tender)
	if err != nil {
		h.logger.Error("failed to update tender", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update tender"})
		return
	}

	c.JSON(http.StatusOK, updatedTender)
}

// DeleteTender godoc
// @Summary      Delete a tender
// @Description  Delete a tender from the database by its ID
// @Tags         tenders
// @Accept       json
// @Produce      json
// @Param        id path     string true "Tender ID"
// @Success      204
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /tenders/{id} [delete]
func (h *TenderHandler) DeleteTender(c *gin.Context) {
	id := c.Param("id")

	err := h.ser.DeleteTender(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to delete tender", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete tender"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
