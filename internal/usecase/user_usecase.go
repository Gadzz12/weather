package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"go_26_10/domain"
	"go_26_10/internal/middleware"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret string
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo domain.UserRepository, jwtSecret string) domain.UserUsecase {
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (u *userUsecase) generateToken(user *domain.User) (string, error) {
	claims := middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.jwtSecret))
}

func (u *userUsecase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	existing, _ := u.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("email sudah terdaftar dalam sistem")
	}

	// 2. Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi kata sandi: %w", err)
	}

	newUser := &domain.User{
		ID:              uuid.New(),
		Name:            strings.TrimSpace(req.Name),
		Email:           strings.ToLower(strings.TrimSpace(req.Email)),
		Password:        string(hashedPassword),
		TemperatureUnit: "Celsius",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// 3. Simpan ke database
	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, fmt.Errorf("gagal menyimpan data pengguna: %w", err)
	}

	// 4. Generate JWT token
	token, err := u.generateToken(newUser)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token autentikasi: %w", err)
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *newUser,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, strings.TrimSpace(req.Email))
	if err != nil {
		return nil, errors.New("email atau kata sandi tidak valid")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("email atau kata sandi tidak valid")
	}

	token, err := u.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token autentikasi: %w", err)
	}

	return &domain.AuthResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *userUsecase) UpdateSettings(ctx context.Context, userID uuid.UUID, req *domain.UpdateSettingsRequest) (*domain.User, error) {
	unit := req.TemperatureUnit
	if !strings.EqualFold(unit, "Celsius") && !strings.EqualFold(unit, "Fahrenheit") {
		return nil, errors.New("unit suhu harus 'Celsius' atau 'Fahrenheit'")
	}

	if strings.EqualFold(unit, "Celsius") {
		unit = "Celsius"
	} else {
		unit = "Fahrenheit"
	}

	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	user.TemperatureUnit = unit
	user.UpdatedAt = time.Now()

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("gagal memperbarui pengaturan preferensi: %w", err)
	}

	return user, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	return user, nil
}
