// Package api API.
//
// @title # MiniTwitter
// @version 1.03.67.83.145
//
// @description API Endpoints for MiniTwitter
// @termsOfService http://swagger.io/terms/
//
// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
//
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host localhost:8080
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package app

import (
	"log/slog"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/zohirovs/internal/config"
	_ "github.com/zohirovs/internal/http/app/docs"
	"github.com/zohirovs/internal/http/handler"
)

func Run(handler *handler.Handler, logger *slog.Logger, config *config.Config, enforcer *casbin.Enforcer) error {
	router := gin.Default()

	// CORS configuration
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	// Swagger documentation setup
	url := ginSwagger.URL("/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url, ginSwagger.PersistAuthorization(true)))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API endpoints
	// User endpoints
	users := router.Group("/")
	{
		users.POST("/register", handler.UserHandler.RegisterUser)
		users.POST("/login", handler.UserHandler.LoginUser)
	}

	// Client endpoints group
	clients := router.Group("api/clients")
	{
		// Tender endpoints
		tenders := clients.Group("/tenders")
		{
			tenders.POST("", handler.TenderHandler.CreateTender)
			tenders.GET("/:id", handler.TenderHandler.GetTender)
			tenders.PUT("/:id/status", handler.TenderHandler.UpdateTenderStatus)
			tenders.DELETE("/:id", handler.TenderHandler.DeleteTender)
		}

		// Bid endpoints for clients (viewing bids)
		bids := clients.Group("/bids")
		{
			bids.GET("/tender/:id", handler.BidHandler.ListBidsForTender) // Changed from /:tender_id/bids to /bids/tender/:id
		}
	}

	// Contractor endpoints group
	contractors := router.Group("api/contractors")
	{
		// Bid endpoints for contractors (submitting bids)
		bids := contractors.Group("/bids")
		{
			bids.POST("", handler.BidHandler.SubmitBid)
		}
	}

	// Start the server
	return router.Run(config.Server.Port)
}
