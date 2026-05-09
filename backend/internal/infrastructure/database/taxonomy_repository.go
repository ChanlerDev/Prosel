package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/chanler/prosel/backend/internal/domain/taxonomy"
)

type CategoryModel struct {
	ID          string  `gorm:"primaryKey;size:36"`
	ParentID    *string `gorm:"size:36;index"`
	Name        string  `gorm:"size:100;not null"`
	Slug        string  `gorm:"size:255;not null;uniqueIndex"`
	Description string
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PostCount   int64 `gorm:"-"`
}

func (CategoryModel) TableName() string { return "categories" }

type TagModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"size:80;not null;uniqueIndex"`
	Slug        string `gorm:"size:255;not null;uniqueIndex"`
	Color       string `gorm:"size:20"`
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PostCount   int64 `gorm:"-"`
}

func (TagModel) TableName() string { return "tags" }

type PostTagModel struct {
	PostID string `gorm:"primaryKey;size:36"`
	TagID  string `gorm:"primaryKey;size:36"`
}

func (PostTagModel) TableName() string { return "post_tags" }

type TopicModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"size:100;not null"`
	Slug        string `gorm:"size:255;not null;uniqueIndex"`
	Description string
	CoverImage  string `gorm:"size:500"`
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (TopicModel) TableName() string { return "topics" }

type TopicItemModel struct {
	TopicID   string `gorm:"primaryKey;size:36"`
	RefType   string `gorm:"primaryKey;size:20"`
	RefID     string `gorm:"primaryKey;size:36"`
	SortOrder int
	Title     string `gorm:"-"`
	Slug      string `gorm:"-"`
}

func (TopicItemModel) TableName() string { return "topic_items" }

type CategoryRepository struct{ db *gorm.DB }

func NewCategoryRepository(db *gorm.DB) *CategoryRepository { return &CategoryRepository{db: db} }

func (r *CategoryRepository) Create(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Create(toCategoryModel(category)).Error
}

func (r *CategoryRepository) Update(ctx context.Context, category *domain.Category) error {
	return r.db.WithContext(ctx).Model(&CategoryModel{}).Where("id = ?", category.ID).Updates(map[string]any{"parent_id": category.ParentID, "name": category.Name, "slug": category.Slug, "description": category.Description, "sort_order": category.SortOrder, "updated_at": category.UpdatedAt}).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	return deleteByID(ctx, r.db, &CategoryModel{}, id)
}

func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*domain.Category, error) {
	var model CategoryModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return categoryFromModel(model, err)
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*domain.Category, error) {
	var model CategoryModel
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error
	return categoryFromModel(model, err)
}

func (r *CategoryRepository) ListTree(ctx context.Context) ([]domain.CategoryNode, error) {
	var models []CategoryModel
	err := r.db.WithContext(ctx).Table("categories").Select("categories.*, COUNT(posts.id) AS post_count").Joins("LEFT JOIN posts ON posts.category_id = categories.id AND posts.status = ?", "published").Group("categories.id").Order("categories.sort_order ASC").Order("categories.name ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	return categoryTree(models), nil
}

func (r *CategoryRepository) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return slugExists(ctx, r.db, "categories", slug, excludeID)
}

type TagRepository struct{ db *gorm.DB }

func NewTagRepository(db *gorm.DB) *TagRepository { return &TagRepository{db: db} }

func (r *TagRepository) Create(ctx context.Context, tag *domain.Tag) error {
	return r.db.WithContext(ctx).Create(toTagModel(tag)).Error
}

func (r *TagRepository) Update(ctx context.Context, tag *domain.Tag) error {
	return r.db.WithContext(ctx).Model(&TagModel{}).Where("id = ?", tag.ID).Updates(map[string]any{"name": tag.Name, "slug": tag.Slug, "color": tag.Color, "description": tag.Description, "updated_at": tag.UpdatedAt}).Error
}

func (r *TagRepository) Delete(ctx context.Context, id string) error {
	return deleteByID(ctx, r.db, &TagModel{}, id)
}

func (r *TagRepository) GetByID(ctx context.Context, id string) (*domain.Tag, error) {
	var model TagModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return tagFromModel(model, err)
}

func (r *TagRepository) List(ctx context.Context) ([]domain.TagWithCount, error) {
	var models []TagModel
	err := r.db.WithContext(ctx).Table("tags").Select("tags.*, COUNT(post_tags.post_id) AS post_count").Joins("LEFT JOIN post_tags ON post_tags.tag_id = tags.id").Joins("LEFT JOIN posts ON posts.id = post_tags.post_id AND posts.status = ?", "published").Group("tags.id").Order("tags.name ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.TagWithCount, 0, len(models))
	for _, model := range models {
		result = append(result, domain.TagWithCount{Tag: *tagFromModelNoError(model), PostCount: model.PostCount})
	}
	return result, nil
}

func (r *TagRepository) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return slugExists(ctx, r.db, "tags", slug, excludeID)
}

