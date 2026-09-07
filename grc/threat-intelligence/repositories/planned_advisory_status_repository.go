package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type PlannedAdvisoryStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryStatus, error)
	GetByID(id uint) (*models.PlannedAdvisoryStatus, error)
	GetByType(statusType string) (*models.PlannedAdvisoryStatus, error)
}

// implementation
type PlannedAdvisoryStatusRepository struct {
	db *gorm.DB
}

func NewPlannedAdvisoryStatusRepository(db *gorm.DB) *PlannedAdvisoryStatusRepository {
	return &PlannedAdvisoryStatusRepository{
		db: db,
	}
}

func (repo *PlannedAdvisoryStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryStatus, error) {

	var plannedAdvisoryStatuses []*models.PlannedAdvisoryStatus

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
	result := query.Find(&plannedAdvisoryStatuses)

	return plannedAdvisoryStatuses, result.Error
}

func (repo *PlannedAdvisoryStatusRepository) GetByID(id uint) (*models.PlannedAdvisoryStatus, error) {
	var status *models.PlannedAdvisoryStatus
	result := repo.db.First(&status, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return status, nil
}

func (repo *PlannedAdvisoryStatusRepository) GetByType(statusType string) (*models.PlannedAdvisoryStatus, error) {
	var plannedAdvisoryStatus *models.PlannedAdvisoryStatus

	result := repo.db.Where("type = ?", statusType).First(&plannedAdvisoryStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return plannedAdvisoryStatus, nil
}
