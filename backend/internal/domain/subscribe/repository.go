package subscribe

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, subscriber *Subscriber) error
	GetByEmail(ctx context.Context, email string) (*Subscriber, error)
	Verify(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
	List(ctx context.Context, filter SubscriberFilter) ([]Subscriber, Pagination, error)
	ListActive(ctx context.Context) ([]Subscriber, error)
	CreateDelivery(ctx context.Context, delivery *EmailDelivery) error
	UpdateDeliveryStatus(ctx context.Context, id string, status EmailDeliveryStatus, errorMessage string, sentAt *time.Time) error
}
