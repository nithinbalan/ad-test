package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlStatus, error)
	GetByID(id uint) (*models.ControlStatus, error)
	GetByType(typeFilter string) (*models.ControlStatus, error)
}

// Service and implementations
type ControlStatusService struct {
	repo repositories.ControlStatusRepositoryInterface
}

func NewControlStatusService(repo repositories.ControlStatusRepositoryInterface) ControlStatusServiceInterface {
	return &ControlStatusService{
		repo: repo,
	}
}

func (css *ControlStatusService) GetAll(sort string, sortBy string, keyword string, status string, typeFilter string) ([]*models.ControlStatus, error) {
	return css.repo.GetAll(sort, sortBy, keyword, status, typeFilter)
}

func (css ControlStatusService) GetByID(id uint) (*models.ControlStatus, error) {
	controlStatus, err := css.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlStatus, nil

}

func (css ControlStatusService) GetByType(typeFilter string) (*models.ControlStatus, error) {
	controlStatus, err := css.repo.GetByType(typeFilter)
	if err != nil {
		return nil, err
	}

	return controlStatus, nil
}
