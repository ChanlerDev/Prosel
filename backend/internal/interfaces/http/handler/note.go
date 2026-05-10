package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/note"
	"github.com/chanler/prosel/backend/internal/interfaces/http/middleware"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/note"
	"github.com/gin-gonic/gin"
)

type NoteService interface {
	CreateNote(ctx context.Context, req usecase.NoteRequest) (*domain.Note, error)
	UpdateNote(ctx context.Context, id string, req usecase.NoteRequest) (*domain.Note, error)
	DeleteNote(ctx context.Context, id string) error
	GetPublicNote(ctx context.Context, slug string) (*domain.Note, error)
	GetAdminNote(ctx context.Context, id string) (*domain.Note, error)
	ListPublicNotes(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error)
	ListAdminNotes(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error)
	PinNote(ctx context.Context, id string, pinned bool) error
}

type NoteHandler struct{ service NoteService }

func NewNoteHandler(service NoteService) *NoteHandler { return &NoteHandler{service: service} }

func (h *NoteHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/notes", h.listPublic)
	group.GET("/notes/:slug", h.getPublic)
}

func (h *NoteHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/notes", h.listAdmin)
	admin.GET("/notes/:id", h.getAdmin)
	admin.POST("/notes", h.create)
	admin.PATCH("/notes/:id", h.update)
	admin.POST("/notes/:id/pin", h.pin)
	admin.DELETE("/notes/:id", h.delete)
}

type noteRequest struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	ContentMarkdown string `json:"contentMarkdown" binding:"required"`
	Mood            string `json:"mood"`
	Weather         string `json:"weather"`
	Location        string `json:"location"`
	Status          string `json:"status"`
}

type notePinRequest struct {
	Pinned bool `json:"pinned"`
}

type noteResponse struct {
	ID              string     `json:"id"`
	AuthorID        string     `json:"authorId"`
	Title           string     `json:"title,omitempty"`
	Slug            string     `json:"slug"`
	ContentMarkdown string     `json:"contentMarkdown,omitempty"`
	ContentText     string     `json:"contentText"`
	Mood            string     `json:"mood,omitempty"`
	Weather         string     `json:"weather,omitempty"`
	Location        string     `json:"location,omitempty"`
	Status          string     `json:"status"`
	PinnedAt        *time.Time `json:"pinnedAt,omitempty"`
	PublishedAt     *time.Time `json:"publishedAt,omitempty"`
	ViewCount       int64      `json:"viewCount"`
	LikeCount       int64      `json:"likeCount"`
	CommentCount    int64      `json:"commentCount"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (h *NoteHandler) listPublic(c *gin.Context) {
	notes, pagination, err := h.service.ListPublicNotes(c.Request.Context(), noteFilter(c))
	if err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OKWithMeta(c, toNoteResponses(notes), pagination)
}

func (h *NoteHandler) getPublic(c *gin.Context) {
	note, err := h.service.GetPublicNote(c.Request.Context(), c.Param("slug"))
	if err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OK(c, toNoteResponse(note))
}

func (h *NoteHandler) listAdmin(c *gin.Context) {
	notes, pagination, err := h.service.ListAdminNotes(c.Request.Context(), noteFilter(c))
	if err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OKWithMeta(c, toNoteResponses(notes), pagination)
}

func (h *NoteHandler) getAdmin(c *gin.Context) {
	note, err := h.service.GetAdminNote(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OK(c, toNoteResponse(note))
}

func (h *NoteHandler) create(c *gin.Context) {
	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid note request", nil)
		return
	}
	note, err := h.service.CreateNote(c.Request.Context(), usecase.NoteRequest{AuthorID: middleware.CurrentUserID(c), Title: req.Title, Slug: req.Slug, ContentMarkdown: req.ContentMarkdown, Mood: req.Mood, Weather: req.Weather, Location: req.Location, Status: req.Status})
	if err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OK(c, toNoteResponse(note))
}

func (h *NoteHandler) update(c *gin.Context) {
	var req noteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid note request", nil)
		return
	}
	note, err := h.service.UpdateNote(c.Request.Context(), c.Param("id"), usecase.NoteRequest{Title: req.Title, Slug: req.Slug, ContentMarkdown: req.ContentMarkdown, Mood: req.Mood, Weather: req.Weather, Location: req.Location, Status: req.Status})
	if err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OK(c, toNoteResponse(note))
}

func (h *NoteHandler) pin(c *gin.Context) {
	var req notePinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid pin request", nil)
		return
	}
	if err := h.service.PinNote(c.Request.Context(), c.Param("id"), req.Pinned); err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *NoteHandler) delete(c *gin.Context) {
	if err := h.service.DeleteNote(c.Request.Context(), c.Param("id")); err != nil {
		h.handleNoteError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *NoteHandler) handleNoteError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrNoteNotFound):
		response.Error(c, http.StatusNotFound, "NOTE_NOT_FOUND", "Note not found", nil)
	case errors.Is(err, domain.ErrSlugTaken):
		response.Error(c, http.StatusConflict, "SLUG_TAKEN", "Note slug already exists", nil)
	case errors.Is(err, domain.ErrInvalidNote):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid note request", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Note request failed", nil)
	}
}

func noteFilter(c *gin.Context) domain.NoteFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	filter := domain.NoteFilter{Page: page, PerPage: perPage, Search: c.Query("search")}
	if status := c.Query("status"); status != "" {
		noteStatus := domain.NoteStatus(status)
		if noteStatus.Valid() {
			filter.Status = &noteStatus
		}
	}
	return filter
}

func toNoteResponses(notes []domain.Note) []noteResponse {
	result := make([]noteResponse, 0, len(notes))
	for _, note := range notes {
		noteCopy := note
		result = append(result, *toNoteResponse(&noteCopy))
	}
	return result
}

func toNoteResponse(note *domain.Note) *noteResponse {
	return &noteResponse{ID: note.ID, AuthorID: note.AuthorID, Title: note.Title, Slug: note.Slug, ContentMarkdown: note.ContentMarkdown, ContentText: note.ContentText, Mood: note.Mood, Weather: note.Weather, Location: note.Location, Status: string(note.Status), PinnedAt: note.PinnedAt, PublishedAt: note.PublishedAt, ViewCount: note.ViewCount, LikeCount: note.LikeCount, CommentCount: note.CommentCount, CreatedAt: note.CreatedAt, UpdatedAt: note.UpdatedAt}
}
