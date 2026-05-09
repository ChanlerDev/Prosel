package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	domain "github.com/chanler/prosel/backend/internal/domain/auth"
)

type JWTService struct {
	secret   []byte
	issuer   string
	duration time.Duration
}

type claims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTService(secret string, issuer string, duration time.Duration) *JWTService {
	return &JWTService{secret: []byte(secret), issuer: issuer, duration: duration}
}

func (s *JWTService) NewAccessToken(user *domain.User) (string, error) {
	now := time.Now().UTC()
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.duration)),
		},
	}).SignedString(s.secret)
}

func (s *JWTService) NewRefreshToken() (string, string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", "", err
	}
	token := base64.RawURLEncoding.EncodeToString(bytes)
	return token, s.HashRefreshToken(token), nil
}

func (s *JWTService) HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (s *JWTService) ParseAccessToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &claims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	}, jwt.WithIssuer(s.issuer))
	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(*claims)
	if !ok || !token.Valid || claims.Subject == "" {
		return "", errors.New("invalid token")
	}
	return claims.Subject, nil
}
