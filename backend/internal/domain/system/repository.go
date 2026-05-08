package system

import (
	"context"
	"errors"
)

var ErrSettingNotFound = errors.New("setting not found")

type SettingRepository interface {
	GetByKey(ctx context.Context, key string) (*SiteSetting, error)
	List(ctx context.Context) ([]SiteSetting, error)
	Upsert(ctx context.Context, setting *SiteSetting) error
}

type HealthChecker interface {
	Check(ctx context.Context) (*HealthStatus, error)
}
