package ai

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/ai"
	postdomain "github.com/chanler/prosel/backend/internal/domain/post"
)

type fakeAIRepo struct {
	summary     *domain.AISummary
	translation *domain.AITranslation
}

func (r *fakeAIRepo) GetSummary(ctx context.Context, refType string, refID string, language string) (*domain.AISummary, error) {
	if r.summary == nil {
		return nil, domain.ErrAINotFound
	}
	return r.summary, nil
}
func (r *fakeAIRepo) UpsertSummary(ctx context.Context, summary *domain.AISummary) error {
	r.summary = summary
	return nil
}
func (r *fakeAIRepo) GetTranslation(ctx context.Context, refType string, refID string, targetLanguage string) (*domain.AITranslation, error) {
	if r.translation == nil {
		return nil, domain.ErrAINotFound
	}
	return r.translation, nil
}
func (r *fakeAIRepo) UpsertTranslation(ctx context.Context, translation *domain.AITranslation) error {
	r.translation = translation
	return nil
}

type fakeAIClient struct {
	summaryOutput     *domain.SummarizeOutput
	translationOutput *domain.TranslateOutput
}

func (c *fakeAIClient) Summarize(ctx context.Context, input domain.SummarizeInput) (*domain.SummarizeOutput, error) {
	return c.summaryOutput, nil
}
func (c *fakeAIClient) Translate(ctx context.Context, input domain.TranslateInput) (*domain.TranslateOutput, error) {
	return c.translationOutput, nil
}
func (c *fakeAIClient) Configured() bool { return c != nil }

type fakePostReader struct{ post *postdomain.Post }

func (r fakePostReader) GetAdminPost(ctx context.Context, id string) (*postdomain.Post, error) {
	return r.post, nil
}

func TestGenerateSummaryPersistsAIResult(t *testing.T) {
	repo := &fakeAIRepo{}
	client := &fakeAIClient{summaryOutput: &domain.SummarizeOutput{Summary: "Short summary", Keywords: []string{"go", "blog"}, Provider: "openai", Model: "gpt-test"}}
	uc := NewAIUsecase(repo, client, fakePostReader{post: &postdomain.Post{ID: "post-1", Title: "Hello", ContentMarkdown: "Long markdown", UpdatedAt: time.Now()}})

	summary, err := uc.GenerateSummary(context.Background(), GenerateSummaryRequest{RefType: "post", RefID: "post-1", Language: "en"})
	if err != nil {
		t.Fatalf("GenerateSummary() error = %v", err)
	}
	if summary.Summary != "Short summary" || summary.Language != "en" || summary.ContentHash == "" {
		t.Fatalf("summary = %#v", summary)
	}
	if repo.summary == nil || len(repo.summary.Keywords) != 2 {
		t.Fatalf("summary was not persisted: %#v", repo.summary)
	}
}

func TestGenerateSummaryReturnsUnavailableWhenClientMissing(t *testing.T) {
	uc := NewAIUsecase(&fakeAIRepo{}, nil, fakePostReader{post: &postdomain.Post{ID: "post-1", Title: "Hello"}})

	_, err := uc.GenerateSummary(context.Background(), GenerateSummaryRequest{RefType: "post", RefID: "post-1", Language: "zh"})
	if !errors.Is(err, domain.ErrAIUnavailable) {
		t.Fatalf("GenerateSummary() error = %v, want %v", err, domain.ErrAIUnavailable)
	}
}

func TestGenerateTranslationPersistsAIResult(t *testing.T) {
	repo := &fakeAIRepo{}
	client := &fakeAIClient{translationOutput: &domain.TranslateOutput{Title: "Bonjour", Summary: "Résumé", ContentMarkdown: "Contenu", Provider: "openai", Model: "gpt-test"}}
	uc := NewAIUsecase(repo, client, fakePostReader{post: &postdomain.Post{ID: "post-1", Title: "Hello", Excerpt: "Summary", ContentMarkdown: "Content"}})

	translation, err := uc.GenerateTranslation(context.Background(), GenerateTranslationRequest{RefType: "post", RefID: "post-1", SourceLanguage: "en", TargetLanguage: "fr"})
	if err != nil {
		t.Fatalf("GenerateTranslation() error = %v", err)
	}
	if translation.Title != "Bonjour" || translation.ContentMarkdown != "Contenu" || translation.TargetLanguage != "fr" {
		t.Fatalf("translation = %#v", translation)
	}
}
