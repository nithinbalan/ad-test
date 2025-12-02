package repositories

import (
	"errors"
	"jethur/grc/masters/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type WorkflowContextRepositoryInterface interface {
	GetAll(params dtos.WorkflowContextListQueryParamsDTO) ([]*models.WorkflowContext, error)
	GetByID(id uint) (*models.WorkflowContext, error)
	GetByType(status string) (*models.WorkflowContext, error)
}

// Implementation
type WorkflowContextRepository struct {
	db *gorm.DB
}

func NewWorkflowContextRepository(db *gorm.DB) *WorkflowContextRepository {
	return &WorkflowContextRepository{
		db: db,
	}
}

func (repo *WorkflowContextRepository) GetAll(params dtos.WorkflowContextListQueryParamsDTO) ([]*models.WorkflowContext, error) {
	var workflowContextes []*models.WorkflowContext

	query := repo.db

	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&workflowContextes)

	return workflowContextes, result.Error
}

func (repo *WorkflowContextRepository) GetByID(id uint) (*models.WorkflowContext, error) {
	var workflowContext *models.WorkflowContext
	result := repo.db.First(&workflowContext, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return workflowContext, result.Error
}

func (repo *WorkflowContextRepository) GetByType(status string) (*models.WorkflowContext, error) {
	var workflowContext *models.WorkflowContext

	result := repo.db.Where("type = ?", status).First(&workflowContext)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Workflow Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return workflowContext, nil
}
