package ai

import "context"

type SummarizeInput struct {
	Title           string
	ContentMarkdown string
	Language        string
}

type SummarizeOutput struct {
	Summary  string
	Keywords []string
	Provider string
	Model    string
}

type TranslateInput struct {
	Title           string
	Summary         string
	ContentMarkdown string
	SourceLanguage  string
	TargetLanguage  string
}

type TranslateOutput struct {
	Title           string
	Summary         string
	ContentMarkdown string
	Provider        string
	Model           string
}

type Client interface {
	Summarize(ctx context.Context, input SummarizeInput) (*SummarizeOutput, error)
	Translate(ctx context.Context, input TranslateInput) (*TranslateOutput, error)
	Configured() bool
}
