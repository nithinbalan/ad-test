package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type IncidentStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentStatus, error)
	GetByID(id uint) (*models.IncidentStatus, error)
}

// implementation
type IncidentStatusRepository struct {
	db *gorm.DB
}

func NewIncidentStatusRepository(db *gorm.DB) *IncidentStatusRepository {
	return &IncidentStatusRepository{
		db: db,
	}
}

func (repo *IncidentStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IncidentStatus, error) {

	var incidentStatuses []*models.IncidentStatus

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
	result := query.Find(&incidentStatuses)

	return incidentStatuses, result.Error
}

func (repo *IncidentStatusRepository) GetByID(id uint) (*models.IncidentStatus, error) {
	var status *models.IncidentStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}
