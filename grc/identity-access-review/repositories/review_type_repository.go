package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ReviewTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewType, error)
	GetByID(id uint) (*models.ReviewType, error)
	GetByType(reviewType string) (*models.ReviewType, error)
}

type ReviewTypeRepository struct {
	db *gorm.DB
}

func NewReviewTypeRepository(db *gorm.DB) *ReviewTypeRepository {
	return &ReviewTypeRepository{
		db: db,
	}
}

func (repo *ReviewTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewType, error) {
	var reviewTypes []*models.ReviewType

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if sortBy != "id" && sortBy != "title_en" && sortBy != "title_ar" && sortBy != "type" && sortBy != "status" {
		sortBy = "title_en"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)

	result := query.Find(&reviewTypes)

	return reviewTypes, result.Error
}

func (repo *ReviewTypeRepository) GetByID(id uint) (*models.ReviewType, error) {
	var reviewType *models.ReviewType

	result := repo.db.First(&reviewType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviewType, nil
}

func (repo *ReviewTypeRepository) GetByType(reviewType string) (*models.ReviewType, error) {
	var resultReviewType *models.ReviewType

	result := repo.db.Where("type = ?", reviewType).First(&resultReviewType)

	if result.Error != nil {
		return nil, result.Error
	}

	return resultReviewType, nil
}
