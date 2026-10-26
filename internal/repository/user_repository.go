package repository

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go_26_10/domain"
)

type userRepository struct {
	db        *gorm.DB
	mockUsers map[uuid.UUID]*domain.User
	mu        sync.RWMutex
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db:        db,
		mockUsers: make(map[uuid.UUID]*domain.User),
	}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	if r.db != nil {
		return r.db.WithContext(ctx).Create(user).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.mockUsers {
		if strings.EqualFold(u.Email, user.Email) {
			return errors.New("email sudah terdaftar")
		}
	}

	userCopy := *user
	r.mockUsers[user.ID] = &userCopy
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	if r.db != nil {
		var user domain.User
		if err := r.db.WithContext(ctx).Where("LOWER(email) = LOWER(?)", email).First(&user).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.mockUsers {
		if strings.EqualFold(u.Email, email) {
			userCopy := *u
			return &userCopy, nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if r.db != nil {
		var user domain.User
		if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &user, nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, exists := r.mockUsers[id]
	if !exists {
		return nil, gorm.ErrRecordNotFound
	}

	userCopy := *u
	return &userCopy, nil
}

func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	if r.db != nil {
		return r.db.WithContext(ctx).Save(user).Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.mockUsers[user.ID]; !exists {
		return gorm.ErrRecordNotFound
	}

	userCopy := *user
	r.mockUsers[user.ID] = &userCopy
	return nil
}
