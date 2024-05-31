package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ModuleRepositoryInterface interface {
	GetModuleIDByIdentifier(identifier string) (uint, error)
}

// Implementation
type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{
		db: db,
	}
}

func (repo *ModuleRepository) GetModuleIDByIdentifier(identifier string) (uint, error) {
	var module models.Module
	if err := repo.db.Where("identifier = ?", identifier).First(&module).Error; err != nil {
		return 0, err
	}
	return module.ID, nil
}
