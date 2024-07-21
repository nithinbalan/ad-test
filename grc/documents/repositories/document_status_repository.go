package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentStatus, error)
	GetByID(id uint) (*models.DocumentStatus, error)
	GetByType(status string) (*models.DocumentStatus, error)
}

// implementation
type DocumentStatusRepository struct {
	db *gorm.DB
}

func NewDocumentStatusRepository(db *gorm.DB) *DocumentStatusRepository {
	return &DocumentStatusRepository{
		db: db,
	}
}

func (repo *DocumentStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentStatus, error) {
	var documentStatuses []*models.DocumentStatus

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
	result := query.Find(&documentStatuses)

	return documentStatuses, result.Error
}

func (repo *DocumentStatusRepository) GetByID(id uint) (*models.DocumentStatus, error) {
	var status *models.DocumentStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *DocumentStatusRepository) GetByType(status string) (*models.DocumentStatus, error) {
	var documentStatus *models.DocumentStatus

	result := repo.db.Where("type = ?", status).First(&documentStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentStatus, nil
}
