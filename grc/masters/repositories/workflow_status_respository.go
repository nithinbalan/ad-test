package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type WorkflowStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowStatus, error)
	GetByID(id uint) (*models.WorkflowStatus, error)
	GetByType(status string) (*models.WorkflowStatus, error)
}

// Implementation
type WorkflowStatusRepository struct {
	db *gorm.DB
}

func NewWorkflowStatusRepository(db *gorm.DB) *WorkflowStatusRepository {
	return &WorkflowStatusRepository{
		db: db,
	}
}

func (repo *WorkflowStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.WorkflowStatus, error) {
	var workflowStatuses []*models.WorkflowStatus

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
	result := query.Find(&workflowStatuses)

	return workflowStatuses, result.Error
}

func (repo *WorkflowStatusRepository) GetByID(id uint) (*models.WorkflowStatus, error) {
	var workflowStatus *models.WorkflowStatus
	result := repo.db.First(&workflowStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return workflowStatus, result.Error
}

func (repo *WorkflowStatusRepository) GetByType(status string) (*models.WorkflowStatus, error) {
	var workflowStatus *models.WorkflowStatus

	result := repo.db.Where("type = ?", status).First(&workflowStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Workflow Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return workflowStatus, nil
}
