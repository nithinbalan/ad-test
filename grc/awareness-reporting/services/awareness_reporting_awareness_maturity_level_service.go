package services

import (
	"errors"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessMaturityLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessMaturityLevel, error)
	GetByID(id uint) (*models.AwarenessMaturityLevel, error)
	GetByType(maturityLevel string) (*models.AwarenessMaturityLevel, error)
}

type AwarenessMaturityLevelService struct {
	repo repositories.AwarenessMaturityLevelRepositoryInterface
}

func NewAwarenessMaturityLevelService(
	repo repositories.AwarenessMaturityLevelRepositoryInterface,
) AwarenessMaturityLevelServiceInterface {
	return &AwarenessMaturityLevelService{
		repo: repo,
	}
}

func (amls *AwarenessMaturityLevelService) GetAll(sort string,sortBy string,keyword string,status string,) ([]*models.AwarenessMaturityLevel, error) {
	return amls.repo.GetAll(sort, sortBy, keyword, status)
}

func (amls AwarenessMaturityLevelService) GetByID(id uint) (*models.AwarenessMaturityLevel, error) {
	awarenessMaturityLevel, err := amls.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return awarenessMaturityLevel, nil
}

func (amls *AwarenessMaturityLevelService) GetByType(maturityLevel string) (*models.AwarenessMaturityLevel, error) {
	awarenessMaturityLevel, err := amls.repo.GetByType(maturityLevel)

	if err != nil {
		return nil, err
	}

	return awarenessMaturityLevel, nil
}