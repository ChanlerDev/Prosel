package file

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	domain "github.com/chanler/prosel/backend/internal/domain/file"
)

type Options struct {
	MaxUploadBytes   int64
	AllowedMimeTypes []string
}

type FileUsecase struct {
	files   domain.Repository
	storage domain.StorageProvider
	options Options
}

type UploadRequest struct {
	UploaderID   string
	OriginalName string
	MimeType     string
	Size         int64
	Data         domain.UploadReader
	RefType      string
	RefID        string
}

func NewFileUsecase(files domain.Repository, storage domain.StorageProvider, options Options) *FileUsecase {
	if options.MaxUploadBytes <= 0 {
		options.MaxUploadBytes = 10 << 20
	}
	if len(options.AllowedMimeTypes) == 0 {
		options.AllowedMimeTypes = []string{"image/jpeg", "image/png", "image/gif", "image/webp"}
	}
	return &FileUsecase{files: files, storage: storage, options: options}
}

func (uc *FileUsecase) Upload(ctx context.Context, req UploadRequest) (*domain.FileAsset, error) {
	originalName := strings.TrimSpace(req.OriginalName)
	mimeType := strings.ToLower(strings.TrimSpace(req.MimeType))
	if strings.TrimSpace(req.UploaderID) == "" || originalName == "" || req.Data == nil || !uc.mimeAllowed(mimeType) {
		return nil, domain.ErrInvalidFile
	}
	if req.Size <= 0 || req.Size > uc.options.MaxUploadBytes {
		return nil, domain.ErrFileTooLarge
	}
	if (strings.TrimSpace(req.RefType) == "") != (strings.TrimSpace(req.RefID) == "") {
		return nil, domain.ErrInvalidFile
	}
	if req.RefType != "" && !validRefType(req.RefType) {
		return nil, domain.ErrInvalidFile
	}

	now := time.Now().UTC()
	fileName := newID() + safeExtension(originalName, mimeType)
	objectKey := now.Format("2006/01/02") + "/" + fileName
	stored, err := uc.storage.Put(ctx, objectKey, req.Data, mimeType)
	if err != nil {
		return nil, err
	}
	status := domain.FileStatusOrphan
	if strings.TrimSpace(req.RefType) != "" {
		status = domain.FileStatusAttached
	}
	asset := &domain.FileAsset{
		ID:           newID(),
		UploaderID:   strings.TrimSpace(req.UploaderID),
		OriginalName: originalName,
		FileName:     fileName,
		StorageType:  domain.StorageLocal,
		ObjectKey:    stored.ObjectKey,
		PublicURL:    stored.PublicURL,
		MimeType:     mimeType,
		ByteSize:     req.Size,
		RefType:      strings.TrimSpace(req.RefType),
		RefID:        strings.TrimSpace(req.RefID),
		Status:       status,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.files.Create(ctx, asset); err != nil {
		return nil, err
	}
	return asset, nil
}

func (uc *FileUsecase) AttachToRef(ctx context.Context, fileID string, refType string, refID string) error {
	fileID = strings.TrimSpace(fileID)
	refType = strings.TrimSpace(refType)
	refID = strings.TrimSpace(refID)
	if fileID == "" || refID == "" || !validRefType(refType) {
		return domain.ErrInvalidFile
	}
	return uc.files.UpdateRef(ctx, fileID, refType, refID)
}

func (uc *FileUsecase) ListFiles(ctx context.Context, filter domain.FileFilter) ([]domain.FileAsset, domain.Pagination, error) {
	filter.Search = strings.TrimSpace(filter.Search)
	filter.MimeType = strings.TrimSpace(filter.MimeType)
	filter.Page, filter.PerPage = domain.NormalizePagination(filter.Page, filter.PerPage)
	return uc.files.List(ctx, filter)
}

func (uc *FileUsecase) DeleteFile(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return domain.ErrInvalidFile
	}
	asset, err := uc.files.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if err := uc.storage.Delete(ctx, asset.ObjectKey); err != nil {
		return err
	}
	return uc.files.MarkDeleted(ctx, id)
}

func (uc *FileUsecase) mimeAllowed(mimeType string) bool {
	for _, allowed := range uc.options.AllowedMimeTypes {
		if mimeType == strings.ToLower(strings.TrimSpace(allowed)) {
			return true
		}
	}
	return false
}

func validRefType(refType string) bool {
	return refType == "post" || refType == "note" || refType == "page"
}

func safeExtension(name string, mimeType string) string {
	ext := strings.ToLower(filepath.Ext(name))
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	}
	if ext == "" {
		return ""
	}
	var builder strings.Builder
	for _, r := range ext {
		if r == '.' || unicode.IsLetter(r) || unicode.IsDigit(r) {
			builder.WriteRune(r)
		}
	}
	return builder.String()
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
