package repository

import (
	"context"
	"fmt"
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

	repo.seedDefaultLocations()
	return repo
}

func (r *locationRepository) seedDefaultLocations() {
	stateJKT := "DKI Jakarta"
	stateJatim := "Jawa Timur"
	stateJabar := "Jawa Barat"
	stateDIY := "Daerah Istimewa Yogyakarta"
	stateJateng := "Jawa Tengah"
	stateBali := "Bali"
	stateSumut := "Sumatera Utara"
	stateSulsel := "Sulawesi Selatan"
	stateSumsel := "Sumatera Selatan"
	stateKepri := "Kepulauan Riau"
	stateLampung := "Lampung"
	stateSumbar := "Sumatera Barat"
	stateRiau := "Riau"
	stateAceh := "Aceh"
	stateKaltim := "Kalimantan Timur"
	stateKalsel := "Kalimantan Selatan"
	stateKalbar := "Kalimantan Barat"
	stateSulut := "Sulawesi Utara"
	stateNTB := "Nusa Tenggara Barat"
	stateNTT := "Nusa Tenggara Timur"
	stateMaluku := "Maluku"
	statePapua := "Papua"

	imgJKT := "https://images.unsplash.com/photo-1555899434-94d1368aa7af?auto=format&fit=crop&w=800&q=80"
	imgBDG := "https://images.unsplash.com/photo-1531244605271-e04eb908c68c?auto=format&fit=crop&w=800&q=80"
	imgSBY := "https://images.unsplash.com/photo-1519680916722-d7e00d99cccb?auto=format&fit=crop&w=800&q=80"
	imgJOG := "https://images.unsplash.com/photo-1527742220271-209c3d39f400?auto=format&fit=crop&w=800&q=80"
	imgSMG := "https://images.unsplash.com/photo-1517485677386-57759360363f?auto=format&fit=crop&w=800&q=80"
	imgDPS := "https://images.unsplash.com/photo-1541703009743-c90e30969049?auto=format&fit=crop&w=800&q=80"
	imgMDN := "https://images.unsplash.com/photo-1517485677386-57759360363f?auto=format&fit=crop&w=800&q=80"
	imgMKS := "https://images.unsplash.com/photo-1531244605271-e04eb908c68c?auto=format&fit=crop&w=800&q=80"

	seeds := []*domain.Location{
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			City:      "Jakarta",
			State:     &stateJKT,
			Country:   "Indonesia",
			Latitude:  -6.208763,
			Longitude: 106.845599,
			ImageURL:  &imgJKT,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			City:      "Bandung",
			State:     &stateJabar,
			Country:   "Indonesia",
			Latitude:  -6.917464,
			Longitude: 107.619123,
			ImageURL:  &imgBDG,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000003"),
			City:      "Surabaya",
			State:     &stateJatim,
			Country:   "Indonesia",
			Latitude:  -7.257472,
			Longitude: 112.752090,
			ImageURL:  &imgSBY,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000004"),
			City:      "Yogyakarta",
			State:     &stateDIY,
			Country:   "Indonesia",
			Latitude:  -7.795562,
			Longitude: 110.369492,
			ImageURL:  &imgJOG,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000005"),
			City:      "Semarang",
			State:     &stateJateng,
			Country:   "Indonesia",
			Latitude:  -6.9932,
			Longitude: 110.4204,
			ImageURL:  &imgSMG,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000006"),
			City:      "Denpasar",
			State:     &stateBali,
			Country:   "Indonesia",
			Latitude:  -8.670458,
			Longitude: 115.212629,
			ImageURL:  &imgDPS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000007"),
			City:      "Medan",
			State:     &stateSumut,
			Country:   "Indonesia",
			Latitude:  3.595196,
			Longitude: 98.672226,
			ImageURL:  &imgMDN,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000008"),
			City:      "Makassar",
			State:     &stateSulsel,
			Country:   "Indonesia",
			Latitude:  -5.147665,
			Longitude: 119.432732,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000009"),
			City:      "Palembang",
			State:     &stateSumsel,
			Country:   "Indonesia",
			Latitude:  -2.976074,
			Longitude: 104.775431,
			ImageURL:  &imgJKT,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000010"),
			City:      "Batam",
			State:     &stateKepri,
			Country:   "Indonesia",
			Latitude:  1.1301,
			Longitude: 104.0529,
			ImageURL:  &imgMDN,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000011"),
			City:      "Malang",
			State:     &stateJatim,
			Country:   "Indonesia",
			Latitude:  -7.9797,
			Longitude: 112.6304,
			ImageURL:  &imgSBY,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000012"),
			City:      "Solo",
			State:     &stateJateng,
			Country:   "Indonesia",
			Latitude:  -7.5755,
			Longitude: 110.8243,
			ImageURL:  &imgJOG,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000013"),
			City:      "Bogor",
			State:     &stateJabar,
			Country:   "Indonesia",
			Latitude:  -6.5971,
			Longitude: 106.8060,
			ImageURL:  &imgBDG,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000014"),
			City:      "Bandar Lampung",
			State:     &stateLampung,
			Country:   "Indonesia",
			Latitude:  -5.4352,
			Longitude: 105.2594,
			ImageURL:  &imgJKT,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000015"),
			City:      "Padang",
			State:     &stateSumbar,
			Country:   "Indonesia",
			Latitude:  -0.9471,
			Longitude: 100.4175,
			ImageURL:  &imgMDN,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000016"),
			City:      "Pekanbaru",
			State:     &stateRiau,
			Country:   "Indonesia",
			Latitude:  0.5071,
			Longitude: 101.4478,
			ImageURL:  &imgMDN,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000017"),
			City:      "Banda Aceh",
			State:     &stateAceh,
			Country:   "Indonesia",
			Latitude:  5.5481,
			Longitude: 95.3235,
			ImageURL:  &imgMDN,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000018"),
			City:      "Balikpapan",
			State:     &stateKaltim,
			Country:   "Indonesia",
			Latitude:  -1.2654,
			Longitude: 116.8312,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000019"),
			City:      "Banjarmasin",
			State:     &stateKalsel,
			Country:   "Indonesia",
			Latitude:  -3.3190,
			Longitude: 114.5892,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000020"),
			City:      "Pontianak",
			State:     &stateKalbar,
			Country:   "Indonesia",
			Latitude:  -0.0264,
			Longitude: 109.3425,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000021"),
			City:      "Samarinda",
			State:     &stateKaltim,
			Country:   "Indonesia",
			Latitude:  -0.5019,
			Longitude: 117.1434,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000022"),
			City:      "Manado",
			State:     &stateSulut,
			Country:   "Indonesia",
			Latitude:  1.4748,
			Longitude: 124.8421,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000023"),
			City:      "Mataram",
			State:     &stateNTB,
			Country:   "Indonesia",
			Latitude:  -8.5827,
			Longitude: 116.1154,
			ImageURL:  &imgDPS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000024"),
			City:      "Kupang",
			State:     &stateNTT,
			Country:   "Indonesia",
			Latitude:  -10.1735,
			Longitude: 123.5816,
			ImageURL:  &imgDPS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000025"),
			City:      "Ambon",
			State:     &stateMaluku,
			Country:   "Indonesia",
			Latitude:  -3.6944,
			Longitude: 128.1814,
			ImageURL:  &imgMKS,
		},
		{
			ID:        uuid.MustParse("00000000-0000-0000-0000-000000000026"),
			City:      "Jayapura",
			State:     &statePapua,
			Country:   "Indonesia",
			Latitude:  -2.5330,
			Longitude: 140.7179,
			ImageURL:  &imgMKS,
		},
	}

	if r.db != nil {
		// 1. Hapus seluruh data tersimpan pengguna yang merujuk ke lokasi luar negeri
		_ = r.db.Exec("DELETE FROM user_saved_locations WHERE location_id IN (SELECT id FROM locations WHERE LOWER(TRIM(country)) != 'indonesia')").Error

		// 2. Hapus semua lokasi luar negeri dari database
		_ = r.db.Exec("DELETE FROM locations WHERE LOWER(TRIM(country)) != 'indonesia'").Error
	}

	for _, loc := range seeds {
		r.mockLocations[loc.ID] = loc
		if r.db != nil {
			var existing domain.Location
			if err := r.db.First(&existing, "id = ?", loc.ID).Error; err != nil {
				_ = r.db.Create(loc).Error
			} else {
				_ = r.db.Model(&domain.Location{}).Where("id = ?", loc.ID).Updates(loc).Error
			}
		}
	}

	if r.db != nil {
		// 3. Normalisasi relasi lama jika ID 02 sebelumnya menyimpan Jakarta sebelum reposisi ke ID 01
		_ = r.db.Exec("UPDATE user_saved_locations SET location_id = '00000000-0000-0000-0000-000000000001' WHERE location_id = '00000000-0000-0000-0000-000000000002'").Error
	}
}

