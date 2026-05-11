package post

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
	"unicode"

	domain "github.com/chanler/prosel/backend/internal/domain/post"
)

type SearchIndexer interface {
	IndexPost(ctx context.Context, post domain.Post) error
	DeleteDocument(ctx context.Context, refType string, refID string) error
}

type PostUsecase struct {
	posts  domain.Repository
	search SearchIndexer
}

type CreatePostRequest struct {
	AuthorID        string
	CategoryID      *string
	TagIDs          []string
	Title           string
	Slug            string
	Excerpt         string
	ContentMarkdown string
	CoverImage      string
	Featured        bool
	SEOTitle        string
	SEODescription  string
}

type UpdatePostRequest struct {
	CategoryID      *string
	TagIDs          []string
	Title           string
	Slug            string
	Excerpt         string
	ContentMarkdown string
	CoverImage      string
	Featured        bool
	SEOTitle        string
	SEODescription  string
}

func NewPostUsecase(posts domain.Repository, search ...SearchIndexer) *PostUsecase {
	uc := &PostUsecase{posts: posts}
	if len(search) > 0 {
		uc.search = search[0]
	}
	return uc
}

func (uc *PostUsecase) CreatePost(ctx context.Context, req CreatePostRequest) (*domain.Post, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" || strings.TrimSpace(req.AuthorID) == "" {
		return nil, domain.ErrInvalidPost
	}
	slug := normalizeSlug(req.Slug)
	if slug == "" {
		slug = normalizeSlug(title)
	}
	if slug == "" {
		return nil, domain.ErrInvalidPost
	}
	exists, err := uc.posts.SlugExists(ctx, slug, nil)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSlugTaken
	}

	now := time.Now().UTC()
	post := &domain.Post{
		ID:              newID(),
		AuthorID:        strings.TrimSpace(req.AuthorID),
		CategoryID:      cleanOptional(req.CategoryID),
		Title:           title,
		Slug:            slug,
		Excerpt:         strings.TrimSpace(req.Excerpt),
		ContentMarkdown: strings.TrimSpace(req.ContentMarkdown),
		ContentText:     markdownText(req.ContentMarkdown),
		CoverImage:      strings.TrimSpace(req.CoverImage),
		Status:          domain.PostDraft,
		Featured:        req.Featured,
		SEOTitle:        strings.TrimSpace(req.SEOTitle),
		SEODescription:  strings.TrimSpace(req.SEODescription),
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := uc.posts.Create(ctx, post); err != nil {
		return nil, err
	}
	if err := uc.posts.ReplaceTags(ctx, post.ID, req.TagIDs); err != nil {
		return nil, err
	}
	return post, nil
}

func (uc *PostUsecase) UpdatePost(ctx context.Context, id string, req UpdatePostRequest) (*domain.Post, error) {
	post, err := uc.posts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, domain.ErrInvalidPost
	}
	slug := normalizeSlug(req.Slug)
	if slug == "" {
		slug = normalizeSlug(title)
	}
	exists, err := uc.posts.SlugExists(ctx, slug, &id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSlugTaken
	}
	post.CategoryID = cleanOptional(req.CategoryID)
	post.Title = title
	post.Slug = slug
	post.Excerpt = strings.TrimSpace(req.Excerpt)
	post.ContentMarkdown = strings.TrimSpace(req.ContentMarkdown)
	post.ContentText = markdownText(req.ContentMarkdown)
	post.CoverImage = strings.TrimSpace(req.CoverImage)
	post.Featured = req.Featured
	post.SEOTitle = strings.TrimSpace(req.SEOTitle)
	post.SEODescription = strings.TrimSpace(req.SEODescription)
	post.UpdatedAt = time.Now().UTC()
	if err := uc.posts.Update(ctx, post); err != nil {
		return nil, err
	}
	if err := uc.posts.ReplaceTags(ctx, post.ID, req.TagIDs); err != nil {
		return nil, err
	}
	if err := uc.indexPost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (uc *PostUsecase) PublishPost(ctx context.Context, id string) (*domain.Post, error) {
	post, err := uc.posts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	publishedAt := time.Now().UTC()
	if post.PublishedAt != nil {
		publishedAt = *post.PublishedAt
	}
	if err := uc.posts.SetStatus(ctx, id, domain.PostPublished, &publishedAt); err != nil {
		return nil, err
	}
	post.Status = domain.PostPublished
	post.PublishedAt = &publishedAt
	post.UpdatedAt = time.Now().UTC()
	if err := uc.indexPost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (uc *PostUsecase) UnpublishPost(ctx context.Context, id string) (*domain.Post, error) {
	post, err := uc.posts.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := uc.posts.SetStatus(ctx, id, domain.PostDraft, nil); err != nil {
		return nil, err
	}
	post.Status = domain.PostDraft
	post.PublishedAt = nil
	post.UpdatedAt = time.Now().UTC()
	if err := uc.indexPost(ctx, post); err != nil {
		return nil, err
	}
	return post, nil
}

func (uc *PostUsecase) DeletePost(ctx context.Context, id string) error {
	if err := uc.posts.Delete(ctx, id); err != nil {
		return err
	}
	if uc.search != nil {
		return uc.search.DeleteDocument(ctx, "post", id)
	}
	return nil
}

func (uc *PostUsecase) GetAdminPost(ctx context.Context, id string) (*domain.Post, error) {
	return uc.posts.GetByID(ctx, id)
}

func (uc *PostUsecase) GetPublishedPost(ctx context.Context, slug string) (*domain.Post, error) {
	post, err := uc.posts.GetBySlug(ctx, slug, false)
	if err != nil {
		return nil, err
	}
	if err := uc.posts.IncrementView(ctx, post.ID); err != nil {
		return nil, err
	}
	post.ViewCount++
	return post, nil
}

func (uc *PostUsecase) ListPublishedPosts(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error) {
	status := domain.PostPublished
	filter.Status = &status
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.posts.List(ctx, filter)
}

func (uc *PostUsecase) ListAdminPosts(ctx context.Context, filter domain.PostListFilter) ([]domain.Post, domain.Pagination, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.posts.List(ctx, filter)
}

func (uc *PostUsecase) indexPost(ctx context.Context, post *domain.Post) error {
	if uc.search == nil {
		return nil
	}
	return uc.search.IndexPost(ctx, *post)
}

var markdownMarkup = regexp.MustCompile(`[` + "`" + `*_#>\[\]()!~|{}+-]+`)
var spacedPunctuation = regexp.MustCompile(`\s+([.,!?;:])`)

func markdownText(value string) string {
	withoutMarkup := markdownMarkup.ReplaceAllString(value, " ")
	return spacedPunctuation.ReplaceAllString(strings.Join(strings.Fields(withoutMarkup), " "), "$1")
}

func normalizeSlug(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var builder strings.Builder
	lastDash := false
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && builder.Len() > 0 {
			builder.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(builder.String(), "-")
}

func cleanOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
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
