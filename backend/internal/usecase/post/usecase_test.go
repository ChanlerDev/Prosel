package post

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/post"
)

type fakePostRepo struct {
	post       *domain.Post
	posts      []domain.Post
	pagination domain.Pagination
	err        error
	slugExists bool
	listed     domain.PostListFilter
	tagIDs     []string
}

func (r *fakePostRepo) Create(ctx context.Context, post *domain.Post) error {
	r.post = post
	return r.err
}
func (r *fakePostRepo) Update(ctx context.Context, post *domain.Post) error {
	r.post = post
	return r.err
}
func (r *fakePostRepo) Delete(ctx context.Context, id string) error { return r.err }
func (r *fakePostRepo) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	return r.post, r.err
}
func (r *fakePostRepo) GetBySlug(ctx context.Context, slug string, includeDraft bool) (*domain.Post, error) {
	return r.post, r.err
}
func (r *fakePostRepo) List(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error) {
	r.listed = filter
	return r.posts, r.pagination, r.err
}
func (r *fakePostRepo) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return r.slugExists, r.err
}
func (r *fakePostRepo) IncrementView(ctx context.Context, id string) error { return r.err }
func (r *fakePostRepo) SetStatus(ctx context.Context, id string, status domain.PostStatus, publishedAt *time.Time) error {
	if r.post != nil {
		r.post.Status = status
		r.post.PublishedAt = publishedAt
	}
	return r.err
}
func (r *fakePostRepo) ReplaceTags(ctx context.Context, postID string, tagIDs []string) error {
	r.tagIDs = tagIDs
	return r.err
}

func TestCreatePostGeneratesSlugAndPlainText(t *testing.T) {
	repo := &fakePostRepo{}
	uc := NewPostUsecase(repo)

	post, err := uc.CreatePost(context.Background(), CreatePostRequest{AuthorID: "user-1", Title: "Hello World!", ContentMarkdown: "# Hello\n\nThis is **markdown**."})
	if err != nil {
		t.Fatalf("CreatePost() error = %v", err)
	}
	if post.Slug != "hello-world" {
		t.Fatalf("Slug = %q, want hello-world", post.Slug)
	}
	if post.ContentText != "Hello This is markdown." {
		t.Fatalf("ContentText = %q", post.ContentText)
	}
	if post.Status != domain.PostDraft {
		t.Fatalf("Status = %q, want draft", post.Status)
	}
	if repo.post == nil || repo.post.ID == "" {
		t.Fatalf("post was not persisted: %#v", repo.post)
	}
}

func TestCreatePostRejectsDuplicateSlug(t *testing.T) {
	uc := NewPostUsecase(&fakePostRepo{slugExists: true})

	_, err := uc.CreatePost(context.Background(), CreatePostRequest{AuthorID: "user-1", Title: "Hello", Slug: "hello", ContentMarkdown: "body"})
	if !errors.Is(err, domain.ErrSlugTaken) {
		t.Fatalf("CreatePost() error = %v, want %v", err, domain.ErrSlugTaken)
	}
}

func TestListPublishedPostsForcesPublishedStatus(t *testing.T) {
	repo := &fakePostRepo{posts: []domain.Post{{ID: "post-1"}}, pagination: domain.Pagination{Page: 1, PerPage: 20, Total: 1, TotalPages: 1}}
	uc := NewPostUsecase(repo)

	_, _, err := uc.ListPublishedPosts(context.Background(), domain.PostListFilter{Page: 0, PerPage: 200})
	if err != nil {
		t.Fatalf("ListPublishedPosts() error = %v", err)
	}
	if repo.listed.Status == nil || *repo.listed.Status != domain.PostPublished {
		t.Fatalf("Status filter = %#v, want published", repo.listed.Status)
	}
	if repo.listed.Page != 1 || repo.listed.PerPage != 100 {
		t.Fatalf("pagination = page %d perPage %d", repo.listed.Page, repo.listed.PerPage)
	}
}

func TestPublishPostSetsPublishedAt(t *testing.T) {
	post := &domain.Post{ID: "post-1", Title: "Hello", Slug: "hello", Status: domain.PostDraft}
	uc := NewPostUsecase(&fakePostRepo{post: post})

	updated, err := uc.PublishPost(context.Background(), "post-1")
	if err != nil {
		t.Fatalf("PublishPost() error = %v", err)
	}
	if updated.Status != domain.PostPublished || updated.PublishedAt == nil {
		t.Fatalf("post = %#v, want published with publishedAt", updated)
	}
}
