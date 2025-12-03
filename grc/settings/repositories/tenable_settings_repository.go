package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type TenableSettingsRepository interface {
	GetFirst() (*models.TenableSettings, error)
	Create(settings *models.TenableSettings) error
	Update(settings *models.TenableSettings) error
	Delete(settings *models.TenableSettings) error
}

type tenableSettingsRepo struct{ db *gorm.DB }

func NewTenableSettingsRepository(db *gorm.DB) TenableSettingsRepository {
	return &tenableSettingsRepo{db}
}

func (repo *tenableSettingsRepo) GetFirst() (*models.TenableSettings, error) {
	var settings models.TenableSettings
	err := repo.db.First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &settings, err
}

func (repo *tenableSettingsRepo) Create(s *models.TenableSettings) error {
	return repo.db.Create(s).Error
}

func (repo *tenableSettingsRepo) Update(s *models.TenableSettings) error {
	return repo.db.Save(s).Error
}

func (repo *tenableSettingsRepo) Delete(settings *models.TenableSettings) error {
	return repo.db.Delete(settings).Error
}
