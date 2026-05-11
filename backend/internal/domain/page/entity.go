package page

import (
	"errors"
	"time"
)

type PageTemplate string
type PageStatus string
type FriendStatus string

const (
	TemplateDefault  PageTemplate = "default"
	TemplateAbout    PageTemplate = "about"
	TemplateFriends  PageTemplate = "friends"
	TemplateProjects PageTemplate = "projects"

	PageDraft     PageStatus = "draft"
	PagePublished PageStatus = "published"
	PageArchived  PageStatus = "archived"

	FriendActive  FriendStatus = "active"
	FriendPending FriendStatus = "pending"
	FriendHidden  FriendStatus = "hidden"
)

var (
	ErrPageNotFound   = errors.New("page not found")
	ErrFriendNotFound = errors.New("friend not found")
	ErrSlugTaken      = errors.New("page slug already exists")
	ErrURLTaken       = errors.New("friend url already exists")
	ErrInvalidPage    = errors.New("invalid page")
	ErrInvalidFriend  = errors.New("invalid friend")
)

type Page struct {
	ID              string
	AuthorID        string
	Title           string
	Slug            string
	Subtitle        string
	ContentMarkdown string
	ContentText     string
	Template        PageTemplate
	Status          PageStatus
	SortOrder       int
	SEOTitle        string
	SEODescription  string
	ViewCount       int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Friend struct {
	ID          string
	Name        string
	URL         string
	AvatarURL   string
	Description string
	Status      FriendStatus
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t PageTemplate) Valid() bool {
	return t == TemplateDefault || t == TemplateAbout || t == TemplateFriends || t == TemplateProjects
}

func (s PageStatus) Valid() bool {
	return s == PageDraft || s == PagePublished || s == PageArchived
}

func (s FriendStatus) Valid() bool {
	return s == FriendActive || s == FriendPending || s == FriendHidden
}
