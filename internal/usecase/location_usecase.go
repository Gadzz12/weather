package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"go_26_10/domain"
)

type locationUsecase struct {
	locationRepo domain.LocationRepository
}

// NewLocationUsecase creates a new instance of LocationUsecase
func NewLocationUsecase(locationRepo domain.LocationRepository) domain.LocationUsecase {
	return &locationUsecase{
		locationRepo: locationRepo,
	}
}

func (u *locationUsecase) SearchLocations(ctx context.Context, query string) ([]domain.Location, error) {
	trimmedQuery := strings.TrimSpace(query)
	if trimmedQuery == "" {
		return nil, errors.New("parameter pencarian 'q' tidak boleh kosong")
	}

	locations, err := u.locationRepo.Search(ctx, trimmedQuery)
	if err != nil {
		return nil, err
	}

	// Filter ketat: Hanya daerah di Indonesia yang dikembalikan
	indonesianLocations := make([]domain.Location, 0, len(locations))
	for _, loc := range locations {
		if strings.EqualFold(strings.TrimSpace(loc.Country), "indonesia") {
			indonesianLocations = append(indonesianLocations, loc)
		}
	}

	return indonesianLocations, nil
}

func (u *locationUsecase) SaveLocation(ctx context.Context, userID, locationID uuid.UUID) (*domain.UserSavedLocation, error) {
	// 1. Verifikasi keberadaan lokasi
	loc, err := u.locationRepo.FindByID(ctx, locationID)
	if err != nil || loc == nil {
		return nil, errors.New("lokasi dengan ID tersebut tidak ditemukan")
	}

	// 2. Pastikan lokasi yang disimpan berada di wilayah Indonesia
	if !strings.EqualFold(strings.TrimSpace(loc.Country), "indonesia") {
		return nil, errors.New("hanya daerah di Indonesia yang dapat disimpan ke daftar favorit")
	}

	// 3. Cek apakah pengguna sudah menyimpan lokasi ini sebelumnya
	existing, err := u.locationRepo.GetUserSavedLocations(ctx, userID)
	if err == nil {
		for _, s := range existing {
			if s.LocationID == locationID {
				return nil, errors.New("lokasi ini sudah tersimpan dalam daftar favorit Anda")
			}
		}
	}

	// 4. Jika ini lokasi pertama yang disimpan, jadikan default
	isDefault := len(existing) == 0

	savedLocation := &domain.UserSavedLocation{
		ID:         uuid.New(),
		UserID:     userID,
		LocationID: locationID,
		IsDefault:  isDefault,
		CreatedAt:  time.Now(),
		Location:   loc,
	}

	if err := u.locationRepo.SaveUserLocation(ctx, savedLocation); err != nil {
		return nil, err
	}

	return savedLocation, nil
}

func (u *locationUsecase) GetSavedLocations(ctx context.Context, userID uuid.UUID) ([]domain.SavedLocationItemResponse, error) {
	savedList, err := u.locationRepo.GetUserSavedLocations(ctx, userID)
	if err != nil {
		return nil, err
	}

	results := make([]domain.SavedLocationItemResponse, 0, len(savedList))
	for _, it := range savedList {
		var loc domain.Location
		if it.Location != nil {
			loc = *it.Location
		} else {
			// Fallback jika tidak ter-preload
			if fetched, fetchErr := u.locationRepo.FindByID(ctx, it.LocationID); fetchErr == nil && fetched != nil {
				loc = *fetched
			}
		}

		results = append(results, domain.SavedLocationItemResponse{
			ID:        it.ID,
			IsDefault: it.IsDefault,
			Location:  loc,
			CreatedAt: it.CreatedAt,
		})
	}

	return results, nil
}
