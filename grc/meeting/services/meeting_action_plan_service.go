package services

import (
	"errors"
	"fmt"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingActionPlanServiceInterface interface {
	GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error)
	GetByID(meetingID uint, id uint) (*models.MeetingActionPlan, error)
	GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingActionPlanDetailDTO, error)
	Create(meetingID uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) (uint, error)
	Update(meetingID uint, id uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error
	Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(meetingID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type MeetingActionPlanService struct {
	repo                meetingRepositories.MeetingActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	meetingSvc          MeetingServiceInterface
	severityLevelSvc    masterServices.SeverityLevelServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	meetingMinuteSvc    MeetingMinuteServiceInterface
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface
}

func NewMeetingActionPlanService(
	repo meetingRepositories.MeetingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingSvc MeetingServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	meetingMinuteSvc MeetingMinuteServiceInterface,
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface) MeetingActionPlanServiceInterface {
	return &MeetingActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		meetingSvc:          meetingSvc,
		severityLevelSvc:    severityLevelSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
		meetingMinuteSvc:    meetingMinuteSvc,
		meetingMinuteStatusSvc:   meetingMinuteStatusSvc,
	}
}

func (maps *MeetingActionPlanService) GetAll(meetingID uint, params *meetingDTOs.MeetingActionPlanListQueryParamsDTO) (*meetingDTOs.MeetingActionPlanListResponseDTO, error) {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return maps.repo.GetAll(meetingID, params)
}

func (maps MeetingActionPlanService) GetByID(meetingID uint, id uint) (*models.MeetingActionPlan, error) {
	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	meetingActionPlan, err := maps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingActionPlan, nil
}

func (maps MeetingActionPlanService) GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingActionPlanDetailDTO, error) {
	var meetingActionPlanDTO meetingDTOs.MeetingActionPlanDetailDTO

	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	meetingActionPlan, err := maps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meetingActionPlan, &meetingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &meetingActionPlanDTO, nil
}

func (maps *MeetingActionPlanService) Create(meetingID uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) (uint, error) {
	tx := maps.repo.BeginTransaction()

	_, err := maps.meetingSvc.GetByID(meetingID)
	if err != nil {
		return 0, err
	}

	if err := maps.validateMeetingActionPlanFormDTO(MeetingActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := maps.parseMeetingActionPlanDates(MeetingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	meetingActionPlan := &models.MeetingActionPlan{
		MeetingID:          meetingID,
		Title:              MeetingActionPlanFormDTO.Title,
		Description:        MeetingActionPlanFormDTO.Description,
		SeverityLevelID:    MeetingActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:  MeetingActionPlanFormDTO.ResponsibleUserID,
		MeetingMinuteID:    MeetingActionPlanFormDTO.MeetingMinuteID,
		TargetDate:         dates.TargetDate,
		Percentage:         *MeetingActionPlanFormDTO.Percentage,
		ActionPlanStatusID: MeetingActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:        currentUserID,
	}

	id, err := maps.repo.Create(meetingActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	meetingMinuteStatus, err := maps.meetingMinuteStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}
	if err := maps.meetingMinuteSvc.StatusUpdate(*MeetingActionPlanFormDTO.MeetingMinuteID,currentUserID,uint(meetingMinuteStatus.ID), tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	actionPlanStatusID := MeetingActionPlanFormDTO.ActionPlanStatusID
	meetingActionPlanUpdate := &models.MeetingActionPlanStatusUpdate{
		MeetingActionPlanID: id,
		ActionPlanStatusID:  actionPlanStatusID,
		Comment:             "Initial action plan status",
		CreatedByID:         currentUserID,
	}

	_, err = maps.repo.CreateInitialMeetingActionPlanStatusUpdate(meetingActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := maps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingActionPlan{}.TableName(),
			id,
			MeetingActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (maps *MeetingActionPlanService) Update(meetingID uint, id uint, currentUserID uint, MeetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error {
	tx := maps.repo.BeginTransaction()

	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	existingMeetingActionPlan, err := maps.GetByID(meetingID, id)

	if err != nil {
		return err
	}

	if err := maps.validateMeetingActionPlanFormDTO(MeetingActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := maps.parseMeetingActionPlanDates(MeetingActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Meeting ActionPlan
	existingMeetingActionPlan.Title = MeetingActionPlanFormDTO.Title
	existingMeetingActionPlan.Description = MeetingActionPlanFormDTO.Description
	existingMeetingActionPlan.SeverityLevelID = MeetingActionPlanFormDTO.SeverityLevelID
	existingMeetingActionPlan.ResponsibleUserID = MeetingActionPlanFormDTO.ResponsibleUserID
	existingMeetingActionPlan.TargetDate = dates.TargetDate
	existingMeetingActionPlan.Percentage = *MeetingActionPlanFormDTO.Percentage
	existingMeetingActionPlan.ActionPlanStatusID = MeetingActionPlanFormDTO.ActionPlanStatusID
	existingMeetingActionPlan.MeetingMinuteID = MeetingActionPlanFormDTO.MeetingMinuteID
	existingMeetingActionPlan.UpdatedByID = &currentUserID

	err = maps.repo.Update(existingMeetingActionPlan)

	if err != nil {
		return err
	}
	tx.Commit()

	return nil
}

func (maps *MeetingActionPlanService) Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := maps.repo.BeginTransaction()

	if _, err := maps.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		meetingActionPlan, err := maps.GetFullDetailsByID(meetingID, id)

		if err != nil {
			continue
		}
		if len(meetingActionPlan.ActionPlanStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range meetingActionPlan.ActionPlanStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := maps.repo.DeleteMeetingActionPlanStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		err = maps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		tx.Commit()

		go func() {
			err := maps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.MeetingActionPlan{}.TableName(),
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
	} else {
		tx.Rollback()

	}

	return failedIds, nil
}

func (maps *MeetingActionPlanService) StatusUpdate(meetingID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO, tx ...utility.DatabaseTransaction) error {
	existingMeetingActionPlan, err := maps.GetByID(meetingID, id)

	if err != nil {
		return err
	}

	existingMeetingActionPlan.ActionPlanStatusID = meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingMeetingActionPlan.Percentage = meetingActionPlanStatusUpdateFormDTO.Percentage
	existingMeetingActionPlan.Comments = &meetingActionPlanStatusUpdateFormDTO.Comment

	existingMeetingActionPlan.UpdatedByID = &currentUserID

	err = maps.repo.Update(existingMeetingActionPlan, tx...)

	if err != nil {
		return err
	}

	return nil
}

func (maps *MeetingActionPlanService) validateMeetingActionPlanFormDTO(meetingActionPlanFormDTO *meetingDTOs.MeetingActionPlanFormDTO) error {
	if meetingActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := maps.severityLevelSvc.GetByID(*meetingActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if _, err := maps.userSvc.GetByID(uint(meetingActionPlanFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := maps.actionPlanStatusSvc.GetByID(meetingActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}
	if _, err := maps.meetingMinuteSvc.GetByID(*meetingActionPlanFormDTO.MeetingMinuteID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting minute does not exist"), "not_exist", "meeting minute")
	}

	return nil
}

func (maps *MeetingActionPlanService) parseMeetingActionPlanDates(dto *meetingDTOs.MeetingActionPlanFormDTO) (*meetingDTOs.MeetingActionPlanParsedDates, error) {
	parsed := &meetingDTOs.MeetingActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
