package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/ai"
	usecase "github.com/chanler/prosel/backend/internal/usecase/ai"
	"github.com/gin-gonic/gin"
)

type fakeAIService struct {
	summary     *domain.AISummary
	translation *domain.AITranslation
	req         usecase.GenerateSummaryRequest
}

func (s *fakeAIService) GenerateSummary(ctx context.Context, req usecase.GenerateSummaryRequest) (*domain.AISummary, error) {
	s.req = req
	return s.summary, nil
}
func (s *fakeAIService) GenerateTranslation(ctx context.Context, req usecase.GenerateTranslationRequest) (*domain.AITranslation, error) {
	return s.translation, nil
}
func (s *fakeAIService) GetPublicSummary(ctx context.Context, refType string, refID string, lang string) (*domain.AISummary, error) {
	return s.summary, nil
}
func (s *fakeAIService) GetPublicTranslation(ctx context.Context, refType string, refID string, lang string) (*domain.AITranslation, error) {
	return s.translation, nil
}
func (s *fakeAIService) Configured() bool { return true }

func TestAIHandlerGenerateSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeAIService{summary: &domain.AISummary{ID: "ai-1", RefType: "post", RefID: "post-1", Language: "en", Summary: "Short", Keywords: []string{"go"}, UpdatedAt: time.Now()}}
	router := gin.New()
	NewAIHandler(service).RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/ai/summaries", strings.NewReader(`{"refType":"post","refId":"post-1","language":"en"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	var body struct {
		Data aiSummaryResponse `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Summary != "Short" || service.req.RefID != "post-1" {
		t.Fatalf("body = %#v req = %#v", body.Data, service.req)
	}
}
