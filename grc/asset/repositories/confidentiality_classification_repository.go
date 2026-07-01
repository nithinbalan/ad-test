package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ConfidentialityClassificationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConfidentialityClassification, error)
	GetByID(id uint) (*models.ConfidentialityClassification, error)
	GetByType(status string) (*models.ConfidentialityClassification, error)
}

type ConfidentialityClassificationRepository struct {
	db *gorm.DB
}

func NewConfidentialityClassificationRepository(db *gorm.DB) *ConfidentialityClassificationRepository {
	return &ConfidentialityClassificationRepository{
		db: db,
	}
}

func (repo *ConfidentialityClassificationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ConfidentialityClassification, error) {
	var confidentialityClassifications []*models.ConfidentialityClassification

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Find(&confidentialityClassifications)

	return confidentialityClassifications, result.Error
}

func (repo *ConfidentialityClassificationRepository) GetByID(id uint) (*models.ConfidentialityClassification, error) {
	var confidentialityClassification *models.ConfidentialityClassification

	result := repo.db.First(&confidentialityClassification, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return confidentialityClassification, nil
}

func (repo *ConfidentialityClassificationRepository) GetByType(status string) (*models.ConfidentialityClassification, error) {
	var confidentialityClassification *models.ConfidentialityClassification

	result := repo.db.Where("type = ?", status).First(&confidentialityClassification)

	if result.Error != nil {
		return nil, result.Error
	}

	return confidentialityClassification, nil
}