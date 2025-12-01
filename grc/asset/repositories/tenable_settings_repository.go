package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type TenableSettingsRepository interface {
	Get() (*models.TenableSettings, error)
	Update(settings *models.TenableSettings) error
	Activate() error
	Deactivate() error
}

type tenableSettingsRepo struct {
	db *gorm.DB
}

func NewTenableSettingsRepository(db *gorm.DB) TenableSettingsRepository {
	return &tenableSettingsRepo{db: db}
}

func (r *tenableSettingsRepo) Get() (*models.TenableSettings, error) {
	var settings models.TenableSettings
	err := r.db.First(&settings).Error
	return &settings, err
}

func (r *tenableSettingsRepo) Update(settings *models.TenableSettings) error {
	return r.db.Save(settings).Error
}

func (r *tenableSettingsRepo) Activate() error {
	return r.db.
		Model(&models.TenableSettings{}).
		Where("id = 1").
		Update("status", "active").Error
}

func (r *tenableSettingsRepo) Deactivate() error {
	return r.db.
		Model(&models.TenableSettings{}).
		Where("id = 1").
		Update("status", "inactive").Error
}
