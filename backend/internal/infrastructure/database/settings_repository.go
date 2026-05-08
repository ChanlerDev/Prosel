package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/system"
)

type SiteSettingModel struct {
	ID           string `gorm:"primaryKey;size:36"`
	SettingKey   string `gorm:"column:setting_key;uniqueIndex;size:100;not null"`
	SettingValue string `gorm:"column:setting_value"`
	ValueType    string `gorm:"size:20;not null"`
	Description  string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (SiteSettingModel) TableName() string {
	return "site_settings"
}

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) GetByKey(ctx context.Context, key string) (*domain.SiteSetting, error) {
	var model SiteSettingModel
	err := r.db.WithContext(ctx).Where("setting_key = ?", key).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSettingNotFound
	}
	if err != nil {
		return nil, err
	}
	setting := toDomainSetting(model)
	return &setting, nil
}

func (r *SettingRepository) List(ctx context.Context) ([]domain.SiteSetting, error) {
	var models []SiteSettingModel
	if err := r.db.WithContext(ctx).Order("setting_key asc").Find(&models).Error; err != nil {
		return nil, err
	}

	settings := make([]domain.SiteSetting, len(models))
	for i, model := range models {
		settings[i] = toDomainSetting(model)
	}
	return settings, nil
}

func (r *SettingRepository) Upsert(ctx context.Context, setting *domain.SiteSetting) error {
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO site_settings (id, setting_key, setting_value, value_type, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
		ON CONFLICT (setting_key)
		DO UPDATE SET setting_value = EXCLUDED.setting_value, value_type = EXCLUDED.value_type, description = EXCLUDED.description, updated_at = NOW()
	`, setting.ID, setting.Key, setting.Value, setting.ValueType, setting.Description).Error
}

func toDomainSetting(model SiteSettingModel) domain.SiteSetting {
	return domain.SiteSetting{
		ID:          model.ID,
		Key:         model.SettingKey,
		Value:       model.SettingValue,
		ValueType:   model.ValueType,
		Description: model.Description,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}
