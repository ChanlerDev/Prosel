package search

import (
	"context"
	"testing"
	"time"

	noteDomain "github.com/chanler/prosel/backend/internal/domain/note"
	pageDomain "github.com/chanler/prosel/backend/internal/domain/page"
	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/search"
)

type fakeSearchRepo struct {
	doc        *domain.SearchDocument
	deleted    []string
	query      string
	filter     domain.SearchFilter
	rebuilt    bool
	results    []domain.SearchResult
	pagination domain.Pagination
	err        error
}

func (r *fakeSearchRepo) UpsertDocument(ctx context.Context, doc *domain.SearchDocument) error {
	r.doc = doc
	return r.err
}

func (r *fakeSearchRepo) DeleteDocument(ctx context.Context, refType string, refID string) error {
	r.deleted = append(r.deleted, refType+":"+refID)
	return r.err
}

func (r *fakeSearchRepo) Search(ctx context.Context, query string, filter domain.SearchFilter) ([]domain.SearchResult, domain.Pagination, error) {
	r.query = query
	r.filter = filter
	return r.results, r.pagination, r.err
}

func (r *fakeSearchRepo) Rebuild(ctx context.Context) error {
	r.rebuilt = true
	return r.err
}

func (r *fakeSearchRepo) Status(ctx context.Context) (*domain.IndexStatus, error) {
	return &domain.IndexStatus{Total: 1}, r.err
}

func TestSearchNormalizesQueryAndPagination(t *testing.T) {
	repo := &fakeSearchRepo{pagination: domain.Pagination{Page: 1, PerPage: 20, Total: 0, TotalPages: 0}}
	uc := NewSearchUsecase(repo)

	_, _, err := uc.Search(context.Background(), "  hello world  ", domain.SearchFilter{Type: "post", Page: 0, PerPage: 500})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if repo.query != "hello world" {
		t.Fatalf("query = %q, want normalized query", repo.query)
	}
	if repo.filter.Page != 1 || repo.filter.PerPage != 100 || repo.filter.Type != domain.RefPost {
		t.Fatalf("filter = %#v", repo.filter)
	}
}

func TestIndexPostUpsertsPublishedAndDeletesDraft(t *testing.T) {
	publishedAt := time.Now().UTC()
	repo := &fakeSearchRepo{}
	uc := NewSearchUsecase(repo)

	err := uc.IndexPost(context.Background(), postDomain.Post{ID: "post-1", Title: "Hello", Slug: "hello", Excerpt: "Intro", ContentText: "Body", Status: postDomain.PostPublished, PublishedAt: &publishedAt})
	if err != nil {
		t.Fatalf("IndexPost() error = %v", err)
	}
	if repo.doc == nil || repo.doc.RefType != domain.RefPost || repo.doc.RefID != "post-1" || repo.doc.Title != "Hello" || repo.doc.PublishedAt == nil {
		t.Fatalf("indexed doc = %#v", repo.doc)
	}

	err = uc.IndexPost(context.Background(), postDomain.Post{ID: "post-1", Status: postDomain.PostDraft})
	if err != nil {
		t.Fatalf("IndexPost(draft) error = %v", err)
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "post:post-1" {
		t.Fatalf("deleted = %#v", repo.deleted)
	}
}

func TestIndexNoteAndPageOnlyExposePublishedContent(t *testing.T) {
	repo := &fakeSearchRepo{}
	uc := NewSearchUsecase(repo)

	if err := uc.IndexNote(context.Background(), noteDomain.Note{ID: "note-1", ContentText: "Private", Status: noteDomain.NotePrivate}); err != nil {
		t.Fatalf("IndexNote(private) error = %v", err)
	}
	if len(repo.deleted) != 1 || repo.deleted[0] != "note:note-1" {
		t.Fatalf("deleted note = %#v", repo.deleted)
	}

	if err := uc.IndexPage(context.Background(), pageDomain.Page{ID: "page-1", Title: "About", Slug: "about", Subtitle: "Me", ContentText: "About body", Status: pageDomain.PagePublished}); err != nil {
		t.Fatalf("IndexPage() error = %v", err)
	}
	if repo.doc == nil || repo.doc.RefType != domain.RefPage || repo.doc.Excerpt != "Me" || repo.doc.SearchText != "About body" {
		t.Fatalf("indexed page = %#v", repo.doc)
	}
}
