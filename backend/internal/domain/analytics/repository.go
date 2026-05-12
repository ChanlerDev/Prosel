package analytics

import "context"

type Repository interface {
	Record(ctx context.Context, event *AnalyticsEvent) error
	Overview(ctx context.Context, rangeValue DateRange) (*AnalyticsOverview, error)
	TopPages(ctx context.Context, rangeValue DateRange, limit int) ([]TopPage, error)
	DailyViews(ctx context.Context, days int) ([]DailyView, error)
	ContentStats(ctx context.Context, refType string, refID string) (*ContentAnalytics, error)
}
