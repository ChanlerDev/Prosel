package comment

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/comment"
)

type fakeCommentRepo struct {
	comment       *domain.Comment
	comments      []domain.Comment
	pagination    domain.Pagination
	err           error
	listed        domain.CommentFilter
	statusUpdates []string
	incrementedID string
	deletedID     string
}

func (r *fakeCommentRepo) Create(ctx context.Context, comment *domain.Comment) error {
	r.comment = comment
	return r.err
}
func (r *fakeCommentRepo) UpdateStatus(ctx context.Context, id string, status domain.CommentStatus) error {
	r.statusUpdates = append(r.statusUpdates, id+":"+string(status))
	return r.err
}
func (r *fakeCommentRepo) Delete(ctx context.Context, id string) error {
	r.deletedID = id
	return r.err
}
func (r *fakeCommentRepo) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	return r.comment, r.err
}
func (r *fakeCommentRepo) ListByRef(ctx context.Context, refType domain.RefType, refID string, onlyApproved bool) ([]domain.Comment, error) {
	return r.comments, r.err
}
func (r *fakeCommentRepo) ListAdmin(ctx context.Context, filter domain.CommentFilter) ([]domain.Comment, domain.Pagination, error) {
	r.listed = filter
	return r.comments, r.pagination, r.err
}
func (r *fakeCommentRepo) IncrementReplyCount(ctx context.Context, id string) error {
	r.incrementedID = id
	return r.err
}

func TestSubmitCommentCreatesPendingRootComment(t *testing.T) {
	repo := &fakeCommentRepo{}
	uc := NewCommentUsecase(repo)

	comment, err := uc.SubmitComment(context.Background(), SubmitCommentRequest{RefType: "post", RefID: "post-1", AuthorName: " Ada ", AuthorEmail: " ADA@EXAMPLE.COM ", Content: " Hello "}, ClientMeta{IP: "127.0.0.1", UserAgent: "agent"})
	if err != nil {
		t.Fatalf("SubmitComment() error = %v", err)
	}
	if comment.Status != domain.CommentPending || comment.ID == "" || comment.CreatedAt.IsZero() {
		t.Fatalf("comment = %#v, want pending persisted comment", comment)
	}
	if repo.comment.AuthorName != "Ada" || repo.comment.AuthorEmail != "ada@example.com" || repo.comment.Content != "Hello" {
		t.Fatalf("stored comment = %#v", repo.comment)
	}
	if repo.comment.ParentID != nil || repo.comment.RootID != nil {
		t.Fatalf("root comment parent/root = %#v/%#v, want nil", repo.comment.ParentID, repo.comment.RootID)
	}
}

func TestSubmitCommentRejectsReplyDeeperThanThreeLevels(t *testing.T) {
	rootID := "root-1"
	parentID := "parent-1"
	repo := &fakeCommentRepo{comment: &domain.Comment{ID: parentID, RefType: domain.RefPost, RefID: "post-1", ParentID: strPtr("middle-1"), RootID: &rootID}}
	uc := NewCommentUsecase(repo)

	_, err := uc.SubmitComment(context.Background(), SubmitCommentRequest{RefType: "post", RefID: "post-1", ParentID: parentID, AuthorName: "Ada", AuthorEmail: "ada@example.com", Content: "too deep"}, ClientMeta{})
	if !errors.Is(err, domain.ErrCommentDepthExceeded) {
		t.Fatalf("SubmitComment() error = %v, want %v", err, domain.ErrCommentDepthExceeded)
	}
}

func TestSubmitCommentReplyUsesParentsRootAndIncrementsParent(t *testing.T) {
	rootID := "root-1"
	parentID := "parent-1"
	repo := &fakeCommentRepo{comment: &domain.Comment{ID: parentID, RefType: domain.RefPost, RefID: "post-1", RootID: &rootID}}
	uc := NewCommentUsecase(repo)

	comment, err := uc.SubmitComment(context.Background(), SubmitCommentRequest{RefType: "post", RefID: "post-1", ParentID: parentID, AuthorName: "Ada", AuthorEmail: "ada@example.com", Content: "reply"}, ClientMeta{})
	if err != nil {
		t.Fatalf("SubmitComment() error = %v", err)
	}
	if comment.ParentID == nil || *comment.ParentID != parentID || comment.RootID == nil || *comment.RootID != rootID {
		t.Fatalf("reply parent/root = %#v/%#v", comment.ParentID, comment.RootID)
	}
	if repo.incrementedID != parentID {
		t.Fatalf("incremented = %q, want %q", repo.incrementedID, parentID)
	}
}

func TestListPublicCommentsBuildsTreeAndFiltersApproved(t *testing.T) {
	now := time.Now().UTC()
	root := domain.Comment{ID: "root", RefType: domain.RefPost, RefID: "post-1", AuthorName: "Root", Content: "root", Status: domain.CommentApproved, CreatedAt: now}
	reply := domain.Comment{ID: "reply", RefType: domain.RefPost, RefID: "post-1", ParentID: &root.ID, RootID: &root.ID, AuthorName: "Reply", Content: "reply", Status: domain.CommentApproved, CreatedAt: now.Add(time.Second)}
	repo := &fakeCommentRepo{comments: []domain.Comment{reply, root}}
	uc := NewCommentUsecase(repo)

	nodes, err := uc.ListPublicComments(context.Background(), "post", "post-1")
	if err != nil {
		t.Fatalf("ListPublicComments() error = %v", err)
	}
	if len(nodes) != 1 || nodes[0].ID != "root" || len(nodes[0].Children) != 1 || nodes[0].Children[0].ID != "reply" {
		t.Fatalf("nodes = %#v", nodes)
	}
}

func TestModerateCommentRejectsPendingStatus(t *testing.T) {
	uc := NewCommentUsecase(&fakeCommentRepo{})

	err := uc.ModerateComment(context.Background(), "comment-1", "pending")
	if !errors.Is(err, domain.ErrInvalidComment) {
		t.Fatalf("ModerateComment() error = %v, want invalid comment", err)
	}
}

func strPtr(value string) *string { return &value }
