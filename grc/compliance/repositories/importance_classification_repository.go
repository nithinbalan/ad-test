package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ImportanceClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ImportanceClassification, error)
	GetByID(id uint) (*models.ImportanceClassification, error)
	GetByType(status string) (*models.ImportanceClassification, error)
}

// Implementation
type ImportanceClassificationRepository struct {
	db *gorm.DB
}

func NewImportanceClassificationRepository(db *gorm.DB) *ImportanceClassificationRepository {
	return &ImportanceClassificationRepository{
		db: db,
	}
}

func (repo *ImportanceClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ImportanceClassification, error) {
	var classifications []*models.ImportanceClassification
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ? OR description_en LIKE ? OR description_ar LIKE ?)",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
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

	result := query.Order(sortBy + " " + sort).Find(&classifications)
	return classifications, result.Error
}

func (repo *ImportanceClassificationRepository) GetByID(id uint) (*models.ImportanceClassification, error) {
	var classification *models.ImportanceClassification
	result := repo.db.First(&classification, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return classification, nil
}

func (repo *ImportanceClassificationRepository) GetByType(status string) (*models.ImportanceClassification, error) {
	var classification *models.ImportanceClassification
	result := repo.db.Where("status = ?", status).First(&classification)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Importance Classification with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return classification, nil
}
