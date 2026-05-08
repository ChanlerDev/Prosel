package system

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/system"
)

type fakeSettingRepository struct {
	settings []domain.SiteSetting
	err      error
}

func (r fakeSettingRepository) GetByKey(ctx context.Context, key string) (*domain.SiteSetting, error) {
	for i := range r.settings {
		if r.settings[i].Key == key {
			return &r.settings[i], nil
		}
	}
	return nil, domain.ErrSettingNotFound
}

func (r fakeSettingRepository) List(ctx context.Context) ([]domain.SiteSetting, error) {
	return r.settings, r.err
}

func (r fakeSettingRepository) Upsert(ctx context.Context, setting *domain.SiteSetting) error {
	return r.err
}

type fakeHealthChecker struct {
	status *domain.HealthStatus
	err    error
}

func (c fakeHealthChecker) Check(ctx context.Context) (*domain.HealthStatus, error) {
	return c.status, c.err
}

func TestGetPublicSettingsReturnsOnlyPublicKeys(t *testing.T) {
	uc := NewSystemUsecase(fakeSettingRepository{settings: []domain.SiteSetting{
		{Key: "site_name", Value: "Prosel", ValueType: domain.ValueTypeString},
		{Key: "site_description", Value: "Blog", ValueType: domain.ValueTypeString},
		{Key: "private_token", Value: "secret", ValueType: domain.ValueTypeString},
	}}, fakeHealthChecker{})

	settings, err := uc.GetPublicSettings(context.Background())
	if err != nil {
		t.Fatalf("GetPublicSettings() error = %v", err)
	}

	if len(settings) != 2 {
		t.Fatalf("len(settings) = %d, want 2", len(settings))
	}
	for _, setting := range settings {
		if setting.Key == "private_token" {
			t.Fatalf("private setting leaked: %#v", setting)
		}
	}
}

func TestCheckHealthNormalizesStatus(t *testing.T) {
	uc := NewSystemUsecase(fakeSettingRepository{}, fakeHealthChecker{status: &domain.HealthStatus{
		Status:     domain.StatusHealthy,
		DatabaseOK: true,
		RedisOK:    false,
		CheckedAt:  time.Now(),
	}})

	status, err := uc.CheckHealth(context.Background())
	if err != nil {
		t.Fatalf("CheckHealth() error = %v", err)
	}
	if status.Status != domain.StatusUnhealthy {
		t.Fatalf("Status = %q, want %q", status.Status, domain.StatusUnhealthy)
	}
}

func TestCheckHealthReturnsErrorWhenCheckerFails(t *testing.T) {
	wantErr := errors.New("database down")
	uc := NewSystemUsecase(fakeSettingRepository{}, fakeHealthChecker{err: wantErr})

	_, err := uc.CheckHealth(context.Background())
	if !errors.Is(err, wantErr) {
		t.Fatalf("CheckHealth() error = %v, want %v", err, wantErr)
	}
}
