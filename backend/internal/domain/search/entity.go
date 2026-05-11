package search

import (
	"errors"
	"time"
)

type RefType string

const (
	RefPost RefType = "post"
	RefNote RefType = "note"
	RefPage RefType = "page"
)

var ErrInvalidSearch = errors.New("invalid search")

type SearchDocument struct {
	ID          string
	RefType     RefType
	RefID       string
	Title       string
	Slug        string
	Excerpt     string
	SearchText  string
	Status      string
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SearchResult struct {
	RefType RefType
	RefID   string
	Title   string
	Slug    string
	Excerpt string
	Rank    float64
}

type IndexStatus struct {
	Total     int64
	Posts     int64
	Notes     int64
	Pages     int64
	UpdatedAt *time.Time
}

func (t RefType) Valid() bool {
	return t == RefPost || t == RefNote || t == RefPage
}
