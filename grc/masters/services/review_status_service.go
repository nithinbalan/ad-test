package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type ReviewStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewStatus, error)
	GetByID(id uint) (*models.ReviewStatus, error)
	GetByType(status string) (*models.ReviewStatus, error)
}

// Service and implementations
type ReviewStatusService struct {
	repo repositories.ReviewStatusRepositoryInterface
}

func NewReviewStatusService(repo repositories.ReviewStatusRepositoryInterface) ReviewStatusServiceInterface {
	return &ReviewStatusService{
		repo: repo,
	}
}

func (rss *ReviewStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewStatus, error) {
	return rss.repo.GetAll(sort, sortBy, keyword, status)
}

func (rss *ReviewStatusService) GetByID(id uint) (*models.ReviewStatus, error) {
	return rss.repo.GetByID(id)
}

func (rss *ReviewStatusService) GetByType(status string) (*models.ReviewStatus, error) {
	// Fetch the review status by its type
	statusResult, err := rss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
