package handler

import (
	"context"
	"encoding/xml"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/subscribe"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/subscribe"
	"github.com/gin-gonic/gin"
)

type SubscribeService interface {
	Subscribe(ctx context.Context, req usecase.SubscribeRequest) (*domain.Subscriber, error)
	Verify(ctx context.Context, token string) error
	Unsubscribe(ctx context.Context, token string) error
	ListSubscribers(ctx context.Context, filter domain.SubscriberFilter) ([]domain.Subscriber, domain.Pagination, error)
	NotifyPostPublished(ctx context.Context, postID string) error
}

type FeedPostService interface {
	ListPublishedPosts(ctx context.Context, filter postDomain.PostListFilter) ([]postDomain.Post, postDomain.Pagination, error)
}

type SubscribeHandler struct {
	service SubscribeService
	posts   FeedPostService
	siteURL string
}

func NewSubscribeHandler(service SubscribeService, posts FeedPostService, siteURL string) *SubscribeHandler {
	return &SubscribeHandler{service: service, posts: posts, siteURL: strings.TrimRight(strings.TrimSpace(siteURL), "/")}
}

func (h *SubscribeHandler) RegisterFeedRoute(router *gin.Engine) {
	router.GET("/feed.xml", h.feed)
}

func (h *SubscribeHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.POST("/subscribe", h.subscribe)
	group.GET("/subscribe/verify", h.verify)
	group.GET("/subscribe/unsubscribe", h.unsubscribe)
}

func (h *SubscribeHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/subscribers", h.listAdmin)
	admin.POST("/subscribers/notify-post", h.notifyPost)
}

type subscribeRequest struct {
	Email string `json:"email" binding:"required"`
	Name  string `json:"name"`
}

type notifyPostRequest struct {
	PostID string `json:"postId" binding:"required"`
}

type subscriberResponse struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	Name           string     `json:"name,omitempty"`
	Status         string     `json:"status"`
	VerifiedAt     *time.Time `json:"verifiedAt,omitempty"`
	UnsubscribedAt *time.Time `json:"unsubscribedAt,omitempty"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Items       []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	GUID        string `xml:"guid"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate,omitempty"`
}

func (h *SubscribeHandler) feed(c *gin.Context) {
	posts, _, err := h.posts.ListPublishedPosts(c.Request.Context(), postDomain.PostListFilter{Page: 1, PerPage: 20})
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Feed request failed", nil)
		return
	}
	items := make([]rssItem, 0, len(posts))
	for _, post := range posts {
		item := rssItem{Title: post.Title, Link: h.publicURL("/posts/" + post.Slug), GUID: post.ID, Description: post.Excerpt}
		if post.PublishedAt != nil {
			item.PubDate = post.PublishedAt.Format(time.RFC1123Z)
		}
		items = append(items, item)
	}
	c.Header("Content-Type", "application/rss+xml; charset=utf-8")
	c.XML(http.StatusOK, rssFeed{Version: "2.0", Channel: rssChannel{Title: "Prosel", Link: h.publicURL("/"), Description: "A personal blog powered by Prosel", Items: items}})
}

func (h *SubscribeHandler) subscribe(c *gin.Context) {
	var req subscribeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid subscribe request", nil)
		return
	}
	subscriber, err := h.service.Subscribe(c.Request.Context(), usecase.SubscribeRequest{Email: req.Email, Name: req.Name})
	if err != nil {
		h.handleSubscribeError(c, err)
		return
	}
	response.OK(c, toSubscriberResponse(subscriber))
}

func (h *SubscribeHandler) verify(c *gin.Context) {
	if err := h.service.Verify(c.Request.Context(), c.Query("token")); err != nil {
		h.handleSubscribeError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *SubscribeHandler) unsubscribe(c *gin.Context) {
	if err := h.service.Unsubscribe(c.Request.Context(), c.Query("token")); err != nil {
		h.handleSubscribeError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *SubscribeHandler) listAdmin(c *gin.Context) {
	subscribers, pagination, err := h.service.ListSubscribers(c.Request.Context(), subscriberFilter(c))
	if err != nil {
		h.handleSubscribeError(c, err)
		return
	}
	response.OKWithMeta(c, toSubscriberResponses(subscribers), pagination)
}

func (h *SubscribeHandler) notifyPost(c *gin.Context) {
	var req notifyPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid notify request", nil)
		return
	}
	if err := h.service.NotifyPostPublished(c.Request.Context(), req.PostID); err != nil {
		h.handleSubscribeError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *SubscribeHandler) handleSubscribeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrSubscriberNotFound):
		response.Error(c, http.StatusNotFound, "SUBSCRIBER_NOT_FOUND", "Subscriber not found", nil)
	case errors.Is(err, domain.ErrSubscriberExists):
		response.Error(c, http.StatusConflict, "SUBSCRIBER_EXISTS", "Subscriber already exists", nil)
	case errors.Is(err, domain.ErrInvalidSubscriber):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid subscribe request", nil)
	case errors.Is(err, postDomain.ErrPostNotFound):
		response.Error(c, http.StatusNotFound, "POST_NOT_FOUND", "Post not found", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Subscribe request failed", nil)
	}
}

func subscriberFilter(c *gin.Context) domain.SubscriberFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	filter := domain.SubscriberFilter{Page: page, PerPage: perPage, Search: c.Query("search")}
	if status := c.Query("status"); status != "" {
		subscriberStatus := domain.SubscriberStatus(status)
		if subscriberStatus.Valid() {
			filter.Status = &subscriberStatus
		}
	}
	return filter
}

func toSubscriberResponses(subscribers []domain.Subscriber) []subscriberResponse {
	result := make([]subscriberResponse, 0, len(subscribers))
	for _, subscriber := range subscribers {
		copy := subscriber
		result = append(result, *toSubscriberResponse(&copy))
	}
	return result
}

func toSubscriberResponse(subscriber *domain.Subscriber) *subscriberResponse {
	return &subscriberResponse{ID: subscriber.ID, Email: subscriber.Email, Name: subscriber.Name, Status: string(subscriber.Status), VerifiedAt: subscriber.VerifiedAt, UnsubscribedAt: subscriber.UnsubscribedAt, CreatedAt: subscriber.CreatedAt, UpdatedAt: subscriber.UpdatedAt}
}

func (h *SubscribeHandler) publicURL(path string) string {
	if h.siteURL == "" {
		return path
	}
	if path == "/" {
		return h.siteURL + "/"
	}
	return h.siteURL + path
}
