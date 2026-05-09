package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	domain "github.com/chanler/prosel/backend/internal/domain/auth"
	"github.com/chanler/prosel/backend/internal/interfaces/http/middleware"
	"github.com/chanler/prosel/backend/internal/interfaces/http/response"
	usecase "github.com/chanler/prosel/backend/internal/usecase/auth"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Login(ctx context.Context, req usecase.LoginRequest, meta usecase.ClientMeta) (*usecase.TokenPair, *domain.User, error)
	Refresh(ctx context.Context, refreshToken string, meta usecase.ClientMeta) (*usecase.TokenPair, error)
	Logout(ctx context.Context, refreshToken string) error
	GetMe(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID string, req usecase.UpdateProfileRequest) (*domain.User, error)
	ChangePassword(ctx context.Context, userID string, req usecase.ChangePasswordRequest) error
}

type AuthHandler struct {
	service AuthService
}

func NewAuthHandler(service AuthService) *AuthHandler { return &AuthHandler{service: service} }

func (h *AuthHandler) RegisterPublicRoutes(group *gin.RouterGroup) {
	group.POST("/auth/login", h.login)
	group.POST("/auth/refresh", h.refresh)
}

func (h *AuthHandler) RegisterProtectedRoutes(api *gin.RouterGroup, admin *gin.RouterGroup) {
	api.POST("/auth/logout", h.logout)
	api.GET("/auth/me", h.me)
	admin.PATCH("/profile", h.updateProfile)
	admin.PATCH("/password", h.changePassword)
}

type loginRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

type tokenResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	ExpiresIn    int64         `json:"expiresIn"`
	User         *userResponse `json:"user,omitempty"`
}

type userResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
	Bio         string `json:"bio,omitempty"`
	Role        string `json:"role"`
}

type updateProfileRequest struct {
	DisplayName string `json:"displayName"`
	AvatarURL   string `json:"avatarUrl"`
	Bio         string `json:"bio"`
}

type changePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8"`
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid login request", nil)
		return
	}
	tokens, user, err := h.service.Login(c.Request.Context(), usecase.LoginRequest{Login: req.Login, Password: req.Password}, clientMeta(c))
	if err != nil {
		h.handleAuthError(c, err)
		return
	}
	setRefreshTokenCookie(c, tokens.RefreshToken)
	response.OK(c, tokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, ExpiresIn: tokens.ExpiresIn, User: toUserResponse(user)})
}

func (h *AuthHandler) refresh(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid refresh request", nil)
		return
	}
	if req.RefreshToken == "" {
		req.RefreshToken = refreshTokenCookie(c)
	}
	if req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required", nil)
		return
	}
	tokens, err := h.service.Refresh(c.Request.Context(), req.RefreshToken, clientMeta(c))
	if err != nil {
		h.handleAuthError(c, err)
		return
	}
	setRefreshTokenCookie(c, tokens.RefreshToken)
	response.OK(c, tokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken, ExpiresIn: tokens.ExpiresIn})
}

func (h *AuthHandler) logout(c *gin.Context) {
	var req refreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid logout request", nil)
		return
	}
	if req.RefreshToken == "" {
		req.RefreshToken = refreshTokenCookie(c)
	}
	if req.RefreshToken == "" {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Refresh token is required", nil)
		return
	}
	if err := h.service.Logout(c.Request.Context(), req.RefreshToken); err != nil {
		h.handleAuthError(c, err)
		return
	}
	clearRefreshTokenCookie(c)
	response.OK(c, map[string]bool{"ok": true})
}

func (h *AuthHandler) me(c *gin.Context) {
	user, err := h.service.GetMe(c.Request.Context(), middleware.CurrentUserID(c))
	if err != nil {
		h.handleAuthError(c, err)
		return
	}
	response.OK(c, toUserResponse(user))
}

func (h *AuthHandler) updateProfile(c *gin.Context) {
	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid profile request", nil)
		return
	}
	user, err := h.service.UpdateProfile(c.Request.Context(), middleware.CurrentUserID(c), usecase.UpdateProfileRequest{DisplayName: req.DisplayName, AvatarURL: req.AvatarURL, Bio: req.Bio})
	if err != nil {
		h.handleAuthError(c, err)
		return
	}
	response.OK(c, toUserResponse(user))
}

func (h *AuthHandler) changePassword(c *gin.Context) {
	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid password request", nil)
		return
	}
	if err := h.service.ChangePassword(c.Request.Context(), middleware.CurrentUserID(c), usecase.ChangePasswordRequest{OldPassword: req.OldPassword, NewPassword: req.NewPassword}); err != nil {
		h.handleAuthError(c, err)
		return
	}
	response.OK(c, map[string]bool{"ok": true})
}

func (h *AuthHandler) handleAuthError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidCredentials):
		response.Error(c, http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid login or password", nil)
	case errors.Is(err, domain.ErrSessionExpired):
		response.Error(c, http.StatusUnauthorized, "SESSION_EXPIRED", "Session expired", nil)
	case errors.Is(err, domain.ErrUserDisabled):
		response.Error(c, http.StatusForbidden, "USER_DISABLED", "User is disabled", nil)
	case errors.Is(err, domain.ErrUserNotFound):
		response.Error(c, http.StatusNotFound, "USER_NOT_FOUND", "User not found", nil)
	default:
		response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Authentication request failed", nil)
	}
}

func clientMeta(c *gin.Context) usecase.ClientMeta {
	return usecase.ClientMeta{UserAgent: c.GetHeader("User-Agent"), IPAddress: c.ClientIP()}
}

func toUserResponse(user *domain.User) *userResponse {
	return &userResponse{ID: user.ID, Username: user.Username, Email: user.Email, DisplayName: user.DisplayName, AvatarURL: user.AvatarURL, Bio: strings.TrimSpace(user.Bio), Role: user.Role}
}

func refreshTokenCookie(c *gin.Context) string {
	value, err := c.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return value
}

func setRefreshTokenCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", token, 60*60*24*7, "/api/v1/auth", "", false, true)
}

func clearRefreshTokenCookie(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/api/v1/auth", "", false, true)
}
