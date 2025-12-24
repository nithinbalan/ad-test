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
	UpdateModule(module []models.Module) error

	GetActiveUserCount() (int64, error)
	GetUserCounts() (activeCount, passiveCount, totalActiveCount int, err error)
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
		Preload("Theme").
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

	var preference *models.SystemPreference
	if err := repo.db.First(&preference, systemPreference.ID).Error; err != nil {
		return err
	}

	if err := repo.db.Model(&preference).Updates(systemPreference).Error; err != nil {
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

func (repo *SystemPreferenceRepository) GetUserCounts() (activeCount, passiveCount, totalActiveCount int, err error) {
	type Counts struct {
		ActiveCount      int
		PassiveCount     int
		TotalActiveUsers int
	}

	var counts Counts
	err = repo.db.Table("users").
		Select(`
			SUM(CASE WHEN lt.type = 'active' THEN 1 ELSE 0 END) AS active_count,
			SUM(CASE WHEN lt.type = 'passive' THEN 1 ELSE 0 END) AS passive_count,
			SUM(CASE WHEN users.user_status_id IS NOT NULL AND us.is_active = ? THEN 1 ELSE 0 END) AS total_active_users
		`, true).
		Joins("JOIN license_types lt ON lt.id = users.license_type_id").
		Joins("JOIN user_statuses us ON us.id = users.user_status_id").
		Scan(&counts).Error
	if err != nil {
		return 0, 0, 0, err
	}

	return counts.ActiveCount, counts.PassiveCount, counts.TotalActiveUsers, nil
}

func (repo *SystemPreferenceRepository) UpdateModule(module []models.Module) error {

	if err := repo.db.Save(module).Error; err != nil {
		return err
	}

	return nil
}
