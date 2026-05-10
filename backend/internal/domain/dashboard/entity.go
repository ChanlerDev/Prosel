package dashboard

import "time"

type DashboardStats struct {
	TotalPosts      int64
	PublishedPosts  int64
	DraftPosts      int64
	PendingComments int64
	TodayViews      int64
	TotalViews      int64
	Categories      int64
	Tags            int64
	Topics          int64
}

type PostSummary struct {
	ID          string
	Title       string
	Slug        string
	Status      string
	ViewCount   int64
	PublishedAt *time.Time
	UpdatedAt   time.Time
}

type ActivityLog struct {
	ID         string
	ActorID    *string
	Action     string
	EntityType string
	EntityID   string
	Message    string
	CreatedAt  time.Time
}

type DashboardOverview struct {
	Stats       DashboardStats
	RecentPosts []PostSummary
	Activities  []ActivityLog
}
