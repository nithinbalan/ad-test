package repositories

import (
	"jethur/grc/models" // Ensure this import path is correct for your project structure

	"gorm.io/gorm"
)

// Interfaces
type WorkflowRoleTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowRoleType, error)
	GetByID(id uint) (*models.WorkflowRoleType, error)
	GetByType(roleType string) (*models.WorkflowRoleType, error)
}

// Implementation
type WorkflowRoleTypeRepository struct {
	db *gorm.DB
}

func NewWorkflowRoleTypeRepository(db *gorm.DB) *WorkflowRoleTypeRepository {
	return &WorkflowRoleTypeRepository{
		db: db,
	}
}

func (repo *WorkflowRoleTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowRoleType, error) {
	var workflowRoleTypes []*models.WorkflowRoleType

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
	result := query.Find(&workflowRoleTypes)

	return workflowRoleTypes, result.Error
}

func (repo *WorkflowRoleTypeRepository) GetByID(id uint) (*models.WorkflowRoleType, error) {
	var roleType *models.WorkflowRoleType
	result := repo.db.First(&roleType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return roleType, nil
}

func (repo *WorkflowRoleTypeRepository) GetByType(roleType string) (*models.WorkflowRoleType, error) {
	var workflowRoleType *models.WorkflowRoleType

	result := repo.db.Where("type = ?", roleType).First(&roleType)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflowRoleType, nil
}
