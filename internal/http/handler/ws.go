package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	ws "github.com/zohirovs/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) HandleWebSocket(c *gin.Context) {
	tenderID := c.Query("tender_id")
	if tenderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "tender_id is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := &ws.Client{
		ID:       uuid.New().String(),
		Conn:     conn,
		TenderID: tenderID,
	}

	h.WsManager.RegisterClient(client)

	defer func() {
		h.WsManager.UnregisterClient(client)
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
