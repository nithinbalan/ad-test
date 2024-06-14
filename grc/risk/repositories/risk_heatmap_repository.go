package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type RiskHeatmapRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskHeatmap, error)
	GetByID(id uint) (*models.RiskHeatmap, error)
	GetByType(status string) (*models.RiskHeatmap, error)
}

// Implementation
type RiskHeatmapRepository struct {
	db *gorm.DB
}

func NewRiskHeatmapRepository(db *gorm.DB) *RiskHeatmapRepository {
	return &RiskHeatmapRepository{
		db: db,
	}
}

func (repo *RiskHeatmapRepository) GetProcessRiskHeatmap(sort string, sortBy string, keyword string, status string) ([]*models.RiskHeatmap, error) {
	var riskHeatmaps []*models.RiskHeatmap

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
	result := query.Find(&riskHeatmaps)

	return riskHeatmaps, result.Error
}
