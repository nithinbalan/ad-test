package services

import (
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
)

type LicenseTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LicenseType, error)
	GetByID(id uint) (*models.LicenseType, error)
	GetByType(status string) (*models.LicenseType, error) // Renamed method

}

// Service and implementations
type LicenseTypeService struct {
	repo repositories.LicenseTypeRepositoryInterface
}

func NewLicenseTypeService(repo repositories.LicenseTypeRepositoryInterface) LicenseTypeServiceInterface {
	return &LicenseTypeService{
		repo: repo,
	}
}

func (lts *LicenseTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.LicenseType, error) {
	return lts.repo.GetAll(sort, sortBy, keyword, status)
}

func (lts LicenseTypeService) GetByID(id uint) (*models.LicenseType, error) {
	return lts.repo.GetByID(id)
}

func (lts *LicenseTypeService) GetByType(status string) (*models.LicenseType, error) {
	// Fetch the lifecycle status by its status
	statusResult, err := lts.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
