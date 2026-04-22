package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskRegisterUpdatePlanStatus, error)
	GetByID(id uint) (*models.RiskRegisterUpdatePlanStatus, error)
	GetByType(status string) (*models.RiskRegisterUpdatePlanStatus, error)
}

type RiskRegisterUpdatePlanStatusService struct {
	repo repositories.RiskRegisterUpdatePlanStatusRepositoryInterface
}

func NewRiskRegisterUpdatePlanStatusService(repo repositories.RiskRegisterUpdatePlanStatusRepositoryInterface) RiskRegisterUpdatePlanStatusServiceInterface {
	return &RiskRegisterUpdatePlanStatusService{
		repo: repo,
	}
}

func (rrupss *RiskRegisterUpdatePlanStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskRegisterUpdatePlanStatus, error) {
	return rrupss.repo.GetAll(sort, sortBy, keyword, status)
}

func (rrupss RiskRegisterUpdatePlanStatusService) GetByID(id uint) (*models.RiskRegisterUpdatePlanStatus, error) {
	riskRegisterUpdatePlanStatus, err := rrupss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskRegisterUpdatePlanStatus, nil
}

func (rrupss *RiskRegisterUpdatePlanStatusService) GetByType(status string) (*models.RiskRegisterUpdatePlanStatus, error) {
	statusResult, err := rrupss.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
