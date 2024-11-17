package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/models"
	"github.com/zohirovs/internal/service"
	websocket "github.com/zohirovs/internal/ws"
)

type BidHandler struct {
	logger     *slog.Logger
	bidService *service.BidService
	wsManager  *websocket.Manager
}

func NewBidHandler(logger *slog.Logger, bid *service.BidService, wsManager *websocket.Manager) *BidHandler {
	return &BidHandler{
		logger:     logger,
		bidService: bid,
		wsManager:  wsManager,
	}
}

func (h *BidHandler) CreateBid(c *gin.Context) {
	var req models.CreateBidRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("failed to bind request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bid, err := h.bidService.CreateBid(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to create bid", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	notification := map[string]interface{}{
		"type": "new_bid",
		"data": bid,
	}

	if err := h.wsManager.BroadcastToTender(bid.TenderID, notification); err != nil {
		h.logger.Error("failed to broadcast bid", "error", err)
	}

	c.JSON(http.StatusCreated, bid)
}

func (h *BidHandler) GetBids(c *gin.Context) {
	tenderID := c.Query("tender_id")
	if tenderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tender_id is required"})
		return
	}

	// Get bids using service
	bids, err := h.bidService.GetBids(c.Request.Context(), tenderID)
	if err != nil {
		h.logger.Error("failed to get bids", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}
