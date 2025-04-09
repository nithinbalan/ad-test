package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessCriticalityLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessCriticalityLevel, error)
	GetByID(id uint) (*models.ProcessCriticalityLevel, error)
	GetByType(status string) (*models.ProcessCriticalityLevel, error)
}

// Service and implementations
type ProcessCriticalityLevelService struct {
	repo repositories.ProcessCriticalityLevelRepositoryInterface
}

func NewProcessCriticalityLevelService(repo repositories.ProcessCriticalityLevelRepositoryInterface) ProcessCriticalityLevelServiceInterface {
	return &ProcessCriticalityLevelService{
		repo: repo,
	}
}

func (pcls *ProcessCriticalityLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessCriticalityLevel, error) {
	return pcls.repo.GetAll(sort, sortBy, keyword, status)
}

func (pcls ProcessCriticalityLevelService) GetByID(id uint) (*models.ProcessCriticalityLevel, error) {
	processCriticalityLevel, err := pcls.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return processCriticalityLevel, nil
}

func (pcls *ProcessCriticalityLevelService) GetByType(status string) (*models.ProcessCriticalityLevel, error) {
	processCriticalityLevel, err := pcls.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return processCriticalityLevel, nil
}
