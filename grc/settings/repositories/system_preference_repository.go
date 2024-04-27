package repositories

import (
	"jethur/grc/settings/models"

	"gorm.io/gorm"
)

// Interfaces
type SystemPreferenceRepositoryInterface interface {
	GetByID(id uint) (*models.SystemPreference, error)
	Update(systemPreference *models.SystemPreference) error
}

// implementation
type SystemPreferenceRepository struct {
	db *gorm.DB
}

func NewSystemPreferenceRepository(db *gorm.DB) *SystemPreferenceRepository {
	return &SystemPreferenceRepository{
		db: db,
	}
}

func (repo *SystemPreferenceRepository) GetByID(id uint) (*models.SystemPreference, error) {
	var status *models.SystemPreference
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *SystemPreferenceRepository) Update(systemPreference *models.SystemPreference) error {
	if err := repo.db.Save(systemPreference).Error; err != nil {
		return err
	}

	return nil
}
