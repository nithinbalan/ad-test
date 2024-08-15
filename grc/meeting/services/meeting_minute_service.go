package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type MeetingMinuteServiceInterface interface {
	GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingMinuteListQueryParamsDTO) (*meetingDTOs.MeetingMinuteListResponseDTO, error)
	GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingMinuteDetailDTO, error)
	GetByID(id uint) (*models.MeetingMinute, error)
	Create(meetingID uint, currentUserID uint, meetingMinuteFormDTO *meetingDTOs.MeetingMinuteFormDTO) (uint, error)
	Update(meetingID uint, id uint, currentUserID uint, meetingMinuteFormDTO *meetingDTOs.MeetingMinuteFormDTO) error
	Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type MeetingMinuteService struct {
	repo             meetingRepositories.MeetingMinuteRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	meetingSvc       MeetingServiceInterface
	meetingAgendaSvc MeetingAgendaServiceInterface
}

func NewMeetingMinuteService(
	repo meetingRepositories.MeetingMinuteRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingSvc MeetingServiceInterface,
	meetingAgendaSvc MeetingAgendaServiceInterface,
) MeetingMinuteServiceInterface {
	return &MeetingMinuteService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		meetingSvc:       meetingSvc,
		meetingAgendaSvc: meetingAgendaSvc,
	}
}

func (mms *MeetingMinuteService) GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingMinuteListQueryParamsDTO) (*meetingDTOs.MeetingMinuteListResponseDTO, error) {
	return mms.repo.GetAll(currentUserID, params)
}

func (mms MeetingMinuteService) GetByID(id uint) (*models.MeetingMinute, error) {
	meetingMinute, err := mms.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingMinute, nil
}

func (mms MeetingMinuteService) GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingMinuteDetailDTO, error) {
	var meetingMinuteDTO meetingDTOs.MeetingMinuteDetailDTO
	if _, err := mms.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	meetingMinute, err := mms.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meetingMinute, &meetingMinuteDTO)
	if err != nil {
		return nil, err
	}

	// Add extra fields required in details screen if any

	return &meetingMinuteDTO, nil
}

func (mms *MeetingMinuteService) Create(meetingID uint, currentUserID uint, meetingMinuteFormDTO *meetingDTOs.MeetingMinuteFormDTO) (uint, error) {
	if _, err := mms.meetingSvc.GetByID(meetingID); err != nil {
		return 0, err
	}
	if err := mms.validateMeetingMinuteFormDTO(meetingMinuteFormDTO); err != nil {
		return 0, err
	}

	meetingMinute := &models.MeetingMinute{
		MeetingAgendaID: meetingMinuteFormDTO.MeetingAgendaID,
		Minutes:         &meetingMinuteFormDTO.Minutes,
		CreatedBYID:     currentUserID,
	}

	id, err := mms.repo.Create(meetingMinute)
	if err != nil {
		return 0, err
	}

	go func() {
		err := mms.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingMinute{}.TableName(),
			id,
			meetingMinuteFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingMinute{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (mms *MeetingMinuteService) Update(meetingID uint, id uint, currentUserID uint, meetingMinuteFormDTO *meetingDTOs.MeetingMinuteFormDTO) error {
	if _, err := mms.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}
	if err := mms.validateMeetingMinuteFormDTO(meetingMinuteFormDTO); err != nil {
		return err
	}
	existingMeetingMinute, err := mms.GetByID(id)
	if err != nil {
		return err
	}

	existingMeetingMinute.MeetingAgendaID = meetingMinuteFormDTO.MeetingAgendaID
	existingMeetingMinute.Minutes = &meetingMinuteFormDTO.Minutes
	existingMeetingMinute.UpdatedBYID = &currentUserID

	if err := mms.repo.Update(existingMeetingMinute); err != nil {
		return err
	}

	go func() {
		err := mms.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.MeetingMinute{}.TableName(),
			id,
			meetingMinuteFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.MeetingMinute{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (mms *MeetingMinuteService) Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := mms.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		// Fetch meeting details including participants
		_, err := mms.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete the meeting
		if err := mms.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		go func() {
			err := mms.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.MeetingMinute{}.TableName(),
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

func (mms *MeetingMinuteService) validateMeetingMinuteFormDTO(meetingMinuteFormDTO *meetingDTOs.MeetingMinuteFormDTO) error {
	if _, err := mms.meetingAgendaSvc.GetByID(*meetingMinuteFormDTO.MeetingAgendaID); err != nil {
		return utility.NewValidationError(fmt.Errorf("meeting agenda does not exist"), "not_exist", "meeting agenda")
	}

	return nil
}