func (r *TagRepository) ReplacePostTags(ctx context.Context, postID string, tagIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("post_id = ?", postID).Delete(&PostTagModel{}).Error; err != nil {
			return err
		}
		for _, tagID := range tagIDs {
			if tagID == "" {
				continue
			}
			if err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&PostTagModel{PostID: postID, TagID: tagID}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

type TopicRepository struct{ db *gorm.DB }

func NewTopicRepository(db *gorm.DB) *TopicRepository { return &TopicRepository{db: db} }

func (r *TopicRepository) Create(ctx context.Context, topic *domain.Topic) error {
	return r.db.WithContext(ctx).Create(toTopicModel(topic)).Error
}

func (r *TopicRepository) Update(ctx context.Context, topic *domain.Topic) error {
	return r.db.WithContext(ctx).Model(&TopicModel{}).Where("id = ?", topic.ID).Updates(map[string]any{"name": topic.Name, "slug": topic.Slug, "description": topic.Description, "cover_image": topic.CoverImage, "sort_order": topic.SortOrder, "updated_at": topic.UpdatedAt}).Error
}

func (r *TopicRepository) Delete(ctx context.Context, id string) error {
	return deleteByID(ctx, r.db, &TopicModel{}, id)
}

func (r *TopicRepository) GetByID(ctx context.Context, id string) (*domain.Topic, error) {
	var model TopicModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return topicFromModel(model, err)
}

func (r *TopicRepository) GetBySlug(ctx context.Context, slug string) (*domain.Topic, error) {
	var model TopicModel
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&model).Error
	return topicFromModel(model, err)
}

func (r *TopicRepository) List(ctx context.Context) ([]domain.Topic, error) {
	var models []TopicModel
	if err := r.db.WithContext(ctx).Order("sort_order ASC").Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]domain.Topic, 0, len(models))
	for _, model := range models {
		result = append(result, *topicFromModelNoError(model))
	}
	return result, nil
}

func (r *TopicRepository) ReplaceItems(ctx context.Context, topicID string, items []domain.TopicItem) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("topic_id = ?", topicID).Delete(&TopicItemModel{}).Error; err != nil {
			return err
		}
		for _, item := range items {
			if item.RefType == "" || item.RefID == "" {
				continue
			}
			if err := tx.Create(&TopicItemModel{TopicID: topicID, RefType: item.RefType, RefID: item.RefID, SortOrder: item.SortOrder}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *TopicRepository) ListItems(ctx context.Context, topicID string) ([]domain.TopicItem, error) {
	var models []TopicItemModel
	err := r.db.WithContext(ctx).Table("topic_items").Select("topic_items.*, posts.title, posts.slug").Joins("LEFT JOIN posts ON posts.id = topic_items.ref_id AND topic_items.ref_type = ?", "post").Where("topic_items.topic_id = ?", topicID).Order("topic_items.sort_order ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]domain.TopicItem, 0, len(models))
	for _, model := range models {
		result = append(result, domain.TopicItem{TopicID: model.TopicID, RefType: model.RefType, RefID: model.RefID, Title: model.Title, Slug: model.Slug, SortOrder: model.SortOrder})
	}
	return result, nil
}

func (r *TopicRepository) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	return slugExists(ctx, r.db, "topics", slug, excludeID)
}

func slugExists(ctx context.Context, db *gorm.DB, table string, slug string, excludeID *string) (bool, error) {
	query := db.WithContext(ctx).Table(table).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func deleteByID(ctx context.Context, db *gorm.DB, model any, id string) error {
	result := db.WithContext(ctx).Delete(model, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrTaxonomyNotFound
	}
	return nil
}

func categoryTree(models []CategoryModel) []domain.CategoryNode {
	nodes := make(map[string]*domain.CategoryNode, len(models))
	for _, model := range models {
		nodes[model.ID] = &domain.CategoryNode{Category: *categoryFromModelNoError(model), PostCount: model.PostCount}
	}
	roots := make([]domain.CategoryNode, 0)
	for _, model := range models {
		node := nodes[model.ID]
		if model.ParentID != nil {
			if parent := nodes[*model.ParentID]; parent != nil {
				parent.Children = append(parent.Children, *node)
				continue
			}
		}
		roots = append(roots, *node)
	}
	return roots
}

func categoryFromModel(model CategoryModel, err error) (*domain.Category, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTaxonomyNotFound
	}
	if err != nil {
		return nil, err
	}
	return categoryFromModelNoError(model), nil
}
func categoryFromModelNoError(model CategoryModel) *domain.Category {
	return &domain.Category{ID: model.ID, ParentID: model.ParentID, Name: model.Name, Slug: model.Slug, Description: model.Description, SortOrder: model.SortOrder, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
func toCategoryModel(category *domain.Category) *CategoryModel {
	return &CategoryModel{ID: category.ID, ParentID: category.ParentID, Name: category.Name, Slug: category.Slug, Description: category.Description, SortOrder: category.SortOrder, CreatedAt: category.CreatedAt, UpdatedAt: category.UpdatedAt}
}

func tagFromModel(model TagModel, err error) (*domain.Tag, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTaxonomyNotFound
	}
	if err != nil {
		return nil, err
	}
	return tagFromModelNoError(model), nil
}
func tagFromModelNoError(model TagModel) *domain.Tag {
	return &domain.Tag{ID: model.ID, Name: model.Name, Slug: model.Slug, Color: model.Color, Description: model.Description, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
func toTagModel(tag *domain.Tag) *TagModel {
	return &TagModel{ID: tag.ID, Name: tag.Name, Slug: tag.Slug, Color: tag.Color, Description: tag.Description, CreatedAt: tag.CreatedAt, UpdatedAt: tag.UpdatedAt}
}

func topicFromModel(model TopicModel, err error) (*domain.Topic, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrTaxonomyNotFound
	}
	if err != nil {
		return nil, err
	}
	return topicFromModelNoError(model), nil
}
func topicFromModelNoError(model TopicModel) *domain.Topic {
	return &domain.Topic{ID: model.ID, Name: model.Name, Slug: model.Slug, Description: model.Description, CoverImage: model.CoverImage, SortOrder: model.SortOrder, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}
func toTopicModel(topic *domain.Topic) *TopicModel {
	return &TopicModel{ID: topic.ID, Name: topic.Name, Slug: topic.Slug, Description: topic.Description, CoverImage: topic.CoverImage, SortOrder: topic.SortOrder, CreatedAt: topic.CreatedAt, UpdatedAt: topic.UpdatedAt}
}
