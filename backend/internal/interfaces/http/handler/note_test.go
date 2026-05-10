package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/note"
	usecase "github.com/chanler/prosel/backend/internal/usecase/note"
	"github.com/gin-gonic/gin"
)

type fakeNoteService struct {
	note       *domain.Note
	notes      []domain.Note
	pagination domain.Pagination
	filter     domain.NoteFilter
	created    usecase.NoteRequest
	updated    usecase.NoteRequest
	pinID      string
	pinned     bool
}

func (s *fakeNoteService) CreateNote(ctx context.Context, req usecase.NoteRequest) (*domain.Note, error) {
	s.created = req
	return s.note, nil
}
func (s *fakeNoteService) UpdateNote(ctx context.Context, id string, req usecase.NoteRequest) (*domain.Note, error) {
	s.updated = req
	return s.note, nil
}
func (s *fakeNoteService) DeleteNote(ctx context.Context, id string) error { return nil }
func (s *fakeNoteService) GetPublicNote(ctx context.Context, slug string) (*domain.Note, error) {
	return s.note, nil
}
func (s *fakeNoteService) GetAdminNote(ctx context.Context, id string) (*domain.Note, error) {
	return s.note, nil
}
func (s *fakeNoteService) ListPublicNotes(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	s.filter = filter
	return s.notes, s.pagination, nil
}
func (s *fakeNoteService) ListAdminNotes(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	s.filter = filter
	return s.notes, s.pagination, nil
}
func (s *fakeNoteService) PinNote(ctx context.Context, id string, pinned bool) error {
	s.pinID = id
	s.pinned = pinned
	return nil
}

func TestNoteHandlerListPublicReturnsMeta(t *testing.T) {
	gin.SetMode(gin.TestMode)
	publishedAt := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	service := &fakeNoteService{notes: []domain.Note{{ID: "note-1", Title: "Hello", Slug: "hello", Status: domain.NotePublished, PublishedAt: &publishedAt}}, pagination: domain.Pagination{Page: 2, PerPage: 10, Total: 11, TotalPages: 2}}
	router := gin.New()
	NewNoteHandler(service).RegisterPublicRoutes(router.Group("/api/v1"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/notes?page=2&perPage=10&search=hello", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var body struct {
		Data []noteResponse    `json:"data"`
		Meta domain.Pagination `json:"meta"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data) != 1 || body.Data[0].Slug != "hello" || body.Meta.TotalPages != 2 {
		t.Fatalf("body = %#v", body)
	}
	if service.filter.Search != "hello" {
		t.Fatalf("filter = %#v", service.filter)
	}
}

func TestNoteHandlerCreatePassesCurrentUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now().UTC()
	service := &fakeNoteService{note: &domain.Note{ID: "note-1", AuthorID: "user-1", Slug: "hello", ContentMarkdown: "body", Status: domain.NotePublished, CreatedAt: now, UpdatedAt: now}}
	router := gin.New()
	group := router.Group("/api/v1/admin")
	group.Use(func(c *gin.Context) { c.Set("userID", "user-1"); c.Next() })
	NewNoteHandler(service).RegisterProtectedRoutes(group)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/notes", strings.NewReader(`{"title":"Hello","contentMarkdown":"body","mood":"calm"}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body=%s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if service.created.AuthorID != "user-1" || service.created.Mood != "calm" {
		t.Fatalf("created = %#v", service.created)
	}
}

func TestNoteHandlerPinParsesPinnedFlag(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := &fakeNoteService{}
	router := gin.New()
	NewNoteHandler(service).RegisterProtectedRoutes(router.Group("/api/v1/admin"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/v1/admin/notes/note-1/pin", strings.NewReader(`{"pinned":true}`))
	request.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if service.pinID != "note-1" || !service.pinned {
		t.Fatalf("pin = %q/%v", service.pinID, service.pinned)
	}
}
