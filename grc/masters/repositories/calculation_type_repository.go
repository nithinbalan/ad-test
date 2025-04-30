package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CalculationTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CalculationType, error)
	GetByID(id uint) (*models.CalculationType, error)
}

// implementation
type CalculationTypeRepository struct {
	db *gorm.DB
}

func NewCalculationTypeRepository(db *gorm.DB) *CalculationTypeRepository {
	return &CalculationTypeRepository{
		db: db,
	}
}

func (repo *CalculationTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CalculationType, error) {
	var calculationTypes []*models.CalculationType
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
	result := query.Find(&calculationTypes)

	return calculationTypes, result.Error
}

func (repo *CalculationTypeRepository) GetByID(id uint) (*models.CalculationType, error) {
	var category *models.CalculationType
	result := repo.db.First(&category, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return category, nil
}
