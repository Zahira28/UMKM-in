package repository

import (
	"errors"

	"backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	FindByEmail(email string) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindByEmailOrUsername(identifier string) (*model.User, error)
	FindByID(id uuid.UUID) (*model.User, error)
	FindByGoogleID(googleID string) (*model.User, error)
	Update(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("LOWER(email) = LOWER(?)", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("LOWER(username) = LOWER(?)", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmailOrUsername(identifier string) (*model.User, error) {
	var user model.User
	err := r.db.Where("LOWER(email) = LOWER(?) OR LOWER(username) = LOWER(?)", identifier, identifier).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	err := r.db.Where("google_id = ?", googleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}
