package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/taxonomy"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	postUsecase "github.com/chanler/prosel/backend/internal/usecase/post"
	usecase "github.com/chanler/prosel/backend/internal/usecase/taxonomy"
	"github.com/gin-gonic/gin"
)

type TaxonomyService interface {
	ListCategories(ctx context.Context) ([]domain.CategoryNode, error)
	CreateCategory(ctx context.Context, req usecase.CategoryRequest) (*domain.Category, error)
	UpdateCategory(ctx context.Context, id string, req usecase.CategoryRequest) (*domain.Category, error)
	DeleteCategory(ctx context.Context, id string) error
	ListTags(ctx context.Context) ([]domain.TagWithCount, error)
	CreateTag(ctx context.Context, req usecase.TagRequest) (*domain.Tag, error)
	UpdateTag(ctx context.Context, id string, req usecase.TagRequest) (*domain.Tag, error)
	DeleteTag(ctx context.Context, id string) error
	ListTopics(ctx context.Context) ([]domain.Topic, error)
	GetTopic(ctx context.Context, slug string) (*domain.TopicDetail, error)
	CreateTopic(ctx context.Context, req usecase.TopicRequest) (*domain.Topic, error)
	UpdateTopic(ctx context.Context, id string, req usecase.TopicRequest) (*domain.Topic, error)
	DeleteTopic(ctx context.Context, id string) error
}

type TaxonomyPostService interface {
	ListPublishedPosts(ctx context.Context, filter postDomain.PostListFilter) ([]postDomain.Post, postDomain.Pagination, error)
}

type TaxonomyHandler struct {
	service TaxonomyService
	posts   TaxonomyPostService
}

func NewTaxonomyHandler(service TaxonomyService, posts TaxonomyPostService) *TaxonomyHandler {
	return &TaxonomyHandler{service: service, posts: posts}
}

func (h *TaxonomyHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/categories", h.listCategories)
	group.GET("/categories/:slug/posts", h.categoryPosts)
	group.GET("/tags", h.listTags)
	group.GET("/tags/:slug/posts", h.tagPosts)
	group.GET("/topics", h.listTopics)
	group.GET("/topics/:slug", h.getTopic)
}

func (h *TaxonomyHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.POST("/categories", h.createCategory)
	admin.PATCH("/categories/:id", h.updateCategory)
	admin.DELETE("/categories/:id", h.deleteCategory)
	admin.POST("/tags", h.createTag)
	admin.PATCH("/tags/:id", h.updateTag)
	admin.DELETE("/tags/:id", h.deleteTag)
	admin.POST("/topics", h.createTopic)
	admin.PATCH("/topics/:id", h.updateTopic)
	admin.DELETE("/topics/:id", h.deleteTopic)
}

type categoryRequest struct {
	ParentID    string `json:"parentId"`
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	SortOrder   int    `json:"sortOrder"`
}
type tagRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug"`
	Color       string `json:"color"`
	Description string `json:"description"`
}
type topicItemRequest struct {
	RefType   string `json:"refType"`
	RefID     string `json:"refId"`
	SortOrder int    `json:"sortOrder"`
}
type topicRequest struct {
	Name        string             `json:"name" binding:"required"`
	Slug        string             `json:"slug"`
	Description string             `json:"description"`
	CoverImage  string             `json:"coverImage"`
	SortOrder   int                `json:"sortOrder"`
	Items       []topicItemRequest `json:"items"`
}

