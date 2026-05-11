package page

import (
	"context"
	"errors"
	"testing"

	domain "github.com/chanler/prosel/backend/internal/domain/page"
)

type fakePageRepo struct {
	page       *domain.Page
	pages      []domain.Page
	pagination domain.Pagination
	err        error
	slugExists bool
	listed     domain.PageFilter
	deletedID  string
}

func (r *fakePageRepo) Create(ctx context.Context, page *domain.Page) error {
	r.page = page
	return r.err
}
func (r *fakePageRepo) Update(ctx context.Context, page *domain.Page) error {
	r.page = page
	return r.err
}
func (r *fakePageRepo) Delete(ctx context.Context, id string) error { r.deletedID = id; return r.err }
func (r *fakePageRepo) GetByID(ctx context.Context, id string) (*domain.Page, error) {
	return r.page, r.err
}
func (r *fakePageRepo) GetBySlug(ctx context.Context, slug string, includeDraft bool) (*domain.Page, error) {
	return r.page, r.err
}
func (r *fakePageRepo) List(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error) {
	r.listed = filter
	return r.pages, r.pagination, r.err
}
func (r *fakePageRepo) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return r.slugExists, r.err
}
func (r *fakePageRepo) IncrementView(ctx context.Context, id string) error { return r.err }

type fakeFriendRepo struct {
	friend    *domain.Friend
	friends   []domain.Friend
	err       error
	status    string
	deletedID string
}

func (r *fakeFriendRepo) Create(ctx context.Context, friend *domain.Friend) error {
	r.friend = friend
	return r.err
}
func (r *fakeFriendRepo) Update(ctx context.Context, friend *domain.Friend) error {
	r.friend = friend
	return r.err
}
func (r *fakeFriendRepo) Delete(ctx context.Context, id string) error { r.deletedID = id; return r.err }
func (r *fakeFriendRepo) GetByID(ctx context.Context, id string) (*domain.Friend, error) {
	return r.friend, r.err
}
func (r *fakeFriendRepo) URLExists(ctx context.Context, url string, excludeID *string) (bool, error) {
	return false, r.err
}
func (r *fakeFriendRepo) List(ctx context.Context, status string) ([]domain.Friend, error) {
	r.status = status
	return r.friends, r.err
}

func TestCreatePageGeneratesSlugAndPlainText(t *testing.T) {
	pages := &fakePageRepo{}
	uc := NewPageUsecase(pages, &fakeFriendRepo{})

	page, err := uc.CreatePage(context.Background(), PageRequest{AuthorID: "user-1", Title: "About Me", ContentMarkdown: "# Hello\nI **write**.", Template: "about", Status: "published"})
	if err != nil {
		t.Fatalf("CreatePage() error = %v", err)
	}
	if page.Slug != "about-me" || page.ContentText != "Hello I write." || page.Template != domain.TemplateAbout || page.Status != domain.PagePublished {
		t.Fatalf("page = %#v", page)
	}
	if pages.page == nil || pages.page.ID == "" {
		t.Fatalf("page was not persisted: %#v", pages.page)
	}
}

func TestCreatePageRejectsDuplicateSlug(t *testing.T) {
	uc := NewPageUsecase(&fakePageRepo{slugExists: true}, &fakeFriendRepo{})

	_, err := uc.CreatePage(context.Background(), PageRequest{AuthorID: "user-1", Title: "About", Slug: "about", ContentMarkdown: "body"})
	if !errors.Is(err, domain.ErrSlugTaken) {
		t.Fatalf("CreatePage() error = %v, want %v", err, domain.ErrSlugTaken)
	}
}

func TestListPublicPagesForcesPublishedStatus(t *testing.T) {
	pages := &fakePageRepo{pages: []domain.Page{{ID: "page-1"}}, pagination: domain.Pagination{Page: 1, PerPage: 20, Total: 1, TotalPages: 1}}
	uc := NewPageUsecase(pages, &fakeFriendRepo{})

	_, _, err := uc.ListPublicPages(context.Background(), domain.PageFilter{Page: 0, PerPage: 200})
	if err != nil {
		t.Fatalf("ListPublicPages() error = %v", err)
	}
	if pages.listed.Status == nil || *pages.listed.Status != domain.PagePublished {
		t.Fatalf("Status filter = %#v, want published", pages.listed.Status)
	}
	if pages.listed.Page != 1 || pages.listed.PerPage != 100 {
		t.Fatalf("pagination = page %d perPage %d", pages.listed.Page, pages.listed.PerPage)
	}
}

func TestListFriendsReturnsOnlyActivePublicFriends(t *testing.T) {
	friends := &fakeFriendRepo{friends: []domain.Friend{{ID: "friend-1", Status: domain.FriendActive}}}
	uc := NewPageUsecase(&fakePageRepo{}, friends)

	_, err := uc.ListFriends(context.Background())
	if err != nil {
		t.Fatalf("ListFriends() error = %v", err)
	}
	if friends.status != string(domain.FriendActive) {
		t.Fatalf("status = %q, want active", friends.status)
	}
}

func TestCreateFriendDefaultsToActive(t *testing.T) {
	friends := &fakeFriendRepo{}
	uc := NewPageUsecase(&fakePageRepo{}, friends)

	friend, err := uc.CreateFriend(context.Background(), FriendRequest{Name: "Ada", URL: "https://example.com", Description: "Writes"})
	if err != nil {
		t.Fatalf("CreateFriend() error = %v", err)
	}
	if friend.Status != domain.FriendActive || friends.friend.Name != "Ada" {
		t.Fatalf("friend = %#v stored = %#v", friend, friends.friend)
	}
}
