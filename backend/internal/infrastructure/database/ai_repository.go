package database

import (
	"context"
	"errors"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/chanler/prosel/backend/internal/domain/ai"
)

type AISummaryModel struct {
	ID          string         `gorm:"primaryKey;size:36"`
	RefType     string         `gorm:"size:20;not null;uniqueIndex:uniq_ai_summary"`
	RefID       string         `gorm:"size:36;not null;uniqueIndex:uniq_ai_summary"`
	Language    string         `gorm:"size:20;not null;uniqueIndex:uniq_ai_summary"`
	ContentHash string         `gorm:"size:64;not null"`
	Summary     string         `gorm:"not null"`
	Keywords    pq.StringArray `gorm:"type:text[];not null"`
	Provider    string         `gorm:"size:50"`
	Model       string         `gorm:"size:100"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (AISummaryModel) TableName() string { return "ai_summaries" }

type AITranslationModel struct {
	ID              string `gorm:"primaryKey;size:36"`
	RefType         string `gorm:"size:20;not null;uniqueIndex:uniq_ai_translation"`
	RefID           string `gorm:"size:36;not null;uniqueIndex:uniq_ai_translation"`
	SourceLanguage  string `gorm:"size:20;not null"`
	TargetLanguage  string `gorm:"size:20;not null;uniqueIndex:uniq_ai_translation"`
	ContentHash     string `gorm:"size:64;not null"`
	Title           string `gorm:"size:255"`
	Summary         string
	ContentMarkdown string `gorm:"not null"`
	Provider        string `gorm:"size:50"`
	Model           string `gorm:"size:100"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (AITranslationModel) TableName() string { return "ai_translations" }

type AIRepository struct{ db *gorm.DB }

func NewAIRepository(db *gorm.DB) *AIRepository { return &AIRepository{db: db} }

func (r *AIRepository) GetSummary(ctx context.Context, refType string, refID string, language string) (*domain.AISummary, error) {
	var model AISummaryModel
	err := r.db.WithContext(ctx).Where("ref_type = ? AND ref_id = ? AND language = ?", refType, refID, language).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAINotFound
	}
	if err != nil {
		return nil, err
	}
	return aiSummaryFromModel(model), nil
}

func (r *AIRepository) UpsertSummary(ctx context.Context, summary *domain.AISummary) error {
	model := toAISummaryModel(summary)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "ref_type"}, {Name: "ref_id"}, {Name: "language"}}, DoUpdates: clause.AssignmentColumns([]string{"content_hash", "summary", "keywords", "provider", "model", "updated_at"})}).Create(&model).Error
}

func (r *AIRepository) GetTranslation(ctx context.Context, refType string, refID string, targetLanguage string) (*domain.AITranslation, error) {
	var model AITranslationModel
	err := r.db.WithContext(ctx).Where("ref_type = ? AND ref_id = ? AND target_language = ?", refType, refID, targetLanguage).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrAINotFound
	}
	if err != nil {
		return nil, err
	}
	return aiTranslationFromModel(model), nil
}

func (r *AIRepository) UpsertTranslation(ctx context.Context, translation *domain.AITranslation) error {
	model := toAITranslationModel(translation)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "ref_type"}, {Name: "ref_id"}, {Name: "target_language"}}, DoUpdates: clause.AssignmentColumns([]string{"source_language", "content_hash", "title", "summary", "content_markdown", "provider", "model", "updated_at"})}).Create(&model).Error
}

func aiSummaryFromModel(model AISummaryModel) *domain.AISummary {
	return &domain.AISummary{ID: model.ID, RefType: model.RefType, RefID: model.RefID, Language: model.Language, ContentHash: model.ContentHash, Summary: model.Summary, Keywords: []string(model.Keywords), Provider: model.Provider, Model: model.Model, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toAISummaryModel(summary *domain.AISummary) AISummaryModel {
	return AISummaryModel{ID: summary.ID, RefType: summary.RefType, RefID: summary.RefID, Language: summary.Language, ContentHash: summary.ContentHash, Summary: summary.Summary, Keywords: pq.StringArray(summary.Keywords), Provider: summary.Provider, Model: summary.Model, CreatedAt: summary.CreatedAt, UpdatedAt: summary.UpdatedAt}
}

func aiTranslationFromModel(model AITranslationModel) *domain.AITranslation {
	return &domain.AITranslation{ID: model.ID, RefType: model.RefType, RefID: model.RefID, SourceLanguage: model.SourceLanguage, TargetLanguage: model.TargetLanguage, ContentHash: model.ContentHash, Title: model.Title, Summary: model.Summary, ContentMarkdown: model.ContentMarkdown, Provider: model.Provider, Model: model.Model, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toAITranslationModel(translation *domain.AITranslation) AITranslationModel {
	return AITranslationModel{ID: translation.ID, RefType: translation.RefType, RefID: translation.RefID, SourceLanguage: translation.SourceLanguage, TargetLanguage: translation.TargetLanguage, ContentHash: translation.ContentHash, Title: translation.Title, Summary: translation.Summary, ContentMarkdown: translation.ContentMarkdown, Provider: translation.Provider, Model: translation.Model, CreatedAt: translation.CreatedAt, UpdatedAt: translation.UpdatedAt}
}
