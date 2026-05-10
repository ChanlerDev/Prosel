package note

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/note"
)

type fakeNoteRepo struct {
	note       *domain.Note
	notes      []domain.Note
	pagination domain.Pagination
	err        error
	slugExists bool
	listed     domain.NoteFilter
	pinnedID   string
	pinnedAt   *time.Time
	deletedID  string
}

func (r *fakeNoteRepo) Create(ctx context.Context, note *domain.Note) error {
	r.note = note
	return r.err
}
func (r *fakeNoteRepo) Update(ctx context.Context, note *domain.Note) error {
	r.note = note
	return r.err
}
func (r *fakeNoteRepo) Delete(ctx context.Context, id string) error { r.deletedID = id; return r.err }
func (r *fakeNoteRepo) GetByID(ctx context.Context, id string) (*domain.Note, error) {
	return r.note, r.err
}
func (r *fakeNoteRepo) GetBySlug(ctx context.Context, slug string) (*domain.Note, error) {
	return r.note, r.err
}
func (r *fakeNoteRepo) ListPublic(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	r.listed = filter
	return r.notes, r.pagination, r.err
}
func (r *fakeNoteRepo) ListAdmin(ctx context.Context, filter domain.NoteFilter) ([]domain.Note, domain.Pagination, error) {
	r.listed = filter
	return r.notes, r.pagination, r.err
}
func (r *fakeNoteRepo) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return r.slugExists, r.err
}
func (r *fakeNoteRepo) IncrementView(ctx context.Context, id string) error { return r.err }
func (r *fakeNoteRepo) SetPinned(ctx context.Context, id string, pinnedAt *time.Time) error {
	r.pinnedID = id
	r.pinnedAt = pinnedAt
	return r.err
}

func TestCreateNoteDefaultsToPublishedAndGeneratesSlug(t *testing.T) {
	repo := &fakeNoteRepo{}
	uc := NewNoteUsecase(repo)

	note, err := uc.CreateNote(context.Background(), NoteRequest{AuthorID: "user-1", Title: "Morning Walk", ContentMarkdown: "A **quiet** note", Mood: "calm", Weather: "sunny", Location: "park"})
	if err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}
	if note.Slug != "morning-walk" || note.Status != domain.NotePublished || note.PublishedAt == nil {
		t.Fatalf("note = %#v, want generated slug and published status", note)
	}
	if note.ContentText != "A quiet note" || repo.note.Mood != "calm" || repo.note.Location != "park" {
		t.Fatalf("stored note = %#v", repo.note)
	}
}

func TestCreateNoteAllowsContentOnlySlugFallback(t *testing.T) {
	repo := &fakeNoteRepo{}
	uc := NewNoteUsecase(repo)

	note, err := uc.CreateNote(context.Background(), NoteRequest{AuthorID: "user-1", ContentMarkdown: "Small note without title"})
	if err != nil {
		t.Fatalf("CreateNote() error = %v", err)
	}
	if note.Slug == "" || note.Title != "" {
		t.Fatalf("note = %#v, want generated slug without title", note)
	}
}

func TestCreateNoteRejectsDuplicateSlug(t *testing.T) {
	uc := NewNoteUsecase(&fakeNoteRepo{slugExists: true})

	_, err := uc.CreateNote(context.Background(), NoteRequest{AuthorID: "user-1", Title: "Hello", Slug: "hello", ContentMarkdown: "body"})
	if !errors.Is(err, domain.ErrSlugTaken) {
		t.Fatalf("CreateNote() error = %v, want %v", err, domain.ErrSlugTaken)
	}
}

func TestListPublicNotesForcesPublishedStatus(t *testing.T) {
	repo := &fakeNoteRepo{notes: []domain.Note{{ID: "note-1"}}, pagination: domain.Pagination{Page: 1, PerPage: 20, Total: 1, TotalPages: 1}}
	uc := NewNoteUsecase(repo)

	_, _, err := uc.ListPublicNotes(context.Background(), domain.NoteFilter{Page: 0, PerPage: 200})
	if err != nil {
		t.Fatalf("ListPublicNotes() error = %v", err)
	}
	if repo.listed.Status == nil || *repo.listed.Status != domain.NotePublished {
		t.Fatalf("Status filter = %#v, want published", repo.listed.Status)
	}
	if repo.listed.Page != 1 || repo.listed.PerPage != 100 {
		t.Fatalf("pagination = page %d perPage %d", repo.listed.Page, repo.listed.PerPage)
	}
}

func TestPinNoteSetsAndClearsPinnedAt(t *testing.T) {
	repo := &fakeNoteRepo{}
	uc := NewNoteUsecase(repo)

	if err := uc.PinNote(context.Background(), "note-1", true); err != nil {
		t.Fatalf("PinNote(true) error = %v", err)
	}
	if repo.pinnedID != "note-1" || repo.pinnedAt == nil {
		t.Fatalf("pin = %q/%#v, want pinnedAt", repo.pinnedID, repo.pinnedAt)
	}
	if err := uc.PinNote(context.Background(), "note-1", false); err != nil {
		t.Fatalf("PinNote(false) error = %v", err)
	}
	if repo.pinnedAt != nil {
		t.Fatalf("pinnedAt = %#v, want nil", repo.pinnedAt)
	}
}
