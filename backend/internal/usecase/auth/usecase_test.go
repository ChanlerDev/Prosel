package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	domain "github.com/chanler/prosel/backend/internal/domain/auth"
)

type fakeUserRepo struct {
	user *domain.User
	err  error
}

func (r *fakeUserRepo) Create(ctx context.Context, user *domain.User) error {
	r.user = user
	return r.err
}
func (r *fakeUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return r.user, r.err
}
func (r *fakeUserRepo) GetByUsernameOrEmail(ctx context.Context, login string) (*domain.User, error) {
	return r.user, r.err
}
func (r *fakeUserRepo) UpdateProfile(ctx context.Context, user *domain.User) error {
	r.user = user
	return r.err
}
func (r *fakeUserRepo) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	return r.err
}
func (r *fakeUserRepo) UpdateLastLogin(ctx context.Context, userID string) error { return r.err }

type fakeSessionRepo struct {
	session *domain.Session
	err     error
}

func (r *fakeSessionRepo) Create(ctx context.Context, session *domain.Session) error {
	r.session = session
	return r.err
}
func (r *fakeSessionRepo) GetByRefreshTokenHash(ctx context.Context, hash string) (*domain.Session, error) {
	return r.session, r.err
}
func (r *fakeSessionRepo) Revoke(ctx context.Context, id string) error              { return r.err }
func (r *fakeSessionRepo) RevokeAllByUser(ctx context.Context, userID string) error { return r.err }

type fakePasswordHasher struct{}

func (fakePasswordHasher) Hash(password string) (string, error)      { return "hash:" + password, nil }
func (fakePasswordHasher) Compare(hash string, password string) bool { return hash == "hash:"+password }

type fakeTokenService struct{}

func (fakeTokenService) NewAccessToken(user *domain.User) (string, error) {
	return "access." + user.ID, nil
}
func (fakeTokenService) NewRefreshToken() (string, string, error) {
	return "refresh", "refresh-hash", nil
}
func (fakeTokenService) HashRefreshToken(token string) string          { return token + "-hash" }
func (fakeTokenService) ParseAccessToken(token string) (string, error) { return "user-1", nil }

func TestLoginCreatesSessionAndReturnsUser(t *testing.T) {
	userRepo := &fakeUserRepo{user: &domain.User{ID: "user-1", Username: "admin", Email: "a@example.com", PasswordHash: "hash:secret", Role: domain.RoleAdmin, Status: domain.StatusActive}}
	sessionRepo := &fakeSessionRepo{}
	uc := NewAuthUsecase(userRepo, sessionRepo, fakePasswordHasher{}, fakeTokenService{}, time.Hour, 24*time.Hour)

	tokens, user, err := uc.Login(context.Background(), LoginRequest{Login: "admin", Password: "secret"}, ClientMeta{UserAgent: "test", IPAddress: "127.0.0.1"})
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if tokens.AccessToken != "access.user-1" || tokens.RefreshToken != "refresh" {
		t.Fatalf("tokens = %#v", tokens)
	}
	if user.ID != "user-1" {
		t.Fatalf("user = %#v", user)
	}
	if sessionRepo.session == nil || sessionRepo.session.RefreshTokenHash != "refresh-hash" {
		t.Fatalf("session was not created: %#v", sessionRepo.session)
	}
}

func TestLoginRejectsInvalidPassword(t *testing.T) {
	uc := NewAuthUsecase(&fakeUserRepo{user: &domain.User{ID: "user-1", PasswordHash: "hash:secret", Role: domain.RoleAdmin, Status: domain.StatusActive}}, &fakeSessionRepo{}, fakePasswordHasher{}, fakeTokenService{}, time.Hour, 24*time.Hour)

	_, _, err := uc.Login(context.Background(), LoginRequest{Login: "admin", Password: "wrong"}, ClientMeta{})
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("Login() error = %v, want %v", err, domain.ErrInvalidCredentials)
	}
}

func TestRefreshRejectsRevokedSession(t *testing.T) {
	revokedAt := time.Now()
	uc := NewAuthUsecase(&fakeUserRepo{user: &domain.User{ID: "user-1", Role: domain.RoleAdmin, Status: domain.StatusActive}}, &fakeSessionRepo{session: &domain.Session{ID: "session-1", UserID: "user-1", ExpiresAt: time.Now().Add(time.Hour), RevokedAt: &revokedAt}}, fakePasswordHasher{}, fakeTokenService{}, time.Hour, 24*time.Hour)

	_, err := uc.Refresh(context.Background(), "refresh", ClientMeta{})
	if !errors.Is(err, domain.ErrSessionExpired) {
		t.Fatalf("Refresh() error = %v, want %v", err, domain.ErrSessionExpired)
	}
}
