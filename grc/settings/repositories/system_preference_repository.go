package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type SystemPreferenceRepositoryInterface interface {
	GetSystemPreference() (*models.SystemPreference, error)
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

func (repo *SystemPreferenceRepository) GetSystemPreference() (*models.SystemPreference, error) {
	var systemPreference *models.SystemPreference
	result := repo.db.First(&systemPreference)

	if result.Error != nil {
		return nil, result.Error
	}

	return systemPreference, nil
}

func (repo *SystemPreferenceRepository) Update(systemPreference *models.SystemPreference) error {
	if err := repo.db.Save(systemPreference).Error; err != nil {
		return err
	}

	return nil
}
