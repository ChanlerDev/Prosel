package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/ai"
	"github.com/chanler/prosel/backend/internal/infrastructure/config"
)

type OpenAIClient struct {
	cfg    config.AIConfig
	client *http.Client
}

type chatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Choices []struct {
		Message chatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

type summarizeJSON struct {
	Summary  string   `json:"summary"`
	Keywords []string `json:"keywords"`
}

type translateJSON struct {
	Title           string `json:"title"`
	Summary         string `json:"summary"`
	ContentMarkdown string `json:"contentMarkdown"`
}

func NewOpenAIClient(cfg config.AIConfig) *OpenAIClient {
	return &OpenAIClient{cfg: cfg, client: &http.Client{Timeout: time.Duration(cfg.TimeoutSeconds) * time.Second}}
}

func (c *OpenAIClient) Configured() bool {
	return strings.EqualFold(c.cfg.Provider, "openai") && c.cfg.APIKey != "" && c.cfg.Model != ""
}

func (c *OpenAIClient) Summarize(ctx context.Context, input domain.SummarizeInput) (*domain.SummarizeOutput, error) {
	if !c.Configured() {
		return nil, domain.ErrAIUnavailable
	}
	content := "Title:\n" + input.Title + "\n\nMarkdown:\n" + input.ContentMarkdown
	message, err := c.chat(ctx, "You summarize blog posts. Return strict JSON only with keys summary and keywords. Summary language: "+input.Language+". Keep summary concise and keywords as short strings.", content)
	if err != nil {
		return nil, err
	}
	var parsed summarizeJSON
	if err := json.Unmarshal([]byte(extractJSON(message)), &parsed); err != nil {
		return nil, err
	}
	return &domain.SummarizeOutput{Summary: parsed.Summary, Keywords: parsed.Keywords, Provider: c.cfg.Provider, Model: c.cfg.Model}, nil
}

func (c *OpenAIClient) Translate(ctx context.Context, input domain.TranslateInput) (*domain.TranslateOutput, error) {
	if !c.Configured() {
		return nil, domain.ErrAIUnavailable
	}
	content := "Title:\n" + input.Title + "\n\nSummary:\n" + input.Summary + "\n\nMarkdown:\n" + input.ContentMarkdown
	message, err := c.chat(ctx, "You translate blog posts from "+input.SourceLanguage+" to "+input.TargetLanguage+". Preserve Markdown structure. Return strict JSON only with keys title, summary, and contentMarkdown.", content)
	if err != nil {
		return nil, err
	}
	var parsed translateJSON
	if err := json.Unmarshal([]byte(extractJSON(message)), &parsed); err != nil {
		return nil, err
	}
	return &domain.TranslateOutput{Title: parsed.Title, Summary: parsed.Summary, ContentMarkdown: parsed.ContentMarkdown, Provider: c.cfg.Provider, Model: c.cfg.Model}, nil
}

func (c *OpenAIClient) chat(ctx context.Context, system string, user string) (string, error) {
	body, err := json.Marshal(chatRequest{Model: c.cfg.Model, Temperature: 0.2, Messages: []chatMessage{{Role: "system", Content: system}, {Role: "user", Content: user}}})
	if err != nil {
		return "", err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.cfg.BaseURL, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	request.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	request.Header.Set("Content-Type", "application/json")
	response, err := c.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	var payload chatResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return "", err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if payload.Error != nil && payload.Error.Message != "" {
			return "", errors.New(payload.Error.Message)
		}
		return "", fmt.Errorf("openai request failed: status %d", response.StatusCode)
	}
	if len(payload.Choices) == 0 || strings.TrimSpace(payload.Choices[0].Message.Content) == "" {
		return "", errors.New("openai returned empty response")
	}
	return payload.Choices[0].Message.Content, nil
}

func extractJSON(value string) string {
	value = strings.TrimSpace(value)
	value = strings.TrimPrefix(value, "```json")
	value = strings.TrimPrefix(value, "```")
	value = strings.TrimSuffix(value, "```")
	return strings.TrimSpace(value)
}
