package auth

import (
	"context"
	"errors"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserDisabled       = errors.New("user disabled")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
	GetByUsernameOrEmail(ctx context.Context, login string) (*User, error)
	UpdateProfile(ctx context.Context, user *User) error
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
	UpdateLastLogin(ctx context.Context, userID string) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByRefreshTokenHash(ctx context.Context, hash string) (*Session, error)
	Revoke(ctx context.Context, id string) error
	RevokeAllByUser(ctx context.Context, userID string) error
}
