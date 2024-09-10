package repositories

import (
	"jethur/grc/models" // Ensure this import path is correct for your project structure

	"gorm.io/gorm"
)

// Interfaces
type WorkflowLevelTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowLevelType, error)
	GetByID(id uint) (*models.WorkflowLevelType, error)
	GetByType(levelType string) (*models.WorkflowLevelType, error)
}

// Implementation
type WorkflowLevelTypeRepository struct {
	db *gorm.DB
}

func NewWorkflowLevelTypeRepository(db *gorm.DB) *WorkflowLevelTypeRepository {
	return &WorkflowLevelTypeRepository{
		db: db,
	}
}

func (repo *WorkflowLevelTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowLevelType, error) {
	var workflowLevelTypes []*models.WorkflowLevelType

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&workflowLevelTypes)

	return workflowLevelTypes, result.Error
}

func (repo *WorkflowLevelTypeRepository) GetByID(id uint) (*models.WorkflowLevelType, error) {
	var levelType *models.WorkflowLevelType
	result := repo.db.First(&levelType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return levelType, nil
}

func (repo *WorkflowLevelTypeRepository) GetByType(levelType string) (*models.WorkflowLevelType, error) {
	var workflowLevelType *models.WorkflowLevelType

	result := repo.db.Where("type = ?", levelType).First(&workflowLevelType)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflowLevelType, nil
}
