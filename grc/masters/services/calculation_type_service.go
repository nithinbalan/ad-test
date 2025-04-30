package services

import (
	"errors"
	"jethur/grc/masters/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type CalculationTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CalculationType, error)
	GetByID(id uint) (*models.CalculationType, error)
}

// Service and implementations
type CalculationTypeService struct {
	repo repositories.CalculationTypeRepositoryInterface
}

func NewCalculationTypeService(repo repositories.CalculationTypeRepositoryInterface) CalculationTypeServiceInterface {
	return &CalculationTypeService{
		repo: repo,
	}
}

func (cts *CalculationTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CalculationType, error) {
	return cts.repo.GetAll(sort, sortBy, keyword, status)
}

func (cts CalculationTypeService) GetByID(id uint) (*models.CalculationType, error) {
	calculationType, err := cts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return calculationType, nil
}
