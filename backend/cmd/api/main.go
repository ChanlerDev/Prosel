package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chanler/prosel/backend/internal/infrastructure/cache"
	"github.com/chanler/prosel/backend/internal/infrastructure/config"
	"github.com/chanler/prosel/backend/internal/infrastructure/database"
	infraLogger "github.com/chanler/prosel/backend/internal/infrastructure/logger"
	"github.com/chanler/prosel/backend/internal/infrastructure/password"
	"github.com/chanler/prosel/backend/internal/infrastructure/token"
	"github.com/chanler/prosel/backend/internal/interfaces/http/handler"
	httpRouter "github.com/chanler/prosel/backend/internal/interfaces/http/router"
	authUsecase "github.com/chanler/prosel/backend/internal/usecase/auth"
	dashboardUsecase "github.com/chanler/prosel/backend/internal/usecase/dashboard"
	postUsecase "github.com/chanler/prosel/backend/internal/usecase/post"
	systemUsecase "github.com/chanler/prosel/backend/internal/usecase/system"
	taxonomyUsecase "github.com/chanler/prosel/backend/internal/usecase/taxonomy"
)

func main() {
	cfg := config.Load()
	log := infraLogger.New()

	ctx, cancel := context.WithTimeout(context.Background(), config.StartupTimeout())
	defer cancel()

	db, err := database.Open(cfg.Database)
	if err != nil {
		log.Error("open database", slog.Any("error", err))
		os.Exit(1)
	}
	if err := database.RunMigrations(ctx, db, cfg.Database.MigrationsDir); err != nil {
		log.Error("run migrations", slog.Any("error", err))
		os.Exit(1)
	}

	redisClient := cache.NewRedisClient(cfg.Redis)
	settingsRepo := database.NewSettingRepository(db)
	healthChecker := database.NewHealthChecker(db, redisClient, cfg.App.Version)
	systemUC := systemUsecase.NewSystemUsecase(settingsRepo, healthChecker)
	systemHandler := handler.NewSystemHandler(systemUC)

	userRepo := database.NewUserRepository(db)
	sessionRepo := database.NewSessionRepository(db)
	passwordHasher := password.NewBcryptHasher(cfg.Auth.PasswordBcryptCost)
	accessDuration := time.Duration(cfg.Auth.AccessTokenMinutes) * time.Minute
	refreshDuration := time.Duration(cfg.Auth.RefreshTokenHours) * time.Hour
	tokenService := token.NewJWTService(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer, accessDuration)
	authUC := authUsecase.NewAuthUsecase(userRepo, sessionRepo, passwordHasher, tokenService, accessDuration, refreshDuration)
	authHandler := handler.NewAuthHandler(authUC)

	postRepo := database.NewPostRepository(db)
	postUC := postUsecase.NewPostUsecase(postRepo)
	postHandler := handler.NewPostHandler(postUC)

	categoryRepo := database.NewCategoryRepository(db)
	tagRepo := database.NewTagRepository(db)
	topicRepo := database.NewTopicRepository(db)
	taxonomyUC := taxonomyUsecase.NewTaxonomyUsecase(categoryRepo, tagRepo, topicRepo)
	taxonomyHandler := handler.NewTaxonomyHandler(taxonomyUC, postUC)

	dashboardRepo := database.NewDashboardRepository(db)
	dashboardUC := dashboardUsecase.NewDashboardUsecase(dashboardRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardUC)

	router := httpRouter.New(cfg, systemHandler, authHandler, postHandler, taxonomyHandler, dashboardHandler, tokenService)
	server := &http.Server{Addr: cfg.HTTP.Address(), Handler: router}

	go func() {
		log.Info("starting api", slog.String("addr", server.Addr))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	<-shutdown

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("shutdown server", slog.Any("error", err))
	}
}
