package repositories

import (
	"jethur/grc/models"
	"jethur/grc/settings/dtos"

	"gorm.io/gorm"
)

type AISettingRepositoryInterface interface {
	UpdateSystemPreferenceAI(enabled bool, updatedByID uint) error

	GetAll() ([]*models.AIModule, error)
	GetFeatureListByModuleID(moduleID int) ([]dtos.AIModuleFeatureListDTO, error)
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

func (repo *AISettingRepository) GetFeatureListByModuleID(moduleID int) ([]dtos.AIModuleFeatureListDTO, error) {
	var items []dtos.AIModuleFeatureListDTO

	err := repo.db.Table("ai_module_features").
		Select(`
			ai_module_features.id,
			ai_module_features.title_en,
			ai_module_features.title_ar,
			ai_module_features.description_en,
			ai_module_features.description_ar,
			ai_module_features.[order],
			knowledge_center_resource_statuses.id AS status_id,
			knowledge_center_resource_statuses.title_en AS status_title_en,
			knowledge_center_resource_statuses.title_ar AS status_title_ar,
			knowledge_center_resource_statuses.status AS status,
			knowledge_center_resource_statuses.color_code AS color_code
		`).
		Joins("LEFT JOIN knowledge_center_resource_statuses ON knowledge_center_resource_statuses.id = ai_module_features.knowledge_center_resource_status_id").
		Where("ai_module_features.ai_module_id = ?", moduleID).
		Order("ai_module_features.[order]").
		Scan(&items).Error

	return items, err
}
