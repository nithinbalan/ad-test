package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ConformanceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Conformance, error)
	GetByID(id uint) (*models.Conformance, error)
	GetByType(status string) (*models.Conformance, error)
}

// Implementation
type ConformanceRepository struct {
	db *gorm.DB
}

func NewConformanceRepository(db *gorm.DB) *ConformanceRepository {
	return &ConformanceRepository{
		db: db,
	}
}

func (repo *ConformanceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Conformance, error) {
	var conformances []*models.Conformance

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
	result := query.Find(&conformances)

	return conformances, result.Error
}

func (repo *ConformanceRepository) GetByID(id uint) (*models.Conformance, error) {
	var conformance *models.Conformance
	result := repo.db.First(&conformance, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return conformance, result.Error
}

func (repo *ConformanceRepository) GetByType(strConformance string) (*models.Conformance, error) {
	var conformance *models.Conformance

	result := repo.db.Where("type = ?", strConformance).First(&conformance)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Conformance type with type "+strConformance)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return conformance, nil
}
