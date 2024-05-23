package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingServiceInterface interface {
	GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*meetingDTOs.MeetingDetailDTO, error)
	GetByID(id uint) (*models.Meeting, error)
	Create(currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type MeetingService struct {
	repo                meetingRepositories.MeetingRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	meetingTypeSvc      MeetingTypeServiceInterface
	meetingFrequencySvc MeetingFrequencyServiceInterface
	userSvc             organizationServices.UserServiceInterface
	meetingStatusSvc    MeetingStatusServiceInterface
}

func NewMeetingService(
	repo meetingRepositories.MeetingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingTypeSvc MeetingTypeServiceInterface,
	meetingFrequencySvc MeetingFrequencyServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	meetingStatusSvc MeetingStatusServiceInterface,
) MeetingServiceInterface {
	return &MeetingService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		meetingTypeSvc:      meetingTypeSvc,
		meetingFrequencySvc: meetingFrequencySvc,
		userSvc:             userSvc,
		meetingStatusSvc:    meetingStatusSvc,
	}
}

func (ms *MeetingService) GetAll(currentUserID uint, params *meetingDTOs.MeetingListQueryParamsDTO) (*meetingDTOs.MeetingListResponseDTO, error) {
	return ms.repo.GetAll(currentUserID, params)
}

func (ms MeetingService) GetByID(id uint) (*models.Meeting, error) {
	meeting, err := ms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meeting, nil
}

func (ms MeetingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*meetingDTOs.MeetingDetailDTO, error) {
	var meetingDTO meetingDTOs.MeetingDetailDTO

	meeting, err := ms.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meeting, &meetingDTO)
	if err != nil {
		return nil, err
	}

	// Add extra fields required in details screen if any

	return &meetingDTO, nil
}

