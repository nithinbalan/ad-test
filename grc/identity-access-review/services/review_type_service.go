package services

import (
	"errors"
	"jethur/grc/identity-access-review/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ReviewTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewType, error)
	GetByID(id uint) (*models.ReviewType, error)
	GetByType(reviewType string) (*models.ReviewType, error)
}

type ReviewTypeService struct {
	repo repositories.ReviewTypeRepositoryInterface
}

func NewReviewTypeService(
	repo repositories.ReviewTypeRepositoryInterface,
) ReviewTypeServiceInterface {
	return &ReviewTypeService{
		repo: repo,
	}
}

func (rts *ReviewTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewType, error) {
	return rts.repo.GetAll(sort, sortBy, keyword, status)
}

func (rts *ReviewTypeService) GetByID(id uint) (*models.ReviewType, error) {
	reviewType, err := rts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return reviewType, nil
}

func (rts *ReviewTypeService) GetByType(reviewType string) (*models.ReviewType, error) {
	resultReviewType, err := rts.repo.GetByType(reviewType)

	if err != nil {
		return nil, err
	}

	return resultReviewType, nil
}
