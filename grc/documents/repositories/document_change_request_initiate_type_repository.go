package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type DocumentChangeRequestInitiateTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateType, error)
	GetByID(id uint) (*models.DocumentChangeRequestInitiateType, error)
	GetByType(requestType string) (*models.DocumentChangeRequestInitiateType, error)
}

// Implementation
type DocumentChangeRequestInitiateTypeRepository struct {
	db *gorm.DB
}

func NewDocumentChangeRequestInitiateTypeRepository(db *gorm.DB) *DocumentChangeRequestInitiateTypeRepository {
	return &DocumentChangeRequestInitiateTypeRepository{
		db: db,
	}
}

func (repo *DocumentChangeRequestInitiateTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.DocumentChangeRequestInitiateType, error) {
	var documentChangeRequestInitiateTypes []*models.DocumentChangeRequestInitiateType

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

	// Define the custom order
	customOrder := `
		CASE 
			WHEN type = 'modify' THEN 1
			WHEN type = 'archive' THEN 2
			WHEN type = 'restore' THEN 3
			ELSE 4
		END
	`

	query = query.Order(customOrder)

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&documentChangeRequestInitiateTypes)

	return documentChangeRequestInitiateTypes, result.Error
}

func (repo *DocumentChangeRequestInitiateTypeRepository) GetByID(id uint) (*models.DocumentChangeRequestInitiateType, error) {
	var requestType *models.DocumentChangeRequestInitiateType
	result := repo.db.First(&requestType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return requestType, result.Error
}

func (repo *DocumentChangeRequestInitiateTypeRepository) GetByType(requestType string) (*models.DocumentChangeRequestInitiateType, error) {
	var documentChangeRequestInitiateType *models.DocumentChangeRequestInitiateType

	result := repo.db.Where("type = ?", requestType).First(&documentChangeRequestInitiateType)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentChangeRequestInitiateType, nil
}
