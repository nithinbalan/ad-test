package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CommitteeTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeType, error)
	GetByID(id uint) (*models.CommitteeType, error)
	GetByType(typeFilter string) (*models.CommitteeType, error)
}

// Service and implementations
type CommitteeTypeService struct {
	repo repositories.CommitteeTypeRepositoryInterface
}

func NewCommitteeTypeService(repo repositories.CommitteeTypeRepositoryInterface) CommitteeTypeServiceInterface {
	return &CommitteeTypeService{
		repo: repo,
	}
}

func (cts *CommitteeTypeService) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.CommitteeType, error) {
	return cts.repo.GetAll(sort, sortBy, keyword, status, typeFilter)
}

func (cts CommitteeTypeService) GetByID(id uint) (*models.CommitteeType, error) {
	committeeType, err := cts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return committeeType, nil

}

func (cts *CommitteeTypeService) GetByType(typeFilter string) (*models.CommitteeType, error) {
	committeeType, err := cts.repo.GetByType(typeFilter)
	if err != nil {
		return nil, err
	}

	return committeeType, nil
}
