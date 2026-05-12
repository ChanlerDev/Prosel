package ai

import (
	"errors"
	"time"
)

const (
	RefTypePost = "post"
	RefTypeNote = "note"
	RefTypePage = "page"
)

var (
	ErrAINotFound    = errors.New("ai result not found")
	ErrAIUnavailable = errors.New("ai provider is not configured")
	ErrInvalidAIRef  = errors.New("invalid ai reference")
)

type AISummary struct {
	ID          string
	RefType     string
	RefID       string
	Language    string
	ContentHash string
	Summary     string
	Keywords    []string
	Provider    string
	Model       string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type AITranslation struct {
	ID              string
	RefType         string
	RefID           string
	SourceLanguage  string
	TargetLanguage  string
	ContentHash     string
	Title           string
	Summary         string
	ContentMarkdown string
	Provider        string
	Model           string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func ValidRefType(refType string) bool {
	return refType == RefTypePost || refType == RefTypeNote || refType == RefTypePage
}
