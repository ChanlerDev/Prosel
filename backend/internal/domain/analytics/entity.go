package analytics

import (
	"errors"
	"time"
)

const (
	EventTypePageView = "page_view"
	RefTypePost       = "post"
	RefTypeNote       = "note"
	RefTypePage       = "page"
)

var ErrInvalidAnalyticsEvent = errors.New("invalid analytics event")

type AnalyticsEvent struct {
	ID         string
	EventType  string
	Path       string
	RefType    string
	RefID      string
	Referer    string
	IPHash     string
	UserAgent  string
	Country    string
	DeviceType string
	Browser    string
	OS         string
	CreatedAt  time.Time
}

type DateRange string

const (
	DateRange7d  DateRange = "7d"
	DateRange30d DateRange = "30d"
	DateRange90d DateRange = "90d"
)

func NormalizeDateRange(value string) DateRange {
	switch DateRange(value) {
	case DateRange7d, DateRange30d, DateRange90d:
		return DateRange(value)
	default:
		return DateRange30d
	}
}

func ValidRefType(value string) bool {
	return value == "" || value == RefTypePost || value == RefTypeNote || value == RefTypePage
}

type AnalyticsOverview struct {
	TodayViews  int64
	WeekViews   int64
	MonthViews  int64
	TopPages    []TopPage
	TopReferers []TopReferer
	Devices     []DeviceStat
}

type TopPage struct {
	Path    string
	RefType string
	RefID   string
	Views   int64
}

type TopReferer struct {
	Referer string
	Views   int64
}

type DeviceStat struct {
	DeviceType string
	Views      int64
}

type DailyView struct {
	Date  time.Time
	Views int64
}

type ContentAnalytics struct {
	RefType string
	RefID   string
	Views   int64
}
