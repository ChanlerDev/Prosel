package auth

import "time"

const (
	RoleAdmin      = "admin"
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	DisplayName  string
	AvatarURL    string
	Bio          string
	Role         string
	Status       string
	LastLoginAt  *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (u User) IsActiveAdmin() bool {
	return u.Role == RoleAdmin && u.Status == StatusActive
}
