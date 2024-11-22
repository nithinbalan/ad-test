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
	GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingAgendaListQueryParamsDTO) ([]*meetingDTOs.MeetingAgendaListDTO, error)
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

func (mas *MeetingAgendaService) GetAll(meetingID uint, currentUserID uint, params *meetingDTOs.MeetingAgendaListQueryParamsDTO) ([]*meetingDTOs.MeetingAgendaListDTO, error) {
	if _, err := mas.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mas.repo.GetAll(meetingID, currentUserID, params)
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
		MeetingID:   meetingID,
		Title:       meetingAgendaFormDTO.Title,
		CreatedBYID: currentUserID,
	}

	id, err := mas.repo.Create(meetingAgenda)
	if err != nil {
		return 0, err
	}

	go func() {
		existingMeetingAgenda, err := mas.GetFullDetailsByID(meetingID, id)
		if err != nil {
			utility.LogError("Failed to fetch meeting minute", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = mas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingAgenda{}.TableName(),
			id,
			existingMeetingAgenda,
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

	oldMeetingAgenda, err := mas.GetFullDetailsByID(meetingID, id)
	if err != nil {
		return err
	}

	existingMeetingAgenda, err := mas.GetByID(id)
	if err != nil {
		return err
	}

	existingMeetingAgenda.MeetingID = meetingID
	existingMeetingAgenda.Title = meetingAgendaFormDTO.Title
	existingMeetingAgenda.UpdatedBYID = &currentUserID

	if err := mas.repo.Update(existingMeetingAgenda); err != nil {
		return err
	}

	go func() {
		existingAgenda, err := mas.GetFullDetailsByID(meetingID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title": oldMeetingAgenda.Title,
		}

		newData := map[string]interface{}{
			"title": existingAgenda.Title,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {

			err := mas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.MeetingAgenda{}.TableName(),
				id,
				changedData,
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
	var deletedMeetingAgenda []map[string]interface{}

	for _, id := range ids {
		// Fetch meeting details including participants
		meetingAgenda, err := mas.repo.GetFullDetailsByID(id)
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
		deletedMeetingAgenda = append(deletedMeetingAgenda, map[string]interface{}{
			"id":         id,
			"title":      meetingAgenda.Title,
			"meeting_id": meetingID,
		})
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
				deletedMeetingAgenda,
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
