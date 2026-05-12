package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/ai"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/ai"
	"github.com/gin-gonic/gin"
)

type AIService interface {
	GenerateSummary(ctx context.Context, req usecase.GenerateSummaryRequest) (*domain.AISummary, error)
	GenerateTranslation(ctx context.Context, req usecase.GenerateTranslationRequest) (*domain.AITranslation, error)
	GetPublicSummary(ctx context.Context, refType string, refID string, lang string) (*domain.AISummary, error)
	GetPublicTranslation(ctx context.Context, refType string, refID string, lang string) (*domain.AITranslation, error)
	Configured() bool
}

type AIHandler struct{ service AIService }

func NewAIHandler(service AIService) *AIHandler { return &AIHandler{service: service} }

func (h *AIHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/ai/summaries", h.getSummary)
	group.GET("/ai/translations", h.getTranslation)
}

func (h *AIHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/ai/status", h.status)
	admin.POST("/ai/summaries", h.generateSummary)
	admin.POST("/ai/translations", h.generateTranslation)
}

type generateSummaryRequest struct {
	RefType  string `json:"refType" binding:"required"`
	RefID    string `json:"refId" binding:"required"`
	Language string `json:"language"`
}

type generateTranslationRequest struct {
	RefType        string `json:"refType" binding:"required"`
	RefID          string `json:"refId" binding:"required"`
	SourceLanguage string `json:"sourceLanguage"`
	TargetLanguage string `json:"targetLanguage" binding:"required"`
}

type aiStatusResponse struct {
	Configured bool `json:"configured"`
}

type aiSummaryResponse struct {
	ID          string    `json:"id"`
	RefType     string    `json:"refType"`
	RefID       string    `json:"refId"`
	Language    string    `json:"language"`
	ContentHash string    `json:"contentHash,omitempty"`
	Summary     string    `json:"summary"`
	Keywords    []string  `json:"keywords"`
	Provider    string    `json:"provider,omitempty"`
	Model       string    `json:"model,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type aiTranslationResponse struct {
	ID              string    `json:"id"`
	RefType         string    `json:"refType"`
	RefID           string    `json:"refId"`
	SourceLanguage  string    `json:"sourceLanguage"`
	TargetLanguage  string    `json:"targetLanguage"`
	ContentHash     string    `json:"contentHash,omitempty"`
	Title           string    `json:"title,omitempty"`
	Summary         string    `json:"summary,omitempty"`
	ContentMarkdown string    `json:"contentMarkdown"`
	Provider        string    `json:"provider,omitempty"`
	Model           string    `json:"model,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (h *AIHandler) status(c *gin.Context) {
	response.OK(c, aiStatusResponse{Configured: h.service.Configured()})
}

func (h *AIHandler) getSummary(c *gin.Context) {
	summary, err := h.service.GetPublicSummary(c.Request.Context(), c.Query("refType"), c.Query("refId"), c.Query("lang"))
	if err != nil {
		h.handleAIError(c, err)
		return
	}
	response.OK(c, toAISummaryResponse(summary))
}

func (h *AIHandler) getTranslation(c *gin.Context) {
	translation, err := h.service.GetPublicTranslation(c.Request.Context(), c.Query("refType"), c.Query("refId"), c.Query("lang"))
	if err != nil {
		h.handleAIError(c, err)
		return
	}
	response.OK(c, toAITranslationResponse(translation))
}

func (h *AIHandler) generateSummary(c *gin.Context) {
	var req generateSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid AI summary request", nil)
		return
	}
	summary, err := h.service.GenerateSummary(c.Request.Context(), usecase.GenerateSummaryRequest{RefType: req.RefType, RefID: req.RefID, Language: req.Language})
	if err != nil {
		h.handleAIError(c, err)
		return
	}
	response.OK(c, toAISummaryResponse(summary))
}

func (h *AIHandler) generateTranslation(c *gin.Context) {
	var req generateTranslationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid AI translation request", nil)
		return
	}
	translation, err := h.service.GenerateTranslation(c.Request.Context(), usecase.GenerateTranslationRequest{RefType: req.RefType, RefID: req.RefID, SourceLanguage: req.SourceLanguage, TargetLanguage: req.TargetLanguage})
	if err != nil {
		h.handleAIError(c, err)
		return
	}
	response.OK(c, toAITranslationResponse(translation))
}

func (h *AIHandler) handleAIError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrAINotFound):
		response.Error(c, http.StatusNotFound, "AI_RESULT_NOT_FOUND", "AI result not found", nil)
	case errors.Is(err, domain.ErrAIUnavailable):
		response.Error(c, http.StatusServiceUnavailable, "AI_PROVIDER_UNCONFIGURED", "AI provider is not configured", nil)
	case errors.Is(err, domain.ErrInvalidAIRef):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid AI request", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "AI request failed", nil)
	}
}

func toAISummaryResponse(summary *domain.AISummary) aiSummaryResponse {
	return aiSummaryResponse{ID: summary.ID, RefType: summary.RefType, RefID: summary.RefID, Language: summary.Language, ContentHash: summary.ContentHash, Summary: summary.Summary, Keywords: summary.Keywords, Provider: summary.Provider, Model: summary.Model, CreatedAt: summary.CreatedAt, UpdatedAt: summary.UpdatedAt}
}

func toAITranslationResponse(translation *domain.AITranslation) aiTranslationResponse {
	return aiTranslationResponse{ID: translation.ID, RefType: translation.RefType, RefID: translation.RefID, SourceLanguage: translation.SourceLanguage, TargetLanguage: translation.TargetLanguage, ContentHash: translation.ContentHash, Title: translation.Title, Summary: translation.Summary, ContentMarkdown: translation.ContentMarkdown, Provider: translation.Provider, Model: translation.Model, CreatedAt: translation.CreatedAt, UpdatedAt: translation.UpdatedAt}
}
