package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentCommunicationLogServiceInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentCommunicationLogListQueryParamsDTO) (*incidentDTOs.IncidentCommunicationLogListResponseDTO, error)
	GetByID(incidentID uint, id uint) (*models.IncidentCommunicationLog, error)
	GetFullDetailsByID(incidentID uint, id uint) (*incidentDTOs.IncidentCommunicationLogDetailDTO, error)
	Create(incidentID uint, currentUserID uint, IncidentCommunicationLogFormDTO *incidentDTOs.IncidentCommunicationLogFormDTO) (uint, error)
	Update(incidentID uint, id uint, currentUserID uint, incident *incidentDTOs.IncidentCommunicationLogFormDTO) error
	Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type IncidentCommunicationLogService struct {
	repo                incidentRepositories.IncidentCommunicationLogRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
	incidentTimelineSvc IncidentTimelineServiceInterface
}

func NewIncidentCommunicationLogService(
	repo incidentRepositories.IncidentCommunicationLogRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	incidentTimelineSvc IncidentTimelineServiceInterface,
) IncidentCommunicationLogServiceInterface {
	return &IncidentCommunicationLogService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentRegisterSvc: incidentRegisterSvc,
		incidentTimelineSvc: incidentTimelineSvc,
	}
}

func (icls *IncidentCommunicationLogService) GetAll(incidentID uint, params *incidentDTOs.IncidentCommunicationLogListQueryParamsDTO) (*incidentDTOs.IncidentCommunicationLogListResponseDTO, error) {
	if _, err := icls.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return icls.repo.GetAll(incidentID, params)
}

func (icls IncidentCommunicationLogService) GetByID(incidentID uint, id uint) (*models.IncidentCommunicationLog, error) {
	if _, err := icls.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentCommunicationLog, err := icls.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentCommunicationLog, nil
}

func (icls IncidentCommunicationLogService) GetFullDetailsByID(incidentID uint, id uint) (*incidentDTOs.IncidentCommunicationLogDetailDTO, error) {
	var incidentCommunicationLogDTO incidentDTOs.IncidentCommunicationLogDetailDTO

	if _, err := icls.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentCommunicationLog, err := icls.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(incidentCommunicationLog, &incidentCommunicationLogDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &incidentCommunicationLogDTO, nil
}

func (icls *IncidentCommunicationLogService) Create(incidentID uint, currentUserID uint, incidentCommunicationLogFormDTO *incidentDTOs.IncidentCommunicationLogFormDTO) (uint, error) {
	tx := icls.repo.BeginTransaction()

	if _, err := icls.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}

	dates, err := icls.parseIncidentCommunicationLogDates(incidentCommunicationLogFormDTO)
	if err != nil {
		return 0, err
	}

	incidentCommunicationLog := &models.IncidentCommunicationLog{
		IncidentRegisterID:  incidentID,
		CommunicatedAt:      *dates.CommunicatedAt,
		CommunicatedBy:      incidentCommunicationLogFormDTO.CommunicatedBy,
		CommunicatedTo:      incidentCommunicationLogFormDTO.CommunicatedTo,
		CommunicationMethod: incidentCommunicationLogFormDTO.CommunicationMethod,
		Summary:             &incidentCommunicationLogFormDTO.Summary,
		CreatedByID:         currentUserID,
	}

	id, err := icls.repo.Create(incidentCommunicationLog, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := icls.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentCommunicationLog{}.TableName(),
			id,
			incidentCommunicationLogFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentCommunicationLog{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (icls *IncidentCommunicationLogService) Update(incidentID uint, id uint, currentUserID uint, incidentCommunicationLogFormDTO *incidentDTOs.IncidentCommunicationLogFormDTO) error {
	tx := icls.repo.BeginTransaction()
	if _, err := icls.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	existingIncidentCommunicationLog, err := icls.GetByID(incidentID, id)
	if err != nil {
		return err
	}

	dates, err := icls.parseIncidentCommunicationLogDates(incidentCommunicationLogFormDTO)
	if err != nil {
		return err
	}

	// Update Incident CommunicationLog
	existingIncidentCommunicationLog.CommunicatedAt = *dates.CommunicatedAt
	existingIncidentCommunicationLog.CommunicatedBy = incidentCommunicationLogFormDTO.CommunicatedBy
	existingIncidentCommunicationLog.CommunicatedTo = incidentCommunicationLogFormDTO.CommunicatedTo
	existingIncidentCommunicationLog.CommunicationMethod = incidentCommunicationLogFormDTO.CommunicationMethod
	existingIncidentCommunicationLog.Summary = &incidentCommunicationLogFormDTO.Summary
	existingIncidentCommunicationLog.UpdatedByID = &currentUserID

	err = icls.repo.Update(existingIncidentCommunicationLog, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	//Update timeline
	//todo
	timelineType, err := icls.incidentTimelineSvc.GetByType("communication")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := icls.incidentTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}
	currentTime := time.Now()
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := icls.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		err := icls.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentCommunicationLog{}.TableName(),
			id,
			incidentCommunicationLogFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.IncidentCommunicationLog{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (icls *IncidentCommunicationLogService) Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := icls.repo.BeginTransaction()

	if _, err := icls.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := icls.repo.GetByID(id)

		if err != nil {
			continue
		}
		//Delete-> update timeline performedAt as null
		//todo
		timelineType, err := icls.incidentTimelineSvc.GetByType("communication")
		if err != nil {
			tx.Rollback()
			return failedIds, err
		}
		timelineStatus, err := icls.incidentTimelineSvc.GetByStatusType("pending")
		if err != nil {
			tx.Rollback()
			return failedIds, err
		}
		timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
			PerformedByID:    nil,
			TimelineStatusID: timelineStatus.ID,
			PerformedAt:      nil,
		}
		if err := icls.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
			tx.Rollback()
			return failedIds, err
		}

		err = icls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := icls.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.IncidentCommunicationLog{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (icls *IncidentCommunicationLogService) parseIncidentCommunicationLogDates(dto *incidentDTOs.IncidentCommunicationLogFormDTO) (*incidentDTOs.IncidentCommunicationLogParsedDates, error) {
	parsed := &incidentDTOs.IncidentCommunicationLogParsedDates{}

	var err error
	if dto.CommunicatedAt != "" {
		parsed.CommunicatedAt, err = utility.ParseDateTime(dto.CommunicatedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid communicated date"), "invalid_date", "communicated at")
		}
	}

	return parsed, nil
}
