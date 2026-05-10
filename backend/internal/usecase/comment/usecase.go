package comment

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/mail"
	"strings"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/comment"
)

const maxCommentDepth = 3

type CommentUsecase struct {
	comments domain.Repository
}

type SubmitCommentRequest struct {
	RefType       string
	RefID         string
	ParentID      string
	AuthorName    string
	AuthorEmail   string
	AuthorWebsite string
	Content       string
}

type AdminReplyRequest struct {
	ParentID string
	Content  string
}

type ClientMeta struct {
	IP        string
	UserAgent string
}

func NewCommentUsecase(comments domain.Repository) *CommentUsecase {
	return &CommentUsecase{comments: comments}
}

func (uc *CommentUsecase) SubmitComment(ctx context.Context, req SubmitCommentRequest, meta ClientMeta) (*domain.Comment, error) {
	refType := domain.RefType(strings.TrimSpace(req.RefType))
	refID := strings.TrimSpace(req.RefID)
	authorName := strings.TrimSpace(req.AuthorName)
	authorEmail := strings.ToLower(strings.TrimSpace(req.AuthorEmail))
	content := strings.TrimSpace(req.Content)
	if !refType.Valid() || refID == "" || authorName == "" || authorEmail == "" || content == "" {
		return nil, domain.ErrInvalidComment
	}
	if _, err := mail.ParseAddress(authorEmail); err != nil {
		return nil, domain.ErrInvalidComment
	}

	var parentID *string
	var rootID *string
	if trimmedParentID := strings.TrimSpace(req.ParentID); trimmedParentID != "" {
		parent, err := uc.comments.GetByID(ctx, trimmedParentID)
		if err != nil {
			return nil, err
		}
		if parent.RefType != refType || parent.RefID != refID {
			return nil, domain.ErrInvalidComment
		}
		if commentDepth(parent) >= maxCommentDepth {
			return nil, domain.ErrCommentDepthExceeded
		}
		parentID = &trimmedParentID
		if parent.RootID != nil {
			rootCopy := *parent.RootID
			rootID = &rootCopy
		} else {
			rootCopy := parent.ID
			rootID = &rootCopy
		}
	}

	now := time.Now().UTC()
	comment := &domain.Comment{ID: newID(), RefType: refType, RefID: refID, ParentID: parentID, RootID: rootID, AuthorName: authorName, AuthorEmail: authorEmail, AuthorWebsite: strings.TrimSpace(req.AuthorWebsite), AuthorIP: strings.TrimSpace(meta.IP), UserAgent: strings.TrimSpace(meta.UserAgent), Content: content, Status: domain.CommentPending, CreatedAt: now, UpdatedAt: now}
	if err := uc.comments.Create(ctx, comment); err != nil {
		return nil, err
	}
	if parentID != nil {
		if err := uc.comments.IncrementReplyCount(ctx, *parentID); err != nil {
			return nil, err
		}
	}
	return comment, nil
}

func (uc *CommentUsecase) ReplyAsAdmin(ctx context.Context, req AdminReplyRequest) (*domain.Comment, error) {
	parentID := strings.TrimSpace(req.ParentID)
	content := strings.TrimSpace(req.Content)
	if parentID == "" || content == "" {
		return nil, domain.ErrInvalidComment
	}
	parent, err := uc.comments.GetByID(ctx, parentID)
	if err != nil {
		return nil, err
	}
	if commentDepth(parent) >= maxCommentDepth {
		return nil, domain.ErrCommentDepthExceeded
	}
	rootID := parent.ID
	if parent.RootID != nil {
		rootID = *parent.RootID
	}
	now := time.Now().UTC()
	comment := &domain.Comment{ID: newID(), RefType: parent.RefType, RefID: parent.RefID, ParentID: &parentID, RootID: &rootID, AuthorName: "Admin", AuthorEmail: "admin@local", Content: content, Status: domain.CommentApproved, IsAdminReply: true, CreatedAt: now, UpdatedAt: now}
	if err := uc.comments.Create(ctx, comment); err != nil {
		return nil, err
	}
	if err := uc.comments.IncrementReplyCount(ctx, parentID); err != nil {
		return nil, err
	}
	return comment, nil
}

func (uc *CommentUsecase) ListPublicComments(ctx context.Context, refTypeValue string, refID string) ([]domain.CommentNode, error) {
	refType := domain.RefType(strings.TrimSpace(refTypeValue))
	refID = strings.TrimSpace(refID)
	if !refType.Valid() || refID == "" {
		return nil, domain.ErrInvalidComment
	}
	comments, err := uc.comments.ListByRef(ctx, refType, refID, true)
	if err != nil {
		return nil, err
	}
	return domain.BuildTree(comments), nil
}

func (uc *CommentUsecase) ModerateComment(ctx context.Context, id string, statusValue string) error {
	id = strings.TrimSpace(id)
	status := domain.CommentStatus(strings.TrimSpace(statusValue))
	if id == "" || !status.ValidModerationTarget() {
		return domain.ErrInvalidComment
	}
	return uc.comments.UpdateStatus(ctx, id, status)
}

func (uc *CommentUsecase) DeleteComment(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidComment
	}
	return uc.comments.Delete(ctx, id)
}

func (uc *CommentUsecase) ListAdminComments(ctx context.Context, filter domain.CommentFilter) ([]domain.Comment, domain.Pagination, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	filter.RefID = strings.TrimSpace(filter.RefID)
	if filter.RefType != "" && !filter.RefType.Valid() {
		filter.RefType = ""
	}
	if filter.Status != nil && !filter.Status.Valid() {
		filter.Status = nil
	}
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.comments.ListAdmin(ctx, filter)
}

func commentDepth(comment *domain.Comment) int {
	if comment.ParentID == nil {
		return 1
	}
	if comment.RootID == nil || *comment.RootID == *comment.ParentID {
		return 2
	}
	return 3
}

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))[:32]
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:])
}
