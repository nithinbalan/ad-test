package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AuthRepositoryInterface interface {
	CreateUserLogin(userLogin *models.UserLogin) error
}

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (repo *AuthRepository) CreateUserLogin(userLogin *models.UserLogin) error {
	if err := repo.db.Create(userLogin).Error; err != nil {
		return err
	}
	return nil
}
