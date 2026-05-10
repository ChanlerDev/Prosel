package comment

import "context"

type Repository interface {
	Create(ctx context.Context, comment *Comment) error
	UpdateStatus(ctx context.Context, id string, status CommentStatus) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*Comment, error)
	ListByRef(ctx context.Context, refType RefType, refID string, onlyApproved bool) ([]Comment, error)
	ListAdmin(ctx context.Context, filter CommentFilter) ([]Comment, Pagination, error)
	IncrementReplyCount(ctx context.Context, id string) error
}
