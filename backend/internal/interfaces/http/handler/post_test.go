package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/post"
	usecase "github.com/chanler/prosel/backend/internal/usecase/post"
	"github.com/gin-gonic/gin"
)

type fakePostService struct {
	post       *domain.Post
	posts      []domain.Post
	pagination domain.Pagination
	filter     domain.PostListFilter
}

func (s *fakePostService) CreatePost(ctx context.Context, req usecase.CreatePostRequest) (*domain.Post, error) {
	s.post = &domain.Post{ID: "post-1", AuthorID: req.AuthorID, Title: req.Title, Slug: req.Slug, Status: domain.PostDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	return s.post, nil
}
func (s *fakePostService) UpdatePost(ctx context.Context, id string, req usecase.UpdatePostRequest) (*domain.Post, error) {
	return s.post, nil
}
func (s *fakePostService) PublishPost(ctx context.Context, id string) (*domain.Post, error) {
	return s.post, nil
}
func (s *fakePostService) UnpublishPost(ctx context.Context, id string) (*domain.Post, error) {
	return s.post, nil
}
func (s *fakePostService) DeletePost(ctx context.Context, id string) error { return nil }
func (s *fakePostService) GetAdminPost(ctx context.Context, id string) (*domain.Post, error) {
	return s.post, nil
}
func (s *fakePostService) GetPublishedPost(ctx context.Context, slug string) (*domain.Post, error) {
	return s.post, nil
}
func (s *fakePostService) ListPublishedPosts(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error) {
	s.filter = filter
	return s.posts, s.pagination, nil
}
func (s *fakePostService) ListAdminPosts(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error) {
	s.filter = filter
	return s.posts, s.pagination, nil
}

func TestPostHandlerListPublicReturnsMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	publishedAt := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	service := &fakePostService{posts: []domain.Post{{ID: "post-1", Title: "Hello", Slug: "hello", Status: domain.PostPublished, PublishedAt: &publishedAt}}, pagination: domain.Pagination{Page: 2, PerPage: 10, Total: 11, TotalPages: 2}}
	router := gin.New()
	NewPostHandler(service).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/posts?page=2&perPage=10&featured=true", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body struct {
		Data []postResponse    `json:"data"`
		Meta domain.Pagination `json:"meta"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Slug != "hello" {
		t.Fatalf("data = %#v", body.Data)
	}
	if body.Meta.Page != 2 || body.Meta.TotalPages != 2 {
		t.Fatalf("meta = %#v", body.Meta)
	}
	if service.filter.Featured == nil || !*service.filter.Featured {
		t.Fatalf("featured filter = %#v", service.filter.Featured)
	}
}
