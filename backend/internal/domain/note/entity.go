package note

import (
	"errors"
	"time"
)

type NoteStatus string

const (
	NoteDraft     NoteStatus = "draft"
	NotePublished NoteStatus = "published"
	NotePrivate   NoteStatus = "private"
	NoteArchived  NoteStatus = "archived"
)

var (
	ErrNoteNotFound = errors.New("note not found")
	ErrSlugTaken    = errors.New("note slug already exists")
	ErrInvalidNote  = errors.New("invalid note")
)

type Note struct {
	ID              string
	AuthorID        string
	Title           string
	Slug            string
	ContentMarkdown string
	ContentText     string
	Mood            string
	Weather         string
	Location        string
	Status          NoteStatus
	PinnedAt        *time.Time
	PublishedAt     *time.Time
	ViewCount       int64
	LikeCount       int64
	CommentCount    int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s NoteStatus) Valid() bool {
	return s == NoteDraft || s == NotePublished || s == NotePrivate || s == NoteArchived
}
