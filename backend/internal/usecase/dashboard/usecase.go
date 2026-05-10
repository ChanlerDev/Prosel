package dashboard

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/dashboard"
)

type DashboardUsecase struct {
	repository domain.Repository
}

func NewDashboardUsecase(repository domain.Repository) *DashboardUsecase {
	return &DashboardUsecase{repository: repository}
}

func (uc *DashboardUsecase) GetOverview(ctx context.Context) (*domain.DashboardOverview, error) {
	stats, err := uc.repository.GetStats(ctx)
	if err != nil {
		return nil, err
	}
	posts, err := uc.repository.GetRecentPosts(ctx, 5)
	if err != nil {
		return nil, err
	}
	activities, err := uc.repository.GetRecentActivities(ctx, 10)
	if err != nil {
		return nil, err
	}
	return &domain.DashboardOverview{Stats: *stats, RecentPosts: posts, Activities: activities}, nil
}

func (uc *DashboardUsecase) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	return uc.repository.GetStats(ctx)
}

func (uc *DashboardUsecase) GetActivities(ctx context.Context, limit int) ([]domain.ActivityLog, error) {
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return uc.repository.GetRecentActivities(ctx, limit)
}

func (uc *DashboardUsecase) RecordActivity(ctx context.Context, log domain.ActivityLog) error {
	if log.ID == "" {
		log.ID = newID()
	}
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now().UTC()
	}
	return uc.repository.RecordActivity(ctx, log)
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
