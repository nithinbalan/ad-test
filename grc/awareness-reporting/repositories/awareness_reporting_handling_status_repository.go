package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AwarenessHandlingStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessHandlingStatus, error)
	GetByID(id uint) (*models.AwarenessHandlingStatus, error)
	GetByType(statusType string) (*models.AwarenessHandlingStatus, error)
}

type AwarenessHandlingStatusRepository struct {
	db *gorm.DB
}

func NewAwarenessHandlingStatusRepository(db *gorm.DB) *AwarenessHandlingStatusRepository {
	return &AwarenessHandlingStatusRepository{
		db: db,
	}
}

func (repo *AwarenessHandlingStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessHandlingStatus, error) {
	var awarenessHandlingStatuses []*models.AwarenessHandlingStatus

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

	result := query.Find(&awarenessHandlingStatuses)

	return awarenessHandlingStatuses, result.Error
}

func (repo *AwarenessHandlingStatusRepository) GetByID(id uint) (*models.AwarenessHandlingStatus, error) {
	var awarenessHandlingStatus models.AwarenessHandlingStatus

	result := repo.db.First(&awarenessHandlingStatus, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &awarenessHandlingStatus, nil
}

func (repo *AwarenessHandlingStatusRepository) GetByType(statusType string) (*models.AwarenessHandlingStatus, error) {
	var awarenessHandlingStatus models.AwarenessHandlingStatus

	result := repo.db.Where("type = ?", statusType).First(&awarenessHandlingStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return &awarenessHandlingStatus, nil
}