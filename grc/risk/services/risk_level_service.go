package services

import (
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
)

type RiskLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskLevel, error)
	GetByID(id uint) (*models.RiskLevel, error)
	GetByType(status string) (*models.RiskLevel, error)
}

// Service and implementations
type RiskLevelService struct {
	repo repositories.RiskLevelRepositoryInterface
}

func NewRiskLevelService(repo repositories.RiskLevelRepositoryInterface) RiskLevelServiceInterface {
	return &RiskLevelService{
		repo: repo,
	}
}

func (rls *RiskLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskLevel, error) {
	return rls.repo.GetAll(sort, sortBy, keyword, status)
}

func (rls RiskLevelService) GetByID(id uint) (*models.RiskLevel, error) {
	return rls.repo.GetByID(id)
}

func (rls *RiskLevelService) GetByType(status string) (*models.RiskLevel, error) {
	// Fetch the risk level by its status
	statusResult, err := rls.repo.GetByType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}
