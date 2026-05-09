package taxonomy

import (
	"context"
	"errors"
	"testing"

	domain "github.com/chanler/prosel/backend/internal/domain/taxonomy"
)

type fakeCategoryRepo struct {
	category   *domain.Category
	categories []domain.CategoryNode
	slugExists bool
	err        error
}

func (r *fakeCategoryRepo) Create(ctx context.Context, category *domain.Category) error {
	r.category = category
	return r.err
}
func (r *fakeCategoryRepo) Update(ctx context.Context, category *domain.Category) error {
	r.category = category
	return r.err
}
func (r *fakeCategoryRepo) Delete(ctx context.Context, id string) error { return r.err }
func (r *fakeCategoryRepo) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	return r.category, r.err
}
func (r *fakeCategoryRepo) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	return r.category, r.err
}
func (r *fakeCategoryRepo) ListTree(ctx context.Context) ([]domain.CategoryNode, error) {
	return r.categories, r.err
}
func (r *fakeCategoryRepo) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return r.slugExists, r.err
}

type fakeTagRepo struct {
	tag        *domain.Tag
	tags       []domain.TagWithCount
	slugExists bool
	err        error
}

func (r *fakeTagRepo) Create(ctx context.Context, tag *domain.Tag) error { r.tag = tag; return r.err }
func (r *fakeTagRepo) Update(ctx context.Context, tag *domain.Tag) error { r.tag = tag; return r.err }
func (r *fakeTagRepo) Delete(ctx context.Context, id string) error       { return r.err }
func (r *fakeTagRepo) GetByID(ctx context.Context, id string) (*domain.Tag, error) {
	return r.tag, r.err
}
func (r *fakeTagRepo) List(ctx context.Context) ([]domain.TagWithCount, error) { return r.tags, r.err }
func (r *fakeTagRepo) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return r.slugExists, r.err
}
func (r *fakeTagRepo) ReplacePostTags(ctx context.Context, postID string, tagIDs []string) error {
	return r.err
}

type fakeTopicRepo struct {
	topic      *domain.Topic
	topics     []domain.Topic
	slugExists bool
	err        error
}

func (r *fakeTopicRepo) Create(ctx context.Context, topic *domain.Topic) error {
	r.topic = topic
	return r.err
}
func (r *fakeTopicRepo) Update(ctx context.Context, topic *domain.Topic) error {
	r.topic = topic
	return r.err
}
func (r *fakeTopicRepo) Delete(ctx context.Context, id string) error { return r.err }
func (r *fakeTopicRepo) GetByID(ctx context.Context, id string) (*domain.Topic, error) {
	return r.topic, r.err
}
func (r *fakeTopicRepo) GetBySlug(ctx context.Context, slug string) (*domain.Topic, error) {
	return r.topic, r.err
}
func (r *fakeTopicRepo) List(ctx context.Context) ([]domain.Topic, error) { return r.topics, r.err }
func (r *fakeTopicRepo) ReplaceItems(ctx context.Context, topicID string, items []domain.TopicItem) error {
	return r.err
}
func (r *fakeTopicRepo) ListItems(ctx context.Context, topicID string) ([]domain.TopicItem, error) {
	return nil, r.err
}
func (r *fakeTopicRepo) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return r.slugExists, r.err
}

func TestCreateCategoryGeneratesSlug(t *testing.T) {
	categories := &fakeCategoryRepo{}
	uc := NewTaxonomyUsecase(categories, &fakeTagRepo{}, &fakeTopicRepo{})

	category, err := uc.CreateCategory(context.Background(), CategoryRequest{Name: "Go Notes", Description: "Backend", SortOrder: 2})
	if err != nil {
		t.Fatalf("CreateCategory() error = %v", err)
	}
	if category.Slug != "go-notes" || category.ID == "" {
		t.Fatalf("category = %#v", category)
	}
	if categories.category == nil || categories.category.Name != "Go Notes" {
		t.Fatalf("category not persisted: %#v", categories.category)
	}
}

func TestCreateTagRejectsDuplicateSlug(t *testing.T) {
	uc := NewTaxonomyUsecase(&fakeCategoryRepo{}, &fakeTagRepo{slugExists: true}, &fakeTopicRepo{})

	_, err := uc.CreateTag(context.Background(), TagRequest{Name: "Go", Slug: "go"})
	if !errors.Is(err, domain.ErrSlugTaken) {
		t.Fatalf("CreateTag() error = %v, want %v", err, domain.ErrSlugTaken)
	}
}

func TestListCategoriesReturnsTree(t *testing.T) {
	categories := []domain.CategoryNode{{Category: domain.Category{ID: "root", Name: "Root"}, Children: []domain.CategoryNode{{Category: domain.Category{ID: "child", Name: "Child"}}}}}
	uc := NewTaxonomyUsecase(&fakeCategoryRepo{categories: categories}, &fakeTagRepo{}, &fakeTopicRepo{})

	result, err := uc.ListCategories(context.Background())
	if err != nil {
		t.Fatalf("ListCategories() error = %v", err)
	}
	if len(result) != 1 || len(result[0].Children) != 1 {
		t.Fatalf("result = %#v", result)
	}
}
