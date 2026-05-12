package file

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	domain "github.com/chanler/prosel/backend/internal/domain/file"
)

type fakeFileRepo struct {
	file       *domain.FileAsset
	files      []domain.FileAsset
	pagination domain.Pagination
	err        error
	listed     domain.FileFilter
	refType    string
	refID      string
	deletedID  string
}

func (r *fakeFileRepo) Create(ctx context.Context, file *domain.FileAsset) error {
	r.file = file
	return r.err
}
func (r *fakeFileRepo) UpdateRef(ctx context.Context, id string, refType string, refID string) error {
	r.refType = refType
	r.refID = refID
	return r.err
}
func (r *fakeFileRepo) MarkDeleted(ctx context.Context, id string) error {
	r.deletedID = id
	return r.err
}
func (r *fakeFileRepo) GetByID(ctx context.Context, id string) (*domain.FileAsset, error) {
	return r.file, r.err
}
func (r *fakeFileRepo) List(ctx context.Context, filter domain.FileFilter) ([]domain.FileAsset, domain.Pagination, error) {
	r.listed = filter
	return r.files, r.pagination, r.err
}

type fakeStorage struct {
	objectKey   string
	contentType string
	data        string
	deletedKey  string
	err         error
}

func (s *fakeStorage) Put(ctx context.Context, objectKey string, data domain.UploadReader, contentType string) (*domain.StoredObject, error) {
	s.objectKey = objectKey
	s.contentType = contentType
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(data)
	s.data = buf.String()
	return &domain.StoredObject{ObjectKey: objectKey, PublicURL: "/uploads/" + objectKey}, s.err
}
func (s *fakeStorage) Delete(ctx context.Context, objectKey string) error {
	s.deletedKey = objectKey
	return s.err
}

func TestUploadStoresImageAsOrphanWithMetadata(t *testing.T) {
	repo := &fakeFileRepo{}
	storage := &fakeStorage{}
	uc := NewFileUsecase(repo, storage, Options{MaxUploadBytes: 1024, AllowedMimeTypes: []string{"image/png"}})

	asset, err := uc.Upload(context.Background(), UploadRequest{UploaderID: "user-1", OriginalName: "Hero Image.png", MimeType: "image/png", Size: 7, Data: strings.NewReader("pngdata")})
	if err != nil {
		t.Fatalf("Upload() error = %v", err)
	}
	if asset.ID == "" || asset.Status != domain.FileStatusOrphan || asset.StorageType != domain.StorageLocal {
		t.Fatalf("asset metadata = %#v, want id, orphan local", asset)
	}
	if asset.FileName == "" || !strings.HasSuffix(asset.FileName, ".png") {
		t.Fatalf("FileName = %q, want generated png name", asset.FileName)
	}
	if storage.contentType != "image/png" || storage.data != "pngdata" {
		t.Fatalf("stored content type/data = %q/%q", storage.contentType, storage.data)
	}
	if repo.file == nil || repo.file.PublicURL == "" || repo.file.ObjectKey == "" {
		t.Fatalf("file was not persisted with storage info: %#v", repo.file)
	}
}

func TestUploadRejectsDisallowedMimeType(t *testing.T) {
	uc := NewFileUsecase(&fakeFileRepo{}, &fakeStorage{}, Options{MaxUploadBytes: 1024, AllowedMimeTypes: []string{"image/png"}})

	_, err := uc.Upload(context.Background(), UploadRequest{UploaderID: "user-1", OriginalName: "doc.pdf", MimeType: "application/pdf", Size: 10, Data: strings.NewReader("pdf")})
	if !errors.Is(err, domain.ErrInvalidFile) {
		t.Fatalf("Upload() error = %v, want %v", err, domain.ErrInvalidFile)
	}
}

func TestUploadRejectsOversizedFile(t *testing.T) {
	uc := NewFileUsecase(&fakeFileRepo{}, &fakeStorage{}, Options{MaxUploadBytes: 5, AllowedMimeTypes: []string{"image/png"}})

	_, err := uc.Upload(context.Background(), UploadRequest{UploaderID: "user-1", OriginalName: "big.png", MimeType: "image/png", Size: 6, Data: strings.NewReader("123456")})
	if !errors.Is(err, domain.ErrFileTooLarge) {
		t.Fatalf("Upload() error = %v, want %v", err, domain.ErrFileTooLarge)
	}
}

func TestAttachToRefValidatesRefType(t *testing.T) {
	repo := &fakeFileRepo{}
	uc := NewFileUsecase(repo, &fakeStorage{}, Options{})

	if err := uc.AttachToRef(context.Background(), "file-1", "post", "post-1"); err != nil {
		t.Fatalf("AttachToRef() error = %v", err)
	}
	if repo.refType != "post" || repo.refID != "post-1" {
		t.Fatalf("ref = %q/%q, want post/post-1", repo.refType, repo.refID)
	}
	if err := uc.AttachToRef(context.Background(), "file-1", "comment", "comment-1"); !errors.Is(err, domain.ErrInvalidFile) {
		t.Fatalf("AttachToRef() invalid ref error = %v, want %v", err, domain.ErrInvalidFile)
	}
}

func TestDeleteFileDeletesStorageThenMarksDeleted(t *testing.T) {
	repo := &fakeFileRepo{file: &domain.FileAsset{ID: "file-1", ObjectKey: "uploads/file.png"}}
	storage := &fakeStorage{}
	uc := NewFileUsecase(repo, storage, Options{})

	if err := uc.DeleteFile(context.Background(), "file-1"); err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}
	if storage.deletedKey != "uploads/file.png" || repo.deletedID != "file-1" {
		t.Fatalf("deleted storage/db = %q/%q, want uploads/file.png/file-1", storage.deletedKey, repo.deletedID)
	}
}
