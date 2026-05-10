package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/dashboard"
)

type ActivityLogModel struct {
	ID         string  `gorm:"primaryKey;size:36"`
	ActorID    *string `gorm:"size:36"`
	Action     string  `gorm:"size:80;not null"`
	EntityType string  `gorm:"size:40"`
	EntityID   string  `gorm:"size:36"`
	Message    string
	CreatedAt  time.Time
}

func (ActivityLogModel) TableName() string { return "activity_logs" }

type DashboardRepository struct{ db *gorm.DB }

func NewDashboardRepository(db *gorm.DB) *DashboardRepository { return &DashboardRepository{db: db} }

func (r *DashboardRepository) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	var stats domain.DashboardStats
	if err := r.db.WithContext(ctx).Model(&PostModel{}).Count(&stats.TotalPosts).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&PostModel{}).Where("status = ?", "published").Count(&stats.PublishedPosts).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&PostModel{}).Where("status = ?", "draft").Count(&stats.DraftPosts).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&CategoryModel{}).Count(&stats.Categories).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&TagModel{}).Count(&stats.Tags).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&TopicModel{}).Count(&stats.Topics).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&PostModel{}).Select("COALESCE(SUM(view_count), 0)").Scan(&stats.TotalViews).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

func (r *DashboardRepository) GetRecentPosts(ctx context.Context, limit int) ([]domain.PostSummary, error) {
	var models []PostModel
	if err := r.db.WithContext(ctx).Order("updated_at DESC").Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}
	posts := make([]domain.PostSummary, 0, len(models))
	for _, model := range models {
		posts = append(posts, domain.PostSummary{ID: model.ID, Title: model.Title, Slug: model.Slug, Status: model.Status, ViewCount: model.ViewCount, PublishedAt: model.PublishedAt, UpdatedAt: model.UpdatedAt})
	}
	return posts, nil
}

func (r *DashboardRepository) GetRecentActivities(ctx context.Context, limit int) ([]domain.ActivityLog, error) {
	var models []ActivityLogModel
	if err := r.db.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&models).Error; err != nil {
		return nil, err
	}
	activities := make([]domain.ActivityLog, 0, len(models))
	for _, model := range models {
		activities = append(activities, domain.ActivityLog{ID: model.ID, ActorID: model.ActorID, Action: model.Action, EntityType: model.EntityType, EntityID: model.EntityID, Message: model.Message, CreatedAt: model.CreatedAt})
	}
	return activities, nil
}

func (r *DashboardRepository) RecordActivity(ctx context.Context, log domain.ActivityLog) error {
	return r.db.WithContext(ctx).Create(&ActivityLogModel{ID: log.ID, ActorID: log.ActorID, Action: log.Action, EntityType: log.EntityType, EntityID: log.EntityID, Message: log.Message, CreatedAt: log.CreatedAt}).Error
}
