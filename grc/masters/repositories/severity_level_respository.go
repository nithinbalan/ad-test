package repositories

import (
	"errors"
	"jethur/grc/masters/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SeverityLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SeverityLevel, error)
	GetByID(id uint) (*models.SeverityLevel, error)
}

// implementation
type SeverityLevelRepository struct {
	db *gorm.DB
}

func NewSeverityLevelRepository(db *gorm.DB) SeverityLevelRepositoryInterface {
	return &SeverityLevelRepository{
		db: db,
	}
}

func (repo *SeverityLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SeverityLevel, error) {

	var severityLevels []*models.SeverityLevel

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
	result := query.Find(&severityLevels)

	return severityLevels, result.Error
}

func (repo *SeverityLevelRepository) GetByID(id uint) (*models.SeverityLevel, error) {
	var severityLevel *models.SeverityLevel
	result := repo.db.First(&severityLevel, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return severityLevel, result.Error
}
