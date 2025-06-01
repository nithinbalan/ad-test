package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AISettingRepositoryInterface interface {
	UpdateSystemPreferenceAI(enabled bool, updatedByID uint) error

	GetAll() ([]*models.AIModule, error)
	GetByModuleID(moduleID int) ([]*models.AIModuleFeature, error)
}

type AISettingRepository struct {
	db *gorm.DB
}

func NewAISettingRepository(db *gorm.DB) *AISettingRepository {
	return &AISettingRepository{
		db: db,
	}
}

func (repo *AISettingRepository) UpdateSystemPreferenceAI(enabled bool, updatedByID uint) error {
	return repo.db.Exec(`
		UPDATE system_preferences 
		SET is_ai_enable = ?, updated_by_id = ?, updated_at = GETDATE()
		WHERE id = 1
	`, enabled, updatedByID).Error
}

func (repo *AISettingRepository) GetAll() ([]*models.AIModule, error) {
	var modules []*models.AIModule
	err := repo.db.Order("[order]").Find(&modules).Error
	return modules, err
}

func (repo *AISettingRepository) GetByModuleID(moduleID int) ([]*models.AIModuleFeature, error) {
	var features []*models.AIModuleFeature
	err := repo.db.Where("ai_module_id = ?", moduleID).Order("[order]").Find(&features).Error
	return features, err
}
