package services

import (
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
)

type ControlReviewStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlReviewStatus, error)
	GetByID(id uint) (*models.ControlReviewStatus, error)
	GetByType(status string) (*models.ControlReviewStatus, error)
}

// Service and implementations
type ControlReviewStatusService struct {
	repo repositories.ControlReviewStatusRepositoryInterface
}

func NewControlReviewStatusService(repo repositories.ControlReviewStatusRepositoryInterface) ControlReviewStatusServiceInterface {
	return &ControlReviewStatusService{
		repo: repo,
	}
}

func (crss *ControlReviewStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlReviewStatus, error) {
	return crss.repo.GetAll(sort, sortBy, keyword, status)
}

func (crss ControlReviewStatusService) GetByID(id uint) (*models.ControlReviewStatus, error) {
	return crss.repo.GetByID(id)
}

func (crss *ControlReviewStatusService) GetByType(status string) (*models.ControlReviewStatus, error) {
	statusResult, err := crss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
