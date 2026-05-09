package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/post"
	"github.com/chanler/prosel/backend/internal/interfaces/http/middleware"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/post"
	"github.com/gin-gonic/gin"
)

type PostService interface {
	CreatePost(ctx context.Context, req usecase.CreatePostRequest) (*domain.Post, error)
	UpdatePost(ctx context.Context, id string, req usecase.UpdatePostRequest) (*domain.Post, error)
	PublishPost(ctx context.Context, id string) (*domain.Post, error)
	UnpublishPost(ctx context.Context, id string) (*domain.Post, error)
	DeletePost(ctx context.Context, id string) error
	GetAdminPost(ctx context.Context, id string) (*domain.Post, error)
	GetPublishedPost(ctx context.Context, slug string) (*domain.Post, error)
	ListPublishedPosts(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error)
	ListAdminPosts(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error)
}

type PostHandler struct {
	service PostService
}

func NewPostHandler(service PostService) *PostHandler { return &PostHandler{service: service} }

func (h *PostHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/posts", h.listPublished)
	group.GET("/posts/:slug", h.getPublished)
}

func (h *PostHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/posts", h.listAdmin)
	admin.GET("/posts/:id", h.getAdmin)
	admin.POST("/posts", h.create)
	admin.PATCH("/posts/:id", h.update)
	admin.DELETE("/posts/:id", h.delete)
	admin.POST("/posts/:id/publish", h.publish)
	admin.POST("/posts/:id/unpublish", h.unpublish)
}

type postRequest struct {
	CategoryID      string `json:"categoryId"`
	Title           string `json:"title" binding:"required"`
	Slug            string `json:"slug"`
	Excerpt         string `json:"excerpt"`
	ContentMarkdown string `json:"contentMarkdown"`
	CoverImage      string `json:"coverImage"`
	Featured        bool   `json:"featured"`
	SEOTitle        string `json:"seoTitle"`
	SEODescription  string `json:"seoDescription"`
}

type postResponse struct {
	ID              string     `json:"id"`
	AuthorID        string     `json:"authorId"`
	CategoryID      *string    `json:"categoryId,omitempty"`
	Title           string     `json:"title"`
	Slug            string     `json:"slug"`
	Excerpt         string     `json:"excerpt,omitempty"`
	ContentMarkdown string     `json:"contentMarkdown,omitempty"`
	ContentText     string     `json:"contentText,omitempty"`
	CoverImage      string     `json:"coverImage,omitempty"`
	Status          string     `json:"status"`
	Featured        bool       `json:"featured"`
	PinnedAt        *time.Time `json:"pinnedAt,omitempty"`
	PublishedAt     *time.Time `json:"publishedAt,omitempty"`
	SEOTitle        string     `json:"seoTitle,omitempty"`
	SEODescription  string     `json:"seoDescription,omitempty"`
	ViewCount       int64      `json:"viewCount"`
	LikeCount       int64      `json:"likeCount"`
	CommentCount    int64      `json:"commentCount"`
	CreatedAt       time.Time  `json:"createdAt"`
	UpdatedAt       time.Time  `json:"updatedAt"`
}

func (h *PostHandler) listPublished(c *gin.Context) {
	posts, pagination, err := h.service.ListPublishedPosts(c.Request.Context(), listFilter(c))
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OKWithMeta(c, toPostResponses(posts), pagination)
}

func (h *PostHandler) getPublished(c *gin.Context) {
	post, err := h.service.GetPublishedPost(c.Request.Context(), c.Param("slug"))
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, toPostResponse(post))
}

func (h *PostHandler) listAdmin(c *gin.Context) {
	posts, pagination, err := h.service.ListAdminPosts(c.Request.Context(), listFilter(c))
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OKWithMeta(c, toPostResponses(posts), pagination)
}

func (h *PostHandler) getAdmin(c *gin.Context) {
	post, err := h.service.GetAdminPost(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, toPostResponse(post))
}

func (h *PostHandler) create(c *gin.Context) {
	var req postRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid post request", nil)
		return
	}
	post, err := h.service.CreatePost(c.Request.Context(), usecase.CreatePostRequest{AuthorID: middleware.CurrentUserID(c), CategoryID: optionalString(req.CategoryID), Title: req.Title, Slug: req.Slug, Excerpt: req.Excerpt, ContentMarkdown: req.ContentMarkdown, CoverImage: req.CoverImage, Featured: req.Featured, SEOTitle: req.SEOTitle, SEODescription: req.SEODescription})
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, toPostResponse(post))
}

func (h *PostHandler) update(c *gin.Context) {
	var req postRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid post request", nil)
		return
	}
	post, err := h.service.UpdatePost(c.Request.Context(), c.Param("id"), usecase.UpdatePostRequest{CategoryID: optionalString(req.CategoryID), Title: req.Title, Slug: req.Slug, Excerpt: req.Excerpt, ContentMarkdown: req.ContentMarkdown, CoverImage: req.CoverImage, Featured: req.Featured, SEOTitle: req.SEOTitle, SEODescription: req.SEODescription})
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, toPostResponse(post))
}

func (h *PostHandler) delete(c *gin.Context) {
	if err := h.service.DeletePost(c.Request.Context(), c.Param("id")); err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *PostHandler) publish(c *gin.Context) {
	post, err := h.service.PublishPost(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, toPostResponse(post))
}

func (h *PostHandler) unpublish(c *gin.Context) {
	post, err := h.service.UnpublishPost(c.Request.Context(), c.Param("id"))
	if err != nil {
		h.handlePostError(c, err)
		return
	}
	response.OK(c, toPostResponse(post))
}

func (h *PostHandler) handlePostError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		response.Error(c, http.StatusNotFound, "POST_NOT_FOUND", "Post not found", nil)
	case errors.Is(err, domain.ErrSlugTaken):
		response.Error(c, http.StatusConflict, "SLUG_TAKEN", "Post slug already exists", nil)
	case errors.Is(err, domain.ErrInvalidPost):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid post request", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Post request failed", nil)
	}
}

func listFilter(c *gin.Context) domain.PostListFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	filter := domain.PostListFilter{Page: page, PerPage: perPage, Search: c.Query("search"), CategoryID: c.Query("category")}
	if status := c.Query("status"); status != "" {
		postStatus := domain.PostStatus(status)
		if postStatus.Valid() {
			filter.Status = &postStatus
		}
	}
	if featured := c.Query("featured"); featured != "" {
		value, err := strconv.ParseBool(featured)
		if err == nil {
			filter.Featured = &value
		}
	}
	return filter
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func toPostResponses(posts []domain.Post) []postResponse {
	result := make([]postResponse, 0, len(posts))
	for _, post := range posts {
		postCopy := post
		result = append(result, *toPostResponse(&postCopy))
	}
	return result
}

func toPostResponse(post *domain.Post) *postResponse {
	return &postResponse{ID: post.ID, AuthorID: post.AuthorID, CategoryID: post.CategoryID, Title: post.Title, Slug: post.Slug, Excerpt: post.Excerpt, ContentMarkdown: post.ContentMarkdown, ContentText: post.ContentText, CoverImage: post.CoverImage, Status: string(post.Status), Featured: post.Featured, PinnedAt: post.PinnedAt, PublishedAt: post.PublishedAt, SEOTitle: post.SEOTitle, SEODescription: post.SEODescription, ViewCount: post.ViewCount, LikeCount: post.LikeCount, CommentCount: post.CommentCount, CreatedAt: post.CreatedAt, UpdatedAt: post.UpdatedAt}
}
