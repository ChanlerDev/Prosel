package ai

import "context"

type Repository interface {
	GetSummary(ctx context.Context, refType string, refID string, language string) (*AISummary, error)
	UpsertSummary(ctx context.Context, summary *AISummary) error
	GetTranslation(ctx context.Context, refType string, refID string, targetLanguage string) (*AITranslation, error)
	UpsertTranslation(ctx context.Context, translation *AITranslation) error
}
