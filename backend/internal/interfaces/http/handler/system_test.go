package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/system"
	"github.com/gin-gonic/gin"
)

type fakeSystemUsecase struct{}

func (fakeSystemUsecase) GetPublicSettings(ctx context.Context) ([]domain.SiteSetting, error) {
	return []domain.SiteSetting{{Key: "site_name", Value: "Prosel", ValueType: domain.ValueTypeString}}, nil
}

func (fakeSystemUsecase) CheckHealth(ctx context.Context) (*domain.HealthStatus, error) {
	return &domain.HealthStatus{Status: domain.StatusHealthy, DatabaseOK: true, RedisOK: true, Version: "test", CheckedAt: time.Date(2026, 5, 9, 0, 0, 0, 0, time.UTC)}, nil
}

func TestSystemHandlerHealthUsesUnifiedResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewSystemHandler(fakeSystemUsecase{}).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/health", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var body map[string]any
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body["data"] == nil {
		t.Fatalf("response missing data: %s", recorder.Body.String())
	}
	if body["error"] != nil {
		t.Fatalf("response error = %#v, want nil", body["error"])
	}
}

func TestSystemHandlerPublicSettingsReturnsMap(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	NewSystemHandler(fakeSystemUsecase{}).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/settings/public", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}

	var body struct {
		Data map[string]any `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data["site_name"] != "Prosel" {
		t.Fatalf("site_name = %#v, want Prosel", body.Data["site_name"])
	}
}
