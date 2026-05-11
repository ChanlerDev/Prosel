package page

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unicode"

	domain "github.com/chanler/prosel/backend/internal/domain/page"
)

type PageUsecase struct {
	pages   domain.PageRepository
	friends domain.FriendRepository
}

type PageRequest struct {
	AuthorID        string
	Title           string
	Slug            string
	Subtitle        string
	ContentMarkdown string
	Template        string
	Status          string
	SortOrder       int
	SEOTitle        string
	SEODescription  string
}

type FriendRequest struct {
	Name        string
	URL         string
	AvatarURL   string
	Description string
	Status      string
	SortOrder   int
}

func NewPageUsecase(pages domain.PageRepository, friends domain.FriendRepository) *PageUsecase {
	return &PageUsecase{pages: pages, friends: friends}
}

func (uc *PageUsecase) CreatePage(ctx context.Context, req PageRequest) (*domain.Page, error) {
	page, err := uc.pageFromRequest(ctx, nil, req)
	if err != nil {
		return nil, err
	}
	if err := uc.pages.Create(ctx, page); err != nil {
		return nil, err
	}
	return page, nil
}

func (uc *PageUsecase) UpdatePage(ctx context.Context, id string, req PageRequest) (*domain.Page, error) {
	page, err := uc.pages.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	updated, err := uc.pageFromRequest(ctx, page, req)
	if err != nil {
		return nil, err
	}
	if err := uc.pages.Update(ctx, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *PageUsecase) DeletePage(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidPage
	}
	return uc.pages.Delete(ctx, id)
}

func (uc *PageUsecase) GetPublicPage(ctx context.Context, slug string) (*domain.Page, error) {
	page, err := uc.pages.GetBySlug(ctx, strings.TrimSpace(slug), false)
	if err != nil {
		return nil, err
	}
	if page.Status != domain.PagePublished {
		return nil, domain.ErrPageNotFound
	}
	if err := uc.pages.IncrementView(ctx, page.ID); err != nil {
		return nil, err
	}
	page.ViewCount++
	return page, nil
}

func (uc *PageUsecase) GetAdminPage(ctx context.Context, id string) (*domain.Page, error) {
	return uc.pages.GetByID(ctx, id)
}

func (uc *PageUsecase) ListPublicPages(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error) {
	status := domain.PagePublished
	filter.Status = &status
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.pages.List(ctx, filter)
}

func (uc *PageUsecase) ListAdminPages(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	if filter.Status != nil && !filter.Status.Valid() {
		filter.Status = nil
	}
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.pages.List(ctx, filter)
}

func (uc *PageUsecase) CreateFriend(ctx context.Context, req FriendRequest) (*domain.Friend, error) {
	friend, err := uc.friendFromRequest(ctx, nil, req)
	if err != nil {
		return nil, err
	}
	if err := uc.friends.Create(ctx, friend); err != nil {
		return nil, err
	}
	return friend, nil
}

func (uc *PageUsecase) UpdateFriend(ctx context.Context, id string, req FriendRequest) (*domain.Friend, error) {
	friend, err := uc.friends.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	updated, err := uc.friendFromRequest(ctx, friend, req)
	if err != nil {
		return nil, err
	}
	if err := uc.friends.Update(ctx, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *PageUsecase) DeleteFriend(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidFriend
	}
	return uc.friends.Delete(ctx, id)
}

func (uc *PageUsecase) ListFriends(ctx context.Context) ([]domain.Friend, error) {
	return uc.friends.List(ctx, string(domain.FriendActive))
}

func (uc *PageUsecase) ListAdminFriends(ctx context.Context, status string) ([]domain.Friend, error) {
	friendStatus := domain.FriendStatus(strings.TrimSpace(status))
	if friendStatus != "" && !friendStatus.Valid() {
		friendStatus = ""
	}
	return uc.friends.List(ctx, string(friendStatus))
}

func (uc *PageUsecase) pageFromRequest(ctx context.Context, existing *domain.Page, req PageRequest) (*domain.Page, error) {
	title := strings.TrimSpace(req.Title)
	contentMarkdown := strings.TrimSpace(req.ContentMarkdown)
	if title == "" || contentMarkdown == "" || (existing == nil && strings.TrimSpace(req.AuthorID) == "") {
		return nil, domain.ErrInvalidPage
	}
	template := domain.PageTemplate(strings.TrimSpace(req.Template))
	if template == "" {
		template = domain.TemplateDefault
	}
	if !template.Valid() {
		return nil, domain.ErrInvalidPage
	}
	status := domain.PageStatus(strings.TrimSpace(req.Status))
	if status == "" {
		status = domain.PagePublished
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidPage
	}
	slug := normalizeSlug(req.Slug)
	if slug == "" {
		slug = normalizeSlug(title)
	}
	if slug == "" {
		return nil, domain.ErrInvalidPage
	}
	var excludeID *string
	if existing != nil {
		excludeID = &existing.ID
	}
	exists, err := uc.pages.SlugExists(ctx, slug, excludeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSlugTaken
	}
	now := time.Now().UTC()
	page := &domain.Page{ID: newID(), AuthorID: strings.TrimSpace(req.AuthorID), CreatedAt: now, UpdatedAt: now}
	if existing != nil {
		copy := *existing
		page = &copy
		page.UpdatedAt = now
	}
	page.Title = title
	page.Slug = slug
	page.Subtitle = strings.TrimSpace(req.Subtitle)
	page.ContentMarkdown = contentMarkdown
	page.ContentText = markdownText(contentMarkdown)
	page.Template = template
	page.Status = status
	page.SortOrder = req.SortOrder
	page.SEOTitle = strings.TrimSpace(req.SEOTitle)
	page.SEODescription = strings.TrimSpace(req.SEODescription)
	return page, nil
}

func (uc *PageUsecase) friendFromRequest(ctx context.Context, existing *domain.Friend, req FriendRequest) (*domain.Friend, error) {
	name := strings.TrimSpace(req.Name)
	friendURL := strings.TrimSpace(req.URL)
	if name == "" || friendURL == "" {
		return nil, domain.ErrInvalidFriend
	}
	parsedURL, err := url.ParseRequestURI(friendURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return nil, domain.ErrInvalidFriend
	}
	status := domain.FriendStatus(strings.TrimSpace(req.Status))
	if status == "" {
		status = domain.FriendActive
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidFriend
	}
	var excludeID *string
	if existing != nil {
		excludeID = &existing.ID
	}
	exists, err := uc.friends.URLExists(ctx, friendURL, excludeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrURLTaken
	}
	now := time.Now().UTC()
	friend := &domain.Friend{ID: newID(), CreatedAt: now, UpdatedAt: now}
	if existing != nil {
		copy := *existing
		friend = &copy
		friend.UpdatedAt = now
	}
	friend.Name = name
	friend.URL = friendURL
	friend.AvatarURL = strings.TrimSpace(req.AvatarURL)
	friend.Description = strings.TrimSpace(req.Description)
	friend.Status = status
	friend.SortOrder = req.SortOrder
	return friend, nil
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

func newID() string {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))[:32]
	}
	bytes[6] = (bytes[6] & 0x0f) | 0x40
	bytes[8] = (bytes[8] & 0x3f) | 0x80
	return hex.EncodeToString(bytes[:4]) + "-" + hex.EncodeToString(bytes[4:6]) + "-" + hex.EncodeToString(bytes[6:8]) + "-" + hex.EncodeToString(bytes[8:10]) + "-" + hex.EncodeToString(bytes[10:])
}
