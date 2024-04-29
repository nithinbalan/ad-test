package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
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
	riskLevel, err := rls.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskLevel, nil
}

func (rls *RiskLevelService) GetByType(status string) (*models.RiskLevel, error) {
	// Fetch the risk level by its status
	statusResult, err := rls.repo.GetByType(status)
	if err != nil {
		// Check if the error is a record not found error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, "Risk Level with status "+status)
		} else {
			// Handle other errors
			return nil, err
		}

	}

	return statusResult, nil
}
