package services

import (
	"errors"

	"jethur/grc/models"
	"jethur/grc/settings/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ReminderContextServiceInterface interface {
	GetAllReminderContexts() ([]*models.ReminderContext, error)
	GetByID(id uint) (*models.ReminderContext, error)
}

// Service and implementations
type ReminderContextService struct {
	repo repositories.ReminderContextRepositoryInterface
}

func NewReminderContextService(repo repositories.ReminderContextRepositoryInterface) ReminderContextServiceInterface {
	return &ReminderContextService{
		repo: repo,
	}
}

func (rcs *ReminderContextService) GetAllReminderContexts() ([]*models.ReminderContext, error) {
	return rcs.repo.GetAllReminderContexts()
}

func (rcs *ReminderContextService) GetByID(id uint) (*models.ReminderContext, error) {
	reminder, err := rcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return reminder, nil
}
