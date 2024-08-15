package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ControlReviewStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlReviewStatus, error)
	GetByID(id uint) (*models.ControlReviewStatus, error)
	GetByType(status string) (*models.ControlReviewStatus, error)
}

// Implementation
type ControlReviewStatusRepository struct {
	db *gorm.DB
}

func NewControlReviewStatusRepository(db *gorm.DB) *ControlReviewStatusRepository {
	return &ControlReviewStatusRepository{
		db: db,
	}
}

func (repo *ControlReviewStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlReviewStatus, error) {

	var controlReviewStatuses []*models.ControlReviewStatus

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
	result := query.Find(&controlReviewStatuses)

	return controlReviewStatuses, result.Error
}

func (repo *ControlReviewStatusRepository) GetByID(id uint) (*models.ControlReviewStatus, error) {
	var controlReviewStatus *models.ControlReviewStatus
	result := repo.db.First(&controlReviewStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return controlReviewStatus, result.Error
}

func (repo *ControlReviewStatusRepository) GetByType(status string) (*models.ControlReviewStatus, error) {
	var controlReviewStatus *models.ControlReviewStatus

	result := repo.db.Where("type = ?", status).First(&controlReviewStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Control Review Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return controlReviewStatus, nil
}
