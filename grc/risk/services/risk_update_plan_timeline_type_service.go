package services

import (
	"errors"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskUpdatePlanTimelineTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.RiskUpdatePlanTimelineType, error)
	GetByID(id uint) (*models.RiskUpdatePlanTimelineType, error)
	GetByType(timelineType string) (*models.RiskUpdatePlanTimelineType, error)
}

// Service and implementations
type RiskUpdatePlanTimelineTypeService struct {
	repo repositories.RiskUpdatePlanTimelineTypeRepositoryInterface
}

func NewRiskUpdatePlanTimelineTypeService(repo repositories.RiskUpdatePlanTimelineTypeRepositoryInterface) RiskUpdatePlanTimelineTypeServiceInterface {
	return &RiskUpdatePlanTimelineTypeService{
		repo: repo,
	}
}

func (rts *RiskUpdatePlanTimelineTypeService) GetAll(sort string, sortBy string, keyword string) ([]*models.RiskUpdatePlanTimelineType, error) {
	return rts.repo.GetAll(sort, sortBy, keyword)
}

func (rts RiskUpdatePlanTimelineTypeService) GetByID(id uint) (*models.RiskUpdatePlanTimelineType, error) {
	timelineType, err := rts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return timelineType, nil
}

func (rts *RiskUpdatePlanTimelineTypeService) GetByType(timelineType string) (*models.RiskUpdatePlanTimelineType, error) {
	typeResult, err := rts.repo.GetByType(timelineType)
	if err != nil {
		return nil, err
	}

	return typeResult, nil
}
