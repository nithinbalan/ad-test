package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type FeedbackOptionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.FeedbackOption, error)
	GetByID(id uint) (*models.FeedbackOption, error)
	GetByType(status string) (*models.FeedbackOption, error)
}

// Service and implementations
type FeedbackOptionService struct {
	repo repositories.FeedbackOptionRepositoryInterface
}

func NewFeedbackOptionService(repo repositories.FeedbackOptionRepositoryInterface) FeedbackOptionServiceInterface {
	return &FeedbackOptionService{
		repo: repo,
	}
}

func (tess *FeedbackOptionService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.FeedbackOption, error) {
	return tess.repo.GetAll(sort, sortBy, keyword, status)
}

func (tess FeedbackOptionService) GetByID(id uint) (*models.FeedbackOption, error) {
	feedbackOption, err := tess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return feedbackOption, nil
}

func (tess *FeedbackOptionService) GetByType(status string) (*models.FeedbackOption, error) {
	feedbackOption, err := tess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return feedbackOption, nil
}
