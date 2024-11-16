package handler

import (
	"net/http"

	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/config"
	"github.com/zohirovs/internal/middleware"
	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/service"
)

// TenderHandler handles tender-related HTTP requests
type TenderHandler struct {
	ser    *service.TenderService
	logger *slog.Logger
	cfg    *config.Config
}

// NewTenderHandler creates a new TenderHandler
func NewTenderHandler(logger *slog.Logger, ser *service.TenderService, cfg *config.Config) *TenderHandler {
	return &TenderHandler{
		ser:    ser,
		logger: logger,
		cfg:    cfg,
	}
}

// ErrorResponse represents a generic error response
type ErrorResponse struct {
	Error string `json:"error"`
}
type SuccessResponse struct {
	Message string `json:"message"`
}

// CreateTender godoc
// @Summary      Create a new tender
// @Description  Create a new tender and store it in the database
// @Tags         tenders
// @Accept       json
// @Produce      json
// @Param        tender body     models.CreateTender true "Tender object"
// @Success      201 {object} models.Tender
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /tenders [post]
func (h *TenderHandler) CreateTender(c *gin.Context) {
	var createTender models.CreateTender

	if err := c.ShouldBindJSON(&createTender); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	tender := models.Tender{
		ClientId:      middleware.GetUserId(c, h.cfg),
		Title:         createTender.Title,
		Description:   createTender.Description,
		Deadline:      createTender.Deadline,
		AttachmentUrl: createTender.AttachmentUrl,
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

// UpdateTenderStatus godoc
// @Summary      Update the status of a tender
// @Description  Update the status of an existing tender by its ID
// @Tags         tenders
// @Accept       json
// @Produce      json
// @Param        id      path      string             true "Tender ID"
// @Param        status  body      models.Status      true "Tender Status"
// @Success      200     {object}  SuccessResponse
// @Failure      400     {object}  ErrorResponse
// @Failure      404     {object}  ErrorResponse
// @Failure      500     {object}  ErrorResponse
// @Router       /tenders/{id}/status [put]
func (h *TenderHandler) UpdateTenderStatus(c *gin.Context) {
	// Extract tender ID from URL parameters
	id := c.Param("id")

	var status models.Status
	if err := c.ShouldBindJSON(&status); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request"})
		return
	}

	// Call the service to update the tender status
	if err := h.ser.UpdateTenderStatus(c.Request.Context(), id, status); err != nil {
		h.logger.Error("failed to update tender status", "error", err)
		if err.Error() == "not found" {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Tender not found"})
		} else {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to update tender status"})
		}
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, SuccessResponse{Message: "Tender status updated successfully"})
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
