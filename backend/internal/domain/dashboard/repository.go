package dashboard

import "context"

type Repository interface {
	GetStats(ctx context.Context) (*DashboardStats, error)
	GetRecentPosts(ctx context.Context, limit int) ([]PostSummary, error)
	GetRecentActivities(ctx context.Context, limit int) ([]ActivityLog, error)
	RecordActivity(ctx context.Context, log ActivityLog) error
}
