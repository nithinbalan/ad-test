package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlType, error)
	GetByID(id uint) (*models.ControlType, error)
}

// Service and implementations
type ControlTypeService struct {
	repo repositories.ControlTypeRepositoryInterface
}

func NewControlTypeService(repo repositories.ControlTypeRepositoryInterface) ControlTypeServiceInterface {
	return &ControlTypeService{
		repo: repo,
	}
}

func (cts *ControlTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ControlType, error) {
	return cts.repo.GetAll(sort, sortBy, keyword, status)
}

func (cts *ControlTypeService) GetByID(id uint) (*models.ControlType, error) {
	controlType, err := cts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlType, nil
}
