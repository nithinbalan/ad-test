package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ControlTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlType, error)
	GetByID(id uint) (*models.ControlType, error)
}

// implementation
type ControlTypeRepository struct {
	db *gorm.DB
}

func NewControlTypeRepository(db *gorm.DB) *ControlTypeRepository {
	return &ControlTypeRepository{
		db: db,
	}
}

func (repo *ControlTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlType, error) {
	var controlTypes []*models.ControlType

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
	result := query.Find(&controlTypes)

	return controlTypes, result.Error
}

func (repo *ControlTypeRepository) GetByID(id uint) (*models.ControlType, error) {
	var controlType *models.ControlType
	result := repo.db.First(&controlType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlType, result.Error
}
