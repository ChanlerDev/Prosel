package router

import (
	"github.com/gin-gonic/gin"

	"github.com/chanler/prosel/backend/internal/infrastructure/config"
	"github.com/chanler/prosel/backend/internal/interfaces/http/handler"
	"github.com/chanler/prosel/backend/internal/interfaces/http/middleware"
)

func New(cfg config.Config, systemHandler *handler.SystemHandler, authHandler *handler.AuthHandler, postHandler *handler.PostHandler, taxonomyHandler *handler.TaxonomyHandler, tokenParser middleware.AccessTokenParser) *gin.Engine {
	if cfg.App.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(gin.Recovery(), gin.Logger(), middleware.CORS(cfg.Cors.AllowedOrigins))

	api := r.Group("/api/v1")
	systemHandler.RegisterPublicRoutes(api)
	authHandler.RegisterPublicRoutes(api)
	postHandler.RegisterPublicRoutes(api)
	taxonomyHandler.RegisterPublicRoutes(api)

	protected := api.Group("")
	protected.Use(middleware.Auth(tokenParser))
	admin := protected.Group("/admin")
	authHandler.RegisterProtectedRoutes(protected, admin)
	postHandler.RegisterProtectedRoutes(admin)
	taxonomyHandler.RegisterProtectedRoutes(admin)

	return r
}
