package taxonomy

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
	"unicode"

	domain "github.com/chanler/prosel/backend/internal/domain/taxonomy"
)

type TaxonomyUsecase struct {
	categories domain.CategoryRepository
	tags       domain.TagRepository
	topics     domain.TopicRepository
}

type CategoryRequest struct {
	ParentID    *string
	Name        string
	Slug        string
	Description string
	SortOrder   int
}

type TagRequest struct {
	Name        string
	Slug        string
	Color       string
	Description string
}

type TopicRequest struct {
	Name        string
	Slug        string
	Description string
	CoverImage  string
	SortOrder   int
	Items       []domain.TopicItem
}

func NewTaxonomyUsecase(categories domain.CategoryRepository, tags domain.TagRepository, topics domain.TopicRepository) *TaxonomyUsecase {
	return &TaxonomyUsecase{categories: categories, tags: tags, topics: topics}
}

func (uc *TaxonomyUsecase) ListCategories(ctx context.Context) ([]domain.CategoryNode, error) {
	return uc.categories.ListTree(ctx)
}

func (uc *TaxonomyUsecase) CreateCategory(ctx context.Context, req CategoryRequest) (*domain.Category, error) {
	slug, err := uc.validateSlug(ctx, req.Name, req.Slug, nil, uc.categories.SlugExists)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidTaxonomy
	}
	now := time.Now().UTC()
	category := &domain.Category{ID: newID(), ParentID: cleanOptional(req.ParentID), Name: name, Slug: slug, Description: strings.TrimSpace(req.Description), SortOrder: req.SortOrder, CreatedAt: now, UpdatedAt: now}
	if err := uc.categories.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *TaxonomyUsecase) UpdateCategory(ctx context.Context, id string, req CategoryRequest) (*domain.Category, error) {
	category, err := uc.categories.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	slug, err := uc.validateSlug(ctx, req.Name, req.Slug, &id, uc.categories.SlugExists)
	if err != nil {
		return nil, err
	}
	category.ParentID = cleanOptional(req.ParentID)
	category.Name = strings.TrimSpace(req.Name)
	category.Slug = slug
	category.Description = strings.TrimSpace(req.Description)
	category.SortOrder = req.SortOrder
	category.UpdatedAt = time.Now().UTC()
	if err := uc.categories.Update(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

func (uc *TaxonomyUsecase) DeleteCategory(ctx context.Context, id string) error {
	return uc.categories.Delete(ctx, id)
}

func (uc *TaxonomyUsecase) ListTags(ctx context.Context) ([]domain.TagWithCount, error) {
	return uc.tags.List(ctx)
}

func (uc *TaxonomyUsecase) CreateTag(ctx context.Context, req TagRequest) (*domain.Tag, error) {
	slug, err := uc.validateSlug(ctx, req.Name, req.Slug, nil, uc.tags.SlugExists)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidTaxonomy
	}
	now := time.Now().UTC()
	tag := &domain.Tag{ID: newID(), Name: name, Slug: slug, Color: strings.TrimSpace(req.Color), Description: strings.TrimSpace(req.Description), CreatedAt: now, UpdatedAt: now}
	if err := uc.tags.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (uc *TaxonomyUsecase) UpdateTag(ctx context.Context, id string, req TagRequest) (*domain.Tag, error) {
	tag, err := uc.tags.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	slug, err := uc.validateSlug(ctx, req.Name, req.Slug, &id, uc.tags.SlugExists)
	if err != nil {
		return nil, err
	}
	tag.Name = strings.TrimSpace(req.Name)
	tag.Slug = slug
	tag.Color = strings.TrimSpace(req.Color)
	tag.Description = strings.TrimSpace(req.Description)
	tag.UpdatedAt = time.Now().UTC()
	if err := uc.tags.Update(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (uc *TaxonomyUsecase) DeleteTag(ctx context.Context, id string) error {
	return uc.tags.Delete(ctx, id)
}

func (uc *TaxonomyUsecase) ReplacePostTags(ctx context.Context, postID string, tagIDs []string) error {
	return uc.tags.ReplacePostTags(ctx, postID, tagIDs)
}

func (uc *TaxonomyUsecase) ListTopics(ctx context.Context) ([]domain.Topic, error) {
	return uc.topics.List(ctx)
}

func (uc *TaxonomyUsecase) GetTopic(ctx context.Context, slug string) (*domain.TopicDetail, error) {
	topic, err := uc.topics.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	items, err := uc.topics.ListItems(ctx, topic.ID)
	if err != nil {
		return nil, err
	}
	return &domain.TopicDetail{Topic: *topic, Items: items}, nil
}

func (uc *TaxonomyUsecase) CreateTopic(ctx context.Context, req TopicRequest) (*domain.Topic, error) {
	slug, err := uc.validateSlug(ctx, req.Name, req.Slug, nil, uc.topics.SlugExists)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, domain.ErrInvalidTaxonomy
	}
	now := time.Now().UTC()
	topic := &domain.Topic{ID: newID(), Name: name, Slug: slug, Description: strings.TrimSpace(req.Description), CoverImage: strings.TrimSpace(req.CoverImage), SortOrder: req.SortOrder, CreatedAt: now, UpdatedAt: now}
	if err := uc.topics.Create(ctx, topic); err != nil {
		return nil, err
	}
	return topic, uc.topics.ReplaceItems(ctx, topic.ID, req.Items)
}

func (uc *TaxonomyUsecase) UpdateTopic(ctx context.Context, id string, req TopicRequest) (*domain.Topic, error) {
	topic, err := uc.topics.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	slug, err := uc.validateSlug(ctx, req.Name, req.Slug, &id, uc.topics.SlugExists)
	if err != nil {
		return nil, err
	}
	topic.Name = strings.TrimSpace(req.Name)
	topic.Slug = slug
	topic.Description = strings.TrimSpace(req.Description)
	topic.CoverImage = strings.TrimSpace(req.CoverImage)
	topic.SortOrder = req.SortOrder
	topic.UpdatedAt = time.Now().UTC()
	if err := uc.topics.Update(ctx, topic); err != nil {
		return nil, err
	}
	return topic, uc.topics.ReplaceItems(ctx, topic.ID, req.Items)
}

func (uc *TaxonomyUsecase) DeleteTopic(ctx context.Context, id string) error {
	return uc.topics.Delete(ctx, id)
}

func (uc *TaxonomyUsecase) validateSlug(ctx context.Context, name string, slug string, excludeID *string, exists func(context.Context, string, *string) (bool, error)) (string, error) {
	normalized := normalizeSlug(slug)
	if normalized == "" {
		normalized = normalizeSlug(name)
	}
	if normalized == "" {
		return "", domain.ErrInvalidTaxonomy
	}
	taken, err := exists(ctx, normalized, excludeID)
	if err != nil {
		return "", err
	}
	if taken {
		return "", domain.ErrSlugTaken
	}
	return normalized, nil
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

func cleanOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
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
