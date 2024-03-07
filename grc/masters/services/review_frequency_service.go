package services

import (
	"jethur/grc/masters/models"
	"jethur/grc/masters/repositories"
)

// interfaces
type ReviewFrequencyServiceInterface interface {
	GetAll() ([]models.ReviewFrequency, error)
}

// Service and implementations
type ReviewFrequencyService struct {
	repo repositories.ReviewFrequencyRepositoryInterface
}

func NewReviewFrequencyService(repo repositories.ReviewFrequencyRepositoryInterface) *ReviewFrequencyService {
	return &ReviewFrequencyService{
		repo: repo,
	}
}

func (rs *ReviewFrequencyService) GetAll() ([]models.ReviewFrequency, error) {
	return rs.repo.GetAll()
}
