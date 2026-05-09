package post

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, post *Post) error
	Update(ctx context.Context, post *Post) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Post, error)
	GetBySlug(ctx context.Context, slug string, includeDraft bool) (*Post, error)
	List(ctx context.Context, filter PostListFilter) ([]Post, Pagination, error)
	SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error)
	IncrementView(ctx context.Context, id string) error
	SetStatus(ctx context.Context, id string, status PostStatus, publishedAt *time.Time) error
	ReplaceTags(ctx context.Context, postID string, tagIDs []string) error
}
