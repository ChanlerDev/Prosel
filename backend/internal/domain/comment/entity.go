package comment

import (
	"errors"
	"time"
)

type RefType string

type CommentStatus string

const (
	RefPost RefType = "post"
	RefNote RefType = "note"
	RefPage RefType = "page"

	CommentPending  CommentStatus = "pending"
	CommentApproved CommentStatus = "approved"
	CommentRejected CommentStatus = "rejected"
	CommentSpam     CommentStatus = "spam"
)

var (
	ErrCommentNotFound      = errors.New("comment not found")
	ErrInvalidComment       = errors.New("invalid comment")
	ErrCommentDepthExceeded = errors.New("comment depth exceeded")
)

type Comment struct {
	ID            string
	RefType       RefType
	RefID         string
	ParentID      *string
	RootID        *string
	AuthorName    string
	AuthorEmail   string
	AuthorWebsite string
	AuthorIP      string
	UserAgent     string
	Content       string
	Status        CommentStatus
	IsAdminReply  bool
	IsPinned      bool
	ReplyCount    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CommentNode struct {
	Comment
	Children []CommentNode
}

func (t RefType) Valid() bool {
	return t == RefPost || t == RefNote || t == RefPage
}

func (s CommentStatus) Valid() bool {
	return s == CommentPending || s == CommentApproved || s == CommentRejected || s == CommentSpam
}

func (s CommentStatus) ValidModerationTarget() bool {
	return s == CommentApproved || s == CommentRejected || s == CommentSpam
}
