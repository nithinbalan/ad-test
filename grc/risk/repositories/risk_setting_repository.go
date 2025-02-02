package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type RiskSettingsRepositoryInterface interface {
	Get() (*models.RiskSetting, error)
	Update(setting *models.RiskSetting) (*models.RiskSetting, error)
}

type RiskSettingsRepository struct {
	db *gorm.DB
}

func NewRiskSettingsRepository(db *gorm.DB) *RiskSettingsRepository {
	return &RiskSettingsRepository{
		db: db,
	}
}

func (repo *RiskSettingsRepository) Get() (*models.RiskSetting, error) {
	var riskSetting *models.RiskSetting
	result := repo.db.
		Preload("ReviewRole").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&riskSetting)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskSetting, nil
}

func (repo *RiskSettingsRepository) Update(updatedSetting *models.RiskSetting) (*models.RiskSetting, error) {
	var setting models.RiskSetting
	if err := repo.db.First(&setting).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&setting).Updates(updatedSetting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}
