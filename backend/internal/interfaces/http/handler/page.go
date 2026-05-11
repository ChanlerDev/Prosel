package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/page"
	"github.com/chanler/prosel/backend/internal/interfaces/http/middleware"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/page"
	"github.com/gin-gonic/gin"
)

type PageService interface {
	CreatePage(ctx context.Context, req usecase.PageRequest) (*domain.Page, error)
	UpdatePage(ctx context.Context, id string, req usecase.PageRequest) (*domain.Page, error)
	DeletePage(ctx context.Context, id string) error
	GetPublicPage(ctx context.Context, slug string) (*domain.Page, error)
	GetAdminPage(ctx context.Context, id string) (*domain.Page, error)
	ListPublicPages(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error)
	ListAdminPages(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error)
	CreateFriend(ctx context.Context, req usecase.FriendRequest) (*domain.Friend, error)
	UpdateFriend(ctx context.Context, id string, req usecase.FriendRequest) (*domain.Friend, error)
	DeleteFriend(ctx context.Context, id string) error
	ListFriends(ctx context.Context) ([]domain.Friend, error)
	ListAdminFriends(ctx context.Context, status string) ([]domain.Friend, error)
}

type PageHandler struct{ service PageService }

func NewPageHandler(service PageService) *PageHandler { return &PageHandler{service: service} }

func (h *PageHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/pages", h.listPublicPages)
	group.GET("/pages/:slug", h.getPublicPage)
	group.GET("/friends", h.listFriends)
}

func (h *PageHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/pages", h.listAdminPages)
	admin.GET("/pages/:id", h.getAdminPage)
	admin.POST("/pages", h.createPage)
	admin.PATCH("/pages/:id", h.updatePage)
	admin.DELETE("/pages/:id", h.deletePage)
	admin.GET("/friends", h.listAdminFriends)
	admin.POST("/friends", h.createFriend)
	admin.PATCH("/friends/:id", h.updateFriend)
	admin.DELETE("/friends/:id", h.deleteFriend)
}

type pageRequest struct {
	Title           string `json:"title" binding:"required"`
	Slug            string `json:"slug"`
	Subtitle        string `json:"subtitle"`
	ContentMarkdown string `json:"contentMarkdown" binding:"required"`
	Template        string `json:"template"`
	Status          string `json:"status"`
	SortOrder       int    `json:"sortOrder"`
	SEOTitle        string `json:"seoTitle"`
	SEODescription  string `json:"seoDescription"`
}

type friendRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url" binding:"required"`
	AvatarURL   string `json:"avatarUrl"`
	Description string `json:"description"`
	Status      string `json:"status"`
	SortOrder   int    `json:"sortOrder"`
}

type pageResponse struct {
	ID              string    `json:"id"`
	AuthorID        string    `json:"authorId"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	Subtitle        string    `json:"subtitle,omitempty"`
	ContentMarkdown string    `json:"contentMarkdown,omitempty"`
	ContentText     string    `json:"contentText"`
	Template        string    `json:"template"`
	Status          string    `json:"status"`
	SortOrder       int       `json:"sortOrder"`
	SEOTitle        string    `json:"seoTitle,omitempty"`
	SEODescription  string    `json:"seoDescription,omitempty"`
	ViewCount       int64     `json:"viewCount"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type friendResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	URL         string    `json:"url"`
	AvatarURL   string    `json:"avatarUrl,omitempty"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	SortOrder   int       `json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (h *PageHandler) listPublicPages(c *gin.Context) {
	pages, pagination, err := h.service.ListPublicPages(c.Request.Context(), pageFilter(c))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OKWithMeta(c, toPageResponses(pages), pagination)
}

func (h *PageHandler) getPublicPage(c *gin.Context) {
	page, err := h.service.GetPublicPage(c.Request.Context(), c.Param("slug"))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toPageResponse(page))
}

func (h *PageHandler) listAdminPages(c *gin.Context) {
	pages, pagination, err := h.service.ListAdminPages(c.Request.Context(), pageFilter(c))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OKWithMeta(c, toPageResponses(pages), pagination)
}

func (h *PageHandler) getAdminPage(c *gin.Context) {
	page, err := h.service.GetAdminPage(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toPageResponse(page))
}

func (h *PageHandler) createPage(c *gin.Context) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page request", nil)
		return
	}
	page, err := h.service.CreatePage(c.Request.Context(), usecase.PageRequest{AuthorID: middleware.CurrentUserID(c), Title: req.Title, Slug: req.Slug, Subtitle: req.Subtitle, ContentMarkdown: req.ContentMarkdown, Template: req.Template, Status: req.Status, SortOrder: req.SortOrder, SEOTitle: req.SEOTitle, SEODescription: req.SEODescription})
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toPageResponse(page))
}

