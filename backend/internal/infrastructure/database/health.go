package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/system"
)

type HealthChecker struct {
	db      *gorm.DB
	redis   *redis.Client
	version string
}

func NewHealthChecker(db *gorm.DB, redisClient *redis.Client, version string) *HealthChecker {
	return &HealthChecker{db: db, redis: redisClient, version: version}
}

func (c *HealthChecker) Check(ctx context.Context) (*domain.HealthStatus, error) {
	status := &domain.HealthStatus{Version: c.version, CheckedAt: time.Now().UTC()}

	sqlDB, err := c.db.DB()
	if err == nil {
		status.DatabaseOK = sqlDB.PingContext(ctx) == nil
	}
	status.RedisOK = c.redis.Ping(ctx).Err() == nil
	status.Normalize()

	return status, nil
}
