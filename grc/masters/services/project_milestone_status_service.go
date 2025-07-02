package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type ProjectMilestoneStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectMilestoneStatus, error)
	GetByID(id uint) (*models.ProjectMilestoneStatus, error)
	GetByType(status string) (*models.ProjectMilestoneStatus, error) // Renamed method

}

// Service and implementations
type ProjectMilestoneStatusService struct {
	repo repositories.ProjectMilestoneStatusRepositoryInterface
}

func NewProjectMilestoneStatusService(repo repositories.ProjectMilestoneStatusRepositoryInterface) ProjectMilestoneStatusServiceInterface {
	return &ProjectMilestoneStatusService{
		repo: repo,
	}
}

func (lss *ProjectMilestoneStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProjectMilestoneStatus, error) {
	return lss.repo.GetAll(sort, sortBy, keyword, status)
}

func (lss ProjectMilestoneStatusService) GetByID(id uint) (*models.ProjectMilestoneStatus, error) {
	return lss.repo.GetByID(id)
}

func (lss *ProjectMilestoneStatusService) GetByType(status string) (*models.ProjectMilestoneStatus, error) {
	// Fetch the project milestone status by its status
	statusResult, err := lss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
