package services

import (
	"errors"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskTimelineServiceInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskTimelineListQueryParamsDTO) (*riskDTOs.RiskTimelineListResponseDTO, error)
	GetAllTypes() ([]models.RiskTimelineType, error)
	GetByID(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error)
	GetByType(timelineType string) (*models.RiskTimelineType, error)
	GetByStatusType(status string) (*models.TimelineStatus, error)
	Create(timelineFormDTOs []*riskDTOs.RiskTimelineFormDTO, tx ...utility.DatabaseTransaction) error
	Update(riskID uint, riskTimelineFormDTO *riskDTOs.RiskTimelineFormDTO, timelineTypeID uint, tx ...utility.DatabaseTransaction) error
	Delete(riskID uint, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type RiskTimelineService struct {
	repo             repositories.RiskTimelineRepositoryInterface
	riskRegisterRepo repositories.RiskRegisterRepositoryInterface
}

func NewRiskTimelineService(repo repositories.RiskTimelineRepositoryInterface, riskRegisterRepo repositories.RiskRegisterRepositoryInterface) RiskTimelineServiceInterface {
	return &RiskTimelineService{
		repo:             repo,
		riskRegisterRepo: riskRegisterRepo,
	}
}

func (rts *RiskTimelineService) GetAll(riskID uint, params *riskDTOs.RiskTimelineListQueryParamsDTO) (*riskDTOs.RiskTimelineListResponseDTO, error) {
	if _, err := rts.riskRegisterRepo.GetByID(riskID); err != nil {
		return nil, err
	}
	return rts.repo.GetAll(riskID, params)
}

func (rts *RiskTimelineService) GetAllTypes() ([]models.RiskTimelineType, error) {
	return rts.repo.GetAllTypes()
}

func (rts *RiskTimelineService) GetByID(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error) {
	riskTimeline, err := rts.repo.GetByID(riskID, timelineTypeID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, riskID)
		}
		return nil, err
	}
	return riskTimeline, nil
}

func (rts *RiskTimelineService) GetByType(timelineType string) (*models.RiskTimelineType, error) {
	typeResult, err := rts.repo.GetByType(timelineType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return typeResult, nil
}

func (rts *RiskTimelineService) GetByStatusType(status string) (*models.TimelineStatus, error) {
	statusResult, err := rts.repo.GetByStatusType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}

func (rts *RiskTimelineService) Create(timelineFormDTOs []*riskDTOs.RiskTimelineFormDTO, tx ...utility.DatabaseTransaction) error {
	var timelines []models.RiskTimeline
	for _, dto := range timelineFormDTOs {
		timeline := models.RiskTimeline{
			RiskRegisterID:     dto.RiskRegisterID,
			RiskTimelineTypeID: dto.RiskTimelineTypeID,
			TimelineStatusID:   dto.TimelineStatusID,
			PerformedByID:      dto.PerformedByID,
			PerformedAt:        dto.PerformedAt,
		}
		timelines = append(timelines, timeline)
	}

	if err := rts.repo.Create(timelines, tx...); err != nil {
		return err
	}
	return nil
}

func (rts *RiskTimelineService) Update(riskID uint, riskTimelineFormDTO *riskDTOs.RiskTimelineFormDTO, timelineTypeID uint, tx ...utility.DatabaseTransaction) error {
	existingTimeline, err := rts.GetByID(riskID, timelineTypeID, tx...)
	if err != nil {
		return err
	}
	existingTimeline.PerformedByID = riskTimelineFormDTO.PerformedByID
	existingTimeline.PerformedAt = riskTimelineFormDTO.PerformedAt
	existingTimeline.TimelineStatusID = riskTimelineFormDTO.TimelineStatusID
	err = rts.repo.Update(existingTimeline, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (rts *RiskTimelineService) Delete(riskID uint, tx ...utility.DatabaseTransaction) error {
	if _, err := rts.riskRegisterRepo.GetByID(riskID); err != nil {
		return err
	}

	// Delete the entry
	if err := rts.repo.Delete(riskID); err != nil {
		return err
	}

	return nil
}