func (h *PageHandler) updatePage(c *gin.Context) {
	var req pageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page request", nil)
		return
	}
	page, err := h.service.UpdatePage(c.Request.Context(), c.Param("id"), usecase.PageRequest{Title: req.Title, Slug: req.Slug, Subtitle: req.Subtitle, ContentMarkdown: req.ContentMarkdown, Template: req.Template, Status: req.Status, SortOrder: req.SortOrder, SEOTitle: req.SEOTitle, SEODescription: req.SEODescription})
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toPageResponse(page))
}

func (h *PageHandler) deletePage(c *gin.Context) {
	if err := h.service.DeletePage(c.Request.Context(), c.Param("id")); err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *PageHandler) listFriends(c *gin.Context) {
	friends, err := h.service.ListFriends(c.Request.Context())
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toFriendResponses(friends))
}

func (h *PageHandler) listAdminFriends(c *gin.Context) {
	friends, err := h.service.ListAdminFriends(c.Request.Context(), c.Query("status"))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toFriendResponses(friends))
}

func (h *PageHandler) createFriend(c *gin.Context) {
	var req friendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid friend request", nil)
		return
	}
	friend, err := h.service.CreateFriend(c.Request.Context(), friendRequestToUsecase(req))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toFriendResponse(friend))
}

func (h *PageHandler) updateFriend(c *gin.Context) {
	var req friendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid friend request", nil)
		return
	}
	friend, err := h.service.UpdateFriend(c.Request.Context(), c.Param("id"), friendRequestToUsecase(req))
	if err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, toFriendResponse(friend))
}

func (h *PageHandler) deleteFriend(c *gin.Context) {
	if err := h.service.DeleteFriend(c.Request.Context(), c.Param("id")); err != nil {
		h.handlePageError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *PageHandler) handlePageError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrPageNotFound):
		response.Error(c, http.StatusNotFound, "PAGE_NOT_FOUND", "Page not found", nil)
	case errors.Is(err, domain.ErrFriendNotFound):
		response.Error(c, http.StatusNotFound, "FRIEND_NOT_FOUND", "Friend not found", nil)
	case errors.Is(err, domain.ErrSlugTaken):
		response.Error(c, http.StatusConflict, "SLUG_TAKEN", "Page slug already exists", nil)
	case errors.Is(err, domain.ErrURLTaken):
		response.Error(c, http.StatusConflict, "URL_TAKEN", "Friend URL already exists", nil)
	case errors.Is(err, domain.ErrInvalidPage), errors.Is(err, domain.ErrInvalidFriend):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid page request", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Page request failed", nil)
	}
}

func pageFilter(c *gin.Context) domain.PageFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	filter := domain.PageFilter{Page: page, PerPage: perPage, Search: c.Query("search")}
	if status := c.Query("status"); status != "" {
		pageStatus := domain.PageStatus(status)
		if pageStatus.Valid() {
			filter.Status = &pageStatus
		}
	}
	return filter
}

func friendRequestToUsecase(req friendRequest) usecase.FriendRequest {
	return usecase.FriendRequest{Name: req.Name, URL: req.URL, AvatarURL: req.AvatarURL, Description: req.Description, Status: req.Status, SortOrder: req.SortOrder}
}

func toPageResponses(pages []domain.Page) []pageResponse {
	result := make([]pageResponse, 0, len(pages))
	for _, page := range pages {
		pageCopy := page
		result = append(result, *toPageResponse(&pageCopy))
	}
	return result
}

func toPageResponse(page *domain.Page) *pageResponse {
	return &pageResponse{ID: page.ID, AuthorID: page.AuthorID, Title: page.Title, Slug: page.Slug, Subtitle: page.Subtitle, ContentMarkdown: page.ContentMarkdown, ContentText: page.ContentText, Template: string(page.Template), Status: string(page.Status), SortOrder: page.SortOrder, SEOTitle: page.SEOTitle, SEODescription: page.SEODescription, ViewCount: page.ViewCount, CreatedAt: page.CreatedAt, UpdatedAt: page.UpdatedAt}
}

func toFriendResponses(friends []domain.Friend) []friendResponse {
	result := make([]friendResponse, 0, len(friends))
	for _, friend := range friends {
		friendCopy := friend
		result = append(result, *toFriendResponse(&friendCopy))
	}
	return result
}

func toFriendResponse(friend *domain.Friend) *friendResponse {
	return &friendResponse{ID: friend.ID, Name: friend.Name, URL: friend.URL, AvatarURL: friend.AvatarURL, Description: friend.Description, Status: string(friend.Status), SortOrder: friend.SortOrder, CreatedAt: friend.CreatedAt, UpdatedAt: friend.UpdatedAt}
}
