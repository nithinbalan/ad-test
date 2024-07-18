package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentAccessTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentAccessType, error)
	GetByID(id uint) (*models.DocumentAccessType, error)
	GetByType(accessType string) (*models.DocumentAccessType, error)
}

// Implementation
type DocumentAccessTypeRepository struct {
	db *gorm.DB
}

func NewDocumentAccessTypeRepository(db *gorm.DB) *DocumentAccessTypeRepository {
	return &DocumentAccessTypeRepository{
		db: db,
	}
}

func (repo *DocumentAccessTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentAccessType, error) {
	var documentAccessTypes []*models.DocumentAccessType

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
	result := query.Find(&documentAccessTypes)

	return documentAccessTypes, result.Error
}

func (repo *DocumentAccessTypeRepository) GetByID(id uint) (*models.DocumentAccessType, error) {
	var accessType *models.DocumentAccessType
	result := repo.db.First(&accessType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return accessType, result.Error
}

func (repo *DocumentAccessTypeRepository) GetByType(accessType string) (*models.DocumentAccessType, error) {
	var documentAccessType *models.DocumentAccessType

	result := repo.db.Where("type = ?", accessType).First(&documentAccessType)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentAccessType, nil
}
