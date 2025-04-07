package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ProcessDependencyTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessDependencyType, error)
	GetByID(id uint) (*models.ProcessDependencyType, error)
	GetByType(status string) (*models.ProcessDependencyType, error)
}

// implementation
type ProcessDependencyTypeRepository struct {
	db *gorm.DB
}

func NewProcessDependencyTypeRepository(db *gorm.DB) *ProcessDependencyTypeRepository {
	return &ProcessDependencyTypeRepository{
		db: db,
	}
}

func (repo *ProcessDependencyTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessDependencyType, error) {
	var processDependencyTypes []*models.ProcessDependencyType

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
	result := query.Find(&processDependencyTypes)

	return processDependencyTypes, result.Error
}

func (repo *ProcessDependencyTypeRepository) GetByID(id uint) (*models.ProcessDependencyType, error) {
	var status *models.ProcessDependencyType
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ProcessDependencyTypeRepository) GetByType(status string) (*models.ProcessDependencyType, error) {
	var processDependencyType *models.ProcessDependencyType

	result := repo.db.Where("type = ?", status).First(&processDependencyType)

	if result.Error != nil {
		return nil, result.Error
	}

	return processDependencyType, nil
}
