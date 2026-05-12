package analytics

import (
	"context"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/analytics"
)

type fakeAnalyticsRepo struct {
	recorded *domain.AnalyticsEvent
	daily    []domain.DailyView
}

func (r *fakeAnalyticsRepo) Record(ctx context.Context, event *domain.AnalyticsEvent) error {
	r.recorded = event
	return nil
}
func (r *fakeAnalyticsRepo) Overview(ctx context.Context, rangeValue domain.DateRange) (*domain.AnalyticsOverview, error) {
	return &domain.AnalyticsOverview{TodayViews: 2, WeekViews: 7, MonthViews: 30}, nil
}
func (r *fakeAnalyticsRepo) TopPages(ctx context.Context, rangeValue domain.DateRange, limit int) ([]domain.TopPage, error) {
	return []domain.TopPage{{Path: "/posts/hello", Views: 5}}, nil
}
func (r *fakeAnalyticsRepo) DailyViews(ctx context.Context, days int) ([]domain.DailyView, error) {
	return r.daily, nil
}
func (r *fakeAnalyticsRepo) ContentStats(ctx context.Context, refType string, refID string) (*domain.ContentAnalytics, error) {
	return &domain.ContentAnalytics{RefType: refType, RefID: refID, Views: 3}, nil
}

func TestRecordPageViewNormalizesRequestAndHashesIP(t *testing.T) {
	repo := &fakeAnalyticsRepo{}
	uc := NewAnalyticsUsecase(repo)

	err := uc.RecordPageView(context.Background(), PageViewRequest{Path: " /posts/hello ", RefType: "post", RefID: "post-1", Referer: "https://example.com"}, ClientMeta{IP: "127.0.0.1", UserAgent: "Mozilla/5.0 (iPhone) Safari/605.1.15"})
	if err != nil {
		t.Fatalf("RecordPageView() error = %v", err)
	}
	if repo.recorded == nil {
		t.Fatal("event was not recorded")
	}
	if repo.recorded.Path != "/posts/hello" || repo.recorded.EventType != "page_view" || repo.recorded.RefType != "post" || repo.recorded.RefID != "post-1" {
		t.Fatalf("event = %#v", repo.recorded)
	}
	if repo.recorded.IPHash == "" || repo.recorded.IPHash == "127.0.0.1" {
		t.Fatalf("IPHash = %q", repo.recorded.IPHash)
	}
	if repo.recorded.DeviceType != "mobile" || repo.recorded.Browser != "Safari" || repo.recorded.OS != "iOS" {
		t.Fatalf("client fields = %#v", repo.recorded)
	}
}

func TestRecordPageViewRejectsInvalidInput(t *testing.T) {
	uc := NewAnalyticsUsecase(&fakeAnalyticsRepo{})
	if err := uc.RecordPageView(context.Background(), PageViewRequest{Path: ""}, ClientMeta{}); err != domain.ErrInvalidAnalyticsEvent {
		t.Fatalf("empty path error = %v", err)
	}
	if err := uc.RecordPageView(context.Background(), PageViewRequest{Path: "/x", RefType: "user", RefID: "1"}, ClientMeta{}); err != domain.ErrInvalidAnalyticsEvent {
		t.Fatalf("invalid ref type error = %v", err)
	}
}

func TestGetDailyViewsNormalizesDays(t *testing.T) {
	repo := &fakeAnalyticsRepo{daily: []domain.DailyView{{Date: time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC), Views: 4}}}
	uc := NewAnalyticsUsecase(repo)

	daily, err := uc.GetDailyViews(context.Background(), 0)
	if err != nil {
		t.Fatalf("GetDailyViews() error = %v", err)
	}
	if len(daily) != 1 || daily[0].Views != 4 {
		t.Fatalf("daily = %#v", daily)
	}
}
