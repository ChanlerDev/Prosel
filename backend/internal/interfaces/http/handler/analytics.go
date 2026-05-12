package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/analytics"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/analytics"
	"github.com/gin-gonic/gin"
)

type AnalyticsService interface {
	RecordPageView(ctx context.Context, req usecase.PageViewRequest, meta usecase.ClientMeta) error
	GetOverview(ctx context.Context, rangeValue domain.DateRange) (*domain.AnalyticsOverview, error)
	GetTopPages(ctx context.Context, rangeValue domain.DateRange, limit int) ([]domain.TopPage, error)
	GetDailyViews(ctx context.Context, days int) ([]domain.DailyView, error)
}

type AnalyticsHandler struct{ service AnalyticsService }

func NewAnalyticsHandler(service AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.POST("/analytics/page-view", h.recordPageView)
}

func (h *AnalyticsHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/analytics/overview", h.overview)
	admin.GET("/analytics/top-pages", h.topPages)
	admin.GET("/analytics/daily", h.daily)
}

type pageViewRequest struct {
	Path    string `json:"path" binding:"required"`
	RefType string `json:"refType"`
	RefID   string `json:"refId"`
	Referer string `json:"referer"`
}

type analyticsOverviewResponse struct {
	TodayViews  int64                `json:"todayViews"`
	WeekViews   int64                `json:"weekViews"`
	MonthViews  int64                `json:"monthViews"`
	TopPages    []topPageResponse    `json:"topPages"`
	TopReferers []topRefererResponse `json:"topReferers"`
	Devices     []deviceStatResponse `json:"devices"`
}

type topPageResponse struct {
	Path    string `json:"path"`
	RefType string `json:"refType,omitempty"`
	RefID   string `json:"refId,omitempty"`
	Views   int64  `json:"views"`
}

type topRefererResponse struct {
	Referer string `json:"referer"`
	Views   int64  `json:"views"`
}

type deviceStatResponse struct {
	DeviceType string `json:"deviceType"`
	Views      int64  `json:"views"`
}

type dailyViewResponse struct {
	Date  string `json:"date"`
	Views int64  `json:"views"`
}

func (h *AnalyticsHandler) recordPageView(c *gin.Context) {
	var req pageViewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid analytics request", nil)
		return
	}
	err := h.service.RecordPageView(c.Request.Context(), usecase.PageViewRequest{Path: req.Path, RefType: req.RefType, RefID: req.RefID, Referer: req.Referer}, usecase.ClientMeta{IP: c.ClientIP(), UserAgent: c.Request.UserAgent()})
	if err != nil {
		h.handleAnalyticsError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *AnalyticsHandler) overview(c *gin.Context) {
	overview, err := h.service.GetOverview(c.Request.Context(), domain.NormalizeDateRange(c.DefaultQuery("range", string(domain.DateRange30d))))
	if err != nil {
		h.handleAnalyticsError(c, err)
		return
	}
	response.OK(c, toAnalyticsOverviewResponse(overview))
}

func (h *AnalyticsHandler) topPages(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	pages, err := h.service.GetTopPages(c.Request.Context(), domain.NormalizeDateRange(c.DefaultQuery("range", string(domain.DateRange30d))), limit)
	if err != nil {
		h.handleAnalyticsError(c, err)
		return
	}
	response.OK(c, toTopPageResponses(pages))
}

func (h *AnalyticsHandler) daily(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	views, err := h.service.GetDailyViews(c.Request.Context(), days)
	if err != nil {
		h.handleAnalyticsError(c, err)
		return
	}
	response.OK(c, toDailyViewResponses(views))
}

func (h *AnalyticsHandler) handleAnalyticsError(c *gin.Context, err error) {
	if err == domain.ErrInvalidAnalyticsEvent {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid analytics request", nil)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Analytics request failed", nil)
}

func toAnalyticsOverviewResponse(overview *domain.AnalyticsOverview) analyticsOverviewResponse {
	return analyticsOverviewResponse{TodayViews: overview.TodayViews, WeekViews: overview.WeekViews, MonthViews: overview.MonthViews, TopPages: toTopPageResponses(overview.TopPages), TopReferers: toTopRefererResponses(overview.TopReferers), Devices: toDeviceStatResponses(overview.Devices)}
}

func toTopPageResponses(pages []domain.TopPage) []topPageResponse {
	result := make([]topPageResponse, 0, len(pages))
	for _, page := range pages {
		result = append(result, topPageResponse{Path: page.Path, RefType: page.RefType, RefID: page.RefID, Views: page.Views})
	}
	return result
}

func toTopRefererResponses(referers []domain.TopReferer) []topRefererResponse {
	result := make([]topRefererResponse, 0, len(referers))
	for _, referer := range referers {
		result = append(result, topRefererResponse{Referer: referer.Referer, Views: referer.Views})
	}
	return result
}

func toDeviceStatResponses(devices []domain.DeviceStat) []deviceStatResponse {
	result := make([]deviceStatResponse, 0, len(devices))
	for _, device := range devices {
		result = append(result, deviceStatResponse{DeviceType: device.DeviceType, Views: device.Views})
	}
	return result
}

func toDailyViewResponses(views []domain.DailyView) []dailyViewResponse {
	result := make([]dailyViewResponse, 0, len(views))
	for _, view := range views {
		result = append(result, dailyViewResponse{Date: view.Date.Format(time.DateOnly), Views: view.Views})
	}
	return result
}
