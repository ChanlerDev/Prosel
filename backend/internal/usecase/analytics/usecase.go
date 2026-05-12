package analytics

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/analytics"
)

type AnalyticsUsecase struct {
	repository domain.Repository
}

type PageViewRequest struct {
	Path    string
	RefType string
	RefID   string
	Referer string
}

type ClientMeta struct {
	IP        string
	UserAgent string
}

func NewAnalyticsUsecase(repository domain.Repository) *AnalyticsUsecase {
	return &AnalyticsUsecase{repository: repository}
}

func (uc *AnalyticsUsecase) RecordPageView(ctx context.Context, req PageViewRequest, meta ClientMeta) error {
	path := strings.TrimSpace(req.Path)
	refType := strings.TrimSpace(req.RefType)
	if path == "" || !domain.ValidRefType(refType) {
		return domain.ErrInvalidAnalyticsEvent
	}
	now := time.Now().UTC()
	event := &domain.AnalyticsEvent{
		ID:        newID(),
		EventType: domain.EventTypePageView,
		Path:      truncate(path, 500),
		RefType:   refType,
		RefID:     strings.TrimSpace(req.RefID),
		Referer:   truncate(strings.TrimSpace(req.Referer), 500),
		IPHash:    hashIP(meta.IP),
		UserAgent: meta.UserAgent,
		CreatedAt: now,
	}
	event.DeviceType, event.Browser, event.OS = parseUserAgent(meta.UserAgent)
	return uc.repository.Record(ctx, event)
}

func (uc *AnalyticsUsecase) GetOverview(ctx context.Context, rangeValue domain.DateRange) (*domain.AnalyticsOverview, error) {
	return uc.repository.Overview(ctx, rangeValue)
}

func (uc *AnalyticsUsecase) GetTopPages(ctx context.Context, rangeValue domain.DateRange, limit int) ([]domain.TopPage, error) {
	if limit < 1 || limit > 100 {
		limit = 10
	}
	return uc.repository.TopPages(ctx, rangeValue, limit)
}

func (uc *AnalyticsUsecase) GetDailyViews(ctx context.Context, days int) ([]domain.DailyView, error) {
	if days < 1 || days > 365 {
		days = 30
	}
	return uc.repository.DailyViews(ctx, days)
}

func (uc *AnalyticsUsecase) GetContentStats(ctx context.Context, refType string, refID string) (*domain.ContentAnalytics, error) {
	refType = strings.TrimSpace(refType)
	refID = strings.TrimSpace(refID)
	if refID == "" || !domain.ValidRefType(refType) || refType == "" {
		return nil, domain.ErrInvalidAnalyticsEvent
	}
	return uc.repository.ContentStats(ctx, refType, refID)
}

func hashIP(ip string) string {
	ip = strings.TrimSpace(ip)
	if ip == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(ip))
	return hex.EncodeToString(sum[:])
}

func parseUserAgent(ua string) (device string, browser string, os string) {
	value := strings.ToLower(ua)
	device = "desktop"
	if strings.Contains(value, "mobile") || strings.Contains(value, "iphone") || strings.Contains(value, "android") {
		device = "mobile"
	}
	if strings.Contains(value, "ipad") || strings.Contains(value, "tablet") {
		device = "tablet"
	}
	switch {
	case strings.Contains(value, "edg/"):
		browser = "Edge"
	case strings.Contains(value, "chrome/") || strings.Contains(value, "crios/"):
		browser = "Chrome"
	case strings.Contains(value, "firefox/"):
		browser = "Firefox"
	case strings.Contains(value, "safari/"):
		browser = "Safari"
	default:
		browser = "Other"
	}
	switch {
	case strings.Contains(value, "iphone") || strings.Contains(value, "ipad"):
		os = "iOS"
	case strings.Contains(value, "android"):
		os = "Android"
	case strings.Contains(value, "mac os") || strings.Contains(value, "macintosh"):
		os = "macOS"
	case strings.Contains(value, "windows"):
		os = "Windows"
	case strings.Contains(value, "linux"):
		os = "Linux"
	default:
		os = "Other"
	}
	return device, browser, os
}

func truncate(value string, max int) string {
	if len(value) <= max {
		return value
	}
	return value[:max]
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))[:32]
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:])
}
