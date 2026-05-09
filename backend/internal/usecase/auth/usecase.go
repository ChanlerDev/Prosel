package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/auth"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash string, password string) bool
}

type TokenService interface {
	NewAccessToken(user *domain.User) (string, error)
	NewRefreshToken() (token string, hash string, err error)
	HashRefreshToken(token string) string
	ParseAccessToken(token string) (userID string, err error)
}

type AuthUsecase struct {
	users                domain.UserRepository
	sessions             domain.SessionRepository
	passwords            PasswordHasher
	tokens               TokenService
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

type LoginRequest struct {
	Login    string
	Password string
}

type ClientMeta struct {
	UserAgent string
	IPAddress string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type UpdateProfileRequest struct {
	DisplayName string
	AvatarURL   string
	Bio         string
}

type ChangePasswordRequest struct {
	OldPassword string
	NewPassword string
}

func NewAuthUsecase(users domain.UserRepository, sessions domain.SessionRepository, passwords PasswordHasher, tokens TokenService, accessTokenDuration time.Duration, refreshTokenDuration time.Duration) *AuthUsecase {
	return &AuthUsecase{users: users, sessions: sessions, passwords: passwords, tokens: tokens, accessTokenDuration: accessTokenDuration, refreshTokenDuration: refreshTokenDuration}
}

func (uc *AuthUsecase) Login(ctx context.Context, req LoginRequest, meta ClientMeta) (*TokenPair, *domain.User, error) {
	user, err := uc.users.GetByUsernameOrEmail(ctx, strings.TrimSpace(req.Login))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, nil, domain.ErrInvalidCredentials
		}
		return nil, nil, err
	}
	if !user.IsActiveAdmin() {
		return nil, nil, domain.ErrUserDisabled
	}
	if !uc.passwords.Compare(user.PasswordHash, req.Password) {
		return nil, nil, domain.ErrInvalidCredentials
	}

	tokens, session, err := uc.issueTokens(user, meta)
	if err != nil {
		return nil, nil, err
	}
	if err := uc.sessions.Create(ctx, session); err != nil {
		return nil, nil, err
	}
	if err := uc.users.UpdateLastLogin(ctx, user.ID); err != nil {
		return nil, nil, err
	}
	return tokens, user, nil
}

func (uc *AuthUsecase) Refresh(ctx context.Context, refreshToken string, meta ClientMeta) (*TokenPair, error) {
	hash := uc.tokens.HashRefreshToken(refreshToken)
	session, err := uc.sessions.GetByRefreshTokenHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil, domain.ErrSessionExpired
		}
		return nil, err
	}
	if !session.IsActive(time.Now().UTC()) {
		return nil, domain.ErrSessionExpired
	}

	user, err := uc.users.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, err
	}
	if !user.IsActiveAdmin() {
		return nil, domain.ErrUserDisabled
	}

	if err := uc.sessions.Revoke(ctx, session.ID); err != nil {
		return nil, err
	}
	tokens, nextSession, err := uc.issueTokens(user, meta)
	if err != nil {
		return nil, err
	}
	if err := uc.sessions.Create(ctx, nextSession); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	hash := uc.tokens.HashRefreshToken(refreshToken)
	session, err := uc.sessions.GetByRefreshTokenHash(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrSessionNotFound) {
			return nil
		}
		return err
	}
	return uc.sessions.Revoke(ctx, session.ID)
}

func (uc *AuthUsecase) GetMe(ctx context.Context, userID string) (*domain.User, error) {
	user, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !user.IsActiveAdmin() {
		return nil, domain.ErrUserDisabled
	}
	return user, nil
}

func (uc *AuthUsecase) UpdateProfile(ctx context.Context, userID string, req UpdateProfileRequest) (*domain.User, error) {
	user, err := uc.GetMe(ctx, userID)
	if err != nil {
		return nil, err
	}
	user.DisplayName = strings.TrimSpace(req.DisplayName)
	user.AvatarURL = strings.TrimSpace(req.AvatarURL)
	user.Bio = strings.TrimSpace(req.Bio)
	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}
	if err := uc.users.UpdateProfile(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, userID string, req ChangePasswordRequest) error {
	user, err := uc.GetMe(ctx, userID)
	if err != nil {
		return err
	}
	if !uc.passwords.Compare(user.PasswordHash, req.OldPassword) {
		return domain.ErrInvalidCredentials
	}
	hash, err := uc.passwords.Hash(req.NewPassword)
	if err != nil {
		return err
	}
	if err := uc.users.UpdatePassword(ctx, userID, hash); err != nil {
		return err
	}
	return uc.sessions.RevokeAllByUser(ctx, userID)
}

func (uc *AuthUsecase) issueTokens(user *domain.User, meta ClientMeta) (*TokenPair, *domain.Session, error) {
	accessToken, err := uc.tokens.NewAccessToken(user)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, refreshHash, err := uc.tokens.NewRefreshToken()
	if err != nil {
		return nil, nil, err
	}

	now := time.Now().UTC()
	session := &domain.Session{
		ID:               newID(),
		UserID:           user.ID,
		RefreshTokenHash: refreshHash,
		UserAgent:        meta.UserAgent,
		IPAddress:        meta.IPAddress,
		ExpiresAt:        now.Add(uc.refreshTokenDuration),
		CreatedAt:        now,
	}
	return &TokenPair{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: int64(uc.accessTokenDuration.Seconds())}, session, nil
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
