package database

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	domain "github.com/chanler/prosel/backend/internal/domain/auth"
)

type UserModel struct {
	ID           string `gorm:"primaryKey;size:36"`
	Username     string `gorm:"uniqueIndex;size:50;not null"`
	Email        string `gorm:"uniqueIndex;size:255;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	DisplayName  string `gorm:"size:100;not null"`
	AvatarURL    string `gorm:"size:500"`
	Bio          string
	Role         string `gorm:"size:20;not null"`
	Status       string `gorm:"size:20;not null"`
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (UserModel) TableName() string { return "users" }

type SessionModel struct {
	ID               string `gorm:"primaryKey;size:36"`
	UserID           string `gorm:"size:36;not null;index"`
	RefreshTokenHash string `gorm:"size:255;not null;uniqueIndex"`
	UserAgent        string
	IPAddress        string `gorm:"size:64"`
	ExpiresAt        time.Time
	RevokedAt        *time.Time
	CreatedAt        time.Time
}

func (SessionModel) TableName() string { return "sessions" }

type UserRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Create(toUserModel(user)).Error
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var model UserModel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error
	return userFromModel(model, err)
}

func (r *UserRepository) GetByUsernameOrEmail(ctx context.Context, login string) (*domain.User, error) {
	var model UserModel
	value := strings.ToLower(strings.TrimSpace(login))
	err := r.db.WithContext(ctx).Where("LOWER(username) = ? OR LOWER(email) = ?", value, value).First(&model).Error
	return userFromModel(model, err)
}

func (r *UserRepository) UpdateProfile(ctx context.Context, user *domain.User) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", user.ID).Updates(map[string]any{
		"display_name": user.DisplayName,
		"avatar_url":   user.AvatarURL,
		"bio":          user.Bio,
		"updated_at":   time.Now().UTC(),
	}).Error
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).Updates(map[string]any{"password_hash": passwordHash, "updated_at": time.Now().UTC()}).Error
}

func (r *UserRepository) UpdateLastLogin(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&UserModel{}).Where("id = ?", userID).Updates(map[string]any{"last_login_at": time.Now().UTC(), "updated_at": time.Now().UTC()}).Error
}

type SessionRepository struct{ db *gorm.DB }

func NewSessionRepository(db *gorm.DB) *SessionRepository { return &SessionRepository{db: db} }

func (r *SessionRepository) Create(ctx context.Context, session *domain.Session) error {
	return r.db.WithContext(ctx).Create(toSessionModel(session)).Error
}

func (r *SessionRepository) GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	var model SessionModel
	err := r.db.WithContext(ctx).Where("refresh_token_hash = ?", hash).First(&model).Error
	return sessionFromModel(model, err)
}

func (r *SessionRepository) Revoke(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Model(&SessionModel{}).Where("id = ? AND revoked_at IS NULL", id).Update("revoked_at", time.Now().UTC()).Error
}

func (r *SessionRepository) RevokeAllByUser(ctx context.Context, userID string) error {
	return r.db.WithContext(ctx).Model(&SessionModel{}).Where("user_id = ? AND revoked_at IS NULL", userID).Update("revoked_at", time.Now().UTC()).Error
}

func userFromModel(model UserModel, err error) (*domain.User, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: model.ID, Username: model.Username, Email: model.Email, PasswordHash: model.PasswordHash, DisplayName: model.DisplayName, AvatarURL: model.AvatarURL, Bio: model.Bio, Role: model.Role, Status: model.Status, LastLoginAt: model.LastLoginAt, CreatedAt: model.CreatedAt, UpdatedAt: model.UpdatedAt}, nil
}

func sessionFromModel(model SessionModel, err error) (*domain.Session, error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.Session{ID: model.ID, UserID: model.UserID, RefreshTokenHash: model.RefreshTokenHash, UserAgent: model.UserAgent, IPAddress: model.IPAddress, ExpiresAt: model.ExpiresAt, RevokedAt: model.RevokedAt, CreatedAt: model.CreatedAt}, nil
}

func toUserModel(user *domain.User) *UserModel {
	return &UserModel{ID: user.ID, Username: user.Username, Email: user.Email, PasswordHash: user.PasswordHash, DisplayName: user.DisplayName, AvatarURL: user.AvatarURL, Bio: user.Bio, Role: user.Role, Status: user.Status, LastLoginAt: user.LastLoginAt}
}

func toSessionModel(session *domain.Session) *SessionModel {
	return &SessionModel{ID: session.ID, UserID: session.UserID, RefreshTokenHash: session.RefreshTokenHash, UserAgent: session.UserAgent, IPAddress: session.IPAddress, ExpiresAt: session.ExpiresAt, RevokedAt: session.RevokedAt, CreatedAt: session.CreatedAt}
}
