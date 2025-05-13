package repositories

import (
	"errors"
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ActiveDirectorySettingRepositoryInterface interface {
	GetFirst() (*models.ActiveDirectorySetting, error)
	Create(ad *models.ActiveDirectorySetting) error
	Update(ad *models.ActiveDirectorySetting) error
}

// Implementation
type ActiveDirectorySettingRepository struct {
	db *gorm.DB
}

func NewActiveDirectorySettingRepository(db *gorm.DB) *ActiveDirectorySettingRepository {
	return &ActiveDirectorySettingRepository{
		db: db,
	}
}
func (repo *ActiveDirectorySettingRepository) GetFirst() (*models.ActiveDirectorySetting, error) {
	var ad models.ActiveDirectorySetting

	result := repo.db.First(&ad)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil // <- clean way to indicate "no existing record"
	}
	return &ad, result.Error
}

func (repo *ActiveDirectorySettingRepository) Create(ad *models.ActiveDirectorySetting) error {
	return repo.db.Create(ad).Error
}

func (repo *ActiveDirectorySettingRepository) Update(ad *models.ActiveDirectorySetting) error {
	return repo.db.Save(ad).Error
}
