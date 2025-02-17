package repository

import (
	"github.com/abin-saji-2003/MicroService-Clean-Architecture-Go-/tree/main/user-service/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
}

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{DB: db}
}

func (r *UserRepo) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
