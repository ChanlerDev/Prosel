package ai

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/ai"
	postdomain "github.com/chanler/prosel/backend/internal/domain/post"
)

type PostReader interface {
	GetAdminPost(ctx context.Context, id string) (*postdomain.Post, error)
}

type AIUsecase struct {
	repo   domain.Repository
	client domain.Client
	posts  PostReader
}

type GenerateSummaryRequest struct {
	RefType  string
	RefID    string
	Language string
}

type GenerateTranslationRequest struct {
	RefType        string
	RefID          string
	SourceLanguage string
	TargetLanguage string
}

func NewAIUsecase(repo domain.Repository, client domain.Client, posts PostReader) *AIUsecase {
	return &AIUsecase{repo: repo, client: client, posts: posts}
}

func (uc *AIUsecase) GenerateSummary(ctx context.Context, req GenerateSummaryRequest) (*domain.AISummary, error) {
	if uc.client == nil || !uc.client.Configured() {
		return nil, domain.ErrAIUnavailable
	}
	refType, refID := normalizeRef(req.RefType, req.RefID)
	language := normalizeLanguage(req.Language, "zh")
	post, err := uc.loadPost(ctx, refType, refID)
	if err != nil {
		return nil, err
	}
	output, err := uc.client.Summarize(ctx, domain.SummarizeInput{Title: post.Title, ContentMarkdown: post.ContentMarkdown, Language: language})
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	summary := &domain.AISummary{ID: newID(), RefType: refType, RefID: refID, Language: language, ContentHash: contentHash(post.Title, post.Excerpt, post.ContentMarkdown), Summary: strings.TrimSpace(output.Summary), Keywords: cleanStrings(output.Keywords), Provider: output.Provider, Model: output.Model, CreatedAt: now, UpdatedAt: now}
	if summary.Summary == "" {
		return nil, domain.ErrInvalidAIRef
	}
	if err := uc.repo.UpsertSummary(ctx, summary); err != nil {
		return nil, err
	}
	return summary, nil
}

func (uc *AIUsecase) GenerateTranslation(ctx context.Context, req GenerateTranslationRequest) (*domain.AITranslation, error) {
	if uc.client == nil || !uc.client.Configured() {
		return nil, domain.ErrAIUnavailable
	}
	refType, refID := normalizeRef(req.RefType, req.RefID)
	sourceLanguage := normalizeLanguage(req.SourceLanguage, "zh")
	targetLanguage := normalizeLanguage(req.TargetLanguage, "en")
	post, err := uc.loadPost(ctx, refType, refID)
	if err != nil {
		return nil, err
	}
	output, err := uc.client.Translate(ctx, domain.TranslateInput{Title: post.Title, Summary: firstNonEmpty(post.Excerpt, post.SEODescription), ContentMarkdown: post.ContentMarkdown, SourceLanguage: sourceLanguage, TargetLanguage: targetLanguage})
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	translation := &domain.AITranslation{ID: newID(), RefType: refType, RefID: refID, SourceLanguage: sourceLanguage, TargetLanguage: targetLanguage, ContentHash: contentHash(post.Title, post.Excerpt, post.ContentMarkdown), Title: strings.TrimSpace(output.Title), Summary: strings.TrimSpace(output.Summary), ContentMarkdown: strings.TrimSpace(output.ContentMarkdown), Provider: output.Provider, Model: output.Model, CreatedAt: now, UpdatedAt: now}
	if translation.ContentMarkdown == "" {
		return nil, domain.ErrInvalidAIRef
	}
	if err := uc.repo.UpsertTranslation(ctx, translation); err != nil {
		return nil, err
	}
	return translation, nil
}

func (uc *AIUsecase) GetPublicSummary(ctx context.Context, refType string, refID string, lang string) (*domain.AISummary, error) {
	refType, refID = normalizeRef(refType, refID)
	if !domain.ValidRefType(refType) || refID == "" {
		return nil, domain.ErrInvalidAIRef
	}
	return uc.repo.GetSummary(ctx, refType, refID, normalizeLanguage(lang, "zh"))
}

func (uc *AIUsecase) GetPublicTranslation(ctx context.Context, refType string, refID string, lang string) (*domain.AITranslation, error) {
	refType, refID = normalizeRef(refType, refID)
	if !domain.ValidRefType(refType) || refID == "" {
		return nil, domain.ErrInvalidAIRef
	}
	return uc.repo.GetTranslation(ctx, refType, refID, normalizeLanguage(lang, "en"))
}

func (uc *AIUsecase) Configured() bool {
	return uc.client != nil && uc.client.Configured()
}

func (uc *AIUsecase) loadPost(ctx context.Context, refType string, refID string) (*postdomain.Post, error) {
	if refType != domain.RefTypePost || refID == "" || uc.posts == nil {
		return nil, domain.ErrInvalidAIRef
	}
	return uc.posts.GetAdminPost(ctx, refID)
}

func normalizeRef(refType string, refID string) (string, string) {
	return strings.TrimSpace(refType), strings.TrimSpace(refID)
}

func normalizeLanguage(value string, fallback string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return fallback
	}
	return value
}

func contentHash(parts ...string) string {
	h := sha256.New()
	for _, part := range parts {
		h.Write([]byte(part))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

func cleanStrings(values []string) []string {
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			result = append(result, value)
		}
	}
	return result
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func newID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return hex.EncodeToString(b[0:4]) + "-" + hex.EncodeToString(b[4:6]) + "-" + hex.EncodeToString(b[6:8]) + "-" + hex.EncodeToString(b[8:10]) + "-" + hex.EncodeToString(b[10:16])
}
