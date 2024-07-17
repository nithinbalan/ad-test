package services

import (
	"errors"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentTimelineServiceInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentTimelineListQueryParamsDTO) (*incidentDTOs.IncidentTimelineListResponseDTO, error)
	GetByID(incidentID uint, timelineTypeID uint) (*models.IncidentTimeline, error)
	GetByType(timelineType string) (*models.IncidentTimelineType, error)
	GetByStatusType(status string) (*models.TimelineStatus, error)
	Create(incidentTimelineFormDTO *incidentDTOs.IncidentTimelineFormDTO, tx ...utility.DatabaseTransaction) error
	Update(incidentID uint, incidentTimelineFormDTO *incidentDTOs.IncidentTimelineFormDTO, timelineTypeID uint, tx ...utility.DatabaseTransaction) error
	Delete(incidentID uint, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type IncidentTimelineService struct {
	repo                 repositories.IncidentTimelineRepositoryInterface
	incidentRegisterRepo repositories.IncidentRegisterRepositoryInterface
}

func NewIncidentTimelineService(repo repositories.IncidentTimelineRepositoryInterface, incidentRegisterRepo repositories.IncidentRegisterRepositoryInterface) IncidentTimelineServiceInterface {
	return &IncidentTimelineService{
		repo:                 repo,
		incidentRegisterRepo: incidentRegisterRepo,
	}
}

func (its *IncidentTimelineService) GetAll(incidentID uint, params *incidentDTOs.IncidentTimelineListQueryParamsDTO) (*incidentDTOs.IncidentTimelineListResponseDTO, error) {
	if _, err := its.incidentRegisterRepo.GetByID(incidentID); err != nil {
		return nil, err
	}
	return its.repo.GetAll(incidentID, params)
}

func (its *IncidentTimelineService) GetByID(incidentID uint, timelineTypeID uint) (*models.IncidentTimeline, error) {
	incidentTimeline, err := its.repo.GetByID(incidentID, timelineTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentID)
		}
		return nil, err
	}
	return incidentTimeline, nil
}

func (its *IncidentTimelineService) GetByType(timelineType string) (*models.IncidentTimelineType, error) {
	typeResult, err := its.repo.GetByType(timelineType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return typeResult, nil
}

func (its *IncidentTimelineService) GetByStatusType(status string) (*models.TimelineStatus, error) {
	statusResult, err := its.repo.GetByStatusType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}

func (its *IncidentTimelineService) Create(incidentTimelineFormDTO *incidentDTOs.IncidentTimelineFormDTO, tx ...utility.DatabaseTransaction) error {
	timeline := &models.IncidentTimeline{
		IncidentRegisterID:     incidentTimelineFormDTO.IncidentRegisterID,
		IncidentTimelineTypeID: incidentTimelineFormDTO.IncidentTimelineTypeID,
		TimelineStatusID:       incidentTimelineFormDTO.TimelineStatusID,
		PerformedByID:          incidentTimelineFormDTO.PerformedByID,
		PerformedAt:            incidentTimelineFormDTO.PerformedAt,
	}
	if err := its.repo.Create(timeline, tx...); err != nil {
		return err
	}
	return nil
}

func (its *IncidentTimelineService) Update(incidentID uint, incidentTimelineFormDTO *incidentDTOs.IncidentTimelineFormDTO, timelineTypeID uint, tx ...utility.DatabaseTransaction) error {
	existingTimeline, err := its.GetByID(incidentID, timelineTypeID)
	if err != nil {
		return err
	}
	existingTimeline.PerformedByID = incidentTimelineFormDTO.PerformedByID
	existingTimeline.PerformedAt = incidentTimelineFormDTO.PerformedAt
	err = its.repo.Update(existingTimeline, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (its *IncidentTimelineService) Delete(incidentID uint, tx ...utility.DatabaseTransaction) error {
	if _, err := its.incidentRegisterRepo.GetByID(incidentID); err != nil {
		return err
	}

	// Delete the entry
	if err := its.repo.Delete(incidentID); err != nil {
		return err
	}

	return nil
}
