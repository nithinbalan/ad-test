package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentChangeRequestInitiateStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateStatus, error)
	GetByID(id uint) (*models.DocumentChangeRequestInitiateStatus, error)
	GetByType(status string) (*models.DocumentChangeRequestInitiateStatus, error)
}

// Implementation
type DocumentChangeRequestInitiateStatusRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestInitiateStatusRepository(db *gorm.DB) *DocumentChangeRequestInitiateStatusRepository {
	return &DocumentChangeRequestInitiateStatusRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestInitiateStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateStatus, error) {
	var documentChangeRequestInitiateStatuses []*models.DocumentChangeRequestInitiateStatus

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
	result := query.Find(&documentChangeRequestInitiateStatuses)

	return documentChangeRequestInitiateStatuses, result.Error
}

func (repo *DocumentChangeRequestInitiateStatusRepository) GetByID(id uint) (*models.DocumentChangeRequestInitiateStatus, error) {
	var status *models.DocumentChangeRequestInitiateStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *DocumentChangeRequestInitiateStatusRepository) GetByType(status string) (*models.DocumentChangeRequestInitiateStatus, error) {
	var documentChangeRequestInitiateStatus *models.DocumentChangeRequestInitiateStatus

	result := repo.db.Where("type = ?", status).First(&documentChangeRequestInitiateStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentChangeRequestInitiateStatus, nil
}
