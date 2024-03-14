package services

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
)

type ReviewFrequencyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error)
	GetByID(id uint) (*models.ReviewFrequency, error)
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

func (rfs *ReviewFrequencyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewFrequency, error) {
	return rfs.repo.GetAll(sort, sortBy, keyword, status)
}


func (rfs ReviewFrequencyService) GetByID(id uint) (*models.ReviewFrequency, error) {
	return rfs.repo.GetByID(id)
}