package database

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

type Migration struct {
	Version string `gorm:"primaryKey;size:64"`
	Name    string `gorm:"size:255;not null"`
}

func RunMigrations(ctx context.Context, db *gorm.DB, dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, file := range files {
		version, name, err := parseMigrationFilename(file)
		if err != nil {
			return err
		}

		applied, err := migrationApplied(ctx, db, version)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := tx.Exec(string(content)).Error; err != nil {
				return err
			}
			return tx.Exec("INSERT INTO schema_migrations (version, name) VALUES (?, ?)", version, name).Error
		})
		if err != nil {
			return fmt.Errorf("apply migration %s: %w", filepath.Base(file), err)
		}
	}

	return nil
}

func migrationApplied(ctx context.Context, db *gorm.DB, version string) (bool, error) {
	var exists bool
	err := db.WithContext(ctx).Raw("SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'schema_migrations')").Scan(&exists).Error
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	var count int64
	err = db.WithContext(ctx).Raw("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count).Error
	return count > 0, err
}

func parseMigrationFilename(file string) (string, string, error) {
	base := filepath.Base(file)
	if !strings.HasSuffix(base, ".sql") {
		return "", "", errors.New("migration must be sql")
	}
	name := strings.TrimSuffix(base, ".sql")
	parts := strings.SplitN(name, "_", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid migration filename %q", base)
	}
	return parts[0], parts[1], nil
}
