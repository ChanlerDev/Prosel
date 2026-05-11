package page

import "context"

type PageRepository interface {
	Create(ctx context.Context, page *Page) error
	Update(ctx context.Context, page *Page) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Page, error)
	GetBySlug(ctx context.Context, slug string, includeDraft bool) (*Page, error)
	List(ctx context.Context, filter PageFilter) ([]Page, Pagination, error)
	SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error)
	IncrementView(ctx context.Context, id string) error
}

type FriendRepository interface {
	Create(ctx context.Context, friend *Friend) error
	Update(ctx context.Context, friend *Friend) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Friend, error)
	URLExists(ctx context.Context, url string, excludeID *string) (bool, error)
	List(ctx context.Context, status string) ([]Friend, error)
}
