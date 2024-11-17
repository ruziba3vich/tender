/*
 * @Author: javohir-a abdusamatovjavohir@gmail.com
 * @Date: 2024-11-17 00:47:41
 * @LastEditors: javohir-a abdusamatovjavohir@gmail.com
 * @LastEditTime: 2024-11-17 04:24:18
 * @FilePath: /tender/internal/http/handler/handler.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package handler

import (
	"log/slog"

	"github.com/zohirovs/internal/config"
	"github.com/zohirovs/internal/service"
	websocket "github.com/zohirovs/internal/ws"
)

type Handler struct {
	UserHandler         *UserHandler
	BidHandler          *BidHandler
	NotificationHandler *NotificationHandler
	TenderHandler       *TenderHandler
	WsManager           *websocket.Manager
}

func NewHandler(logger *slog.Logger, service *service.Service, cfg *config.Config) *Handler {
	wsManager := websocket.NewManager()
	return &Handler{
		UserHandler:         NewUserHandler(logger, service.User),
		BidHandler:          NewBidHandler(logger, service.Bid, wsManager),
		NotificationHandler: NewNotificationHandler(logger, service.Notification),
		TenderHandler:       NewTenderHandler(logger, service.Tender, cfg),
		WsManager:           wsManager,
	}
}
