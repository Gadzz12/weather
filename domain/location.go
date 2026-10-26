package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Location represents the locations table entity
type Location struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	City      string    `gorm:"type:varchar(100);not null" json:"city"`
	State     *string   `gorm:"type:varchar(100)" json:"state"`
	Country   string    `gorm:"type:varchar(100);not null" json:"country"`
	Latitude  float64   `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64   `gorm:"type:decimal(11,8);not null" json:"longitude"`
	ImageURL  *string   `gorm:"type:text" json:"image_url"`
}

// UserSavedLocation represents the user_saved_locations table entity
type UserSavedLocation struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	LocationID uuid.UUID `gorm:"type:uuid;not null" json:"location_id"`
	IsDefault  bool      `gorm:"default:false" json:"is_default"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

// SaveLocationRequest represents the payload for adding a saved location
type SaveLocationRequest struct {
	LocationID uuid.UUID `json:"location_id" binding:"required"`
}

// SavedLocationItemResponse represents the formatted response for saved locations
type SavedLocationItemResponse struct {
	ID        uuid.UUID `json:"id"`
	IsDefault bool      `json:"is_default"`
	Location  Location  `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

// LocationRepository defines the database contract for location operations
type LocationRepository interface {
	Search(ctx context.Context, query string) ([]Location, error)
	FindByID(ctx context.Context, id uuid.UUID) (*Location, error)
	Create(ctx context.Context, location *Location) error
	SaveUserLocation(ctx context.Context, savedLocation *UserSavedLocation) error
	GetUserSavedLocations(ctx context.Context, userID uuid.UUID) ([]UserSavedLocation, error)
	GetDefaultUserLocation(ctx context.Context, userID uuid.UUID) (*UserSavedLocation, error)
	DeleteUserSavedLocation(ctx context.Context, userID, locationID uuid.UUID) error
}

// LocationUsecase defines the business logic contract for location operations
type LocationUsecase interface {
	SearchLocations(ctx context.Context, query string) ([]Location, error)
	SaveLocation(ctx context.Context, userID, locationID uuid.UUID) (*UserSavedLocation, error)
	GetSavedLocations(ctx context.Context, userID uuid.UUID) ([]SavedLocationItemResponse, error)
}

