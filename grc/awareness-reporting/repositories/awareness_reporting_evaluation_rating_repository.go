package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AwarenessEvaluationRatingRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessEvaluationRating, error)
	GetByID(id uint) (*models.AwarenessEvaluationRating, error)
	GetByType(rating string) (*models.AwarenessEvaluationRating, error)
}

type AwarenessEvaluationRatingRepository struct {
	db *gorm.DB
}

func NewAwarenessEvaluationRatingRepository(db *gorm.DB) *AwarenessEvaluationRatingRepository {
	return &AwarenessEvaluationRatingRepository{
		db: db,
	}
}

func (repo *AwarenessEvaluationRatingRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessEvaluationRating, error) {
	var awarenessEvaluationRatings []*models.AwarenessEvaluationRating

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Find(&awarenessEvaluationRatings)

	return awarenessEvaluationRatings, result.Error
}

func (repo *AwarenessEvaluationRatingRepository) GetByID(id uint) (*models.AwarenessEvaluationRating, error) {
	var awarenessEvaluationRating *models.AwarenessEvaluationRating

	result := repo.db.First(&awarenessEvaluationRating, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessEvaluationRating, nil
}

func (repo *AwarenessEvaluationRatingRepository) GetByType(rating string) (*models.AwarenessEvaluationRating, error) {
	var awarenessEvaluationRating *models.AwarenessEvaluationRating

	result := repo.db.Where("type = ?", rating).First(&awarenessEvaluationRating)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessEvaluationRating, nil
}