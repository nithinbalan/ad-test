package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRIRiskToleranceStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KRIRiskToleranceStatus, error)
	GetByID(id uint) (*models.KRIRiskToleranceStatus, error)
	GetByType(statusType string) (*models.KRIRiskToleranceStatus, error)
}

type KRIRiskToleranceStatusService struct {
	repo repositories.KRIRiskToleranceStatusRepositoryInterface
}

func NewKRIRiskToleranceStatusService(repo repositories.KRIRiskToleranceStatusRepositoryInterface) KRIRiskToleranceStatusServiceInterface {
	return &KRIRiskToleranceStatusService{
		repo: repo,
	}
}

func (krtss *KRIRiskToleranceStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.KRIRiskToleranceStatus, error) {
	return krtss.repo.GetAll(sort, sortBy, keyword, status)
}

func (krtss KRIRiskToleranceStatusService) GetByID(id uint) (*models.KRIRiskToleranceStatus, error) {
	status, err := krtss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return status, nil
}

func (krtss *KRIRiskToleranceStatusService) GetByType(statusType string) (*models.KRIRiskToleranceStatus, error) {
	status, err := krtss.repo.GetByType(statusType)
	if err != nil {
		return nil, err
	}

	return status, nil
}
