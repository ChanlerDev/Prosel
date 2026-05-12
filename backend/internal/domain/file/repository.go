package file

import "context"

type Repository interface {
	Create(ctx context.Context, file *FileAsset) error
	UpdateRef(ctx context.Context, id string, refType string, refID string) error
	MarkDeleted(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*FileAsset, error)
	List(ctx context.Context, filter FileFilter) ([]FileAsset, Pagination, error)
}

type StorageProvider interface {
	Put(ctx context.Context, objectKey string, data UploadReader, contentType string) (*StoredObject, error)
	Delete(ctx context.Context, objectKey string) error
}
