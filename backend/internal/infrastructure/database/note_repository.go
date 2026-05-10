package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/note"
)

type NoteModel struct {
	ID              string `gorm:"primaryKey;size:36"`
	AuthorID        string `gorm:"size:36;index"`
	Title           string `gorm:"size:200"`
	Slug            string `gorm:"size:255;uniqueIndex"`
	ContentMarkdown string `gorm:"not null"`
	ContentText     string `gorm:"not null"`
	Mood            string `gorm:"size:80"`
	Weather         string `gorm:"size:80"`
	Location        string `gorm:"size:120"`
	Status          string `gorm:"size:20;not null;index"`
	PinnedAt        *time.Time
	PublishedAt     *time.Time
	ViewCount       int64
	LikeCount       int64
	CommentCount    int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (NoteModel) TableName() string { return "notes" }

type NoteRepository struct{ db *gorm.DB }

func NewNoteRepository(db *gorm.DB) *NoteRepository { return &NoteRepository{db: db} }

func (r *NoteRepository) Create(ctx context.Context, note *domain.Note) error {
	return r.db.WithContext(ctx).Create(toNoteModel(note)).Error
}

func (r *NoteRepository) Update(ctx context.Context, note *domain.Note) error {
	return r.db.WithContext(ctx).Model(&NoteModel{}).Where("id = ?", note.ID).Updates(map[string]any{
		"title":            note.Title,
		"slug":             note.Slug,
		"content_markdown": note.ContentMarkdown,
		"content_text":     note.ContentText,
		"mood":             note.Mood,
		"weather":          note.Weather,
		"location":         note.Location,
		"status":           string(note.Status),
		"published_at":     note.PublishedAt,
		"updated_at":       note.UpdatedAt,
	}).Error
}

func (r *NoteRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&NoteModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNoteNotFound
	}
	return nil
}

func (r *NoteRepository) GetByID(ctx context.Context, id string) (*domain.Note, error) {
	var model NoteModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return noteFromModel(model, err)
}

func (r *NoteRepository) GetBySlug(ctx context.Context, slug string) (*domain.Note, error) {
	var model NoteModel
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error
	return noteFromModel(model, err)
}

func (r *NoteRepository) ListPublic(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	status := domain.NotePublished
	filter.Status = &status
	return r.list(ctx, filter)
}

func (r *NoteRepository) ListAdmin(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	return r.list(ctx, filter)
}

func (r *NoteRepository) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	query := r.db.WithContext(ctx).Model(&NoteModel{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *NoteRepository) IncrementView(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&NoteModel{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *NoteRepository) SetPinned(ctx context.Context, id string, pinnedAt *time.Time) error {
	result := r.db.WithContext(ctx).Model(&NoteModel{}).Where("id = ?", id).Updates(map[string]any{"pinned_at": pinnedAt, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNoteNotFound
	}
	return nil
}

func (r *NoteRepository) list(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	query := r.db.WithContext(ctx).Model(&NoteModel{})
	query = applyNoteFilter(query, filter)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	var models []NoteModel
	if err := query.Order("pinned_at DESC NULLS LAST").Order("published_at DESC NULLS LAST").Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&models).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	notes := make([]domain.Note, 0, len(models))
	for _, model := range models {
		notes = append(notes, *noteFromModelNoError(model))
	}
	return notes, domain.NewPagination(page, perPage, total), nil
}

func applyNoteFilter(query *gorm.DB, filter domain.NoteFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR content_text ILIKE ? OR mood ILIKE ? OR location ILIKE ?", search, search, search, search)
	}
	return query
}

func noteFromModel(model NoteModel, err error) (*domain.Note, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrNoteNotFound
	}
	if err != nil {
		return nil, err
	}
	return noteFromModelNoError(model), nil
}

func noteFromModelNoError(model NoteModel) *domain.Note {
	return &domain.Note{ID: model.ID, AuthorID: model.AuthorID, Title: model.Title, Slug: model.Slug, ContentMarkdown: model.ContentMarkdown, ContentText: model.ContentText, Mood: model.Mood, Weather: model.Weather, Location: model.Location, Status: domain.NoteStatus(model.Status), PinnedAt: model.PinnedAt, PublishedAt: model.PublishedAt, ViewCount: model.ViewCount, LikeCount: model.LikeCount, CommentCount: model.CommentCount, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toNoteModel(note *domain.Note) *NoteModel {
	return &NoteModel{ID: note.ID, AuthorID: note.AuthorID, Title: note.Title, Slug: note.Slug, ContentMarkdown: note.ContentMarkdown, ContentText: note.ContentText, Mood: note.Mood, Weather: note.Weather, Location: note.Location, Status: string(note.Status), PinnedAt: note.PinnedAt, PublishedAt: note.PublishedAt, ViewCount: note.ViewCount, LikeCount: note.LikeCount, CommentCount: note.CommentCount, CreatedAt: note.CreatedAt, UpdatedAt: note.UpdatedAt}
}
