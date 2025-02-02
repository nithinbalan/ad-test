package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interface for RiskSettingsRepository
type RiskSettingsRepositoryInterface interface {
	Get() (*models.RiskSetting, error)
	Update(setting *models.RiskSetting) (*models.RiskSetting, error)
}

// Implementation of RiskSettingsRepository
type RiskSettingsRepository struct {
	db *gorm.DB
}

// Constructor
func NewRiskSettingsRepository(db *gorm.DB) *RiskSettingsRepository {
	return &RiskSettingsRepository{
		db: db,
	}
}

// GET: Fetch the Single Risk Setting with Role Title
func (repo *RiskSettingsRepository) Get() (*models.RiskSetting, error) {
	var riskSetting *models.RiskSetting
	result := repo.db.
		Preload("ReviewRole"). // Preload role details
		First(&riskSetting)    // Fetch the first (and only) risk setting

	if result.Error != nil {
		return nil, result.Error
	}

	return riskSetting, nil
}

// UPDATE: Update the Single Risk Setting
func (repo *RiskSettingsRepository) Update(updatedSetting *models.RiskSetting) (*models.RiskSetting, error) {
	// Assuming there's only one record, we'll fetch the first one
	var setting models.RiskSetting
	if err := repo.db.First(&setting).Error; err != nil {
		return nil, err
	}

	// Apply updates
	if err := repo.db.Model(&setting).Updates(updatedSetting).Error; err != nil {
		return nil, err
	}

	return &setting, nil
}
