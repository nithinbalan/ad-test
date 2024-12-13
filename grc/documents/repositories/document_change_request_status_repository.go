package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentChangeRequestStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestStatus, error)
	GetByID(id uint) (*models.DocumentChangeRequestStatus, error)
	GetByType(status string) (*models.DocumentChangeRequestStatus, error)
}

// Implementation
type DocumentChangeRequestStatusRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestStatusRepository(db *gorm.DB) *DocumentChangeRequestStatusRepository {
	return &DocumentChangeRequestStatusRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestStatus, error) {
	var documentChangeRequestStatuses []*models.DocumentChangeRequestStatus

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
	result := query.Find(&documentChangeRequestStatuses)

	return documentChangeRequestStatuses, result.Error
}

func (repo *DocumentChangeRequestStatusRepository) GetByID(id uint) (*models.DocumentChangeRequestStatus, error) {
	var status *models.DocumentChangeRequestStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *DocumentChangeRequestStatusRepository) GetByType(status string) (*models.DocumentChangeRequestStatus, error) {
	var documentChangeRequestStatus *models.DocumentChangeRequestStatus

	result := repo.db.Where("type = ?", status).First(&documentChangeRequestStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentChangeRequestStatus, nil
}
