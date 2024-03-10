package services

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
)

type ReviewFrequencyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error)
}

// Service and implementations
type ReviewFrequencyService struct {
	repo repositories.ReviewFrequencyRepositoryInterface
}

func NewReviewFrequencyService(repo repositories.ReviewFrequencyRepositoryInterface) ReviewFrequencyServiceInterface {
	return &ReviewFrequencyService{
		repo: repo,
	}
}

func (sls *ReviewFrequencyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error) {
	return sls.repo.GetAll(sort, sortBy, keyword, status)
}
