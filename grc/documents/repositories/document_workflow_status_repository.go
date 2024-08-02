package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentWorkflowStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentWorkflowStatus, error)
	GetByID(id uint) (*models.DocumentWorkflowStatus, error)
	GetByType(status string) (*models.DocumentWorkflowStatus, error)
}

// Implementation
type DocumentWorkflowStatusRepository struct {
	db *gorm.DB
}

func NewDocumentWorkflowStatusRepository(db *gorm.DB) *DocumentWorkflowStatusRepository {
	return &DocumentWorkflowStatusRepository{
		db: db,
	}
}

func (repo *DocumentWorkflowStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentWorkflowStatus, error) {
	var statuses []*models.DocumentWorkflowStatus

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

func (repo *DocumentWorkflowStatusRepository) GetByID(id uint) (*models.DocumentWorkflowStatus, error) {
	var status *models.DocumentWorkflowStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *DocumentWorkflowStatusRepository) GetByType(status string) (*models.DocumentWorkflowStatus, error) {
	var reviewStatus *models.DocumentWorkflowStatus

	result := repo.db.Where("type = ?", status).First(&reviewStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviewStatus, nil
}
