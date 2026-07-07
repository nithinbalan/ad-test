package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type QualysSettingsRepository interface {
	GetFirst() (*models.QualysSettings, error)
	Create(settings *models.QualysSettings) error
	Update(settings *models.QualysSettings) error
	Delete(settings *models.QualysSettings) error
}

type qualysSettingsRepo struct{ db *gorm.DB }

func NewQualysSettingsRepository(db *gorm.DB) QualysSettingsRepository {
	return &qualysSettingsRepo{db}
}

func (repo *qualysSettingsRepo) GetFirst() (*models.QualysSettings, error) {
	var settings models.QualysSettings
	err := repo.db.First(&settings).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return &settings, err
}

func (repo *qualysSettingsRepo) Create(s *models.QualysSettings) error {
	return repo.db.Create(s).Error
}

func (repo *qualysSettingsRepo) Update(s *models.QualysSettings) error {
	return repo.db.Save(s).Error
}

func (repo *qualysSettingsRepo) Delete(settings *models.QualysSettings) error {
	return repo.db.Delete(settings).Error
}