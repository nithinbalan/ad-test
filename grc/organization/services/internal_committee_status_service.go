package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type InternalCommitteeStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.InternalCommitteeStatus, error)
	GetByID(id uint) (*models.InternalCommitteeStatus, error)
	GetByType(typeFilter string) (*models.InternalCommitteeStatus, error)
}

// Service and implementations
type InternalCommitteeStatusService struct {
	repo repositories.InternalCommitteeStatusRepositoryInterface
}

func NewInternalCommitteeStatusService(repo repositories.InternalCommitteeStatusRepositoryInterface) InternalCommitteeStatusServiceInterface {
	return &InternalCommitteeStatusService{
		repo: repo,
	}
}

func (icss *InternalCommitteeStatusService) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.InternalCommitteeStatus, error) {
	return icss.repo.GetAll(sort, sortBy, keyword, status, typeFilter)
}

func (icss InternalCommitteeStatusService) GetByID(id uint) (*models.InternalCommitteeStatus, error) {
	internalCommitteeStatus, err := icss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return internalCommitteeStatus, nil

}

func (icss InternalCommitteeStatusService) GetByType(typeFilter string) (*models.InternalCommitteeStatus, error) {
	internalCommitteeStatus, err := icss.repo.GetByType(typeFilter)
	if err != nil {
		return nil, err
	}

	return internalCommitteeStatus, nil
}
