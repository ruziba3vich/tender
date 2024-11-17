// internal/handler/bid.go
package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/config"
	"github.com/zohirovs/internal/middleware"
	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/service"
)

type BidHandler struct {
	ser    *service.BidService
	logger *slog.Logger
	cfg    *config.Config
}

func NewBidHandler(logger *slog.Logger, ser *service.BidService, cfg *config.Config) *BidHandler {
	return &BidHandler{
		ser:    ser,
		logger: logger,
		cfg:    cfg,
	}
}

// SubmitBid godoc
// @Summary      Submit a new bid for a tender
// @Description  Allows contractors to submit a bid for an open tender
// @Tags         bids
// @Accept       json
// @Produce      json
// @Param        bid body     models.CreateBid true "Bid object"
// @Success      201  {object}  models.Bid
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/contractors/bids [post]
func (h *BidHandler) SubmitBid(c *gin.Context) {
	var createBid models.CreateBid
	if err := c.ShouldBindJSON(&createBid); err != nil {
		h.logger.Error("failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate price and delivery time
	if createBid.Price <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Price must be greater than 0"})
		return
	}
	if createBid.DeliveryTime <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Delivery time must be greater than 0"})
		return
	}

	bid := models.Bid{
		TenderId:     createBid.TenderId,
		ContractorId: middleware.GetUserId(c, h.cfg),
		Price:        createBid.Price,
		DeliveryTime: createBid.DeliveryTime,
		Comments:     createBid.Comments,
	}

	createdBid, err := h.ser.CreateBid(c.Request.Context(), &bid)
	if err != nil {
		h.logger.Error("failed to create bid", "error", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create bid"})
		return
	}

	c.JSON(http.StatusCreated, createdBid)
}

// ListBidsForTender godoc
// @Summary      List all bids for a tender
// @Description  Retrieve all bids submitted for a specific tender with optional filtering
// @Tags         bids
// @Accept       json
// @Produce      json
// @Param        tender_id    path      string  true  "Tender ID"
// @Param        min_price    query     number  false "Minimum price filter"
// @Param        max_price    query     number  false "Maximum price filter"
// @Param        max_delivery query     number  false "Maximum delivery time filter"
// @Success      200  {array}   models.Bid
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/clients/tenders/{tender_id}/bids [get]
func (h *BidHandler) ListBidsForTender(c *gin.Context) {
	tenderId := c.Param("tender_id")
	if tenderId == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Tender ID is required"})
		return
	}

	// Parse filter parameters
	filter := make(map[string]interface{})

	// Handle price range filter
	if minPrice := c.Query("min_price"); minPrice != "" {
		price, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			h.logger.Error("invalid min_price parameter", "error", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid min_price parameter"})
			return
		}
		if price > 0 {
			filter["min_price"] = price
		}
	}

	if maxPrice := c.Query("max_price"); maxPrice != "" {
		price, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			h.logger.Error("invalid max_price parameter", "error", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid max_price parameter"})
			return
		}
		if price > 0 {
			filter["max_price"] = price
		}
	}

	// Handle delivery time filter
	if maxDelivery := c.Query("max_delivery"); maxDelivery != "" {
		days, err := strconv.ParseInt(maxDelivery, 10, 64)
		if err != nil {
			h.logger.Error("invalid max_delivery parameter", "error", err)
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid max_delivery parameter"})
			return
		}
		if days > 0 {
			filter["max_delivery"] = days
		}
	}

	// Call service with separate tenderId parameter
	bids, err := h.ser.ListBidsForTender(c.Request.Context(), tenderId, filter)
	if err != nil {
		h.logger.Error("failed to list bids", "error", err, "tender_id", tenderId)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to retrieve bids"})
		return
	}

	c.JSON(http.StatusOK, bids)
}
