package handler

import (
	"context"
	"encoding/xml"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/subscribe"
	usecase "github.com/chanler/prosel/backend/internal/usecase/subscribe"
	"github.com/gin-gonic/gin"
)

type fakeSubscribeService struct {
	subscriber *domain.Subscriber
	subs       []domain.Subscriber
	pagination domain.Pagination
	verify     string
	unsub      string
	notified   string
	filter     domain.SubscriberFilter
}

func (s *fakeSubscribeService) Subscribe(ctx context.Context, req usecase.SubscribeRequest) (*domain.Subscriber, error) {
	s.subscriber = &domain.Subscriber{ID: "sub-1", Email: req.Email, Name: req.Name, Status: domain.SubscriberPending, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return s.subscriber, nil
}
func (s *fakeSubscribeService) Verify(ctx context.Context, token string) error {
	s.verify = token
	return nil
}
func (s *fakeSubscribeService) Unsubscribe(ctx context.Context, token string) error {
	s.unsub = token
	return nil
}
func (s *fakeSubscribeService) ListSubscribers(ctx context.Context, filter domain.SubscriberFilter) ([]domain.Subscriber, domain.Pagination, error) {
	s.filter = filter
	return s.subs, s.pagination, nil
}
func (s *fakeSubscribeService) NotifyPostPublished(ctx context.Context, postID string) error {
	s.notified = postID
	return nil
}

type fakeFeedPostService struct {
	posts []postDomain.Post
}

func (s fakeFeedPostService) ListPublishedPosts(ctx context.Context, filter postDomain.PostListFilter) ([]postDomain.Post, postDomain.Pagination, error) {
	return s.posts, postDomain.Pagination{Page: 1, PerPage: 20, Total: int64(len(s.posts)), TotalPages: 1}, nil
}

func TestSubscribeHandlerFeedReturnsRSSXML(t *testing.T) {
	gin.SetMode(gin.TestMode)
	publishedAt := time.Date(2026, 5, 12, 0, 0, 0, 0, time.UTC)
	router := gin.New()
	NewSubscribeHandler(&fakeSubscribeService{}, fakeFeedPostService{posts: []postDomain.Post{{ID: "post-1", Title: "Hello", Slug: "hello", Excerpt: "Intro", Status: postDomain.PostPublished, PublishedAt: &publishedAt}}}, "https://example.com").RegisterFeedRoute(router)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/feed.xml", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var feed rssFeed
	if err := xml.Unmarshal(recorder.Body.Bytes(), &feed); err != nil {
		t.Fatalf("decode rss: %v", err)
	}
	if feed.Channel.Title != "Prosel" || len(feed.Channel.Items) != 1 || feed.Channel.Items[0].Link != "https://example.com/posts/hello" {
		t.Fatalf("feed = %#v", feed)
	}
}

func TestSubscribeHandlerVerifyUsesToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeSubscribeService{}
	router := gin.New()
	NewSubscribeHandler(service, fakeFeedPostService{}, "").RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/subscribe/verify?token=tok", nil))

	if recorder.Code != http.StatusOK || service.verify != "tok" {
		t.Fatalf("status=%d token=%q", recorder.Code, service.verify)
	}
}

func TestSubscribeHandlerAdminListFiltersStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeSubscribeService{subs: []domain.Subscriber{{ID: "sub-1", Email: "a@example.com", Status: domain.SubscriberActive}}, pagination: domain.Pagination{Page: 1, PerPage: 20, Total: 1, TotalPages: 1}}
	router := gin.New()
	NewSubscribeHandler(service, fakeFeedPostService{}, "").RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/admin/subscribers?status=active", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d", recorder.Code)
	}
	if service.filter.Status == nil || *service.filter.Status != domain.SubscriberActive {
		t.Fatalf("filter = %#v", service.filter.Status)
	}
}
