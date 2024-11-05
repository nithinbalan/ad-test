package services

import (
	"errors"
	"jethur/grc/asset/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatStatus, error)
	GetByID(id uint) (*models.ThreatStatus, error)
	GetByType(status string) (*models.ThreatStatus, error)
}

// Service and implementations
type ThreatStatusService struct {
	repo repositories.ThreatStatusRepositoryInterface
}

func NewThreatStatusService(repo repositories.ThreatStatusRepositoryInterface) ThreatStatusServiceInterface {
	return &ThreatStatusService{
		repo: repo,
	}
}

func (tss *ThreatStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThreatStatus, error) {
	return tss.repo.GetAll(sort, sortBy, keyword, status)
}

func (tss ThreatStatusService) GetByID(id uint) (*models.ThreatStatus, error) {
	threatStatus, err := tss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return threatStatus, nil
}

func (tss *ThreatStatusService) GetByType(status string) (*models.ThreatStatus, error) {
	threatStatus, err := tss.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return threatStatus, nil
}
