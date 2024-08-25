package services

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditTimelineServiceInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditTimelineListQueryParamsDTO) (*auditDTOs.AuditTimelineListResponseDTO, error)
	GetAllTypes() ([]models.AuditTimelineType, error)
	GetByID(auditID uint, timelineTypeID uint) (*models.AuditTimeline, error)
	GetByType(timelineType string) (*models.AuditTimelineType, error)
	GetByStatusType(status string) (*models.TimelineStatus, error)
	Create(timelineFormDTOs []*auditDTOs.AuditTimelineFormDTO, tx ...utility.DatabaseTransaction) error
	Update(auditID uint, auditTimelineFormDTO *auditDTOs.AuditTimelineFormDTO, timelineTypeID uint, tx ...utility.DatabaseTransaction) error
	Delete(auditID uint, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type AuditTimelineService struct {
	repo      repositories.AuditTimelineRepositoryInterface
	auditRepo repositories.AuditRepositoryInterface
}

func NewAuditTimelineService(repo repositories.AuditTimelineRepositoryInterface, auditRepo repositories.AuditRepositoryInterface) AuditTimelineServiceInterface {
	return &AuditTimelineService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (ats *AuditTimelineService) GetAll(auditID uint, params *auditDTOs.AuditTimelineListQueryParamsDTO) (*auditDTOs.AuditTimelineListResponseDTO, error) {
	if _, err := ats.auditRepo.GetByID(auditID); err != nil {
		return nil, err
	}
	return ats.repo.GetAll(auditID, params)
}

func (ats *AuditTimelineService) GetAllTypes() ([]models.AuditTimelineType, error) {
	return ats.repo.GetAllTypes()
}

func (ats *AuditTimelineService) GetByID(auditID uint, timelineTypeID uint) (*models.AuditTimeline, error) {
	auditTimeline, err := ats.repo.GetByID(auditID, timelineTypeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, auditID)
		}
		return nil, err
	}
	return auditTimeline, nil
}

func (ats *AuditTimelineService) GetByType(timelineType string) (*models.AuditTimelineType, error) {
	typeResult, err := ats.repo.GetByType(timelineType)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return typeResult, nil
}

func (ats *AuditTimelineService) GetByStatusType(status string) (*models.TimelineStatus, error) {
	statusResult, err := ats.repo.GetByStatusType(status)
	if err != nil {
		// Handle error appropriately
		return nil, err
	}

	return statusResult, nil
}

func (ats *AuditTimelineService) Create(timelineFormDTOs []*auditDTOs.AuditTimelineFormDTO, tx ...utility.DatabaseTransaction) error {
	var timelines []models.AuditTimeline
	for _, dto := range timelineFormDTOs {
		timeline := models.AuditTimeline{
			AuditID:             dto.AuditID,
			AuditTimelineTypeID: dto.AuditTimelineTypeID,
			TimelineStatusID:    dto.TimelineStatusID,
			PerformedByID:       dto.PerformedByID,
			PerformedAt:         dto.PerformedAt,
		}
		timelines = append(timelines, timeline)
	}

	if err := ats.repo.Create(timelines, tx...); err != nil {
		return err
	}
	return nil
}

func (ats *AuditTimelineService) Update(auditID uint, auditTimelineFormDTO *auditDTOs.AuditTimelineFormDTO, timelineTypeID uint, tx ...utility.DatabaseTransaction) error {
	existingTimeline, err := ats.GetByID(auditID, timelineTypeID)
	if err != nil {
		return err
	}
	existingTimeline.PerformedByID = auditTimelineFormDTO.PerformedByID
	existingTimeline.PerformedAt = auditTimelineFormDTO.PerformedAt
	existingTimeline.TimelineStatusID = auditTimelineFormDTO.TimelineStatusID
	err = ats.repo.Update(existingTimeline, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (ats *AuditTimelineService) Delete(auditID uint, tx ...utility.DatabaseTransaction) error {
	if _, err := ats.auditRepo.GetByID(auditID); err != nil {
		return err
	}

	// Delete the entry
	if err := ats.repo.Delete(auditID); err != nil {
		return err
	}

	return nil
}
