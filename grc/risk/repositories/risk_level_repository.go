package repositories

import (
	"jethur/grc/risk/models"

	"gorm.io/gorm"
)

// Interfaces
type RiskLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskLevel, error)
	GetByID(id uint) (*models.RiskLevel, error)
	GetByType(status string) (*models.RiskLevel, error)
}

// Implementation
type RiskLevelRepository struct {
	db *gorm.DB
}

func NewRiskLevelRepository(db *gorm.DB) *RiskLevelRepository {
	return &RiskLevelRepository{
		db: db,
	}
}

func (repo *RiskLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskLevel, error) {
	var riskLevels []*models.RiskLevel

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
	result := query.Find(&riskLevels)

	return riskLevels, result.Error
}

func (repo *RiskLevelRepository) GetByID(id uint) (*models.RiskLevel, error) {
	var riskLevel *models.RiskLevel
	result := repo.db.First(&riskLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskLevel, nil

}

func (repo *RiskLevelRepository) GetByType(status string) (*models.RiskLevel, error) {
	var riskLevel *models.RiskLevel

	result := repo.db.Where("type = ?", status).First(&riskLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskLevel, nil
}
