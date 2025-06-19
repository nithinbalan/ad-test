package services

import (
	"jethur/grc/models"
	"jethur/grc/strategy/repositories"
)

type KPIStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIStatus, error)
	GetByID(id uint) (*models.KPIStatus, error)
	GetByType(status string) (*models.KPIStatus, error) // Renamed method

}

// Service and implementations
type KPIStatusService struct {
	repo repositories.KPIStatusRepositoryInterface
}

func NewKPIStatusService(repo repositories.KPIStatusRepositoryInterface) KPIStatusServiceInterface {
	return &KPIStatusService{
		repo: repo,
	}
}

func (lss *KPIStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIStatus, error) {
	return lss.repo.GetAll(sort, sortBy, keyword, status)
}

func (lss KPIStatusService) GetByID(id uint) (*models.KPIStatus, error) {
	return lss.repo.GetByID(id)
}

func (lss *KPIStatusService) GetByType(status string) (*models.KPIStatus, error) {
	// Fetch the kpi status by its status
	statusResult, err := lss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
