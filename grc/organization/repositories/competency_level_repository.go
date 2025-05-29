package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type CompetencyLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyLevel, error)
	GetByID(id uint) (*models.CompetencyLevel, error)
}

// implementation
type CompetencyLevelRepository struct {
	db *gorm.DB
}

func NewCompetencyLevelRepository(db *gorm.DB) *CompetencyLevelRepository {
	return &CompetencyLevelRepository{
		db: db,
	}
}

func (repo *CompetencyLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyLevel, error) {
	var competencyLevels []*models.CompetencyLevel

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
	result := query.Find(&competencyLevels)

	return competencyLevels, result.Error
}

func (repo *CompetencyLevelRepository) GetByID(id uint) (*models.CompetencyLevel, error) {
	var status *models.CompetencyLevel
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}
