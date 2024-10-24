package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentChangeRequestTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestType, error)
	GetByID(id uint) (*models.DocumentChangeRequestType, error)
	GetByType(requestType string) (*models.DocumentChangeRequestType, error)
}

// Implementation
type DocumentChangeRequestTypeRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestTypeRepository(db *gorm.DB) *DocumentChangeRequestTypeRepository {
	return &DocumentChangeRequestTypeRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestType, error) {
	var documentChangeRequestTypes []*models.DocumentChangeRequestType

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
	result := query.Find(&documentChangeRequestTypes)

	return documentChangeRequestTypes, result.Error
}

func (repo *DocumentChangeRequestTypeRepository) GetByID(id uint) (*models.DocumentChangeRequestType, error) {
	var requestType *models.DocumentChangeRequestType
	result := repo.db.First(&requestType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return requestType, result.Error
}

func (repo *DocumentChangeRequestTypeRepository) GetByType(requestType string) (*models.DocumentChangeRequestType, error) {
	var documentChangeRequestType *models.DocumentChangeRequestType

	result := repo.db.Where("type = ?", requestType).First(&documentChangeRequestType)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentChangeRequestType, nil
}
