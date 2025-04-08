package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ProcessCriticalityLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessCriticalityLevel, error)
	GetByID(id uint) (*models.ProcessCriticalityLevel, error)
	GetByType(status string) (*models.ProcessCriticalityLevel, error)
}

// implementation
type ProcessCriticalityLevelRepository struct {
	db *gorm.DB
}

func NewProcessCriticalityLevelRepository(db *gorm.DB) *ProcessCriticalityLevelRepository {
	return &ProcessCriticalityLevelRepository{
		db: db,
	}
}

func (repo *ProcessCriticalityLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessCriticalityLevel, error) {
	var processCriticalityLevels []*models.ProcessCriticalityLevel

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
	result := query.Find(&processCriticalityLevels)

	return processCriticalityLevels, result.Error
}

func (repo *ProcessCriticalityLevelRepository) GetByID(id uint) (*models.ProcessCriticalityLevel, error) {
	var status *models.ProcessCriticalityLevel
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ProcessCriticalityLevelRepository) GetByType(status string) (*models.ProcessCriticalityLevel, error) {
	var processCriticalityLevel *models.ProcessCriticalityLevel

	result := repo.db.Where("type = ?", status).First(&processCriticalityLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return processCriticalityLevel, nil
}
