package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type MeetingActionPlanStatusUpdateServiceInterface interface {
	GetByID(meetingID uint, meetingActionPlanID uint, id uint) (*models.MeetingActionPlanStatusUpdate, error)
	GetAll(meetingID uint, meetingActionPlanID uint, sort string, sortBy string, keyword string) ([]*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error)
	Create(meetingID uint, meetingActionPlanID uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) (uint, error)
	Update(meetingID uint, meetingActionPlanID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) error
	GetLatestMeetingActionPlanStatusUpdate(meetingID uint, actionPlanID uint) (*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error)
}

type MeetingActionPlanStatusUpdateService struct {
	repo                   repositories.MeetingActionPlanStatusUpdateRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	meetingActionPlanSvc   MeetingActionPlanServiceInterface
	actionPlanStatusSvc    masterServices.ActionPlanStatusServiceInterface
	meetingMinuteSvc       MeetingMinuteServiceInterface
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface
}

func NewMeetingActionPlanStatusUpdateService(
	repo repositories.MeetingActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingActionPlanSvc MeetingActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	meetingMinuteSvc MeetingMinuteServiceInterface,
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface,
) MeetingActionPlanStatusUpdateServiceInterface {
	return &MeetingActionPlanStatusUpdateService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		meetingActionPlanSvc:   meetingActionPlanSvc,
		actionPlanStatusSvc:    actionPlanStatusSvc,
		meetingMinuteSvc:       meetingMinuteSvc,
		meetingMinuteStatusSvc: meetingMinuteStatusSvc,
	}
}

func (masus *MeetingActionPlanStatusUpdateService) GetAll(meetingID uint, meetingActionPlanID uint, sort string, sortBy string, keyword string) ([]*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error) {
	if _, err := masus.meetingActionPlanSvc.GetByID(meetingID, meetingActionPlanID); err != nil {
		return nil, err
	}

	return masus.repo.GetAll(meetingActionPlanID, sort, sortBy, keyword)
}

func (masus MeetingActionPlanStatusUpdateService) GetByID(meetingID uint, meetingActionPlanID uint, id uint) (*models.MeetingActionPlanStatusUpdate, error) {
	if _, err := masus.meetingActionPlanSvc.GetByID(meetingID, meetingActionPlanID); err != nil {
		return nil, err
	}

	meetingActionPlanStatusUpdate, err := masus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return meetingActionPlanStatusUpdate, nil
}

func (masus *MeetingActionPlanStatusUpdateService) Create(meetingID uint, meetingActionPlanID uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := masus.repo.BeginTransaction()

	meetingActionPlan, err := masus.meetingActionPlanSvc.GetByID(meetingID, meetingActionPlanID)
	if err != nil {
		return 0, err
	}

	if err := masus.validateMeetingActionPlanStatusUpdateFormDTO(meetingActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	meetingActionPlanUpdate := &models.MeetingActionPlanStatusUpdate{
		MeetingActionPlanID: meetingActionPlanID,
		ActionPlanStatusID:  meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:             meetingActionPlanStatusUpdateFormDTO.Comment,
		Percentage:          meetingActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:         currentUserID,
	}

	id, err := masus.repo.Create(meetingActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = masus.meetingActionPlanSvc.StatusUpdate(meetingID, meetingActionPlanID, currentUserID, meetingActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.MeetingActionPlanStatusUpdateDocument

	for _, doc := range meetingActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("meeting", "meeting", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), encryptedFileName)

		document := &models.MeetingActionPlanStatusUpdateDocument{
			MeetingActionPlanStatusUpdateID: id,
			URL:                             url,
			Name:                            doc.FileName,
			Size:                            doc.Size,
			Token:                           encryptedFileName,
			OriginalName:                    doc.OriginalName,
			Type:                            doc.Type,
			CreatedByID:                     currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := masus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("meeting", "meeting", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	// Check if all related action plans are completed or terminated
	completedStatuses := []string{"completed", "terminated", "on_hold"}

	meetingMinuteID := meetingActionPlan.MeetingMinuteID
	meetingActionPlanStatus, err := masus.actionPlanStatusSvc.GetByID(meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID)
	if err != nil {
		return 0, err
	}
	// Check if meetingActionPlanStatus.Type exists in completedStatuses array
	statusExists := false
	for _, status := range completedStatuses {
		if meetingActionPlanStatus.Type == status {
			statusExists = true
			break
		}
	}

	if statusExists {
		// Get the count of completed/terminated action plans excluding the current one
		allMeetingMinutesCompleted, err := masus.repo.CheckAllActionPlansCompletedStatusCount(*meetingMinuteID, meetingActionPlanID, completedStatuses, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// If the total count equals the number of action plans, update the meeting minute status to "closed"
		if allMeetingMinutesCompleted == 0 {
			meetingMinuteStatus, err := masus.meetingMinuteStatusSvc.GetByType("closed")
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			if err := masus.meetingMinuteSvc.StatusUpdate(*meetingMinuteID, currentUserID, uint(meetingMinuteStatus.ID), tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := masus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingActionPlanStatusUpdate{}.TableName(),
			id,
			meetingActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.MeetingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (masus *MeetingActionPlanStatusUpdateService) Update(meetingID uint, meetingActionPlanID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) error {
	tx := masus.repo.BeginTransaction()

	if _, err := masus.meetingActionPlanSvc.GetByID(meetingID, meetingActionPlanID); err != nil {
		return err
	}

	existingMeetingActionPlanStatusUpdate, err := masus.GetByID(meetingID, meetingActionPlanID, id)

	if err != nil {
		return err
	}

	if err := masus.validateMeetingActionPlanStatusUpdateFormDTO(meetingActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingMeetingActionPlanStatusUpdate.Comment = meetingActionPlanStatusUpdateFormDTO.Comment
	existingMeetingActionPlanStatusUpdate.Percentage = meetingActionPlanStatusUpdateFormDTO.Percentage
	existingMeetingActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = masus.repo.Update(existingMeetingActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = masus.meetingActionPlanSvc.StatusUpdate(meetingID, meetingActionPlanID, currentUserID, meetingActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.MeetingActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.MeetingActionPlanStatusUpdateDocument

	for _, doc := range meetingActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("meeting", "meeting", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), encryptedFileName)

			document := &models.MeetingActionPlanStatusUpdateDocument{
				MeetingActionPlanStatusUpdateID: id,
				URL:                             url,
				Name:                            doc.FileName,
				Size:                            doc.Size,
				Token:                           encryptedFileName,
				OriginalName:                    doc.OriginalName,
				Type:                            doc.Type,
				CreatedByID:                     currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := masus.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := masus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := masus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("meeting", "meeting", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("meeting", "meeting", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (masus *MeetingActionPlanStatusUpdateService) validateMeetingActionPlanStatusUpdateFormDTO(meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) error {

	if _, err := masus.actionPlanStatusSvc.GetByID(meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (masus MeetingActionPlanStatusUpdateService) GetLatestMeetingActionPlanStatusUpdate(meetingID uint, actionPlanID uint) (*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error) {
	if _, err := masus.meetingActionPlanSvc.GetByID(meetingID, actionPlanID); err != nil {
		return nil, err
	}

	return masus.repo.GetLatestMeetingActionPlanStatusUpdate(actionPlanID)
}
