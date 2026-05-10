package note

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, note *Note) error
	Update(ctx context.Context, note *Note) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Note, error)
	GetBySlug(ctx context.Context, slug string) (*Note, error)
	ListPublic(ctx context.Context, filter NoteFilter) ([]Note, Pagination, error)
	ListAdmin(ctx context.Context, filter NoteFilter) ([]Note, Pagination, error)
	SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error)
	IncrementView(ctx context.Context, id string) error
	SetPinned(ctx context.Context, id string, pinnedAt *time.Time) error
}
