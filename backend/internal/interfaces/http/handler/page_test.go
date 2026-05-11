package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/page"
	usecase "github.com/chanler/prosel/backend/internal/usecase/page"
	"github.com/gin-gonic/gin"
)

type fakePageService struct {
	page       *domain.Page
	pages      []domain.Page
	pagination domain.Pagination
	friend     *domain.Friend
	friends    []domain.Friend
	filter     domain.PageFilter
	created    usecase.PageRequest
	friendReq  usecase.FriendRequest
	status     string
}

func (s *fakePageService) CreatePage(ctx context.Context, req usecase.PageRequest) (*domain.Page, error) {
	s.created = req
	return s.page, nil
}
func (s *fakePageService) UpdatePage(ctx context.Context, id string, req usecase.PageRequest) (*domain.Page, error) {
	s.created = req
	return s.page, nil
}
func (s *fakePageService) DeletePage(ctx context.Context, id string) error { return nil }
func (s *fakePageService) GetPublicPage(ctx context.Context, slug string) (*domain.Page, error) {
	return s.page, nil
}
func (s *fakePageService) GetAdminPage(ctx context.Context, id string) (*domain.Page, error) {
	return s.page, nil
}
func (s *fakePageService) ListPublicPages(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error) {
	s.filter = filter
	return s.pages, s.pagination, nil
}
func (s *fakePageService) ListAdminPages(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error) {
	s.filter = filter
	return s.pages, s.pagination, nil
}
func (s *fakePageService) CreateFriend(ctx context.Context, req usecase.FriendRequest) (*domain.Friend, error) {
	s.friendReq = req
	return s.friend, nil
}
func (s *fakePageService) UpdateFriend(ctx context.Context, id string, req usecase.FriendRequest) (*domain.Friend, error) {
	s.friendReq = req
	return s.friend, nil
}
func (s *fakePageService) DeleteFriend(ctx context.Context, id string) error { return nil }
func (s *fakePageService) ListFriends(ctx context.Context) ([]domain.Friend, error) {
	return s.friends, nil
}
func (s *fakePageService) ListAdminFriends(ctx context.Context, status string) ([]domain.Friend, error) {
	s.status = status
	return s.friends, nil
}

func TestPageHandlerListPublicReturnsMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakePageService{pages: []domain.Page{{ID: "page-1", Title: "About", Slug: "about", Status: domain.PagePublished, CreatedAt: now, UpdatedAt: now}}, pagination: domain.Pagination{Page: 2, PerPage: 10, Total: 11, TotalPages: 2}}
	router := gin.New()
	NewPageHandler(service).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/pages?page=2&perPage=10&search=about", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body struct {
		Data []pageResponse    `json:"data"`
		Meta domain.Pagination `json:"meta"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Slug != "about" || body.Meta.TotalPages != 2 {
		t.Fatalf("body = %#v", body)
	}
	if service.filter.Search != "about" {
		t.Fatalf("filter = %#v", service.filter)
	}
}

func TestPageHandlerCreatePassesCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakePageService{page: &domain.Page{ID: "page-1", AuthorID: "user-1", Title: "About", Slug: "about", Status: domain.PagePublished, Template: domain.TemplateAbout, CreatedAt: now, UpdatedAt: now}}
	router := gin.New()
	group := router.Group("/api/v1/admin")
	group.Use(func(c *gin.Context) { c.Set("userID", "user-1"); c.Next() })
	NewPageHandler(service).RegisterProtectedRoutes(group)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/pages", strings.NewReader(`{"title":"About","contentMarkdown":"body","template":"about"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if service.created.AuthorID != "user-1" || service.created.Template != "about" {
		t.Fatalf("created = %#v", service.created)
	}
}

func TestPageHandlerListFriendsPublic(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakePageService{friends: []domain.Friend{{ID: "friend-1", Name: "Ada", URL: "https://example.com", Status: domain.FriendActive, CreatedAt: now, UpdatedAt: now}}}
	router := gin.New()
	NewPageHandler(service).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/friends", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body struct {
		Data []friendResponse `json:"data"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Name != "Ada" {
		t.Fatalf("body = %#v", body)
	}
}

func TestPageHandlerCreateFriend(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakePageService{friend: &domain.Friend{ID: "friend-1", Name: "Ada", URL: "https://example.com", Status: domain.FriendActive, CreatedAt: now, UpdatedAt: now}}
	router := gin.New()
	NewPageHandler(service).RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/friends", strings.NewReader(`{"name":"Ada","url":"https://example.com"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if service.friendReq.Name != "Ada" || service.friendReq.URL != "https://example.com" {
		t.Fatalf("friend req = %#v", service.friendReq)
	}
}
