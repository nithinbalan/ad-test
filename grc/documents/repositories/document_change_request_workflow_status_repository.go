package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentChangeRequestWorkflowStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestWorkflowStatus, error)
	GetByID(id uint) (*models.DocumentChangeRequestWorkflowStatus, error)
	GetByType(status string) (*models.DocumentChangeRequestWorkflowStatus, error)
}

// Implementation
type DocumentChangeRequestWorkflowStatusRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestWorkflowStatusRepository(db *gorm.DB) *DocumentChangeRequestWorkflowStatusRepository {
	return &DocumentChangeRequestWorkflowStatusRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestWorkflowStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestWorkflowStatus, error) {
	var statuses []*models.DocumentChangeRequestWorkflowStatus

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
	result := query.Find(&statuses)

	return statuses, result.Error
}

func (repo *DocumentChangeRequestWorkflowStatusRepository) GetByID(id uint) (*models.DocumentChangeRequestWorkflowStatus, error) {
	var status *models.DocumentChangeRequestWorkflowStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *DocumentChangeRequestWorkflowStatusRepository) GetByType(status string) (*models.DocumentChangeRequestWorkflowStatus, error) {
	var workflowStatus *models.DocumentChangeRequestWorkflowStatus

	result := repo.db.Where("type = ?", status).First(&workflowStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflowStatus, nil
}
