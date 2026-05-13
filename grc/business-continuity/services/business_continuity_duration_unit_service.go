package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BusinessContinuityDurationUnitServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BusinessContinuityDurationUnit, error)
	GetByID(id uint) (*models.BusinessContinuityDurationUnit, error)
	GetByType(durationType string) (*models.BusinessContinuityDurationUnit, error)
}

// Service and implementations
type BusinessContinuityDurationUnitService struct {
	repo repositories.BusinessContinuityDurationUnitRepositoryInterface
}

func NewBusinessContinuityDurationUnitService(repo repositories.BusinessContinuityDurationUnitRepositoryInterface) BusinessContinuityDurationUnitServiceInterface {
	return &BusinessContinuityDurationUnitService{
		repo: repo,
	}
}

func (bcdus *BusinessContinuityDurationUnitService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BusinessContinuityDurationUnit, error) {
	return bcdus.repo.GetAll(sort, sortBy, keyword, status)
}

func (bcdus BusinessContinuityDurationUnitService) GetByID(id uint) (*models.BusinessContinuityDurationUnit, error) {
	businessContinuityDurationUnit, err := bcdus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return businessContinuityDurationUnit, nil
}

func (bcdus *BusinessContinuityDurationUnitService) GetByType(durationType string) (*models.BusinessContinuityDurationUnit, error) {
	businessContinuityDurationUnit, err := bcdus.repo.GetByType(durationType)

	if err != nil {
		return nil, err
	}

	return businessContinuityDurationUnit, nil
}
