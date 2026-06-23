package services

import (
	"errors"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessHandlingStatusServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessHandlingStatus, error)
	GetByID(id uint) (*models.AwarenessHandlingStatus, error)
	GetByType(statusType string) (*models.AwarenessHandlingStatus, error)
}

type AwarenessHandlingStatusService struct {
	repo repositories.AwarenessHandlingStatusRepositoryInterface
}

func NewAwarenessHandlingStatusService(
	repo repositories.AwarenessHandlingStatusRepositoryInterface,
) AwarenessHandlingStatusServiceInterface {
	return &AwarenessHandlingStatusService{
		repo: repo,
	}
}

func (ahss *AwarenessHandlingStatusService) GetAll(sort string,sortBy string,keyword string,status string,) ([]*models.AwarenessHandlingStatus, error) {
	return ahss.repo.GetAll(sort, sortBy, keyword, status)
}

func (ahss *AwarenessHandlingStatusService) GetByID(id uint) (*models.AwarenessHandlingStatus, error) {
	awarenessHandlingStatus, err := ahss.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return awarenessHandlingStatus, nil
}

func (ahss *AwarenessHandlingStatusService) GetByType(statusType string) (*models.AwarenessHandlingStatus, error) {
	awarenessHandlingStatus, err := ahss.repo.GetByType(statusType)

	if err != nil {
		return nil, err
	}

	return awarenessHandlingStatus, nil
}