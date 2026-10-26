package repository

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go_26_10/domain"
)

type locationRepository struct {
	db                *gorm.DB
	mockLocations     map[uuid.UUID]*domain.Location
	mockUserLocations map[uuid.UUID][]*domain.UserSavedLocation
	mu                sync.RWMutex
}

// NewLocationRepository creates a new instance of LocationRepository
func NewLocationRepository(db *gorm.DB) domain.LocationRepository {
	repo := &locationRepository{
		db:                db,
		mockLocations:     make(map[uuid.UUID]*domain.Location),
		mockUserLocations: make(map[uuid.UUID][]*domain.UserSavedLocation),
	}

	// Inisialisasi daftar kota bawaan untuk fallback pencarian dan demo
	repo.seedDefaultLocations()
	return repo
}

func (r *locationRepository) seedDefaultLocations() {
	stateCA := "CA"
	stateNY := "NY"
	stateJKT := "DKI Jakarta"
	stateJatim := "Jawa Timur"
	stateJabar := "Jawa Barat"
	imgLA := "https://images.unsplash.com/photo-1580655653885-65763b2597d0?auto=format&fit=crop&w=800&q=80"
	imgJKT := "https://images.unsplash.com/photo-1555899434-94d1368aa7af?auto=format&fit=crop&w=800&q=80"
	imgTKO := "https://images.unsplash.com/photo-1503899036084-c55cdd92da26?auto=format&fit=crop&w=800&q=80"
	imgLDN := "https://images.unsplash.com/photo-1513635269975-59663e0ac1ad?auto=format&fit=crop&w=800&q=80"
	imgNYC := "https://images.unsplash.com/photo-1496442226666-8d4d0e62e6e9?auto=format&fit=crop&w=800&q=80"

	seeds := []*domain.Location{
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			City:      "Los Angeles",
			State:     &stateCA,
			Country:   "USA",
			Latitude:  34.052235,
			Longitude: -118.243683,
			ImageURL:  &imgLA,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			City:      "Jakarta",
			State:     &stateJKT,
			Country:   "Indonesia",
			Latitude:  -6.208763,
			Longitude: 106.845599,
			ImageURL:  &imgJKT,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			City:      "Tokyo",
			State:     nil,
			Country:   "Japan",
			Latitude:  35.676192,
			Longitude: 139.650311,
			ImageURL:  &imgTKO,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			City:      "London",
			State:     nil,
			Country:   "United Kingdom",
			Latitude:  51.507351,
			Longitude: -0.127758,
			ImageURL:  &imgLDN,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			City:      "New York",
			State:     &stateNY,
			Country:   "USA",
			Latitude:  40.712776,
			Longitude: -74.005974,
			ImageURL:  &imgNYC,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			City:      "Surabaya",
			State:     &stateJatim,
			Country:   "Indonesia",
			Latitude:  -7.257472,
			Longitude: 112.752090,
			ImageURL:  &imgJKT,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000007"),
			City:      "Bandung",
			State:     &stateJabar,
			Country:   "Indonesia",
			Latitude:  -6.917464,
			Longitude: 107.619123,
			ImageURL:  &imgJKT,
		},
	}

	for _, s := range seeds {
		r.mockLocations[s.ID] = s
		if r.db != nil {
			var existing domain.Location
			if err := r.db.First(&existing, "id = ?", s.ID).Error; err != nil {
				_ = r.db.Create(s).Error
			}
		}
	}
}

func (r *locationRepository) Search(ctx context.Context, query string) ([]domain.Location, error) {
	if r.db != nil {
		var locations []domain.Location
		err := r.db.WithContext(ctx).Where("city ILIKE ? OR country ILIKE ?", "%"+query+"%", "%"+query+"%").Find(&locations).Error
		if err == nil && len(locations) > 0 {
			return locations, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []domain.Location
	lowerQ := strings.ToLower(strings.TrimSpace(query))

	for _, loc := range r.mockLocations {
		matchCity := strings.Contains(strings.ToLower(loc.City), lowerQ)
		matchCountry := strings.Contains(strings.ToLower(loc.Country), lowerQ)
		matchState := loc.State != nil && strings.Contains(strings.ToLower(*loc.State), lowerQ)

		if matchCity || matchCountry || matchState {
			results = append(results, *loc)
		}
	}

	return results, nil
}

func (r *locationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	if r.db != nil {
		var loc domain.Location
		if err := r.db.WithContext(ctx).First(&loc, "id = ?", id).Error; err == nil {
			return &loc, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	loc, exists := r.mockLocations[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	locCopy := *loc
	return &locCopy, nil
}

func (r *locationRepository) Create(ctx context.Context, location *domain.Location) error {
	if r.db != nil {
		if err := r.db.WithContext(ctx).Create(location).Error; err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	locCopy := *location
	r.mockLocations[location.ID] = &locCopy
	return nil
}

func (r *locationRepository) SaveUserLocation(ctx context.Context, savedLocation *domain.UserSavedLocation) error {
	if r.db != nil {
		if err := r.db.WithContext(ctx).Create(savedLocation).Error; err == nil {
			return nil
		}
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	savedCopy := *savedLocation
	if savedCopy.Location == nil {
		if loc, exists := r.mockLocations[savedCopy.LocationID]; exists {
			savedCopy.Location = loc
		}
	}

	r.mockUserLocations[savedCopy.UserID] = append(r.mockUserLocations[savedCopy.UserID], &savedCopy)
	return nil
}

func (r *locationRepository) GetUserSavedLocations(ctx context.Context, userID uuid.UUID) ([]domain.UserSavedLocation, error) {
	if r.db != nil {
		var list []domain.UserSavedLocation
		if err := r.db.WithContext(ctx).Preload("Location").Where("user_id = ?", userID).Find(&list).Error; err == nil && len(list) > 0 {
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	savedList, exists := r.mockUserLocations[userID]
	if !exists {
		return []domain.UserSavedLocation{}, nil
	}

	var results []domain.UserSavedLocation
	for _, it := range savedList {
		results = append(results, *it)
	}
	return results, nil
}

func (r *locationRepository) GetDefaultUserLocation(ctx context.Context, userID uuid.UUID) (*domain.UserSavedLocation, error) {
	if r.db != nil {
		var saved domain.UserSavedLocation
		if err := r.db.WithContext(ctx).Preload("Location").Where("user_id = ? AND is_default = ?", userID, true).First(&saved).Error; err == nil {
			return &saved, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	savedList, exists := r.mockUserLocations[userID]
	if !exists || len(savedList) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, it := range savedList {
		if it.IsDefault {
			itCopy := *it
			return &itCopy, nil
		}
	}

	first := *savedList[0]
	return &first, nil
}

func (r *locationRepository) DeleteUserSavedLocation(ctx context.Context, userID, locationID uuid.UUID) error {
	if r.db != nil {
		_ = r.db.WithContext(ctx).Where("user_id = ? AND location_id = ?", userID, locationID).Delete(&domain.UserSavedLocation{}).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	savedList, exists := r.mockUserLocations[userID]
	if !exists {
		return nil
	}

	var updated []*domain.UserSavedLocation
	for _, it := range savedList {
		if it.LocationID != locationID {
			updated = append(updated, it)
		}
	}
	r.mockUserLocations[userID] = updated
	return nil
}
