package file

import (
	"errors"
	"io"
	"time"
)

type StorageType string

type FileStatus string

const (
	StorageLocal StorageType = "local"
	StorageS3    StorageType = "s3"

	FileStatusAttached FileStatus = "attached"
	FileStatusOrphan   FileStatus = "orphan"
	FileStatusDeleted  FileStatus = "deleted"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrInvalidFile  = errors.New("invalid file")
	ErrFileTooLarge = errors.New("file too large")
)

type FileAsset struct {
	ID           string
	UploaderID   string
	OriginalName string
	FileName     string
	StorageType  StorageType
	ObjectKey    string
	PublicURL    string
	MimeType     string
	ByteSize     int64
	Width        *int
	Height       *int
	RefType      string
	RefID        string
	Status       FileStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type StoredObject struct {
	ObjectKey string
	PublicURL string
}

type UploadReader interface {
	io.Reader
}

func (s FileStatus) Valid() bool {
	return s == FileStatusAttached || s == FileStatusOrphan || s == FileStatusDeleted
}
