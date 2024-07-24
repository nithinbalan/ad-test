package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskStatus, error)
	GetByID(id uint) (*models.RiskStatus, error)
	GetByType(status string) (*models.RiskStatus, error)
}

// Service and implementations
type RiskStatusService struct {
	repo repositories.RiskStatusRepositoryInterface
}

func NewRiskStatusService(repo repositories.RiskStatusRepositoryInterface) RiskStatusServiceInterface {
	return &RiskStatusService{
		repo: repo,
	}
}

func (rss *RiskStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskStatus, error) {
	return rss.repo.GetAll(sort, sortBy, keyword, status)
}

func (rss RiskStatusService) GetByID(id uint) (*models.RiskStatus, error) {
	riskStatus, err := rss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskStatus, nil
}

func (rss *RiskStatusService) GetByType(status string) (*models.RiskStatus, error) {
	statusResult, err := rss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
