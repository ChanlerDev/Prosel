package subscribe

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/mail"
	"strings"
	"time"

	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/subscribe"
)

type PostReader interface {
	GetAdminPost(ctx context.Context, id string) (*postDomain.Post, error)
}

type SubscribeUsecase struct {
	subscribers domain.Repository
	mailer      domain.Mailer
	posts       PostReader
	siteURL     string
}

type Options struct {
	SiteURL string
}

type SubscribeRequest struct {
	Email string
	Name  string
}

func NewSubscribeUsecase(subscribers domain.Repository, mailer domain.Mailer, posts PostReader, opts Options) *SubscribeUsecase {
	return &SubscribeUsecase{subscribers: subscribers, mailer: mailer, posts: posts, siteURL: strings.TrimRight(strings.TrimSpace(opts.SiteURL), "/")}
}

func (uc *SubscribeUsecase) Subscribe(ctx context.Context, req SubscribeRequest) (*domain.Subscriber, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))
	name := strings.TrimSpace(req.Name)
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, domain.ErrInvalidSubscriber
	}

	if existing, err := uc.subscribers.GetByEmail(ctx, email); err == nil {
		if existing.Status == domain.SubscriberUnsubscribed {
			return nil, domain.ErrSubscriberExists
		}
		_ = uc.sendVerification(ctx, existing)
		return existing, nil
	} else if err != domain.ErrSubscriberNotFound {
		return nil, err
	}

	now := time.Now().UTC()
	subscriber := &domain.Subscriber{ID: newID(), Email: email, Name: name, Status: domain.SubscriberPending, VerifyToken: newToken(), UnsubscribeToken: newToken(), CreatedAt: now, UpdatedAt: now}
	if err := uc.subscribers.Create(ctx, subscriber); err != nil {
		return nil, err
	}
	if err := uc.sendVerification(ctx, subscriber); err != nil {
		return nil, err
	}
	return subscriber, nil
}

func (uc *SubscribeUsecase) Verify(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return domain.ErrInvalidSubscriber
	}
	return uc.subscribers.Verify(ctx, token)
}

func (uc *SubscribeUsecase) Unsubscribe(ctx context.Context, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return domain.ErrInvalidSubscriber
	}
	return uc.subscribers.Unsubscribe(ctx, token)
}

func (uc *SubscribeUsecase) ListSubscribers(ctx context.Context, filter domain.SubscriberFilter) ([]domain.Subscriber, domain.Pagination, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	if filter.Status != nil && !filter.Status.Valid() {
		filter.Status = nil
	}
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.subscribers.List(ctx, filter)
}

func (uc *SubscribeUsecase) NotifyPostPublished(ctx context.Context, postID string) error {
	postID = strings.TrimSpace(postID)
	if postID == "" || uc.posts == nil {
		return domain.ErrInvalidSubscriber
	}
	post, err := uc.posts.GetAdminPost(ctx, postID)
	if err != nil {
		return err
	}
	if post.Status != postDomain.PostPublished {
		return domain.ErrInvalidSubscriber
	}
	subscribers, err := uc.subscribers.ListActive(ctx)
	if err != nil {
		return err
	}
	for _, subscriber := range subscribers {
		subscriberCopy := subscriber
		if err := uc.notifySubscriber(ctx, post, &subscriberCopy); err != nil {
			return err
		}
	}
	return nil
}

func (uc *SubscribeUsecase) sendVerification(ctx context.Context, subscriber *domain.Subscriber) error {
	if uc.mailer == nil {
		return nil
	}
	link := uc.publicURL("/subscribe/verify?token=" + subscriber.VerifyToken)
	body := "Confirm your Prosel subscription:\n\n" + link
	return uc.mailer.Send(ctx, domain.MailMessage{To: subscriber.Email, Subject: "Confirm your Prosel subscription", Body: body})
}

func (uc *SubscribeUsecase) notifySubscriber(ctx context.Context, post *postDomain.Post, subscriber *domain.Subscriber) error {
	subID := subscriber.ID
	delivery := &domain.EmailDelivery{ID: newID(), SubscriberID: &subID, Subject: "New post: " + post.Title, RefType: "post", RefID: post.ID, Status: domain.EmailDeliveryPending, CreatedAt: time.Now().UTC()}
	if err := uc.subscribers.CreateDelivery(ctx, delivery); err != nil {
		return err
	}
	body := post.Title + "\n\n" + strings.TrimSpace(post.Excerpt) + "\n\n" + uc.publicURL("/posts/"+post.Slug) + "\n\nUnsubscribe: " + uc.publicURL("/subscribe/unsubscribe?token="+subscriber.UnsubscribeToken)
	err := uc.mailer.Send(ctx, domain.MailMessage{To: subscriber.Email, Subject: delivery.Subject, Body: body})
	if err != nil {
		_ = uc.subscribers.UpdateDeliveryStatus(ctx, delivery.ID, domain.EmailDeliveryFailed, err.Error(), nil)
		return err
	}
	sentAt := time.Now().UTC()
	return uc.subscribers.UpdateDeliveryStatus(ctx, delivery.ID, domain.EmailDeliverySent, "", &sentAt)
}

func (uc *SubscribeUsecase) publicURL(path string) string {
	if uc.siteURL == "" {
		return path
	}
	return uc.siteURL + path
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))[:32]
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:])
}

func newToken() string {
	var bytes [32]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(bytes[:])
}
