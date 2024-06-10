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

type MeetingAgendaServiceInterface interface {
	GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingAgendaListQueryParamsDTO) (*meetingDTOs.MeetingAgendaListResponseDTO, error)
	GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingAgendaDetailDTO, error)
	GetByID(id uint) (*models.MeetingAgenda, error)
	Create(meetingID uint, currentUserID uint, meetingAgendaFormDTO *meetingDTOs.MeetingAgendaFormDTO) (uint, error)
	Update(meetingID uint, id uint, currentUserID uint, meetingAgendaFormDTO *meetingDTOs.MeetingAgendaFormDTO) error
	Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type MeetingAgendaService struct {
	repo           meetingRepositories.MeetingAgendaRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	meetingSvc     MeetingServiceInterface
}

func NewMeetingAgendaService(
	repo meetingRepositories.MeetingAgendaRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingSvc MeetingServiceInterface,
) MeetingAgendaServiceInterface {
	return &MeetingAgendaService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		meetingSvc:     meetingSvc,
	}
}

func (mas *MeetingAgendaService) GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingAgendaListQueryParamsDTO) (*meetingDTOs.MeetingAgendaListResponseDTO, error) {
	return mas.repo.GetAll(currentUserID, params)
}

func (mas MeetingAgendaService) GetByID(id uint) (*models.MeetingAgenda, error) {
	meetingAgenda, err := mas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingAgenda, nil
}

func (mas MeetingAgendaService) GetFullDetailsByID(meetingID uint, id uint) (*meetingDTOs.MeetingAgendaDetailDTO, error) {
	var meetingAgendaDTO meetingDTOs.MeetingAgendaDetailDTO
	if _, err := mas.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	meetingAgenda, err := mas.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(meetingAgenda, &meetingAgendaDTO)
	if err != nil {
		return nil, err
	}

	// Add extra fields required in details screen if any

	return &meetingAgendaDTO, nil
}

func (mas *MeetingAgendaService) Create(meetingID uint, currentUserID uint, meetingAgendaFormDTO *meetingDTOs.MeetingAgendaFormDTO) (uint, error) {
	if _, err := mas.meetingSvc.GetByID(meetingID); err != nil {
		return 0, err
	}
	meetingAgenda := &models.MeetingAgenda{
		MeetingID: meetingID,
		Title:     meetingAgendaFormDTO.Title,
		TimeFrom:  &meetingAgendaFormDTO.TimeFrom,
		TimeTo:    &meetingAgendaFormDTO.TimeTo,
		CreatedBy: currentUserID,
	}

	id, err := mas.repo.Create(meetingAgenda)
	if err != nil {
		return 0, err
	}

	go func() {
		err := mas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingAgenda{}.TableName(),
			id,
			meetingAgendaFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingAgenda{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (mas *MeetingAgendaService) Update(meetingID uint, id uint, currentUserID uint, meetingAgendaFormDTO *meetingDTOs.MeetingAgendaFormDTO) error {
	if _, err := mas.meetingSvc.GetByID(meetingID); err != nil {
		return err
	}

	existingMeetingAgenda, err := mas.GetByID(id)
	if err != nil {
		return err
	}

	existingMeetingAgenda.MeetingID = meetingID
	existingMeetingAgenda.Title = meetingAgendaFormDTO.Title
	existingMeetingAgenda.TimeFrom = &meetingAgendaFormDTO.TimeFrom
	existingMeetingAgenda.TimeTo = &meetingAgendaFormDTO.TimeTo
	existingMeetingAgenda.UpdatedBy = &currentUserID

	if err := mas.repo.Update(existingMeetingAgenda); err != nil {
		return err
	}

	go func() {
		err := mas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.MeetingAgenda{}.TableName(),
			id,
			meetingAgendaFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.MeetingAgenda{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (mas *MeetingAgendaService) Delete(meetingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := mas.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		// Fetch meeting details including participants
		_, err := mas.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete the meeting
		if err := mas.repo.Delete(id); err != nil {
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
			err := mas.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.MeetingAgenda{}.TableName(),
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
