package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpPhaseServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpPhase, error)
	GetByID(id uint) (*models.BcpPhase, error)
}

// Service and implementations
type BcpPhaseService struct {
	repo repositories.BcpPhaseRepositoryInterface
}

func NewBcpPhaseService(repo repositories.BcpPhaseRepositoryInterface) BcpPhaseServiceInterface {
	return &BcpPhaseService{
		repo: repo,
	}
}

func (ics *BcpPhaseService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BcpPhase, error) {
	return ics.repo.GetAll(sort, sortBy, keyword, status)
}

func (ics BcpPhaseService) GetByID(id uint) (*models.BcpPhase, error) {
	bcpPhase, err := ics.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return bcpPhase, nil
}
