package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type StatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Status, error)
	GetByID(id uint) (*models.Status, error)
	GetByType(status string) (*models.Status, error)
}

// Service and implementations
type StatusService struct {
	repo repositories.StatusRepositoryInterface
}

func NewStatusService(repo repositories.StatusRepositoryInterface) StatusServiceInterface {
	return &StatusService{
		repo: repo,
	}
}

func (ss *StatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Status, error) {
	return ss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ss StatusService) GetByID(id uint) (*models.Status, error) {
	return ss.repo.GetByID(id)
}

func (ss *StatusService) GetByType(status string) (*models.Status, error) {
	// Fetch the status by its status
	statusResult, err := ss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
