package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/page"
)

type PageModel struct {
	ID              string `gorm:"primaryKey;size:36"`
	AuthorID        string `gorm:"size:36;index"`
	Title           string `gorm:"size:200;not null"`
	Slug            string `gorm:"size:255;not null;uniqueIndex"`
	Subtitle        string `gorm:"size:300"`
	ContentMarkdown string `gorm:"not null"`
	ContentText     string `gorm:"not null"`
	Template        string `gorm:"size:40;not null"`
	Status          string `gorm:"size:20;not null;index"`
	SortOrder       int
	SEOTitle        string `gorm:"size:200"`
	SEODescription  string `gorm:"size:500"`
	ViewCount       int64
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (PageModel) TableName() string { return "pages" }

type FriendModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	Name        string `gorm:"size:100;not null"`
	URL         string `gorm:"size:500;not null;uniqueIndex"`
	AvatarURL   string `gorm:"size:500"`
	Description string `gorm:"size:500"`
	Status      string `gorm:"size:20;not null;index"`
	SortOrder   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (FriendModel) TableName() string { return "friends" }

type PageRepository struct{ db *gorm.DB }

func NewPageRepository(db *gorm.DB) *PageRepository { return &PageRepository{db: db} }

func (r *PageRepository) Create(ctx context.Context, page *domain.Page) error {
	return r.db.WithContext(ctx).Create(toPageModel(page)).Error
}

func (r *PageRepository) Update(ctx context.Context, page *domain.Page) error {
	return r.db.WithContext(ctx).Model(&PageModel{}).Where("id = ?", page.ID).Updates(map[string]any{"title": page.Title, "slug": page.Slug, "subtitle": page.Subtitle, "content_markdown": page.ContentMarkdown, "content_text": page.ContentText, "template": string(page.Template), "status": string(page.Status), "sort_order": page.SortOrder, "seo_title": page.SEOTitle, "seo_description": page.SEODescription, "updated_at": page.UpdatedAt}).Error
}

func (r *PageRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&PageModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrPageNotFound
	}
	return nil
}

func (r *PageRepository) GetByID(ctx context.Context, id string) (*domain.Page, error) {
	var model PageModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return pageFromModel(model, err)
}

func (r *PageRepository) GetBySlug(ctx context.Context, slug string, includeDraft bool) (*domain.Page, error) {
	var model PageModel
	query := r.db.WithContext(ctx).Where("slug = ?", slug)
	if !includeDraft {
		query = query.Where("status = ?", string(domain.PagePublished))
	}
	err := query.First(&model).Error
	return pageFromModel(model, err)
}

func (r *PageRepository) List(ctx context.Context, filter domain.PageFilter) ([]domain.Page, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	query := r.db.WithContext(ctx).Model(&PageModel{})
	query = applyPageFilter(query, filter)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	var models []PageModel
	if err := query.Order("sort_order ASC").Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&models).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	pages := make([]domain.Page, 0, len(models))
	for _, model := range models {
		pages = append(pages, *pageFromModelNoError(model))
	}
	return pages, domain.NewPagination(page, perPage, total), nil
}

func (r *PageRepository) SlugExists(ctx context.Context, slug string, excludeID *string) (bool, error) {
	query := r.db.WithContext(ctx).Model(&PageModel{}).Where("slug = ?", slug)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PageRepository) IncrementView(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&PageModel{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

type FriendRepository struct{ db *gorm.DB }

func NewFriendRepository(db *gorm.DB) *FriendRepository { return &FriendRepository{db: db} }

func (r *FriendRepository) Create(ctx context.Context, friend *domain.Friend) error {
	return r.db.WithContext(ctx).Create(toFriendModel(friend)).Error
}

func (r *FriendRepository) Update(ctx context.Context, friend *domain.Friend) error {
	return r.db.WithContext(ctx).Model(&FriendModel{}).Where("id = ?", friend.ID).Updates(map[string]any{"name": friend.Name, "url": friend.URL, "avatar_url": friend.AvatarURL, "description": friend.Description, "status": string(friend.Status), "sort_order": friend.SortOrder, "updated_at": friend.UpdatedAt}).Error
}

func (r *FriendRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&FriendModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrFriendNotFound
	}
	return nil
}

func (r *FriendRepository) GetByID(ctx context.Context, id string) (*domain.Friend, error) {
	var model FriendModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return friendFromModel(model, err)
}

func (r *FriendRepository) URLExists(ctx context.Context, url string, excludeID *string) (bool, error) {
	query := r.db.WithContext(ctx).Model(&FriendModel{}).Where("url = ?", url)
	if excludeID != nil {
		query = query.Where("id <> ?", *excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FriendRepository) List(ctx context.Context, status string) ([]domain.Friend, error) {
	query := r.db.WithContext(ctx).Model(&FriendModel{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var models []FriendModel
	if err := query.Order("sort_order ASC").Order("name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	friends := make([]domain.Friend, 0, len(models))
	for _, model := range models {
		friends = append(friends, *friendFromModelNoError(model))
	}
	return friends, nil
}

func applyPageFilter(query *gorm.DB, filter domain.PageFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR subtitle ILIKE ? OR content_text ILIKE ?", search, search, search)
	}
	return query
}

func pageFromModel(model PageModel, err error) (*domain.Page, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrPageNotFound
	}
	if err != nil {
		return nil, err
	}
	return pageFromModelNoError(model), nil
}

func pageFromModelNoError(model PageModel) *domain.Page {
	return &domain.Page{ID: model.ID, AuthorID: model.AuthorID, Title: model.Title, Slug: model.Slug, Subtitle: model.Subtitle, ContentMarkdown: model.ContentMarkdown, ContentText: model.ContentText, Template: domain.PageTemplate(model.Template), Status: domain.PageStatus(model.Status), SortOrder: model.SortOrder, SEOTitle: model.SEOTitle, SEODescription: model.SEODescription, ViewCount: model.ViewCount, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toPageModel(page *domain.Page) *PageModel {
	return &PageModel{ID: page.ID, AuthorID: page.AuthorID, Title: page.Title, Slug: page.Slug, Subtitle: page.Subtitle, ContentMarkdown: page.ContentMarkdown, ContentText: page.ContentText, Template: string(page.Template), Status: string(page.Status), SortOrder: page.SortOrder, SEOTitle: page.SEOTitle, SEODescription: page.SEODescription, ViewCount: page.ViewCount, CreatedAt: page.CreatedAt, UpdatedAt: page.UpdatedAt}
}

func friendFromModel(model FriendModel, err error) (*domain.Friend, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrFriendNotFound
	}
	if err != nil {
		return nil, err
	}
	return friendFromModelNoError(model), nil
}

func friendFromModelNoError(model FriendModel) *domain.Friend {
	return &domain.Friend{ID: model.ID, Name: model.Name, URL: model.URL, AvatarURL: model.AvatarURL, Description: model.Description, Status: domain.FriendStatus(model.Status), SortOrder: model.SortOrder, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toFriendModel(friend *domain.Friend) *FriendModel {
	return &FriendModel{ID: friend.ID, Name: friend.Name, URL: friend.URL, AvatarURL: friend.AvatarURL, Description: friend.Description, Status: string(friend.Status), SortOrder: friend.SortOrder, CreatedAt: friend.CreatedAt, UpdatedAt: friend.UpdatedAt}
}
