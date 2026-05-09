package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/post"
)

type PostModel struct {
	ID              string  `gorm:"primaryKey;size:36"`
	AuthorID        string  `gorm:"size:36;index"`
	CategoryID      *string `gorm:"size:36;index"`
	Title           string  `gorm:"size:200;not null"`
	Slug            string  `gorm:"size:255;not null;uniqueIndex"`
	Excerpt         string  `gorm:"size:500"`
	ContentMarkdown string  `gorm:"not null"`
	ContentText     string  `gorm:"not null"`
	CoverImage      string  `gorm:"size:500"`
	Status          string  `gorm:"size:20;not null;index"`
	Featured        bool
	PinnedAt        *time.Time
	PublishedAt     *time.Time
	SEOTitle        string `gorm:"size:200"`
	SEODescription  string `gorm:"size:500"`
	ViewCount       int64
	LikeCount       int64
	CommentCount    int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (PostModel) TableName() string { return "posts" }

type PostRepository struct{ db *gorm.DB }

func NewPostRepository(db *gorm.DB) *PostRepository { return &PostRepository{db: db} }

func (r *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Create(toPostModel(post)).Error
}

func (r *PostRepository) Update(ctx context.Context, post *domain.Post) error {
	return r.db.WithContext(ctx).Model(&PostModel{}).Where("id = ?", post.ID).Updates(map[string]any{
		"category_id":      post.CategoryID,
		"title":            post.Title,
		"slug":             post.Slug,
		"excerpt":          post.Excerpt,
		"content_markdown": post.ContentMarkdown,
		"content_text":     post.ContentText,
		"cover_image":      post.CoverImage,
		"featured":         post.Featured,
		"seo_title":        post.SEOTitle,
		"seo_description":  post.SEODescription,
		"updated_at":       post.UpdatedAt,
	}).Error
}

func (r *PostRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&PostModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPostNotFound
	}
	return nil
}

func (r *PostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	var model PostModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return postFromModel(model, err)
}

func (r *PostRepository) GetBySlug(ctx context.Context, slug string, includeDraft bool) (*domain.Post, error) {
	var model PostModel
	query := r.db.WithContext(ctx).Where("slug = ?", slug)
	if !includeDraft {
		query = query.Where("status = ?", string(domain.PostPublished))
	}
	err := query.First(&model).Error
	return postFromModel(model, err)
}

func (r *PostRepository) List(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	query := r.db.WithContext(ctx).Model(&PostModel{})
	query = applyPostFilter(query, filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}

	var models []PostModel
	err := query.Order("pinned_at DESC NULLS LAST").Order("published_at DESC NULLS LAST").Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&models).Error
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	posts := make([]domain.Post, 0, len(models))
	for _, model := range models {
		posts = append(posts, *postFromModelNoError(model))
	}
	return posts, domain.NewPagination(page, perPage, total), nil
}

func (r *PostRepository) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	query := r.db.WithContext(ctx).Model(&PostModel{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostRepository) IncrementView(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&PostModel{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *PostRepository) SetStatus(ctx context.Context, id string, status domain.PostStatus, publishedAt *time.Time) error {
	result := r.db.WithContext(ctx).Model(&PostModel{}).Where("id = ?", id).Updates(map[string]any{"status": string(status), "published_at": publishedAt, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPostNotFound
	}
	return nil
}

func applyPostFilter(query *gorm.DB, filter domain.PostListFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}
	if filter.Featured != nil {
		query = query.Where("featured = ?", *filter.Featured)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR excerpt ILIKE ? OR content_text ILIKE ?", search, search, search)
	}
	return query
}

func postFromModel(model PostModel, err error) (*domain.Post, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPostNotFound
	}
	if err != nil {
		return nil, err
	}
	return postFromModelNoError(model), nil
}

func postFromModelNoError(model PostModel) *domain.Post {
	return &domain.Post{ID: model.ID, AuthorID: model.AuthorID, CategoryID: model.CategoryID, Title: model.Title, Slug: model.Slug, Excerpt: model.Excerpt, ContentMarkdown: model.ContentMarkdown, ContentText: model.ContentText, CoverImage: model.CoverImage, Status: domain.PostStatus(model.Status), Featured: model.Featured, PinnedAt: model.PinnedAt, PublishedAt: model.PublishedAt, SEOTitle: model.SEOTitle, SEODescription: model.SEODescription, ViewCount: model.ViewCount, LikeCount: model.LikeCount, CommentCount: model.CommentCount, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toPostModel(post *domain.Post) *PostModel {
	return &PostModel{ID: post.ID, AuthorID: post.AuthorID, CategoryID: post.CategoryID, Title: post.Title, Slug: post.Slug, Excerpt: post.Excerpt, ContentMarkdown: post.ContentMarkdown, ContentText: post.ContentText, CoverImage: post.CoverImage, Status: string(post.Status), Featured: post.Featured, PinnedAt: post.PinnedAt, PublishedAt: post.PublishedAt, SEOTitle: post.SEOTitle, SEODescription: post.SEODescription, ViewCount: post.ViewCount, LikeCount: post.LikeCount, CommentCount: post.CommentCount, CreatedAt: post.CreatedAt, UpdatedAt: post.UpdatedAt}
}
