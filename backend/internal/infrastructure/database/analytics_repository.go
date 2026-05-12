package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/analytics"
)

type AnalyticsEventModel struct {
	ID         string `gorm:"primaryKey;size:36"`
	EventType  string `gorm:"size:40;not null"`
	Path       string `gorm:"size:500;not null"`
	RefType    string `gorm:"size:20"`
	RefID      string `gorm:"size:36"`
	Referer    string `gorm:"size:500"`
	IPHash     string `gorm:"size:128"`
	UserAgent  string
	Country    string `gorm:"size:80"`
	DeviceType string `gorm:"size:40"`
	Browser    string `gorm:"size:80"`
	OS         string `gorm:"size:80"`
	CreatedAt  time.Time
}

func (AnalyticsEventModel) TableName() string { return "analytics_events" }

type AnalyticsRepository struct{ db *gorm.DB }

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository { return &AnalyticsRepository{db: db} }

func (r *AnalyticsRepository) Record(ctx context.Context, event *domain.AnalyticsEvent) error {
	return r.db.WithContext(ctx).Create(&AnalyticsEventModel{ID: event.ID, EventType: event.EventType, Path: event.Path, RefType: event.RefType, RefID: event.RefID, Referer: event.Referer, IPHash: event.IPHash, UserAgent: event.UserAgent, Country: event.Country, DeviceType: event.DeviceType, Browser: event.Browser, OS: event.OS, CreatedAt: event.CreatedAt}).Error
}

func (r *AnalyticsRepository) Overview(ctx context.Context, rangeValue domain.DateRange) (*domain.AnalyticsOverview, error) {
	now := time.Now().UTC()
	startToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	overview := &domain.AnalyticsOverview{}
	if err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).Where("created_at >= ?", startToday).Count(&overview.TodayViews).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).Where("created_at >= ?", startToday.AddDate(0, 0, -6)).Count(&overview.WeekViews).Error; err != nil {
		return nil, err
	}
	if err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).Where("created_at >= ?", startToday.AddDate(0, 0, -29)).Count(&overview.MonthViews).Error; err != nil {
		return nil, err
	}
	var err error
	overview.TopPages, err = r.TopPages(ctx, rangeValue, 10)
	if err != nil {
		return nil, err
	}
	overview.TopReferers, err = r.topReferers(ctx, rangeValue, 10)
	if err != nil {
		return nil, err
	}
	overview.Devices, err = r.devices(ctx, rangeValue)
	if err != nil {
		return nil, err
	}
	return overview, nil
}

func (r *AnalyticsRepository) TopPages(ctx context.Context, rangeValue domain.DateRange, limit int) ([]domain.TopPage, error) {
	var rows []struct {
		Path    string
		RefType string
		RefID   string
		Views   int64
	}
	err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).
		Select("path, COALESCE(ref_type, '') AS ref_type, COALESCE(ref_id, '') AS ref_id, COUNT(*) AS views").
		Where("created_at >= ?", rangeStart(rangeValue)).
		Group("path, ref_type, ref_id").
		Order("views DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	pages := make([]domain.TopPage, 0, len(rows))
	for _, row := range rows {
		pages = append(pages, domain.TopPage{Path: row.Path, RefType: row.RefType, RefID: row.RefID, Views: row.Views})
	}
	return pages, nil
}

func (r *AnalyticsRepository) DailyViews(ctx context.Context, days int) ([]domain.DailyView, error) {
	var rows []struct {
		Date  time.Time
		Views int64
	}
	err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).
		Select("DATE_TRUNC('day', created_at) AS date, COUNT(*) AS views").
		Where("created_at >= ?", time.Now().UTC().AddDate(0, 0, -days+1)).
		Group("date").
		Order("date ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	views := make([]domain.DailyView, 0, len(rows))
	for _, row := range rows {
		views = append(views, domain.DailyView{Date: row.Date, Views: row.Views})
	}
	return views, nil
}

func (r *AnalyticsRepository) ContentStats(ctx context.Context, refType string, refID string) (*domain.ContentAnalytics, error) {
	stats := &domain.ContentAnalytics{RefType: refType, RefID: refID}
	if err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).Where("ref_type = ? AND ref_id = ?", refType, refID).Count(&stats.Views).Error; err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *AnalyticsRepository) topReferers(ctx context.Context, rangeValue domain.DateRange, limit int) ([]domain.TopReferer, error) {
	var rows []domain.TopReferer
	err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).
		Select("referer, COUNT(*) AS views").
		Where("created_at >= ? AND referer <> ''", rangeStart(rangeValue)).
		Group("referer").
		Order("views DESC").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (r *AnalyticsRepository) devices(ctx context.Context, rangeValue domain.DateRange) ([]domain.DeviceStat, error) {
	var rows []domain.DeviceStat
	err := r.db.WithContext(ctx).Model(&AnalyticsEventModel{}).
		Select("COALESCE(NULLIF(device_type, ''), 'unknown') AS device_type, COUNT(*) AS views").
		Where("created_at >= ?", rangeStart(rangeValue)).
		Group("device_type").
		Order("views DESC").
		Scan(&rows).Error
	return rows, err
}

func rangeStart(rangeValue domain.DateRange) time.Time {
	now := time.Now().UTC()
	days := 30
	switch rangeValue {
	case domain.DateRange7d:
		days = 7
	case domain.DateRange90d:
		days = 90
	}
	return now.AddDate(0, 0, -days+1)
}
