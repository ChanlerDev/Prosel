package note

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strings"
	"time"
	"unicode"

	domain "github.com/chanler/prosel/backend/internal/domain/note"
)

type NoteUsecase struct {
	notes domain.Repository
}

type NoteRequest struct {
	AuthorID        string
	Title           string
	Slug            string
	ContentMarkdown string
	Mood            string
	Weather         string
	Location        string
	Status          string
}

func NewNoteUsecase(notes domain.Repository) *NoteUsecase {
	return &NoteUsecase{notes: notes}
}

func (uc *NoteUsecase) CreateNote(ctx context.Context, req NoteRequest) (*domain.Note, error) {
	note, err := uc.noteFromRequest(ctx, nil, req)
	if err != nil {
		return nil, err
	}
	if note.Status == domain.NotePublished && note.PublishedAt == nil {
		publishedAt := note.CreatedAt
		note.PublishedAt = &publishedAt
	}
	if err := uc.notes.Create(ctx, note); err != nil {
		return nil, err
	}
	return note, nil
}

func (uc *NoteUsecase) UpdateNote(ctx context.Context, id string, req NoteRequest) (*domain.Note, error) {
	note, err := uc.notes.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	updated, err := uc.noteFromRequest(ctx, note, req)
	if err != nil {
		return nil, err
	}
	if note.Status != domain.NotePublished && updated.Status == domain.NotePublished && updated.PublishedAt == nil {
		publishedAt := updated.UpdatedAt
		updated.PublishedAt = &publishedAt
	}
	if updated.Status != domain.NotePublished {
		updated.PublishedAt = nil
	}
	if err := uc.notes.Update(ctx, updated); err != nil {
		return nil, err
	}
	return updated, nil
}

func (uc *NoteUsecase) DeleteNote(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidNote
	}
	return uc.notes.Delete(ctx, id)
}

func (uc *NoteUsecase) GetPublicNote(ctx context.Context, slug string) (*domain.Note, error) {
	note, err := uc.notes.GetBySlug(ctx, strings.TrimSpace(slug))
	if err != nil {
		return nil, err
	}
	if note.Status != domain.NotePublished {
		return nil, domain.ErrNoteNotFound
	}
	if err := uc.notes.IncrementView(ctx, note.ID); err != nil {
		return nil, err
	}
	note.ViewCount++
	return note, nil
}

func (uc *NoteUsecase) GetAdminNote(ctx context.Context, id string) (*domain.Note, error) {
	return uc.notes.GetByID(ctx, id)
}

func (uc *NoteUsecase) ListPublicNotes(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	status := domain.NotePublished
	filter.Status = &status
	filter.Search = strings.TrimSpace(filter.Search)
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.notes.ListPublic(ctx, filter)
}

func (uc *NoteUsecase) ListAdminNotes(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	if filter.Status != nil && !filter.Status.Valid() {
		filter.Status = nil
	}
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.notes.ListAdmin(ctx, filter)
}

func (uc *NoteUsecase) PinNote(ctx context.Context, id string, pinned bool) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidNote
	}
	var pinnedAt *time.Time
	if pinned {
		now := time.Now().UTC()
		pinnedAt = &now
	}
	return uc.notes.SetPinned(ctx, id, pinnedAt)
}

func (uc *NoteUsecase) noteFromRequest(ctx context.Context, existing *domain.Note, req NoteRequest) (*domain.Note, error) {
	contentMarkdown := strings.TrimSpace(req.ContentMarkdown)
	if strings.TrimSpace(req.AuthorID) == "" && existing == nil || contentMarkdown == "" {
		return nil, domain.ErrInvalidNote
	}
	status := domain.NoteStatus(strings.TrimSpace(req.Status))
	if status == "" {
		status = domain.NotePublished
	}
	if !status.Valid() {
		return nil, domain.ErrInvalidNote
	}
	baseForSlug := req.Title
	if strings.TrimSpace(baseForSlug) == "" {
		baseForSlug = contentMarkdown
	}
	slug := normalizeSlug(req.Slug)
	if slug == "" {
		slug = normalizeSlug(baseForSlug)
	}
	if slug == "" {
		return nil, domain.ErrInvalidNote
	}
	var excludeID *string
	if existing != nil {
		excludeID = &existing.ID
	}
	exists, err := uc.notes.SlugExists(ctx, slug, excludeID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrSlugTaken
	}
	now := time.Now().UTC()
	note := &domain.Note{ID: newID(), AuthorID: strings.TrimSpace(req.AuthorID), CreatedAt: now, UpdatedAt: now}
	if existing != nil {
		copy := *existing
		note = &copy
		note.UpdatedAt = now
	}
	note.Title = strings.TrimSpace(req.Title)
	note.Slug = slug
	note.ContentMarkdown = contentMarkdown
	note.ContentText = markdownText(contentMarkdown)
	note.Mood = strings.TrimSpace(req.Mood)
	note.Weather = strings.TrimSpace(req.Weather)
	note.Location = strings.TrimSpace(req.Location)
	note.Status = status
	return note, nil
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
