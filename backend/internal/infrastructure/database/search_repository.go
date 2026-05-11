package database

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	domain "github.com/chanler/prosel/backend/internal/domain/search"
)

type SearchDocumentModel struct {
	ID          string `gorm:"primaryKey;size:36"`
	RefType     string `gorm:"size:20;not null;uniqueIndex:idx_search_ref_unique"`
	RefID       string `gorm:"size:36;not null;uniqueIndex:idx_search_ref_unique"`
	Title       string `gorm:"size:255;not null"`
	Slug        string `gorm:"size:255"`
	Excerpt     string `gorm:"size:500"`
	SearchText  string `gorm:"not null"`
	Status      string `gorm:"size:20;not null"`
	PublishedAt *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (SearchDocumentModel) TableName() string { return "search_documents" }

type SearchRepository struct{ db *gorm.DB }

func NewSearchRepository(db *gorm.DB) *SearchRepository { return &SearchRepository{db: db} }

func (r *SearchRepository) UpsertDocument(ctx context.Context, doc *domain.SearchDocument) error {
	now := time.Now().UTC()
	model := toSearchDocumentModel(doc)
	if model.ID == "" {
		model.ID = searchDocumentID(doc.RefType, doc.RefID)
	}
	model.Status = "published"
	model.UpdatedAt = now
	if model.CreatedAt.IsZero() {
		model.CreatedAt = now
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ref_type"}, {Name: "ref_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"title", "slug", "excerpt", "search_text", "status", "published_at", "updated_at"}),
	}).Create(model).Error
}

func (r *SearchRepository) DeleteDocument(ctx context.Context, refType string, refID string) error {
	return r.db.WithContext(ctx).Where("ref_type = ? AND ref_id = ?", refType, refID).Delete(&SearchDocumentModel{}).Error
}

func (r *SearchRepository) Search(ctx context.Context, query string, filter domain.SearchFilter) ([]domain.SearchResult, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	base := r.db.WithContext(ctx).Model(&SearchDocumentModel{}).Where("status = ?", "published")
	if filter.Type.Valid() {
		base = base.Where("ref_type = ?", string(filter.Type))
	}
	base = base.Where("search_vector @@ plainto_tsquery('simple', ?)", query)

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}

	type row struct {
		RefType string
		RefID   string
		Title   string
		Slug    string
		Excerpt string
		Rank    float64
	}
	var rows []row
	err := base.Select("ref_type, ref_id, title, slug, excerpt, ts_rank(search_vector, plainto_tsquery('simple', ?)) AS rank", query).
		Order("rank DESC").Order("published_at DESC NULLS LAST").Order("updated_at DESC").
		Limit(perPage).Offset((page - 1) * perPage).Scan(&rows).Error
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	results := make([]domain.SearchResult, 0, len(rows))
	for _, item := range rows {
		results = append(results, domain.SearchResult{RefType: domain.RefType(item.RefType), RefID: item.RefID, Title: item.Title, Slug: item.Slug, Excerpt: item.Excerpt, Rank: item.Rank})
	}
	return results, domain.NewPagination(page, perPage, total), nil
}

func (r *SearchRepository) Rebuild(ctx context.Context) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM search_documents").Error; err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO search_documents (id, ref_type, ref_id, title, slug, excerpt, search_text, status, published_at, created_at, updated_at)
SELECT substring(md5('post:' || id) from 1 for 8) || '-' || substring(md5('post:' || id) from 9 for 4) || '-' || substring(md5('post:' || id) from 13 for 4) || '-' || substring(md5('post:' || id) from 17 for 4) || '-' || substring(md5('post:' || id) from 21 for 12), 'post', id, title, slug, excerpt, content_text, 'published', published_at, NOW(), NOW()
FROM posts WHERE status = 'published'`).Error; err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO search_documents (id, ref_type, ref_id, title, slug, excerpt, search_text, status, published_at, created_at, updated_at)
SELECT substring(md5('note:' || id) from 1 for 8) || '-' || substring(md5('note:' || id) from 9 for 4) || '-' || substring(md5('note:' || id) from 13 for 4) || '-' || substring(md5('note:' || id) from 17 for 4) || '-' || substring(md5('note:' || id) from 21 for 12), 'note', id, CASE WHEN title = '' THEN left(content_text, 120) ELSE title END, slug, left(content_text, 300), content_text, 'published', published_at, NOW(), NOW()
FROM notes WHERE status = 'published'`).Error; err != nil {
			return err
		}
		return tx.Exec(`INSERT INTO search_documents (id, ref_type, ref_id, title, slug, excerpt, search_text, status, published_at, created_at, updated_at)
SELECT substring(md5('page:' || id) from 1 for 8) || '-' || substring(md5('page:' || id) from 9 for 4) || '-' || substring(md5('page:' || id) from 13 for 4) || '-' || substring(md5('page:' || id) from 17 for 4) || '-' || substring(md5('page:' || id) from 21 for 12), 'page', id, title, slug, CASE WHEN subtitle = '' THEN left(content_text, 300) ELSE subtitle END, content_text, 'published', NULL, NOW(), NOW()
FROM pages WHERE status = 'published'`).Error
	})
}

func (r *SearchRepository) Status(ctx context.Context) (*domain.IndexStatus, error) {
	var status domain.IndexStatus
	if err := r.db.WithContext(ctx).Model(&SearchDocumentModel{}).Count(&status.Total).Error; err != nil {
		return nil, err
	}
	counts := []struct {
		RefType string
		Count   int64
	}{}
	if err := r.db.WithContext(ctx).Model(&SearchDocumentModel{}).Select("ref_type, count(*) AS count").Group("ref_type").Scan(&counts).Error; err != nil {
		return nil, err
	}
	for _, count := range counts {
		switch domain.RefType(count.RefType) {
		case domain.RefPost:
			status.Posts = count.Count
		case domain.RefNote:
			status.Notes = count.Count
		case domain.RefPage:
			status.Pages = count.Count
		}
	}
	var updatedAt *time.Time
	if err := r.db.WithContext(ctx).Model(&SearchDocumentModel{}).Select("max(updated_at)").Scan(&updatedAt).Error; err != nil {
		return nil, err
	}
	status.UpdatedAt = updatedAt
	return &status, nil
}

func toSearchDocumentModel(doc *domain.SearchDocument) *SearchDocumentModel {
	return &SearchDocumentModel{ID: doc.ID, RefType: string(doc.RefType), RefID: doc.RefID, Title: doc.Title, Slug: doc.Slug, Excerpt: doc.Excerpt, SearchText: doc.SearchText, Status: doc.Status, PublishedAt: doc.PublishedAt, CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt}
}

func searchDocumentID(refType domain.RefType, refID string) string {
	sum := md5.Sum([]byte(string(refType) + ":" + refID))
	hexValue := hex.EncodeToString(sum[:])
	return hexValue[:8] + "-" + hexValue[8:12] + "-" + hexValue[12:16] + "-" + hexValue[16:20] + "-" + hexValue[20:32]
}
