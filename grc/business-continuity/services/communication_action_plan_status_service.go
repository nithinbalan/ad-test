package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type CommunicationActionPlanStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommunicationActionPlanStatus, error)
	GetByID(id uint) (*models.CommunicationActionPlanStatus, error)
	GetByType(status string) (*models.CommunicationActionPlanStatus, error)
}

// Service and implementations
type CommunicationActionPlanStatusService struct {
	repo repositories.CommunicationActionPlanStatusRepositoryInterface
}

func NewCommunicationActionPlanStatusService(repo repositories.CommunicationActionPlanStatusRepositoryInterface) CommunicationActionPlanStatusServiceInterface {
	return &CommunicationActionPlanStatusService{
		repo: repo,
	}
}

func (tess *CommunicationActionPlanStatusService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CommunicationActionPlanStatus, error) {
	return tess.repo.GetAll(sort, sortBy, keyword, status)
}

func (tess CommunicationActionPlanStatusService) GetByID(id uint) (*models.CommunicationActionPlanStatus, error) {
	communicationActionPlanStatus, err := tess.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return communicationActionPlanStatus, nil
}

func (tess *CommunicationActionPlanStatusService) GetByType(status string) (*models.CommunicationActionPlanStatus, error) {
	communicationActionPlanStatus, err := tess.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return communicationActionPlanStatus, nil
}
