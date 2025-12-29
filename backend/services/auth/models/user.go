package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID               string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email            string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash     string         `gorm:"" json:"-"`
	GoogleID         *string        `gorm:"index:idx_users_google_id,unique,where:google_id IS NOT NULL" json:"google_id,omitempty"`
	TelegramID       *int64         `gorm:"index:idx_users_telegram_id,unique,where:telegram_id IS NOT NULL" json:"telegram_id,omitempty"`
	TelegramUsername *string        `gorm:"size:100" json:"telegram_username,omitempty"`
	FirstName        string         `gorm:"size:100" json:"first_name"`
	LastName         string         `gorm:"size:100" json:"last_name"`
	BaseCurrency     string         `gorm:"size:3;default:'USD'" json:"base_currency"`
	AvatarURL        string         `gorm:"size:500" json:"avatar_url,omitempty"`
	IsActive         bool           `gorm:"default:true" json:"is_active"`
	LastLoginAt      *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName returns the table name for GORM
func (User) TableName() string {
	return "users"
}

// RefreshToken represents a refresh token stored in the database
type RefreshToken struct {
	ID        string         `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time      `gorm:"not null" json:"expires_at"`
	IsRevoked bool           `gorm:"default:false" json:"is_revoked"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName returns the table name for GORM
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// OAuthState stores OAuth state for CSRF protection
type OAuthState struct {
	ID        string    `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	State     string    `gorm:"uniqueIndex;not null" json:"state"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName returns the table name for GORM
func (OAuthState) TableName() string {
	return "oauth_states"
}

