package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ThreatStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatStatus, error)
	GetByID(id uint) (*models.ThreatStatus, error)
	GetByType(status string) (*models.ThreatStatus, error)
}

// Implementation
type ThreatStatusRepository struct {
	db *gorm.DB
}

func NewThreatStatusRepository(db *gorm.DB) *ThreatStatusRepository {
	return &ThreatStatusRepository{
		db: db,
	}
}

func (repo *ThreatStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatStatus, error) {
	var threatStatuses []*models.ThreatStatus

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
	result := query.Find(&threatStatuses)

	return threatStatuses, result.Error
}

func (repo *ThreatStatusRepository) GetByID(id uint) (*models.ThreatStatus, error) {
	var status *models.ThreatStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ThreatStatusRepository) GetByType(status string) (*models.ThreatStatus, error) {
	var threatStatus *models.ThreatStatus

	result := repo.db.Where("type = ?", status).First(&threatStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return threatStatus, nil
}
