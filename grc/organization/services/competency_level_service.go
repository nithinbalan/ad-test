package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CompetencyLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyLevel, error)
	GetByID(id uint) (*models.CompetencyLevel, error)
}

// Service and implementations
type CompetencyLevelService struct {
	repo repositories.CompetencyLevelRepositoryInterface
}

func NewCompetencyLevelService(repo repositories.CompetencyLevelRepositoryInterface) CompetencyLevelServiceInterface {
	return &CompetencyLevelService{
		repo: repo,
	}
}

func (cls *CompetencyLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyLevel, error) {
	return cls.repo.GetAll(sort, sortBy, keyword, status)
}

func (cls CompetencyLevelService) GetByID(id uint) (*models.CompetencyLevel, error) {
	controlEfficiencyStatus, err := cls.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return controlEfficiencyStatus, nil
}
