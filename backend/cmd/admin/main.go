package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/auth"
	"github.com/chanler/prosel/backend/internal/infrastructure/config"
	"github.com/chanler/prosel/backend/internal/infrastructure/database"
	infraLogger "github.com/chanler/prosel/backend/internal/infrastructure/logger"
	"github.com/chanler/prosel/backend/internal/infrastructure/password"
)

func main() {
	username := flag.String("username", "admin", "admin username")
	email := flag.String("email", "admin@example.com", "admin email")
	displayName := flag.String("display-name", "Admin", "admin display name")
	plainPassword := flag.String("password", "", "admin password")
	flag.Parse()

	if strings.TrimSpace(*plainPassword) == "" {
		fmt.Fprintln(os.Stderr, "-password is required")
		os.Exit(2)
	}

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

	hasher := password.NewBcryptHasher(cfg.Auth.PasswordBcryptCost)
	hash, err := hasher.Hash(*plainPassword)
	if err != nil {
		log.Error("hash password", slog.Any("error", err))
		os.Exit(1)
	}

	repo := database.NewUserRepository(db)
	user := &domain.User{ID: newID(), Username: strings.TrimSpace(*username), Email: strings.TrimSpace(*email), PasswordHash: hash, DisplayName: strings.TrimSpace(*displayName), Role: domain.RoleAdmin, Status: domain.StatusActive}
	if err := repo.Create(ctx, user); err != nil {
		log.Error("create admin", slog.Any("error", err))
		os.Exit(1)
	}
	log.Info("admin created", slog.String("username", user.Username), slog.String("email", user.Email))
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))[:32]
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:])
}
