package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/comment"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/comment"
	"github.com/gin-gonic/gin"
)

type CommentService interface {
	SubmitComment(ctx context.Context, req usecase.SubmitCommentRequest, meta usecase.ClientMeta) (*domain.Comment, error)
	ReplyAsAdmin(ctx context.Context, req usecase.AdminReplyRequest) (*domain.Comment, error)
	ListPublicComments(ctx context.Context, refType string, refID string) ([]domain.CommentNode, error)
	ModerateComment(ctx context.Context, id string, status string) error
	DeleteComment(ctx context.Context, id string) error
	ListAdminComments(ctx context.Context, filter domain.CommentFilter) ([]domain.Comment, domain.Pagination, error)
}

type CommentHandler struct{ service CommentService }

func NewCommentHandler(service CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

func (h *CommentHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.GET("/comments", h.listPublic)
	group.POST("/comments", h.submit)
}

func (h *CommentHandler) RegisterProtectedRoutes(admin *gin.RouterGroup) {
	admin.GET("/comments", h.listAdmin)
	admin.PATCH("/comments/:id/status", h.moderate)
	admin.POST("/comments/:id/reply", h.reply)
	admin.DELETE("/comments/:id", h.delete)
}

type submitCommentRequest struct {
	RefType       string `json:"refType" binding:"required"`
	RefID         string `json:"refId" binding:"required"`
	ParentID      string `json:"parentId"`
	AuthorName    string `json:"authorName" binding:"required"`
	AuthorEmail   string `json:"authorEmail" binding:"required"`
	AuthorWebsite string `json:"authorWebsite"`
	Content       string `json:"content" binding:"required"`
}

type moderateCommentRequest struct {
	Status string `json:"status" binding:"required"`
}

type adminReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

type commentResponse struct {
	ID            string            `json:"id"`
	RefType       string            `json:"refType"`
	RefID         string            `json:"refId"`
	ParentID      *string           `json:"parentId,omitempty"`
	RootID        *string           `json:"rootId,omitempty"`
	AuthorName    string            `json:"authorName"`
	AuthorEmail   string            `json:"authorEmail,omitempty"`
	AuthorWebsite string            `json:"authorWebsite,omitempty"`
	Content       string            `json:"content"`
	Status        string            `json:"status"`
	IsAdminReply  bool              `json:"isAdminReply"`
	IsPinned      bool              `json:"isPinned"`
	ReplyCount    int               `json:"replyCount"`
	Children      []commentResponse `json:"children,omitempty"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

func (h *CommentHandler) listPublic(c *gin.Context) {
	nodes, err := h.service.ListPublicComments(c.Request.Context(), c.Query("refType"), c.Query("refId"))
	if err != nil {
		h.handleCommentError(c, err)
		return
	}
	response.OK(c, toCommentNodeResponses(nodes))
}

func (h *CommentHandler) submit(c *gin.Context) {
	var req submitCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid comment request", nil)
		return
	}
	comment, err := h.service.SubmitComment(c.Request.Context(), usecase.SubmitCommentRequest{RefType: req.RefType, RefID: req.RefID, ParentID: req.ParentID, AuthorName: req.AuthorName, AuthorEmail: req.AuthorEmail, AuthorWebsite: req.AuthorWebsite, Content: req.Content}, usecase.ClientMeta{IP: c.ClientIP(), UserAgent: c.Request.UserAgent()})
	if err != nil {
		h.handleCommentError(c, err)
		return
	}
	response.OK(c, toCommentResponse(comment, true))
}

func (h *CommentHandler) listAdmin(c *gin.Context) {
	comments, pagination, err := h.service.ListAdminComments(c.Request.Context(), commentFilter(c))
	if err != nil {
		h.handleCommentError(c, err)
		return
	}
	response.OKWithMeta(c, toCommentResponses(comments, true), pagination)
}

func (h *CommentHandler) moderate(c *gin.Context) {
	var req moderateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid moderation request", nil)
		return
	}
	if err := h.service.ModerateComment(c.Request.Context(), c.Param("id"), req.Status); err != nil {
		h.handleCommentError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *CommentHandler) reply(c *gin.Context) {
	var req adminReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid reply request", nil)
		return
	}
	comment, err := h.service.ReplyAsAdmin(c.Request.Context(), usecase.AdminReplyRequest{ParentID: c.Param("id"), Content: req.Content})
	if err != nil {
		h.handleCommentError(c, err)
		return
	}
	response.OK(c, toCommentResponse(comment, true))
}

func (h *CommentHandler) delete(c *gin.Context) {
	if err := h.service.DeleteComment(c.Request.Context(), c.Param("id")); err != nil {
		h.handleCommentError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *CommentHandler) handleCommentError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrCommentNotFound):
		response.Error(c, http.StatusNotFound, "COMMENT_NOT_FOUND", "Comment not found", nil)
	case errors.Is(err, domain.ErrInvalidComment):
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid comment request", nil)
	case errors.Is(err, domain.ErrCommentDepthExceeded):
		response.Error(c, http.StatusBadRequest, "COMMENT_DEPTH_EXCEEDED", "Comment nesting depth exceeded", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Comment request failed", nil)
	}
}

func commentFilter(c *gin.Context) domain.CommentFilter {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "20"))
	filter := domain.CommentFilter{Page: page, PerPage: perPage, RefType: domain.RefType(c.Query("refType")), RefID: c.Query("refId"), Search: c.Query("search")}
	if status := c.Query("status"); status != "" {
		commentStatus := domain.CommentStatus(status)
		if commentStatus.Valid() {
			filter.Status = &commentStatus
		}
	}
	return filter
}

func toCommentNodeResponses(nodes []domain.CommentNode) []commentResponse {
	result := make([]commentResponse, 0, len(nodes))
	for _, node := range nodes {
		response := *toCommentResponse(&node.Comment, false)
		response.Children = toCommentNodeResponses(node.Children)
		result = append(result, response)
	}
	return result
}

func toCommentResponses(comments []domain.Comment, includeEmail bool) []commentResponse {
	result := make([]commentResponse, 0, len(comments))
	for _, comment := range comments {
		copy := comment
		result = append(result, *toCommentResponse(&copy, includeEmail))
	}
	return result
}

func toCommentResponse(comment *domain.Comment, includeEmail bool) *commentResponse {
	res := &commentResponse{ID: comment.ID, RefType: string(comment.RefType), RefID: comment.RefID, ParentID: comment.ParentID, RootID: comment.RootID, AuthorName: comment.AuthorName, AuthorWebsite: comment.AuthorWebsite, Content: comment.Content, Status: string(comment.Status), IsAdminReply: comment.IsAdminReply, IsPinned: comment.IsPinned, ReplyCount: comment.ReplyCount, CreatedAt: comment.CreatedAt, UpdatedAt: comment.UpdatedAt}
	if includeEmail {
		res.AuthorEmail = comment.AuthorEmail
	}
	return res
}
