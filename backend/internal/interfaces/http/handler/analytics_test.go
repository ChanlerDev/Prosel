package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/analytics"
	usecase "github.com/chanler/prosel/backend/internal/usecase/analytics"
	"github.com/gin-gonic/gin"
)

type fakeAnalyticsService struct {
	pageView usecase.PageViewRequest
	meta     usecase.ClientMeta
}

func (s *fakeAnalyticsService) RecordPageView(ctx context.Context, req usecase.PageViewRequest, meta usecase.ClientMeta) error {
	s.pageView = req
	s.meta = meta
	return nil
}
func (s *fakeAnalyticsService) GetOverview(ctx context.Context, rangeValue domain.DateRange) (*domain.AnalyticsOverview, error) {
	return &domain.AnalyticsOverview{TodayViews: 1, WeekViews: 7, MonthViews: 30, TopPages: []domain.TopPage{{Path: "/posts/hello", Views: 5}}, TopReferers: []domain.TopReferer{{Referer: "https://example.com", Views: 2}}, Devices: []domain.DeviceStat{{DeviceType: "desktop", Views: 6}}}, nil
}
func (s *fakeAnalyticsService) GetTopPages(ctx context.Context, rangeValue domain.DateRange, limit int) ([]domain.TopPage, error) {
	return []domain.TopPage{{Path: "/", Views: int64(limit)}}, nil
}
func (s *fakeAnalyticsService) GetDailyViews(ctx context.Context, days int) ([]domain.DailyView, error) {
	return []domain.DailyView{{Date: time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC), Views: int64(days)}}, nil
}

func TestAnalyticsHandlerRecordsPageView(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAnalyticsService{}
	router := gin.New()
	NewAnalyticsHandler(service).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/analytics/page-view", strings.NewReader(`{"path":"/posts/hello","refType":"post","refId":"post-1","referer":"https://example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "Go test")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if service.pageView.Path != "/posts/hello" || service.pageView.RefType != "post" || service.pageView.RefID != "post-1" || service.meta.UserAgent != "Go test" {
		t.Fatalf("recorded = %#v meta=%#v", service.pageView, service.meta)
	}
}

func TestAnalyticsHandlerOverviewReturnsStats(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewAnalyticsHandler(&fakeAnalyticsService{}).RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/admin/analytics/overview?range=7d", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"todayViews":1`) || !strings.Contains(recorder.Body.String(), `"deviceType":"desktop"`) {
		t.Fatalf("body = %s", recorder.Body.String())
	}
}

func TestAnalyticsHandlerDailyUsesDaysQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewAnalyticsHandler(&fakeAnalyticsService{}).RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/admin/analytics/daily?days=14", nil))

	if recorder.Code != http.StatusOK || !strings.Contains(recorder.Body.String(), `"views":14`) {
		t.Fatalf("status = %d, body = %s", recorder.Code, recorder.Body.String())
	}
}
