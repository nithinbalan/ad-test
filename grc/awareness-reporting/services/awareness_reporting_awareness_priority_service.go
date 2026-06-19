package services

import (
	"errors"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessPriorityServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessPriority, error)
	GetByID(id uint) (*models.AwarenessPriority, error)
	GetByType(priority string) (*models.AwarenessPriority, error)
}

type AwarenessPriorityService struct {
	repo repositories.AwarenessPriorityRepositoryInterface
}

func NewAwarenessPriorityService(
	repo repositories.AwarenessPriorityRepositoryInterface,
) AwarenessPriorityServiceInterface {
	return &AwarenessPriorityService{
		repo: repo,
	}
}

func (aps *AwarenessPriorityService) GetAll(sort string,sortBy string,keyword string,status string,) ([]*models.AwarenessPriority, error) {
	return aps.repo.GetAll(sort, sortBy, keyword, status)
}

func (aps AwarenessPriorityService) GetByID(id uint) (*models.AwarenessPriority, error) {
	awarenessPriority, err := aps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return awarenessPriority, nil
}

func (aps *AwarenessPriorityService) GetByType(priority string) (*models.AwarenessPriority, error) {
	awarenessPriority, err := aps.repo.GetByType(priority)

	if err != nil {
		return nil, err
	}

	return awarenessPriority, nil
}