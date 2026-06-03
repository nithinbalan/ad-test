package services

import (
	"errors"
	"jethur/grc/business-continuity/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BCMImpactAnalysisLevelServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMImpactAnalysisLevel, error)
	GetByID(id uint) (*models.BCMImpactAnalysisLevel, error)
	GetByType(status string) (*models.BCMImpactAnalysisLevel, error)
}

// Service and implementations
type BCMImpactAnalysisLevelService struct {
	repo repositories.BCMImpactAnalysisLevelRepositoryInterface
}

func NewBCMImpactAnalysisLevelService(repo repositories.BCMImpactAnalysisLevelRepositoryInterface) BCMImpactAnalysisLevelServiceInterface {
	return &BCMImpactAnalysisLevelService{
		repo: repo,
	}
}

func (bials *BCMImpactAnalysisLevelService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.BCMImpactAnalysisLevel, error) {
	return bials.repo.GetAll(sort, sortBy, keyword, status)
}

func (bials BCMImpactAnalysisLevelService) GetByID(id uint) (*models.BCMImpactAnalysisLevel, error) {
	bcmImpactAnalysisLevel, err := bials.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return bcmImpactAnalysisLevel, nil
}

func (bials *BCMImpactAnalysisLevelService) GetByType(status string) (*models.BCMImpactAnalysisLevel, error) {
	bcmImpactAnalysisLevel, err := bials.repo.GetByType(status)

	if err != nil {
		return nil, err
	}

	return bcmImpactAnalysisLevel, nil
}
