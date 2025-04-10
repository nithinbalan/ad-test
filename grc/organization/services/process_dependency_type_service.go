package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessDependencyTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessDependencyType, error)
	GetByID(id uint) (*models.ProcessDependencyType, error)
	GetByType(status string) (*models.ProcessDependencyType, error)
}

// Service and implementations
type ProcessDependencyTypeService struct {
	repo repositories.ProcessDependencyTypeRepositoryInterface
}

func NewProcessDependencyTypeService(repo repositories.ProcessDependencyTypeRepositoryInterface) ProcessDependencyTypeServiceInterface {
	return &ProcessDependencyTypeService{
		repo: repo,
	}
}

func (pdts *ProcessDependencyTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessDependencyType, error) {
	return pdts.repo.GetAll(sort, sortBy, keyword, status)
}

func (pdts ProcessDependencyTypeService) GetByID(id uint) (*models.ProcessDependencyType, error) {
	processDependencyType, err := pdts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return processDependencyType, nil
}

func (pdts *ProcessDependencyTypeService) GetByType(status string) (*models.ProcessDependencyType, error) {
	processDependencyType, err := pdts.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return processDependencyType, nil
}
