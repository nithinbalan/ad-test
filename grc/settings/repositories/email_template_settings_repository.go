package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type EmailTemplateSettingsRepositoryInterface interface {
	GetEmailTemplateSettings() (*models.EmailTemplateSettings, error)
	Update(emailTemplateSettings *models.EmailTemplateSettings) error
}

// Implementation
type EmailTemplateSettingsRepository struct {
	db *gorm.DB
}

func NewEmailTemplateSettingsRepository(db *gorm.DB) *EmailTemplateSettingsRepository {
	return &EmailTemplateSettingsRepository{
		db: db,
	}
}

func (repo *EmailTemplateSettingsRepository) GetEmailTemplateSettings() (*models.EmailTemplateSettings, error) {
	var emailTemplateSettings *models.EmailTemplateSettings
	result := repo.db.First(&emailTemplateSettings)
	if result.Error != nil {
		return nil, result.Error
	}
	return emailTemplateSettings, nil
}

func (repo *EmailTemplateSettingsRepository) Update(emailTemplateSettings *models.EmailTemplateSettings) error {
	if err := repo.db.Save(emailTemplateSettings).Error; err != nil {
		return err
	}
	return nil
}
