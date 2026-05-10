package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/comment"
)

type CommentModel struct {
	ID            string  `gorm:"primaryKey;size:36"`
	RefType       string  `gorm:"size:20;not null;index"`
	RefID         string  `gorm:"size:36;not null;index"`
	ParentID      *string `gorm:"size:36;index"`
	RootID        *string `gorm:"size:36;index"`
	AuthorName    string  `gorm:"size:80;not null"`
	AuthorEmail   string  `gorm:"size:255;not null"`
	AuthorWebsite string  `gorm:"size:500"`
	AuthorIP      string  `gorm:"size:64"`
	UserAgent     string
	Content       string `gorm:"not null"`
	Status        string `gorm:"size:20;not null;index"`
	IsAdminReply  bool
	IsPinned      bool
	ReplyCount    int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (CommentModel) TableName() string { return "comments" }

type CommentRepository struct{ db *gorm.DB }

func NewCommentRepository(db *gorm.DB) *CommentRepository { return &CommentRepository{db: db} }

func (r *CommentRepository) Create(ctx context.Context, comment *domain.Comment) error {
	return r.db.WithContext(ctx).Create(toCommentModel(comment)).Error
}

func (r *CommentRepository) UpdateStatus(ctx context.Context, id string, status domain.CommentStatus) error {
	result := r.db.WithContext(ctx).Model(&CommentModel{}).Where("id = ?", id).Updates(map[string]any{"status": string(status), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

func (r *CommentRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&CommentModel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrCommentNotFound
	}
	return nil
}

func (r *CommentRepository) GetByID(ctx context.Context, id string) (*domain.Comment, error) {
	var model CommentModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return commentFromModel(model, err)
}

func (r *CommentRepository) ListByRef(ctx context.Context, refType domain.RefType, refID string, onlyApproved bool) ([]domain.Comment, error) {
	query := r.db.WithContext(ctx).Where("ref_type = ? AND ref_id = ?", string(refType), refID)
	if onlyApproved {
		query = query.Where("status = ?", string(domain.CommentApproved))
	}
	var models []CommentModel
	if err := query.Order("is_pinned DESC").Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	return commentsFromModels(models), nil
}

func (r *CommentRepository) ListAdmin(ctx context.Context, filter domain.CommentFilter) ([]domain.Comment, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	query := applyCommentFilter(r.db.WithContext(ctx).Model(&CommentModel{}), filter)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	var models []CommentModel
	if err := query.Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&models).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	return commentsFromModels(models), domain.NewPagination(page, perPage, total), nil
}

func (r *CommentRepository) IncrementReplyCount(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&CommentModel{}).Where("id = ?", id).UpdateColumn("reply_count", gorm.Expr("reply_count + 1")).Error
}

func applyCommentFilter(query *gorm.DB, filter domain.CommentFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	}
	if filter.RefType != "" {
		query = query.Where("ref_type = ?", string(filter.RefType))
	}
	if filter.RefID != "" {
		query = query.Where("ref_id = ?", filter.RefID)
	}
	if filter.Search != "" {
		search := "%" + filter.Search + "%"
		query = query.Where("author_name ILIKE ? OR author_email ILIKE ? OR content ILIKE ?", search, search, search)
	}
	return query
}

func commentFromModel(model CommentModel, err error) (*domain.Comment, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrCommentNotFound
	}
	if err != nil {
		return nil, err
	}
	return commentFromModelNoError(model), nil
}

func commentsFromModels(models []CommentModel) []domain.Comment {
	comments := make([]domain.Comment, 0, len(models))
	for _, model := range models {
		comments = append(comments, *commentFromModelNoError(model))
	}
	return comments
}

func commentFromModelNoError(model CommentModel) *domain.Comment {
	return &domain.Comment{ID: model.ID, RefType: domain.RefType(model.RefType), RefID: model.RefID, ParentID: model.ParentID, RootID: model.RootID, AuthorName: model.AuthorName, AuthorEmail: model.AuthorEmail, AuthorWebsite: model.AuthorWebsite, AuthorIP: model.AuthorIP, UserAgent: model.UserAgent, Content: model.Content, Status: domain.CommentStatus(model.Status), IsAdminReply: model.IsAdminReply, IsPinned: model.IsPinned, ReplyCount: model.ReplyCount, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toCommentModel(comment *domain.Comment) *CommentModel {
	return &CommentModel{ID: comment.ID, RefType: string(comment.RefType), RefID: comment.RefID, ParentID: comment.ParentID, RootID: comment.RootID, AuthorName: comment.AuthorName, AuthorEmail: comment.AuthorEmail, AuthorWebsite: comment.AuthorWebsite, AuthorIP: comment.AuthorIP, UserAgent: comment.UserAgent, Content: comment.Content, Status: string(comment.Status), IsAdminReply: comment.IsAdminReply, IsPinned: comment.IsPinned, ReplyCount: comment.ReplyCount, CreatedAt: comment.CreatedAt, UpdatedAt: comment.UpdatedAt}
}
