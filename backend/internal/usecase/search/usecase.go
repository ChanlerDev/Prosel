package search

import (
	"context"
	"strings"

	noteDomain "github.com/chanler/prosel/backend/internal/domain/note"
	pageDomain "github.com/chanler/prosel/backend/internal/domain/page"
	postDomain "github.com/chanler/prosel/backend/internal/domain/post"
	domain "github.com/chanler/prosel/backend/internal/domain/search"
)

type SearchUsecase struct {
	search domain.Repository
}

func NewSearchUsecase(search domain.Repository) *SearchUsecase {
	return &SearchUsecase{search: search}
}

func (uc *SearchUsecase) Search(ctx context.Context, query string, filter domain.SearchFilter) ([]domain.SearchResult, domain.Pagination, error) {
	query = strings.TrimSpace(query)
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	if filter.Type != "" && !filter.Type.Valid() {
		filter.Type = ""
	}
	if query == "" {
		return []domain.SearchResult{}, domain.NewPagination(filter.Page, filter.PerPage, 0), nil
	}
	return uc.search.Search(ctx, query, filter)
}

func (uc *SearchUsecase) IndexPost(ctx context.Context, post postDomain.Post) error {
	if post.Status != postDomain.PostPublished {
		return uc.search.DeleteDocument(ctx, string(domain.RefPost), post.ID)
	}
	return uc.search.UpsertDocument(ctx, &domain.SearchDocument{RefType: domain.RefPost, RefID: post.ID, Title: post.Title, Slug: post.Slug, Excerpt: post.Excerpt, SearchText: post.ContentText, Status: string(post.Status), PublishedAt: post.PublishedAt})
}

func (uc *SearchUsecase) IndexNote(ctx context.Context, note noteDomain.Note) error {
	if note.Status != noteDomain.NotePublished {
		return uc.search.DeleteDocument(ctx, string(domain.RefNote), note.ID)
	}
	title := strings.TrimSpace(note.Title)
	if title == "" {
		title = firstWords(note.ContentText, 12)
	}
	excerpt := firstWords(note.ContentText, 40)
	return uc.search.UpsertDocument(ctx, &domain.SearchDocument{RefType: domain.RefNote, RefID: note.ID, Title: title, Slug: note.Slug, Excerpt: excerpt, SearchText: note.ContentText, Status: string(note.Status), PublishedAt: note.PublishedAt})
}

func (uc *SearchUsecase) IndexPage(ctx context.Context, page pageDomain.Page) error {
	if page.Status != pageDomain.PagePublished {
		return uc.search.DeleteDocument(ctx, string(domain.RefPage), page.ID)
	}
	excerpt := strings.TrimSpace(page.Subtitle)
	if excerpt == "" {
		excerpt = firstWords(page.ContentText, 40)
	}
	return uc.search.UpsertDocument(ctx, &domain.SearchDocument{RefType: domain.RefPage, RefID: page.ID, Title: page.Title, Slug: page.Slug, Excerpt: excerpt, SearchText: page.ContentText, Status: string(page.Status), PublishedAt: nil})
}

func (uc *SearchUsecase) DeleteDocument(ctx context.Context, refType string, refID string) error {
	return uc.search.DeleteDocument(ctx, refType, refID)
}

func (uc *SearchUsecase) RebuildIndex(ctx context.Context) error {
	return uc.search.Rebuild(ctx)
}

func (uc *SearchUsecase) IndexStatus(ctx context.Context) (*domain.IndexStatus, error) {
	return uc.search.Status(ctx)
}

func firstWords(value string, limit int) string {
	words := strings.Fields(value)
	if len(words) <= limit {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:limit], " ")
}
