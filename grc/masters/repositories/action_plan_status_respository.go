package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ActionPlanStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ActionPlanStatus, error)
	GetByID(id uint) (*models.ActionPlanStatus, error)
	GetByType(status string) (*models.ActionPlanStatus, error)
}

// Implementation
type ActionPlanStatusRepository struct {
	db *gorm.DB
}

func NewActionPlanStatusRepository(db *gorm.DB) *ActionPlanStatusRepository {
	return &ActionPlanStatusRepository{
		db: db,
	}
}

func (repo *ActionPlanStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ActionPlanStatus, error) {
	var actionPlanStatuses []*models.ActionPlanStatus

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
	result := query.Find(&actionPlanStatuses)

	return actionPlanStatuses, result.Error
}

func (repo *ActionPlanStatusRepository) GetByID(id uint) (*models.ActionPlanStatus, error) {
	var actionPlanStatus *models.ActionPlanStatus
	result := repo.db.First(&actionPlanStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return actionPlanStatus, result.Error
}

func (repo *ActionPlanStatusRepository) GetByType(status string) (*models.ActionPlanStatus, error) {
	var actionPlanStatus *models.ActionPlanStatus

	result := repo.db.Where("type = ?", status).First(&actionPlanStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Action Plan Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return actionPlanStatus, nil
}
