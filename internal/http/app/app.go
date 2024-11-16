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

	// _ "github.com/abdulazizax/mini-twitter/api-service/internal/items/http/app/docs"
	// "github.com/abdulazizax/mini-twitter/api-service/internal/items/middleware"

	// casbin "github.com/casbin/casbin/v2"
	// "github.com/gin-contrib/cors"

	// "github.com/abdulazizax/mini-twitter/api-service/internal/items/http/handler"
	// "github.com/abdulazizax/mini-twitter/api-service/internal/pkg/config"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/zohirovs/internal/config"
	_ "github.com/zohirovs/internal/http/app/docs"
	"github.com/zohirovs/internal/http/handler"
	"github.com/zohirovs/internal/middleware"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Run initializes and starts the HTTP server for the MiniTwitter API.
// It sets up routing, middleware, and Swagger documentation.
//
// Parameters:
// - handler: Pointer to the Handler struct containing all route handlers
// - logger: Structured logger for logging
// - config: Application configuration
// - enforcer: Casbin enforcer for authorization
//
// Returns:
// - error: Any error that occurs during server startup
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

	// API ednpoints
	tenders := router.Group("/tenders")
	tenders.Use(middleware.AuthzMiddleware("/tenders", enforcer, config))
	{
		tenders.POST("", handler.TenderHandler.CreateTender)
		tenders.GET("", handler.TenderHandler.GetTender)
		tenders.PUT(":id/status", handler.TenderHandler.UpdateTenderStatus)
		tenders.DELETE("", handler.TenderHandler.DeleteTender)
	}
	// Start the server
	return router.Run(config.Server.Port)
}
