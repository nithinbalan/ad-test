package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/threat-intelligence/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type PlannedAdvisoryStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryStatus, error)
	GetByID(id uint) (*models.PlannedAdvisoryStatus, error)
	GetByType(statusType string) (*models.PlannedAdvisoryStatus, error)
}

// Service and implementations
type PlannedAdvisoryStatusService struct {
	repo repositories.PlannedAdvisoryStatusRepositoryInterface
}

func NewPlannedAdvisoryStatusService(repo repositories.PlannedAdvisoryStatusRepositoryInterface) PlannedAdvisoryStatusServiceInterface {
	return &PlannedAdvisoryStatusService{
		repo: repo,
	}
}

func (pass *PlannedAdvisoryStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PlannedAdvisoryStatus, error) {
	return pass.repo.GetAll(sort, sortBy, keyword, status)
}

func (pass PlannedAdvisoryStatusService) GetByID(id uint) (*models.PlannedAdvisoryStatus, error) {
	plannedAdvisoryStatus, err := pass.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return plannedAdvisoryStatus, nil
}

func (pass *PlannedAdvisoryStatusService) GetByType(statusType string) (*models.PlannedAdvisoryStatus, error) {
	statusResult, err := pass.repo.GetByType(statusType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
