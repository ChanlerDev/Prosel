package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/dashboard"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	GetOverview(ctx context.Context) (*domain.DashboardOverview, error)
	GetStats(ctx context.Context) (*domain.DashboardStats, error)
	GetActivities(ctx context.Context, limit int) ([]domain.ActivityLog, error)
}

type DashboardHandler struct{ service DashboardService }

func NewDashboardHandler(service DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/dashboard/overview", h.overview)
	admin.GET("/dashboard/stats", h.stats)
	admin.GET("/activity-logs", h.activities)
}

type dashboardStatsResponse struct {
	TotalPosts      int64 `json:"totalPosts"`
	PublishedPosts  int64 `json:"publishedPosts"`
	DraftPosts      int64 `json:"draftPosts"`
	PendingComments int64 `json:"pendingComments"`
	TodayViews      int64 `json:"todayViews"`
	TotalViews      int64 `json:"totalViews"`
	Categories      int64 `json:"categories"`
	Tags            int64 `json:"tags"`
	Topics          int64 `json:"topics"`
}

type postSummaryResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Status      string     `json:"status"`
	ViewCount   int64      `json:"viewCount"`
	PublishedAt *time.Time `json:"publishedAt,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type activityLogResponse struct {
	ID         string    `json:"id"`
	ActorID    *string   `json:"actorId,omitempty"`
	Action     string    `json:"action"`
	EntityType string    `json:"entityType,omitempty"`
	EntityID   string    `json:"entityId,omitempty"`
	Message    string    `json:"message,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type overviewResponse struct {
	Stats       dashboardStatsResponse `json:"stats"`
	RecentPosts []postSummaryResponse  `json:"recentPosts"`
	Activities  []activityLogResponse  `json:"activities"`
}

func (h *DashboardHandler) overview(c *gin.Context) {
	overview, err := h.service.GetOverview(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Dashboard overview failed", nil)
		return
	}
	response.OK(c, toOverviewResponse(overview))
}

func (h *DashboardHandler) stats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Dashboard stats failed", nil)
		return
	}
	response.OK(c, toStatsResponse(*stats))
}

func (h *DashboardHandler) activities(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	activities, err := h.service.GetActivities(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Activity logs failed", nil)
		return
	}
	response.OK(c, toActivityResponses(activities))
}

func toOverviewResponse(overview *domain.DashboardOverview) overviewResponse {
	return overviewResponse{Stats: toStatsResponse(overview.Stats), RecentPosts: toPostSummaryResponses(overview.RecentPosts), Activities: toActivityResponses(overview.Activities)}
}

func toStatsResponse(stats domain.DashboardStats) dashboardStatsResponse {
	return dashboardStatsResponse{TotalPosts: stats.TotalPosts, PublishedPosts: stats.PublishedPosts, DraftPosts: stats.DraftPosts, PendingComments: stats.PendingComments, TodayViews: stats.TodayViews, TotalViews: stats.TotalViews, Categories: stats.Categories, Tags: stats.Tags, Topics: stats.Topics}
}

func toPostSummaryResponses(posts []domain.PostSummary) []postSummaryResponse {
	result := make([]postSummaryResponse, 0, len(posts))
	for _, post := range posts {
		result = append(result, postSummaryResponse{ID: post.ID, Title: post.Title, Slug: post.Slug, Status: post.Status, ViewCount: post.ViewCount, PublishedAt: post.PublishedAt, UpdatedAt: post.UpdatedAt})
	}
	return result
}

func toActivityResponses(activities []domain.ActivityLog) []activityLogResponse {
	result := make([]activityLogResponse, 0, len(activities))
	for _, activity := range activities {
		result = append(result, activityLogResponse{ID: activity.ID, ActorID: activity.ActorID, Action: activity.Action, EntityType: activity.EntityType, EntityID: activity.EntityID, Message: activity.Message, CreatedAt: activity.CreatedAt})
	}
	return result
}
