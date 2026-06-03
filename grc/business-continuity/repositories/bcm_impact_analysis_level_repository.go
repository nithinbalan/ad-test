package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type BCMImpactAnalysisLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMImpactAnalysisLevel, error)
	GetByID(id uint) (*models.BCMImpactAnalysisLevel, error)
	GetByType(status string) (*models.BCMImpactAnalysisLevel, error)
}

// implementation
type BCMImpactAnalysisLevelRepository struct {
	db *gorm.DB
}

func NewBCMImpactAnalysisLevelRepository(db *gorm.DB) *BCMImpactAnalysisLevelRepository {
	return &BCMImpactAnalysisLevelRepository{
		db: db,
	}
}

func (repo *BCMImpactAnalysisLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMImpactAnalysisLevel, error) {
	var bcmImpactAnalysisLevels []*models.BCMImpactAnalysisLevel

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
	result := query.Find(&bcmImpactAnalysisLevels)

	return bcmImpactAnalysisLevels, result.Error
}

func (repo *BCMImpactAnalysisLevelRepository) GetByID(id uint) (*models.BCMImpactAnalysisLevel, error) {
	var status *models.BCMImpactAnalysisLevel
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *BCMImpactAnalysisLevelRepository) GetByType(status string) (*models.BCMImpactAnalysisLevel, error) {
	var bcmImpactAnalysisLevel *models.BCMImpactAnalysisLevel

	result := repo.db.Where("type = ?", status).First(&bcmImpactAnalysisLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return bcmImpactAnalysisLevel, nil
}
