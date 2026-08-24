package services

import (
	"jethur/grc/kpi-register/repositories"
	"jethur/grc/models"
)

type KPIReviewStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIReviewStatus, error)
	GetByID(id uint) (*models.KPIReviewStatus, error)
	GetByType(status string) (*models.KPIReviewStatus, error) // Renamed method

}

// Service and implementations
type KPIReviewStatusService struct {
	repo repositories.KPIReviewStatusRepositoryInterface
}

func NewKPIReviewStatusService(repo repositories.KPIReviewStatusRepositoryInterface) KPIReviewStatusServiceInterface {
	return &KPIReviewStatusService{
		repo: repo,
	}
}

func (krss *KPIReviewStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KPIReviewStatus, error) {
	return krss.repo.GetAll(sort, sortBy, keyword, status)
}

func (krss KPIReviewStatusService) GetByID(id uint) (*models.KPIReviewStatus, error) {
	return krss.repo.GetByID(id)
}

func (krss *KPIReviewStatusService) GetByType(status string) (*models.KPIReviewStatus, error) {
	// Fetch the kpi status by its status
	statusResult, err := krss.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
