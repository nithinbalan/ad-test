package services

import (
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
)

type RiskAppetiteToleranceLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskAppetiteToleranceLevel, error)
	GetByID(id uint) (*models.RiskAppetiteToleranceLevel, error)
}

// Service and implementations
type RiskAppetiteToleranceLevelService struct {
	repo repositories.RiskAppetiteToleranceLevelRepositoryInterface
}

func NewRiskAppetiteToleranceLevelService(repo repositories.RiskAppetiteToleranceLevelRepositoryInterface) RiskAppetiteToleranceLevelServiceInterface {
	return &RiskAppetiteToleranceLevelService{
		repo: repo,
	}
}

func (rats *RiskAppetiteToleranceLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskAppetiteToleranceLevel, error) {
	return rats.repo.GetAll(sort, sortBy, keyword, status)
}

func (rats RiskAppetiteToleranceLevelService) GetByID(id uint) (*models.RiskAppetiteToleranceLevel, error) {
	return rats.repo.GetByID(id)
}
