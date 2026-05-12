package subscribe

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/subscribe"
)

type fakeSubscriberRepo struct {
	subscriber *domain.Subscriber
	subs       []domain.Subscriber
	deliveries []domain.EmailDelivery
	pagination domain.Pagination
	err        error
	verified   string
	unsub      string
	listed     domain.SubscriberFilter
}

func (r *fakeSubscriberRepo) Create(ctx context.Context, subscriber *domain.Subscriber) error {
	r.subscriber = subscriber
	return r.err
}
func (r *fakeSubscriberRepo) GetByEmail(ctx context.Context, email string) (*domain.Subscriber, error) {
	if r.subscriber == nil {
		return nil, domain.ErrSubscriberNotFound
	}
	return r.subscriber, r.err
}
func (r *fakeSubscriberRepo) Verify(ctx context.Context, token string) error {
	r.verified = token
	return r.err
}
func (r *fakeSubscriberRepo) Unsubscribe(ctx context.Context, token string) error {
	r.unsub = token
	return r.err
}
func (r *fakeSubscriberRepo) List(ctx context.Context, filter domain.SubscriberFilter) ([]domain.Subscriber, domain.Pagination, error) {
	r.listed = filter
	return r.subs, r.pagination, r.err
}
func (r *fakeSubscriberRepo) ListActive(ctx context.Context) ([]domain.Subscriber, error) {
	return r.subs, r.err
}
func (r *fakeSubscriberRepo) CreateDelivery(ctx context.Context, delivery *domain.EmailDelivery) error {
	r.deliveries = append(r.deliveries, *delivery)
	return r.err
}
func (r *fakeSubscriberRepo) UpdateDeliveryStatus(ctx context.Context, id string, status domain.EmailDeliveryStatus, errorMessage string, sentAt *time.Time) error {
	for i := range r.deliveries {
		if r.deliveries[i].ID == id {
			r.deliveries[i].Status = status
			r.deliveries[i].ErrorMessage = errorMessage
			r.deliveries[i].SentAt = sentAt
		}
	}
	return r.err
}

type fakeMailer struct {
	messages []domain.MailMessage
	err      error
}

func (m *fakeMailer) Send(ctx context.Context, message domain.MailMessage) error {
	m.messages = append(m.messages, message)
	return m.err
}

type fakePostReader struct {
	post *postDomain.Post
	err  error
}

func (r fakePostReader) GetAdminPost(ctx context.Context, id string) (*postDomain.Post, error) {
	return r.post, r.err
}

func TestSubscribeCreatesPendingSubscriberAndSendsVerificationEmail(t *testing.T) {
	repo := &fakeSubscriberRepo{}
	mailer := &fakeMailer{}
	uc := NewSubscribeUsecase(repo, mailer, nil, Options{SiteURL: "https://example.com"})

	subscriber, err := uc.Subscribe(context.Background(), SubscribeRequest{Email: " USER@Example.COM ", Name: " Chanler "})
	if err != nil {
		t.Fatalf("Subscribe() error = %v", err)
	}
	if subscriber.Email != "user@example.com" || subscriber.Name != "Chanler" || subscriber.Status != domain.SubscriberPending {
		t.Fatalf("subscriber = %#v", subscriber)
	}
	if subscriber.VerifyToken == "" || subscriber.UnsubscribeToken == "" {
		t.Fatalf("tokens were not generated: %#v", subscriber)
	}
	if len(mailer.messages) != 1 {
		t.Fatalf("messages = %d, want 1", len(mailer.messages))
	}
	if !strings.Contains(mailer.messages[0].Body, "/subscribe/verify?token="+subscriber.VerifyToken) {
		t.Fatalf("verification body = %q", mailer.messages[0].Body)
	}
}

func TestSubscribeRejectsInvalidEmail(t *testing.T) {
	uc := NewSubscribeUsecase(&fakeSubscriberRepo{}, &fakeMailer{}, nil, Options{})

	_, err := uc.Subscribe(context.Background(), SubscribeRequest{Email: "not-email"})
	if !errors.Is(err, domain.ErrInvalidSubscriber) {
		t.Fatalf("Subscribe() error = %v, want %v", err, domain.ErrInvalidSubscriber)
	}
}

func TestVerifyAndUnsubscribeRequireToken(t *testing.T) {
	repo := &fakeSubscriberRepo{}
	uc := NewSubscribeUsecase(repo, &fakeMailer{}, nil, Options{})

	if err := uc.Verify(context.Background(), " token-1 "); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if repo.verified != "token-1" {
		t.Fatalf("verified token = %q", repo.verified)
	}
	if err := uc.Unsubscribe(context.Background(), " token-2 "); err != nil {
		t.Fatalf("Unsubscribe() error = %v", err)
	}
	if repo.unsub != "token-2" {
		t.Fatalf("unsubscribe token = %q", repo.unsub)
	}
	if err := uc.Verify(context.Background(), " "); !errors.Is(err, domain.ErrInvalidSubscriber) {
		t.Fatalf("Verify(empty) error = %v", err)
	}
}

func TestNotifyPostPublishedSendsToActiveSubscribersAndRecordsDeliveries(t *testing.T) {
	publishedAt := time.Date(2026, 5, 12, 1, 2, 3, 0, time.UTC)
	repo := &fakeSubscriberRepo{subs: []domain.Subscriber{{ID: "sub-1", Email: "one@example.com", Status: domain.SubscriberActive, UnsubscribeToken: "unsub-1"}, {ID: "sub-2", Email: "two@example.com", Status: domain.SubscriberActive, UnsubscribeToken: "unsub-2"}}}
	mailer := &fakeMailer{}
	postReader := fakePostReader{post: &postDomain.Post{ID: "post-1", Title: "New Post", Slug: "new-post", Excerpt: "Intro", Status: postDomain.PostPublished, PublishedAt: &publishedAt}}
	uc := NewSubscribeUsecase(repo, mailer, postReader, Options{SiteURL: "https://example.com"})

	if err := uc.NotifyPostPublished(context.Background(), "post-1"); err != nil {
		t.Fatalf("NotifyPostPublished() error = %v", err)
	}
	if len(mailer.messages) != 2 || len(repo.deliveries) != 2 {
		t.Fatalf("messages=%d deliveries=%d, want 2/2", len(mailer.messages), len(repo.deliveries))
	}
	if repo.deliveries[0].Status != domain.EmailDeliverySent || repo.deliveries[0].SentAt == nil {
		t.Fatalf("delivery = %#v", repo.deliveries[0])
	}
	if !strings.Contains(mailer.messages[0].Body, "https://example.com/posts/new-post") || !strings.Contains(mailer.messages[0].Body, "unsubscribe?token=unsub-1") {
		t.Fatalf("message body = %q", mailer.messages[0].Body)
	}
}

func TestNotifyPostPublishedRejectsDraftPost(t *testing.T) {
	uc := NewSubscribeUsecase(&fakeSubscriberRepo{}, &fakeMailer{}, fakePostReader{post: &postDomain.Post{ID: "post-1", Status: postDomain.PostDraft}}, Options{})

	err := uc.NotifyPostPublished(context.Background(), "post-1")
	if !errors.Is(err, domain.ErrInvalidSubscriber) {
		t.Fatalf("NotifyPostPublished() error = %v, want %v", err, domain.ErrInvalidSubscriber)
	}
}
