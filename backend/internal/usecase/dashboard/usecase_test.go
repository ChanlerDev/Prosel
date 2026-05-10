package dashboard

import (
	"context"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/dashboard"
)

type fakeDashboardRepo struct {
	stats      *domain.DashboardStats
	posts      []domain.PostSummary
	activities []domain.ActivityLog
}

func (r fakeDashboardRepo) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	return r.stats, nil
}
func (r fakeDashboardRepo) GetRecentPosts(ctx context.Context, limit int) ([]domain.PostSummary, error) {
	return r.posts[:min(limit, len(r.posts))], nil
}
func (r fakeDashboardRepo) GetRecentActivities(ctx context.Context, limit int) ([]domain.ActivityLog, error) {
	return r.activities[:min(limit, len(r.activities))], nil
}
func (r fakeDashboardRepo) RecordActivity(ctx context.Context, log domain.ActivityLog) error {
	return nil
}

func TestGetOverviewCombinesStatsPostsAndActivities(t *testing.T) {
	now := time.Now().UTC()
	uc := NewDashboardUsecase(fakeDashboardRepo{
		stats:      &domain.DashboardStats{TotalPosts: 3, PublishedPosts: 2, DraftPosts: 1, TotalViews: 20},
		posts:      []domain.PostSummary{{ID: "post-1", Title: "Hello", Status: "published", UpdatedAt: now}},
		activities: []domain.ActivityLog{{ID: "activity-1", Action: "post.published", Message: "Published Hello", CreatedAt: now}},
	})

	overview, err := uc.GetOverview(context.Background())
	if err != nil {
		t.Fatalf("GetOverview() error = %v", err)
	}
	if overview.Stats.TotalPosts != 3 || len(overview.RecentPosts) != 1 || len(overview.Activities) != 1 {
		t.Fatalf("overview = %#v", overview)
	}
}
