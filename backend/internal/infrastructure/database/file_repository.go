package database

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/file"
)

type FileModel struct {
	ID           string `gorm:"primaryKey;size:36"`
	UploaderID   string `gorm:"size:36;index"`
	OriginalName string `gorm:"size:255;not null"`
	FileName     string `gorm:"size:255;not null"`
	StorageType  string `gorm:"size:20;not null"`
	ObjectKey    string `gorm:"size:500;not null"`
	PublicURL    string `gorm:"size:500;not null"`
	MimeType     string `gorm:"size:120;not null"`
	ByteSize     int64  `gorm:"not null"`
	Width        *int
	Height       *int
	RefType      string `gorm:"size:20;index"`
	RefID        string `gorm:"size:36;index"`
	Status       string `gorm:"size:20;not null;index"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (FileModel) TableName() string { return "files" }

type FileRepository struct{ db *gorm.DB }

func NewFileRepository(db *gorm.DB) *FileRepository { return &FileRepository{db: db} }

func (r *FileRepository) Create(ctx context.Context, file *domain.FileAsset) error {
	return r.db.WithContext(ctx).Create(toFileModel(file)).Error
}

func (r *FileRepository) UpdateRef(ctx context.Context, id string, refType string, refID string) error {
	result := r.db.WithContext(ctx).Model(&FileModel{}).Where("id = ? AND status <> ?", id, string(domain.FileStatusDeleted)).Updates(map[string]any{"ref_type": refType, "ref_id": refID, "status": string(domain.FileStatusAttached), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrFileNotFound
	}
	return nil
}

func (r *FileRepository) MarkDeleted(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&FileModel{}).Where("id = ?", id).Updates(map[string]any{"status": string(domain.FileStatusDeleted), "updated_at": time.Now().UTC()})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrFileNotFound
	}
	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id string) (*domain.FileAsset, error) {
	var model FileModel
	err := r.db.WithContext(ctx).Where("id = ? AND status <> ?", id, string(domain.FileStatusDeleted)).First(&model).Error
	return fileFromModel(model, err)
}

func (r *FileRepository) List(ctx context.Context, filter domain.FileFilter) ([]domain.FileAsset, domain.Pagination, error) {
	page, perPage := domain.NormalizePagination(filter.Page, filter.PerPage)
	query := r.db.WithContext(ctx).Model(&FileModel{})
	query = applyFileFilter(query, filter)

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, domain.Pagination{}, err
	}
	var models []FileModel
	err := query.Order("created_at DESC").Limit(perPage).Offset((page - 1) * perPage).Find(&models).Error
	if err != nil {
		return nil, domain.Pagination{}, err
	}
	files := make([]domain.FileAsset, 0, len(models))
	for _, model := range models {
		files = append(files, *fileFromModelNoError(model))
	}
	return files, domain.NewPagination(page, perPage, total), nil
}

func applyFileFilter(query *gorm.DB, filter domain.FileFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", string(*filter.Status))
	} else {
		query = query.Where("status <> ?", string(domain.FileStatusDeleted))
	}
	if filter.MimeType != "" {
		query = query.Where("mime_type LIKE ?", filter.MimeType+"%")
	}
	if filter.Search != "" {
		query = query.Where("original_name ILIKE ? OR file_name ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}
	return query
}

func fileFromModel(model FileModel, err error) (*domain.FileAsset, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrFileNotFound
	}
	if err != nil {
		return nil, err
	}
	return fileFromModelNoError(model), nil
}

func fileFromModelNoError(model FileModel) *domain.FileAsset {
	return &domain.FileAsset{ID: model.ID, UploaderID: model.UploaderID, OriginalName: model.OriginalName, FileName: model.FileName, StorageType: domain.StorageType(model.StorageType), ObjectKey: model.ObjectKey, PublicURL: model.PublicURL, MimeType: model.MimeType, ByteSize: model.ByteSize, Width: model.Width, Height: model.Height, RefType: model.RefType, RefID: model.RefID, Status: domain.FileStatus(model.Status), CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}
}

func toFileModel(file *domain.FileAsset) *FileModel {
	return &FileModel{ID: file.ID, UploaderID: file.UploaderID, OriginalName: file.OriginalName, FileName: file.FileName, StorageType: string(file.StorageType), ObjectKey: file.ObjectKey, PublicURL: file.PublicURL, MimeType: file.MimeType, ByteSize: file.ByteSize, Width: file.Width, Height: file.Height, RefType: file.RefType, RefID: file.RefID, Status: string(file.Status), CreatedAt: file.CreatedAt, UpdatedAt: file.UpdatedAt}
}