func (r *locationRepository) Search(ctx context.Context, query string) ([]domain.Location, error) {
	cleanQ := strings.TrimSpace(query)
	if r.db != nil {
		var locations []domain.Location
		// Hanya cari kota/wilayah di Indonesia
		err := r.db.WithContext(ctx).
			Where("(city ILIKE ? OR COALESCE(state, '') ILIKE ?) AND LOWER(TRIM(country)) = ?", "%"+cleanQ+"%", "%"+cleanQ+"%", "indonesia").
			Find(&locations).Error
		if err == nil {
			return locations, nil
		}
	}

	var results []domain.Location
	r.mu.RLock()
	defer r.mu.RUnlock()
	q := strings.ToLower(cleanQ)
	for _, loc := range r.mockLocations {
		if !strings.EqualFold(loc.Country, "indonesia") {
			continue
		}
		cityName := strings.ToLower(loc.City)
		stateName := ""
		if loc.State != nil {
			stateName = strings.ToLower(*loc.State)
		}
		if q == "" || strings.Contains(cityName, q) || strings.Contains(stateName, q) || q == "indonesia" {
			results = append(results, *loc)
		}
	}
	return results, nil
}

func (r *locationRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Location, error) {
	if r.db != nil {
		var loc domain.Location
		if err := r.db.WithContext(ctx).Where("LOWER(TRIM(country)) = ?", "indonesia").First(&loc, "id = ?", id).Error; err == nil {
			return &loc, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	if loc, ok := r.mockLocations[id]; ok && strings.EqualFold(loc.Country, "indonesia") {
		return loc, nil
	}
	return nil, fmt.Errorf("location not found")
}

func (r *locationRepository) Create(ctx context.Context, location *domain.Location) error {
	if location.ID == uuid.Nil {
		location.ID = uuid.New()
	}

	// Validasi integritas: hanya daerah di Indonesia yang didukung
	if !strings.EqualFold(strings.TrimSpace(location.Country), "indonesia") {
		return fmt.Errorf("hanya daerah di wilayah Indonesia yang didukung")
	}

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
	// Pastikan lokasi yang disimpan berada di Indonesia
	if savedLocation.Location != nil && !strings.EqualFold(savedLocation.Location.Country, "indonesia") {
		return fmt.Errorf("hanya daerah di Indonesia yang dapat disimpan")
	}

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
	if savedCopy.Location != nil && !strings.EqualFold(savedCopy.Location.Country, "indonesia") {
		return fmt.Errorf("hanya daerah di Indonesia yang dapat disimpan")
	}
	r.mockUserLocations[savedLocation.UserID] = append(r.mockUserLocations[savedLocation.UserID], &savedCopy)
	return nil
}

func (r *locationRepository) GetUserSavedLocations(ctx context.Context, userID uuid.UUID) ([]domain.UserSavedLocation, error) {
	if r.db != nil {
		var list []domain.UserSavedLocation
		if err := r.db.WithContext(ctx).
			Joins("JOIN locations ON locations.id = user_saved_locations.location_id").
			Where("user_saved_locations.user_id = ? AND LOWER(TRIM(locations.country)) = ?", userID, "indonesia").
			Preload("Location").
			Find(&list).Error; err == nil && len(list) > 0 {
			return list, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	ptrSlice := r.mockUserLocations[userID]
	var results []domain.UserSavedLocation
	for _, saved := range ptrSlice {
		if saved.Location != nil && !strings.EqualFold(saved.Location.Country, "indonesia") {
			continue
		}
		results = append(results, *saved)
	}
	return results, nil
}

func (r *locationRepository) GetDefaultUserLocation(ctx context.Context, userID uuid.UUID) (*domain.UserSavedLocation, error) {
	if r.db != nil {
		var saved domain.UserSavedLocation
		if err := r.db.WithContext(ctx).
			Joins("JOIN locations ON locations.id = user_saved_locations.location_id").
			Where("user_saved_locations.user_id = ? AND user_saved_locations.is_default = ? AND LOWER(TRIM(locations.country)) = ?", userID, true, "indonesia").
			Preload("Location").
			First(&saved).Error; err == nil {
			return &saved, nil
		}
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, saved := range r.mockUserLocations[userID] {
		if saved.IsDefault {
			if saved.Location != nil && !strings.EqualFold(saved.Location.Country, "indonesia") {
				continue
			}
			itCopy := *saved
			return &itCopy, nil
		}
	}
	return nil, fmt.Errorf("no default location found")
}

func (r *locationRepository) DeleteUserSavedLocation(ctx context.Context, userID, locationID uuid.UUID) error {
	if r.db != nil {
		_ = r.db.WithContext(ctx).Where("user_id = ? AND location_id = ?", userID, locationID).Delete(&domain.UserSavedLocation{}).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	locs := r.mockUserLocations[userID]
	for i, saved := range locs {
		if saved.LocationID == locationID {
			r.mockUserLocations[userID] = append(locs[:i], locs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("user location not found")
}
