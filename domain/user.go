package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents the users table entity
type User struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"type:varchar(100);not null" json:"name"`
	Email           string    `gorm:"type:varchar(150);unique;not null" json:"email"`
	Password        string    `gorm:"type:varchar(255);not null" json:"-"`
	TemperatureUnit string    `gorm:"type:varchar(10);default:'Celsius'" json:"temperature_unit"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// RegisterRequest represents the payload for user registration
type RegisterRequest struct {
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UpdateSettingsRequest represents the payload for updating user preferences
type UpdateSettingsRequest struct {
	TemperatureUnit string `json:"temperature_unit" binding:"required,oneof=Celsius Fahrenheit"`
}

// AuthResponse represents the response containing user info and JWT token
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

// UserRepository defines the database contract for user operations
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
}

// UserUsecase defines the business logic contract for user operations
type UserUsecase interface {
	Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
	UpdateSettings(ctx context.Context, userID uuid.UUID, req *UpdateSettingsRequest) (*User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
}

