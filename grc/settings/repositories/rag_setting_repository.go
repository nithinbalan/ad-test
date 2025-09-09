package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type RAGSettingRepositoryInterface interface {
	GetAllRAGSettings() ([]*models.RagSetting, error)
}

// implementation
type RAGSettingRepository struct {
	db *gorm.DB
}

func NewRAGSettingRepository(db *gorm.DB) *RAGSettingRepository {
	return &RAGSettingRepository{
		db: db,
	}
}

func (repo *RAGSettingRepository) GetAllRAGSettings() ([]*models.RagSetting, error) {
	var contexts []*models.RagSetting

	query := repo.db.Model(&models.RagSetting{})

	err := query.Find(&contexts).Error
	return contexts, err
}
