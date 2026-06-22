package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AwarenessMaturityLevelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessMaturityLevel, error)
	GetByID(id uint) (*models.AwarenessMaturityLevel, error)
	GetByType(maturityLevel string) (*models.AwarenessMaturityLevel, error)
}

type AwarenessMaturityLevelRepository struct {
	db *gorm.DB
}

func NewAwarenessMaturityLevelRepository(db *gorm.DB) *AwarenessMaturityLevelRepository {
	return &AwarenessMaturityLevelRepository{
		db: db,
	}
}

func (repo *AwarenessMaturityLevelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessMaturityLevel, error) {
	var awarenessMaturityLevels []*models.AwarenessMaturityLevel

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Find(&awarenessMaturityLevels)

	return awarenessMaturityLevels, result.Error
}

func (repo *AwarenessMaturityLevelRepository) GetByID(id uint) (*models.AwarenessMaturityLevel, error) {
	var awarenessMaturityLevel *models.AwarenessMaturityLevel

	result := repo.db.First(&awarenessMaturityLevel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessMaturityLevel, nil
}

func (repo *AwarenessMaturityLevelRepository) GetByType(maturityLevel string) (*models.AwarenessMaturityLevel, error) {
	var awarenessMaturityLevel *models.AwarenessMaturityLevel

	result := repo.db.Where("type = ?", maturityLevel).First(&awarenessMaturityLevel)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessMaturityLevel, nil
}