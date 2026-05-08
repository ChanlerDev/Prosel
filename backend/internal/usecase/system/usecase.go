package system

import (
	"context"

	domain "github.com/chanler/prosel/backend/internal/domain/system"
)

var publicSettingKeys = map[string]struct{}{
	"site_name":          {},
	"site_description":   {},
	"site_url":           {},
	"posts_per_page":     {},
	"comment_moderation": {},
	"analytics_enabled":  {},
}

type SystemUsecase struct {
	settings domain.SettingRepository
	health   domain.HealthChecker
}

func NewSystemUsecase(settings domain.SettingRepository, health domain.HealthChecker) *SystemUsecase {
	return &SystemUsecase{settings: settings, health: health}
}

func (uc *SystemUsecase) GetPublicSettings(ctx context.Context) ([]domain.SiteSetting, error) {
	settings, err := uc.settings.List(ctx)
	if err != nil {
		return nil, err
	}

	publicSettings := make([]domain.SiteSetting, 0, len(settings))
	for _, setting := range settings {
		if _, ok := publicSettingKeys[setting.Key]; ok {
			publicSettings = append(publicSettings, setting)
		}
	}
	return publicSettings, nil
}

func (uc *SystemUsecase) CheckHealth(ctx context.Context) (*domain.HealthStatus, error) {
	status, err := uc.health.Check(ctx)
	if err != nil {
		return nil, err
	}
	status.Normalize()
	return status, nil
}
