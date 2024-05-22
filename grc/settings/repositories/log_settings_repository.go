package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type LogSettingsRepositoryInterface interface {
	GetLogSettings() (*models.LogSettings, error)
	Update(logSettings *models.LogSettings) error
}

// Implementation
type LogSettingsRepository struct {
	db *gorm.DB
}

func NewLogSettingsRepository(db *gorm.DB) *LogSettingsRepository {
	return &LogSettingsRepository{
		db: db,
	}
}

func (repo *LogSettingsRepository) GetLogSettings() (*models.LogSettings, error) {
	var logSettings *models.LogSettings
	result := repo.db.First(&logSettings)

	if result.Error != nil {
		return nil, result.Error
	}

	return logSettings, nil
}

func (repo *LogSettingsRepository) Update(logSettings *models.LogSettings) error {
	if err := repo.db.Save(logSettings).Error; err != nil {
		return err
	}

	return nil
}
