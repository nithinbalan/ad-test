package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ControlEfficiencyStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlEfficiencyStatus, error)
	GetByID(id uint) (*models.ControlEfficiencyStatus, error)
	GetByType(status string) (*models.ControlEfficiencyStatus, error)
}

// implementation
type ControlEfficiencyStatusRepository struct {
	db *gorm.DB
}

func NewControlEfficiencyStatusRepository(db *gorm.DB) *ControlEfficiencyStatusRepository {
	return &ControlEfficiencyStatusRepository{
		db: db,
	}
}

func (repo *ControlEfficiencyStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlEfficiencyStatus, error) {
	var controlEfficiencyStatuses []*models.ControlEfficiencyStatus

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
	result := query.Find(&controlEfficiencyStatuses)

	return controlEfficiencyStatuses, result.Error
}

func (repo *ControlEfficiencyStatusRepository) GetByID(id uint) (*models.ControlEfficiencyStatus, error) {
	var status *models.ControlEfficiencyStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, result.Error
}

func (repo *ControlEfficiencyStatusRepository) GetByType(status string) (*models.ControlEfficiencyStatus, error) {
	var controlEfficiencyStatus *models.ControlEfficiencyStatus

	result := repo.db.Where("type = ?", status).First(&controlEfficiencyStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlEfficiencyStatus, nil
}