func (ms *MeetingService) Create(currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) (uint, error) {
	tx := ms.repo.BeginTransaction()

	if err := ms.validateMeetingFormDTO(meetingFormDTO); err != nil {
		return 0, err
	}

	dates, err := ms.parseMeetingDates(meetingFormDTO)
	if err != nil {
		return 0, err
	}

	meeting := &models.Meeting{
		Title:              meetingFormDTO.Title,
		Location:           meetingFormDTO.Location,
		Date:               dates.Date,
		StartTime:          dates.StartTime,
		EndTime:            dates.EndTime,
		OrganizerID:        meetingFormDTO.OrganizerID,
		IsPublic:           meetingFormDTO.IsPublic,
		CreatedBy:          currentUserID,
		MeetingTypeID:      meetingFormDTO.MeetingTypeID,
		MeetingFrequencyID: meetingFormDTO.MeetingFrequencyID,
		MeetingID:          meetingFormDTO.MeetingID,
		MeetingStatusID:    meetingFormDTO.MeetingStatusID,
	}

	id, err := ms.repo.Create(meeting, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var participants []*models.MeetingParticipant
	for _, participantID := range meetingFormDTO.UserIDs {
		participant := &models.MeetingParticipant{
			MeetingID: id,
			UserID:    *participantID,
			CreatedBy: currentUserID,
		}
		participants = append(participants, participant)
	}

	if len(participants) > 0 {
		if err := ms.repo.CreateParticipants(participants, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := ms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Meeting{}.TableName(),
			id,
			meetingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Meeting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ms *MeetingService) Update(id uint, currentUserID uint, meetingFormDTO *meetingDTOs.MeetingFormDTO) error {
	tx := ms.repo.BeginTransaction()

	if err := ms.validateMeetingFormDTO(meetingFormDTO); err != nil {
		return err
	}

	existingMeeting, err := ms.GetByID(id)
	if err != nil {
		return err
	}

	dates, err := ms.parseMeetingDates(meetingFormDTO)
	if err != nil {
		return err
	}

	existingMeeting.MeetingID = meetingFormDTO.MeetingID
	existingMeeting.Title = meetingFormDTO.Title
	existingMeeting.Location = meetingFormDTO.Location
	existingMeeting.Date = dates.Date
	existingMeeting.StartTime = dates.StartTime
	existingMeeting.EndTime = dates.EndTime
	existingMeeting.OrganizerID = meetingFormDTO.OrganizerID
	existingMeeting.IsPublic = meetingFormDTO.IsPublic
	existingMeeting.UpdatedBy = &currentUserID
	existingMeeting.MeetingTypeID = meetingFormDTO.MeetingTypeID
	existingMeeting.MeetingFrequencyID = meetingFormDTO.MeetingFrequencyID
	existingMeeting.MeetingStatusID = meetingFormDTO.MeetingStatusID

	if err := ms.repo.Update(existingMeeting, tx); err != nil {
		tx.Rollback()
		return err
	}

	meeting, err := ms.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingParticipantIDs := make(map[uint]bool)
	for _, participant := range meeting.Participants {
		existingParticipantIDs[participant.UserID] = true
	}

	newParticipantIDs := make(map[uint]bool)
	for _, userID := range meetingFormDTO.UserIDs {
		newParticipantIDs[*userID] = true
	}

	var participantsToAdd []*models.MeetingParticipant
	for _, userID := range meetingFormDTO.UserIDs {
		if _, exists := existingParticipantIDs[*userID]; !exists {
			participant := &models.MeetingParticipant{
				MeetingID: id,
				UserID:    *userID,
				CreatedBy: currentUserID,
			}
			participantsToAdd = append(participantsToAdd, participant)
		}
	}

	if len(participantsToAdd) > 0 {
		if err := ms.repo.CreateParticipants(participantsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var participantIDsToDelete []uint
	for _, existingParticipant := range meeting.Participants {
		if _, exists := newParticipantIDs[existingParticipant.UserID]; !exists {
			participantIDsToDelete = append(participantIDsToDelete, existingParticipant.ID)
		}
	}

	// Delete participants not in the new list
	if len(participantIDsToDelete) > 0 {
		if err := ms.repo.DeleteParticipants(participantIDsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := ms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Meeting{}.TableName(),
			id,
			meetingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.Meeting{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ms *MeetingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ms.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		// Fetch meeting details including participants
		meeting, err := ms.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete participants associated with the meeting
		if len(meeting.Participants) > 0 {
			var deleteParticipantIDs []uint
			for _, participant := range meeting.Participants {
				deleteParticipantIDs = append(deleteParticipantIDs, participant.ID)
			}
			if err := ms.repo.DeleteParticipants(deleteParticipantIDs, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Delete the meeting
		if err := ms.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := ms.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Meeting{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ms *MeetingService) validateMeetingFormDTO(meetingFormDTO *meetingDTOs.MeetingFormDTO) error {
	if _, err := ms.userSvc.GetByID(meetingFormDTO.OrganizerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("organizer does not exist"), "not_exist", "organizer")
	}

	if _, err := ms.meetingTypeSvc.GetByID(meetingFormDTO.MeetingTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting type does not exist"), "not_exist", "meeting type")
	}

	if _, err := ms.meetingStatusSvc.GetByID(meetingFormDTO.MeetingStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting status does not exist"), "not_exist", "meeting status")
	}

	if meetingFormDTO.MeetingFrequencyID != nil {
		if _, err := ms.meetingFrequencySvc.GetByID(*meetingFormDTO.MeetingFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("meeting frequency does not exist"), "not_exist", "meeting frequency")
		}
	}

	if meetingFormDTO.MeetingID != nil {
		if _, err := ms.GetByID(*meetingFormDTO.MeetingID); err != nil {
			return utility.NewValidationError(fmt.Errorf("meeting does not exist"), "not_exist", "meeting")
		}
	}

	return nil
}

func (ms *MeetingService) parseMeetingDates(dto *meetingDTOs.MeetingFormDTO) (*meetingDTOs.MeetingParsedDates, error) {
	parsed := &meetingDTOs.MeetingParsedDates{}

	var err error

	parsed.Date, err = utility.ParseMandatoryDate(dto.Date)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid date"), "invalid_date", "date")
	}

	parsed.StartTime, err = utility.ParseMandatoryTime(dto.StartTime)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start time"), "invalid_time", "start time")
	}

	parsed.EndTime, err = utility.ParseMandatoryTime(dto.EndTime)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end time"), "invalid_time", "end time")
	}

	return parsed, nil
}