type categoryResponse struct {
	ID          string             `json:"id"`
	ParentID    *string            `json:"parentId,omitempty"`
	Name        string             `json:"name"`
	Slug        string             `json:"slug"`
	Description string             `json:"description,omitempty"`
	SortOrder   int                `json:"sortOrder"`
	PostCount   int64              `json:"postCount"`
	Children    []categoryResponse `json:"children"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}
type tagResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Color       string    `json:"color,omitempty"`
	Description string    `json:"description,omitempty"`
	PostCount   int64     `json:"postCount,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
type topicResponse struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Slug        string              `json:"slug"`
	Description string              `json:"description,omitempty"`
	CoverImage  string              `json:"coverImage,omitempty"`
	SortOrder   int                 `json:"sortOrder"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
	Items       []topicItemResponse `json:"items,omitempty"`
}
type topicItemResponse struct {
	RefType   string `json:"refType"`
	RefID     string `json:"refId"`
	Title     string `json:"title"`
	Slug      string `json:"slug,omitempty"`
	SortOrder int    `json:"sortOrder"`
}

func (h *TaxonomyHandler) listCategories(c *gin.Context) {
	nodes, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, toCategoryResponses(nodes))
}
func (h *TaxonomyHandler) listTags(c *gin.Context) {
	tags, err := h.service.ListTags(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, toTagResponses(tags))
}
func (h *TaxonomyHandler) listTopics(c *gin.Context) {
	topics, err := h.service.ListTopics(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, toTopicResponses(topics))
}
func (h *TaxonomyHandler) getTopic(c *gin.Context) {
	topic, err := h.service.GetTopic(c.Request.Context(), c.Param("slug"))
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, toTopicDetailResponse(topic))
}

func (h *TaxonomyHandler) categoryPosts(c *gin.Context) {
	categories, err := h.service.ListCategories(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	category := findCategory(categories, c.Param("slug"))
	if category == nil {
		h.handleError(c, domain.ErrTaxonomyNotFound)
		return
	}
	posts, pagination, err := h.posts.ListPublishedPosts(c.Request.Context(), postListFilter(c, category.ID, ""))
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OKWithMeta(c, toPostResponses(posts), pagination)
}

func (h *TaxonomyHandler) tagPosts(c *gin.Context) {
	tags, err := h.service.ListTags(c.Request.Context())
	if err != nil {
		h.handleError(c, err)
		return
	}
	var tagID string
	for _, tag := range tags {
		if tag.Slug == c.Param("slug") {
			tagID = tag.ID
			break
		}
	}
	if tagID == "" {
		h.handleError(c, domain.ErrTaxonomyNotFound)
		return
	}
	posts, pagination, err := h.posts.ListPublishedPosts(c.Request.Context(), postListFilter(c, "", tagID))
	if err != nil {
		h.handleError(c, err)
		return
	}
	response.OKWithMeta(c, toPostResponses(posts), pagination)
}

func (h *TaxonomyHandler) createCategory(c *gin.Context) {
	var req categoryRequest
	if bind(c, &req) {
		category, err := h.service.CreateCategory(c.Request.Context(), usecase.CategoryRequest{ParentID: optionalString(req.ParentID), Name: req.Name, Slug: req.Slug, Description: req.Description, SortOrder: req.SortOrder})
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.OK(c, toCategoryResponse(domain.CategoryNode{Category: *category}))
	}
}
func (h *TaxonomyHandler) updateCategory(c *gin.Context) {
	var req categoryRequest
	if bind(c, &req) {
		category, err := h.service.UpdateCategory(c.Request.Context(), c.Param("id"), usecase.CategoryRequest{ParentID: optionalString(req.ParentID), Name: req.Name, Slug: req.Slug, Description: req.Description, SortOrder: req.SortOrder})
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.OK(c, toCategoryResponse(domain.CategoryNode{Category: *category}))
	}
}
func (h *TaxonomyHandler) deleteCategory(c *gin.Context) {
	if err := h.service.DeleteCategory(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}
func (h *TaxonomyHandler) createTag(c *gin.Context) {
	var req tagRequest
	if bind(c, &req) {
		tag, err := h.service.CreateTag(c.Request.Context(), usecase.TagRequest{Name: req.Name, Slug: req.Slug, Color: req.Color, Description: req.Description})
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.OK(c, toTagResponse(domain.TagWithCount{Tag: *tag}))
	}
}
func (h *TaxonomyHandler) updateTag(c *gin.Context) {
	var req tagRequest
	if bind(c, &req) {
		tag, err := h.service.UpdateTag(c.Request.Context(), c.Param("id"), usecase.TagRequest{Name: req.Name, Slug: req.Slug, Color: req.Color, Description: req.Description})
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.OK(c, toTagResponse(domain.TagWithCount{Tag: *tag}))
	}
}
func (h *TaxonomyHandler) deleteTag(c *gin.Context) {
	if err := h.service.DeleteTag(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}
func (h *TaxonomyHandler) createTopic(c *gin.Context) {
	var req topicRequest
	if bind(c, &req) {
		topic, err := h.service.CreateTopic(c.Request.Context(), topicUsecaseRequest(req))
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.OK(c, toTopicResponse(*topic))
	}
}
func (h *TaxonomyHandler) updateTopic(c *gin.Context) {
	var req topicRequest
	if bind(c, &req) {
		topic, err := h.service.UpdateTopic(c.Request.Context(), c.Param("id"), topicUsecaseRequest(req))
		if err != nil {
			h.handleError(c, err)
			return
		}
		response.OK(c, toTopicResponse(*topic))
	}
}
func (h *TaxonomyHandler) deleteTopic(c *gin.Context) {
	if err := h.service.DeleteTopic(c.Request.Context(), c.Param("id")); err != nil {
		h.handleError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid taxonomy request", nil)
		return false
	}
	return true
}
func (h *TaxonomyHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrTaxonomyNotFound):
		response.Error(c, http.StatusNotFound, "TAXONOMY_NOT_FOUND", "Taxonomy item not found", nil)
	case errors.Is(err, domain.ErrSlugTaken):
		response.Error(c, http.StatusConflict, "SLUG_TAKEN", "Taxonomy slug already exists", nil)
	case errors.Is(err, domain.ErrInvalidTaxonomy):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid taxonomy request", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Taxonomy request failed", nil)
	}
}

func postListFilter(c *gin.Context, categoryID string, tagID string) postDomain.PostListFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	return postDomain.PostListFilter{Page: page, PerPage: perPage, CategoryID: categoryID, TagID: tagID}
}
func findCategory(nodes []domain.CategoryNode, slug string) *domain.Category {
	for _, node := range nodes {
		if node.Slug == slug {
			return &node.Category
		}
		if child := findCategory(node.Children, slug); child != nil {
			return child
		}
	}
	return nil
}
func topicUsecaseRequest(req topicRequest) usecase.TopicRequest {
	items := make([]domain.TopicItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, domain.TopicItem{RefType: item.RefType, RefID: item.RefID, SortOrder: item.SortOrder})
	}
	return usecase.TopicRequest{Name: req.Name, Slug: req.Slug, Description: req.Description, CoverImage: req.CoverImage, SortOrder: req.SortOrder, Items: items}
}
func toCategoryResponses(nodes []domain.CategoryNode) []categoryResponse {
	result := make([]categoryResponse, 0, len(nodes))
	for _, node := range nodes {
		result = append(result, toCategoryResponse(node))
	}
	return result
}
func toCategoryResponse(node domain.CategoryNode) categoryResponse {
	return categoryResponse{ID: node.ID, ParentID: node.ParentID, Name: node.Name, Slug: node.Slug, Description: node.Description, SortOrder: node.SortOrder, PostCount: node.PostCount, Children: toCategoryResponses(node.Children), CreatedAt: node.CreatedAt, UpdatedAt: node.UpdatedAt}
}
func toTagResponses(tags []domain.TagWithCount) []tagResponse {
	result := make([]tagResponse, 0, len(tags))
	for _, tag := range tags {
		result = append(result, toTagResponse(tag))
	}
	return result
}
func toTagResponse(tag domain.TagWithCount) tagResponse {
	return tagResponse{ID: tag.ID, Name: tag.Name, Slug: tag.Slug, Color: tag.Color, Description: tag.Description, PostCount: tag.PostCount, CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt}
}
func toTopicResponses(topics []domain.Topic) []topicResponse {
	result := make([]topicResponse, 0, len(topics))
	for _, topic := range topics {
		result = append(result, toTopicResponse(topic))
	}
	return result
}
func toTopicResponse(topic domain.Topic) topicResponse {
	return topicResponse{ID: topic.ID, Name: topic.Name, Slug: topic.Slug, Description: topic.Description, CoverImage: topic.CoverImage, SortOrder: topic.SortOrder, CreatedAt: topic.CreatedAt, UpdatedAt: topic.UpdatedAt}
}
func toTopicDetailResponse(topic *domain.TopicDetail) topicResponse {
	response := toTopicResponse(topic.Topic)
	for _, item := range topic.Items {
		response.Items = append(response.Items, topicItemResponse{RefType: item.RefType, RefID: item.RefID, Title: item.Title, Slug: item.Slug, SortOrder: item.SortOrder})
	}
	return response
}

var _ TaxonomyPostService = (*postUsecase.PostUsecase)(nil)
