package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type TriggerConditionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TriggerCondition, error)
	GetByID(id uint) (*models.TriggerCondition, error)
	GetByType(status string) (*models.TriggerCondition, error) // Renamed method
}

// Service and implementations
type TriggerConditionService struct {
	repo repositories.TriggerConditionRepositoryInterface
}

func NewTriggerConditionService(repo repositories.TriggerConditionRepositoryInterface) TriggerConditionServiceInterface {
	return &TriggerConditionService{
		repo: repo,
	}
}

func (tcs *TriggerConditionService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.TriggerCondition, error) {
	return tcs.repo.GetAll(sort, sortBy, keyword, status)
}

func (tcs TriggerConditionService) GetByID(id uint) (*models.TriggerCondition, error) {
	triggerCondition, err := tcs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return triggerCondition, nil

}

func (tcs *TriggerConditionService) GetByType(status string) (*models.TriggerCondition, error) {
	statusResult, err := tcs.repo.GetByType(status)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}
