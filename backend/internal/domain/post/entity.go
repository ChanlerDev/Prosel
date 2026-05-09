package post

import (
	"errors"
	"time"
)

type PostStatus string

const (
	PostDraft     PostStatus = "draft"
	PostPublished PostStatus = "published"
	PostArchived  PostStatus = "archived"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrSlugTaken    = errors.New("post slug already exists")
	ErrInvalidPost  = errors.New("invalid post")
)

type Post struct {
	ID              string
	AuthorID        string
	CategoryID      *string
	Title           string
	Slug            string
	Excerpt         string
	ContentMarkdown string
	ContentText     string
	CoverImage      string
	Status          PostStatus
	Featured        bool
	PinnedAt        *time.Time
	PublishedAt     *time.Time
	SEOTitle        string
	SEODescription  string
	ViewCount       int64
	LikeCount       int64
	CommentCount    int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (s PostStatus) Valid() bool {
	return s == PostDraft || s == PostPublished || s == PostArchived
}
