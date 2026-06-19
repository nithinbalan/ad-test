package services

import (
	"errors"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessEvaluationRatingServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessEvaluationRating, error)
	GetByID(id uint) (*models.AwarenessEvaluationRating, error)
	GetByType(rating string) (*models.AwarenessEvaluationRating, error)
}

type AwarenessEvaluationRatingService struct {
	repo repositories.AwarenessEvaluationRatingRepositoryInterface
}

func NewAwarenessEvaluationRatingService(
	repo repositories.AwarenessEvaluationRatingRepositoryInterface,
) AwarenessEvaluationRatingServiceInterface {
	return &AwarenessEvaluationRatingService{
		repo: repo,
	}
}

func (aers *AwarenessEvaluationRatingService) GetAll(sort string,sortBy string,keyword string,status string,) ([]*models.AwarenessEvaluationRating, error) {
	return aers.repo.GetAll(sort, sortBy, keyword, status)
}

func (aers AwarenessEvaluationRatingService) GetByID(id uint) (*models.AwarenessEvaluationRating, error) {
	awarenessEvaluationRating, err := aers.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return awarenessEvaluationRating, nil
}

func (aers *AwarenessEvaluationRatingService) GetByType(rating string) (*models.AwarenessEvaluationRating, error) {
	awarenessEvaluationRating, err := aers.repo.GetByType(rating)

	if err != nil {
		return nil, err
	}

	return awarenessEvaluationRating, nil
}