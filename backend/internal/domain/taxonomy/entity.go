package taxonomy

import (
	"errors"
	"time"
)

var (
	ErrTaxonomyNotFound = errors.New("taxonomy item not found")
	ErrSlugTaken        = errors.New("taxonomy slug already exists")
	ErrInvalidTaxonomy  = errors.New("invalid taxonomy item")
)

type Category struct {
	ID          string
	ParentID    *string
	Name        string
	Slug        string
	Description string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CategoryNode struct {
	Category
	PostCount int64
	Children  []CategoryNode
}

type Tag struct {
	ID          string
	Name        string
	Slug        string
	Color       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TagWithCount struct {
	Tag
	PostCount int64
}

type Topic struct {
	ID          string
	Name        string
	Slug        string
	Description string
	CoverImage  string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TopicItem struct {
	TopicID   string
	RefType   string
	RefID     string
	Title     string
	Slug      string
	SortOrder int
}

type TopicDetail struct {
	Topic
	Items []TopicItem
}
