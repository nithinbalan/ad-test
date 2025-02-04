package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ApplicabilityStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ApplicabilityStatus, error)
	GetByID(id uint) (*models.ApplicabilityStatus, error)
	GetByType(status string) (*models.ApplicabilityStatus, error)
}

// Implementation
type ApplicabilityStatusRepository struct {
	db *gorm.DB
}

func NewApplicabilityStatusRepository(db *gorm.DB) *ApplicabilityStatusRepository {
	return &ApplicabilityStatusRepository{
		db: db,
	}
}

func (repo *ApplicabilityStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ApplicabilityStatus, error) {
	var applicabilityStatuses []*models.ApplicabilityStatus

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
	result := query.Find(&applicabilityStatuses)

	return applicabilityStatuses, result.Error
}

func (repo *ApplicabilityStatusRepository) GetByID(id uint) (*models.ApplicabilityStatus, error) {
	var applicabilityStatus *models.ApplicabilityStatus
	result := repo.db.First(&applicabilityStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return applicabilityStatus, result.Error
}

func (repo *ApplicabilityStatusRepository) GetByType(status string) (*models.ApplicabilityStatus, error) {
	var applicabilityStatus *models.ApplicabilityStatus

	result := repo.db.Where("type = ?", status).First(&applicabilityStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Applicability Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return applicabilityStatus, nil
}
