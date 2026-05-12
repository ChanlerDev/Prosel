package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	infraAI "github.com/chanler/prosel/backend/internal/infrastructure/ai"
	"github.com/chanler/prosel/backend/internal/infrastructure/storage"
	"github.com/chanler/prosel/backend/internal/infrastructure/mailer"
	"github.com/chanler/prosel/backend/internal/infrastructure/cache"
	"github.com/chanler/prosel/backend/internal/infrastructure/config"
	"github.com/chanler/prosel/backend/internal/infrastructure/database"
	infraLogger "github.com/chanler/prosel/backend/internal/infrastructure/logger"
	"github.com/chanler/prosel/backend/internal/infrastructure/password"
	"github.com/chanler/prosel/backend/internal/infrastructure/token"
	"github.com/chanler/prosel/backend/internal/interfaces/http/handler"
	httpRouter "github.com/chanler/prosel/backend/internal/interfaces/http/router"
	aiUsecase "github.com/chanler/prosel/backend/internal/usecase/ai"
	fileUsecase "github.com/chanler/prosel/backend/internal/usecase/file"
	subscribeUsecase "github.com/chanler/prosel/backend/internal/usecase/subscribe"
	authUsecase "github.com/chanler/prosel/backend/internal/usecase/auth"
	commentUsecase "github.com/chanler/prosel/backend/internal/usecase/comment"
	dashboardUsecase "github.com/chanler/prosel/backend/internal/usecase/dashboard"
	noteUsecase "github.com/chanler/prosel/backend/internal/usecase/note"
	pageUsecase "github.com/chanler/prosel/backend/internal/usecase/page"
	postUsecase "github.com/chanler/prosel/backend/internal/usecase/post"
	searchUsecase "github.com/chanler/prosel/backend/internal/usecase/search"
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

	searchRepo := database.NewSearchRepository(db)
	searchUC := searchUsecase.NewSearchUsecase(searchRepo)
	searchHandler := handler.NewSearchHandler(searchUC)

	postRepo := database.NewPostRepository(db)
	postUC := postUsecase.NewPostUsecase(postRepo, searchUC)
	postHandler := handler.NewPostHandler(postUC)

	categoryRepo := database.NewCategoryRepository(db)
	tagRepo := database.NewTagRepository(db)
	topicRepo := database.NewTopicRepository(db)
	taxonomyUC := taxonomyUsecase.NewTaxonomyUsecase(categoryRepo, tagRepo, topicRepo)
	taxonomyHandler := handler.NewTaxonomyHandler(taxonomyUC, postUC)

	dashboardRepo := database.NewDashboardRepository(db)
	dashboardUC := dashboardUsecase.NewDashboardUsecase(dashboardRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardUC)

	commentRepo := database.NewCommentRepository(db)
	commentUC := commentUsecase.NewCommentUsecase(commentRepo)
	commentHandler := handler.NewCommentHandler(commentUC)

	noteRepo := database.NewNoteRepository(db)
	noteUC := noteUsecase.NewNoteUsecase(noteRepo, searchUC)
	noteHandler := handler.NewNoteHandler(noteUC)

	pageRepo := database.NewPageRepository(db)
	friendRepo := database.NewFriendRepository(db)
	pageUC := pageUsecase.NewPageUsecase(pageRepo, friendRepo, searchUC)
	pageHandler := handler.NewPageHandler(pageUC)

	fileRepo := database.NewFileRepository(db)
	localStorage := storage.NewLocalStorage(cfg.File.UploadDir, cfg.File.UploadPublicURL)
	fileUC := fileUsecase.NewFileUsecase(fileRepo, localStorage, fileUsecase.Options{MaxUploadBytes: int64(cfg.File.MaxUploadMB) << 20})
	fileHandler := handler.NewFileHandler(fileUC)

	subscriberRepo := database.NewSubscriberRepository(db)
	mailService := mailer.NewSMTPMailer(cfg.Mail, log)
	subscribeUC := subscribeUsecase.NewSubscribeUsecase(subscriberRepo, mailService, postUC, subscribeUsecase.Options{SiteURL: cfg.Site.URL})
	subscribeHandler := handler.NewSubscribeHandler(subscribeUC, postUC, cfg.Site.URL)

	aiRepo := database.NewAIRepository(db)
	aiClient := infraAI.NewOpenAIClient(cfg.AI)
	aiUC := aiUsecase.NewAIUsecase(aiRepo, aiClient, postUC)
	aiHandler := handler.NewAIHandler(aiUC)

	router := httpRouter.New(cfg, systemHandler, authHandler, postHandler, taxonomyHandler, dashboardHandler, commentHandler, noteHandler, pageHandler, searchHandler, fileHandler, subscribeHandler, aiHandler, tokenService)
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
