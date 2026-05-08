package handler

import (
	"context"
	"net/http"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/system"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

type SystemService interface {
	GetPublicSettings(ctx context.Context) ([]domain.SiteSetting, error)
	CheckHealth(ctx context.Context) (*domain.HealthStatus, error)
}

type SystemHandler struct {
	service SystemService
}

func NewSystemHandler(service SystemService) *SystemHandler {
	return &SystemHandler{service: service}
}

func (h *SystemHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/health", h.health)
	group.GET("/settings/public", h.publicSettings)
}

type healthResponse struct {
	Status     string    `json:"status"`
	DatabaseOK bool      `json:"databaseOk"`
	RedisOK    bool      `json:"redisOk"`
	Version    string    `json:"version"`
	CheckedAt  time.Time `json:"checkedAt"`
}

func (h *SystemHandler) health(c *gin.Context) {
	status, err := h.service.CheckHealth(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusServiceUnavailable, "SERVICE_UNAVAILABLE", "Health check failed", nil)
		return
	}

	response.OK(c, healthResponse{
		Status:     status.Status,
		DatabaseOK: status.DatabaseOK,
		RedisOK:    status.RedisOK,
		Version:    status.Version,
		CheckedAt:  status.CheckedAt,
	})
}

type publicSettingResponse struct {
	Key       string `json:"key"`
	Value     any    `json:"value"`
	ValueType string `json:"valueType"`
}

func (h *SystemHandler) publicSettings(c *gin.Context) {
	settings, err := h.service.GetPublicSettings(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to load settings", nil)
		return
	}

	result := make(map[string]any, len(settings))
	for _, setting := range settings {
		result[setting.Key] = setting.PublicValue()
	}
	response.OK(c, result)
}
