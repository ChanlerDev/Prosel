package search

import "context"

type Repository interface {
	UpsertDocument(ctx context.Context, doc *SearchDocument) error
	DeleteDocument(ctx context.Context, refType string, refID string) error
	Search(ctx context.Context, query string, filter SearchFilter) ([]SearchResult, Pagination, error)
	Rebuild(ctx context.Context) error
	Status(ctx context.Context) (*IndexStatus, error)
}
