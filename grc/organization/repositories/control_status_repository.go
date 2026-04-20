package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

// Interfaces
type ControlStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlStatus, error)
	GetByID(id uint) (*models.ControlStatus, error)
	GetByType(typeFilter string) (*models.ControlStatus, error)
}

// implementation
type ControlStatusRepository struct {
	db *gorm.DB
}

func NewControlStatusRepository(db *gorm.DB) *ControlStatusRepository {
	return &ControlStatusRepository{
		db: db,
	}
}

func (repo *ControlStatusRepository) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlStatus, error) {
	var controlStatuses []*models.ControlStatus

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

	if typeFilter == "needs_improvement" || typeFilter == "doesnt_need_improvement" {
		query = query.Where("type = ?", typeFilter)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&controlStatuses)

	return controlStatuses, result.Error
}

func (repo *ControlStatusRepository) GetByID(id uint) (*models.ControlStatus, error) {
	var controlStatus *models.ControlStatus

	result := repo.db.First(&controlStatus, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlStatus, result.Error
}

func (repo *ControlStatusRepository) GetByType(typeFilter string) (*models.ControlStatus, error) {
	var controlStatus *models.ControlStatus

	result := repo.db.Where("type = ?", typeFilter).First(&controlStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlStatus, nil
}
