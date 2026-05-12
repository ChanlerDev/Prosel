package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	domain "github.com/chanler/prosel/backend/internal/domain/file"
)

type LocalStorage struct {
	rootDir   string
	publicURL string
}

func NewLocalStorage(rootDir string, publicURL string) *LocalStorage {
	rootDir = strings.TrimSpace(rootDir)
	if rootDir == "" {
		rootDir = "uploads"
	}
	publicURL = strings.TrimRight(strings.TrimSpace(publicURL), "/")
	if publicURL == "" {
		publicURL = "/uploads"
	}
	return &LocalStorage{rootDir: rootDir, publicURL: publicURL}
}

func (s *LocalStorage) Put(ctx context.Context, objectKey string, data domain.UploadReader, contentType string) (*domain.StoredObject, error) {
	cleanKey := filepath.Clean(objectKey)
	if strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return nil, domain.ErrInvalidFile
	}
	path := filepath.Join(s.rootDir, cleanKey)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if _, err := io.Copy(file, data); err != nil {
		return nil, err
	}
	publicPath := strings.TrimLeft(filepath.ToSlash(cleanKey), "/")
	return &domain.StoredObject{ObjectKey: publicPath, PublicURL: s.publicURL + "/" + publicPath}, nil
}

func (s *LocalStorage) Delete(ctx context.Context, objectKey string) error {
	cleanKey := filepath.Clean(objectKey)
	if strings.HasPrefix(cleanKey, "..") || filepath.IsAbs(cleanKey) {
		return domain.ErrInvalidFile
	}
	path := filepath.Join(s.rootDir, cleanKey)
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
