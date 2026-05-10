package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/comment"
	usecase "github.com/chanler/prosel/backend/internal/usecase/comment"
	"github.com/gin-gonic/gin"
)

type fakeCommentService struct {
	comment    *domain.Comment
	nodes      []domain.CommentNode
	comments   []domain.Comment
	pagination domain.Pagination
	filter     domain.CommentFilter
	submitted  usecase.SubmitCommentRequest
	meta       usecase.ClientMeta
	statusID   string
	status     string
}

func (s *fakeCommentService) SubmitComment(ctx context.Context, req usecase.SubmitCommentRequest, meta usecase.ClientMeta) (*domain.Comment, error) {
	s.submitted = req
	s.meta = meta
	return s.comment, nil
}
func (s *fakeCommentService) ReplyAsAdmin(ctx context.Context, req usecase.AdminReplyRequest) (*domain.Comment, error) {
	return s.comment, nil
}
func (s *fakeCommentService) ListPublicComments(ctx context.Context, refType string, refID string) ([]domain.CommentNode, error) {
	return s.nodes, nil
}
func (s *fakeCommentService) ModerateComment(ctx context.Context, id string, status string) error {
	s.statusID = id
	s.status = status
	return nil
}
func (s *fakeCommentService) DeleteComment(ctx context.Context, id string) error { return nil }
func (s *fakeCommentService) ListAdminComments(ctx context.Context, filter domain.CommentFilter) ([]domain.Comment, domain.Pagination, error) {
	s.filter = filter
	return s.comments, s.pagination, nil
}

func TestCommentHandlerSubmitUsesClientMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakeCommentService{comment: &domain.Comment{ID: "comment-1", RefType: domain.RefPost, RefID: "post-1", AuthorName: "Ada", Content: "Hello", Status: domain.CommentPending, CreatedAt: now, UpdatedAt: now}}
	router := gin.New()
	NewCommentHandler(service).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/comments", strings.NewReader(`{"refType":"post","refId":"post-1","authorName":"Ada","authorEmail":"ada@example.com","content":"Hello"}`))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("User-Agent", "comment-test")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if service.submitted.RefType != "post" || service.meta.UserAgent != "comment-test" {
		t.Fatalf("submitted = %#v meta = %#v", service.submitted, service.meta)
	}
}

func TestCommentHandlerListAdminParsesFiltersAndReturnsMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakeCommentService{comments: []domain.Comment{{ID: "comment-1", RefType: domain.RefPost, RefID: "post-1", AuthorName: "Ada", Content: "Hello", Status: domain.CommentPending, CreatedAt: now, UpdatedAt: now}}, pagination: domain.Pagination{Page: 2, PerPage: 10, Total: 11, TotalPages: 2}}
	router := gin.New()
	NewCommentHandler(service).RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/admin/comments?page=2&perPage=10&status=pending&refType=post&refId=post-1&search=Ada", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body struct {
		Data []commentResponse `json:"data"`
		Meta domain.Pagination `json:"meta"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Meta.TotalPages != 2 {
		t.Fatalf("body = %#v", body)
	}
	if service.filter.Status == nil || *service.filter.Status != domain.CommentPending || service.filter.RefType != domain.RefPost || service.filter.RefID != "post-1" {
		t.Fatalf("filter = %#v", service.filter)
	}
}
