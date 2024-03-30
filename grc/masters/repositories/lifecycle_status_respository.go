package repositories

import (
	"errors"
	"jethur/grc/masters/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type LifecycleStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LifecycleStatus, error)
	GetByID(id uint) (*models.LifecycleStatus, error)
	GetByStatus(status string) (*models.LifecycleStatus, error)
}

// implementation
type LifecycleStatusRepository struct {
	db *gorm.DB
}

func NewLifecycleStatusRepository(db *gorm.DB) *LifecycleStatusRepository {
	return &LifecycleStatusRepository{
		db: db,
	}
}

func (repo *LifecycleStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LifecycleStatus, error) {

	var lifecycleStatuses []*models.LifecycleStatus

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
	result := query.Find(&lifecycleStatuses)

	return lifecycleStatuses, result.Error
}

func (repo *LifecycleStatusRepository) GetByID(id uint) (*models.LifecycleStatus, error) {
	var lifecycleStatus *models.LifecycleStatus
	result := repo.db.First(&lifecycleStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return lifecycleStatus, result.Error
}

func (repo *LifecycleStatusRepository) GetByStatus(status string) (*models.LifecycleStatus, error) {
	var lifecycleStatus *models.LifecycleStatus

	result := repo.db.Where("type = ?", status).First(&lifecycleStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Lifecycle Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return lifecycleStatus, nil
}
