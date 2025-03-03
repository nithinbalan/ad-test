package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type SystemPreferenceRepositoryInterface interface {
	GetSystemPreference() (*models.SystemPreference, error)
	GetByID(id uint) (*models.SystemPreference, error)
	Update(systemPreference *models.SystemPreference) error

	GetActiveUserCount() (int64, error)
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
	result := repo.db.
		Preload("Timezone").
		Preload("Currency").
		First(&systemPreference)

	if result.Error != nil {
		return nil, result.Error
	}

	return systemPreference, nil
}

func (repo *SystemPreferenceRepository) GetByID(id uint) (*models.SystemPreference, error) {
	var systemPreference *models.SystemPreference

	result := repo.db.First(&systemPreference, id)
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

func (repo *SystemPreferenceRepository) GetActiveUserCount() (int64, error) {
	var count int64
	err := repo.db.Model(&models.User{}).
		Joins("JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Where("user_statuses.is_active = ?", true).
		Count(&count).
		Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
