package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/search"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	"github.com/gin-gonic/gin"
)

type SearchService interface {
	Search(ctx context.Context, query string, filter domain.SearchFilter) ([]domain.SearchResult, domain.Pagination, error)
	RebuildIndex(ctx context.Context) error
	IndexStatus(ctx context.Context) (*domain.IndexStatus, error)
}

type SearchHandler struct{ service SearchService }

func NewSearchHandler(service SearchService) *SearchHandler { return &SearchHandler{service: service} }

func (h *SearchHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/search", h.search)
}

func (h *SearchHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/search", h.status)
	admin.POST("/search/rebuild", h.rebuild)
}

type searchResultResponse struct {
	RefType string  `json:"refType"`
	RefID   string  `json:"refId"`
	Title   string  `json:"title"`
	Slug    string  `json:"slug,omitempty"`
	Excerpt string  `json:"excerpt,omitempty"`
	Rank    float64 `json:"rank"`
}

type searchStatusResponse struct {
	Total     int64      `json:"total"`
	Posts     int64      `json:"posts"`
	Notes     int64      `json:"notes"`
	Pages     int64      `json:"pages"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (h *SearchHandler) search(c *gin.Context) {
	results, pagination, err := h.service.Search(c.Request.Context(), c.Query("q"), searchFilter(c))
	if err != nil {
		h.handleSearchError(c, err)
		return
	}
	response.OKWithMeta(c, toSearchResultResponses(results), pagination)
}

func (h *SearchHandler) status(c *gin.Context) {
	status, err := h.service.IndexStatus(c.Request.Context())
	if err != nil {
		h.handleSearchError(c, err)
		return
	}
	response.OK(c, toSearchStatusResponse(status))
}

func (h *SearchHandler) rebuild(c *gin.Context) {
	if err := h.service.RebuildIndex(c.Request.Context()); err != nil {
		h.handleSearchError(c, err)
		return
	}
	status, err := h.service.IndexStatus(c.Request.Context())
	if err != nil {
		h.handleSearchError(c, err)
		return
	}
	response.OK(c, toSearchStatusResponse(status))
}

func (h *SearchHandler) handleSearchError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrInvalidSearch) {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid search request", nil)
		return
	}
	response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Search request failed", nil)
}

func searchFilter(c *gin.Context) domain.SearchFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	return domain.SearchFilter{Type: domain.RefType(c.Query("type")), Page: page, PerPage: perPage}
}

func toSearchResultResponses(results []domain.SearchResult) []searchResultResponse {
	items := make([]searchResultResponse, 0, len(results))
	for _, result := range results {
		items = append(items, searchResultResponse{RefType: string(result.RefType), RefID: result.RefID, Title: result.Title, Slug: result.Slug, Excerpt: result.Excerpt, Rank: result.Rank})
	}
	return items
}

func toSearchStatusResponse(status *domain.IndexStatus) searchStatusResponse {
	return searchStatusResponse{Total: status.Total, Posts: status.Posts, Notes: status.Notes, Pages: status.Pages, UpdatedAt: status.UpdatedAt}
}
