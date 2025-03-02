package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ReviewStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewStatus, error)
	GetByID(id uint) (*models.ReviewStatus, error)
	GetByType(status string) (*models.ReviewStatus, error)
}

// Implementation
type ReviewStatusRepository struct {
	db *gorm.DB
}

func NewReviewStatusRepository(db *gorm.DB) *ReviewStatusRepository {
	return &ReviewStatusRepository{
		db: db,
	}
}

func (repo *ReviewStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewStatus, error) {
	var reviewStatuses []*models.ReviewStatus

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
	result := query.Find(&reviewStatuses)

	return reviewStatuses, result.Error
}

func (repo *ReviewStatusRepository) GetByID(id uint) (*models.ReviewStatus, error) {
	var reviewStatus *models.ReviewStatus
	result := repo.db.First(&reviewStatus, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return reviewStatus, result.Error
}

func (repo *ReviewStatusRepository) GetByType(status string) (*models.ReviewStatus, error) {
	var reviewStatus *models.ReviewStatus

	result := repo.db.Where("type = ?", status).First(&reviewStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Review Status with type "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return reviewStatus, nil
}
