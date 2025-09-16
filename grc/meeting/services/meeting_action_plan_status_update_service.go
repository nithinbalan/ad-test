package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/meeting/repositories"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingActionPlanStatusUpdateServiceInterface interface {
	GetByID(meetingID uint, meetingActionPlanID uint, id uint) (*models.MeetingActionPlanStatusUpdate, error)
	GetAll(meetingID uint, meetingActionPlanID uint, sort string, sortBy string, keyword string) ([]*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error)
	Create(meetingID uint, meetingActionPlanID uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) (uint, error)
	Update(meetingID uint, meetingActionPlanID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) error
	GetLatestMeetingActionPlanStatusUpdate(meetingID uint, actionPlanID uint) (*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(meetingID uint, meetingActionPlanID uint, id uint, token string) (string, error)
	Approve(meetingId uint, actionPlanID uint, currentUserID uint, MeetingActionPlanReviewUserStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(meetingId uint, actionPlanID uint, currentUserID uint, MeetingActionPlanReviewUserStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanReviewUserStatusUpdateFormDTO) error
}

type MeetingActionPlanStatusUpdateService struct {
	repo                   repositories.MeetingActionPlanStatusUpdateRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	meetingActionPlanSvc   MeetingActionPlanServiceInterface
	actionPlanStatusSvc    masterServices.ActionPlanStatusServiceInterface
	meetingMinuteSvc       MeetingMinuteServiceInterface
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface
	meetingSvc             MeetingServiceInterface
	taskSvc                taskServices.TaskServiceInterface
	taskStatusSvc          taskServices.TaskStatusServiceInterface
	userSvc                organizationServices.UserServiceInterface
	notificationSvc        commonServices.NotificationServiceInterface
}

func NewMeetingActionPlanStatusUpdateService(
	repo repositories.MeetingActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	meetingActionPlanSvc MeetingActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	meetingMinuteSvc MeetingMinuteServiceInterface,
	meetingMinuteStatusSvc MeetingMinuteStatusServiceInterface,
	meetingSvc MeetingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) MeetingActionPlanStatusUpdateServiceInterface {
	return &MeetingActionPlanStatusUpdateService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		meetingActionPlanSvc:   meetingActionPlanSvc,
		actionPlanStatusSvc:    actionPlanStatusSvc,
		meetingMinuteSvc:       meetingMinuteSvc,
		meetingMinuteStatusSvc: meetingMinuteStatusSvc,
		meetingSvc:             meetingSvc,
		taskSvc:                taskSvc,
		taskStatusSvc:          taskStatusSvc,
		userSvc:                userSvc,
		notificationSvc:        notificationSvc,
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
		Percentage:          &meetingActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:         currentUserID,
	}

	id, err := masus.repo.Create(meetingActionPlanUpdate, tx)
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

		url := utility.GenerateDocumentURL("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	// Check if all related action plans are completed or terminated
	completedStatuses := []string{"completed", "terminated", "close"}

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

	meeting, err := masus.meetingSvc.GetByID(meetingID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := masus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && meetingActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, meetingID, meetingActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:    "meeting",
			MeetingActionPlanID: &meetingActionPlanID,
			Title:               meetingActionPlan.Title,
			Description:         meetingActionPlan.Description,
			URLParam:            &urlParam,
			QueryParam:          nil,
			SeverityLevelID:     meetingActionPlan.SeverityLevelID,
			ResponsibleUserID:   *meetingActionPlan.ReviewUserID,
		}

		if _, err := masus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingMeetingActionPlanStatusUpdate, err := masus.GetByID(meetingID, meetingActionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch meeting action plan status update", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = masus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.MeetingActionPlanStatusUpdate{}.TableName(),
			id,
			existingMeetingActionPlanStatusUpdate,
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

		if meetingActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && meetingActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"meeting_register_id":%d, "action_plan_id":%d}`, meetingID, meetingActionPlanID)
			queryParam := ""
			entity := "meeting-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*meetingActionPlan.ReviewUserID},
				ModuleIdentifier:     "meeting",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Meeting Action Plan Status Update",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
				TitleEn:              meetingActionPlan.Title,
				TitleAr:              meetingActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := masus.userSvc.GetByID(*meetingActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := masus.userSvc.GetByID(meetingActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Meeting Action Plan %s: Resolved", meetingActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            reviewUser.FirstNameEn,
				"LastNameEn":             reviewUser.LastNameEn,
				"MeetingReferenceNumber": meeting.ReferenceNumber,
				"MeetingTitle":           meeting.Title,
				"ActionPlanTitle":        meetingActionPlan.Title,
				"ActionPlanDescription":  meetingActionPlan.Description,
				"ResolvedByFirstNameEn":  resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":   resolvedBy.LastNameEn,
				"ResolvedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/review-meeting/%d/action-plans?action_plan_id=%d", appURL, meetingID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "meeting/templates/email/meeting_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = masus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (masus *MeetingActionPlanStatusUpdateService) Update(meetingID uint, meetingActionPlanID uint, id uint, currentUserID uint, meetingActionPlanStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanStatusUpdateFormDTO) error {
	tx := masus.repo.BeginTransaction()

	if _, err := masus.meetingActionPlanSvc.GetByID(meetingID, meetingActionPlanID); err != nil {
		return err
	}

	oldActionPlanStatusUpdate, err := masus.GetByID(meetingID, meetingActionPlanID, id)
	if err != nil {
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
	existingMeetingActionPlanStatusUpdate.Percentage = &meetingActionPlanStatusUpdateFormDTO.Percentage
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

			// genereate URL
			url := utility.GenerateDocumentURL("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlanStatusUpdate, err := masus.GetByID(meetingID, meetingActionPlanID, id)
		if err != nil {
			tx.Rollback()
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldActionPlanStatusUpdate.Comment,
			"percentage": oldActionPlanStatusUpdate.Percentage,
		}

		newData := map[string]interface{}{
			"comment":    existingActionPlanStatusUpdate.Comment,
			"percentage": existingActionPlanStatusUpdate.Percentage,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := masus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.MeetingActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.MeetingActionPlanStatusUpdate{}.TableName(),
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

func (masus *MeetingActionPlanStatusUpdateService) GetDocumentPath(meetingID uint, meetingActionPlanID uint, meetingActionPlanUpdateID uint, token string) (string, error) {
	existingMeetingActionPlanStatusUpdate, err := masus.GetByID(meetingID, meetingActionPlanID, meetingActionPlanUpdateID)
	if err != nil {
		return "", err
	}

	if len(existingMeetingActionPlanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the meeting action plan update")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("meeting", "meetings", strconv.FormatUint(uint64(meetingID), 10), "action-plans", strconv.FormatUint(uint64(meetingActionPlanID), 10), "updates", strconv.FormatUint(uint64(meetingActionPlanUpdateID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (masus *MeetingActionPlanStatusUpdateService) Approve(meetingId uint, actionPlanID uint, currentUserID uint, MeetingActionPlanReviewUserStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanReviewUserStatusUpdateFormDTO) error {
	meeting, err := masus.meetingSvc.GetByID(meetingId)
	if err != nil {
		return err
	}
	existingMeetingActionPlan, err := masus.meetingActionPlanSvc.GetByID(meetingId, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingMeetingActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this Meeting action plan.")
	}

	if currentUserID != *existingMeetingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	meetingActionPlanStatusResolved, err := masus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	meetingActionPlanClosedStatus, err := masus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := masus.repo.BeginTransaction()

	if existingMeetingActionPlan.ActionPlanStatusID != meetingActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The Meeting action plan status must be resolved.")
	}

	existingMeetingActionPlan.ActionPlanStatusID = meetingActionPlanClosedStatus.ID

	if err := masus.repo.MeetingActionPlanUpdateStatus(existingMeetingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	meetingActionPlanUpdate := &models.MeetingActionPlanStatusUpdate{
		MeetingActionPlanID: actionPlanID,
		ActionPlanStatusID:  meetingActionPlanClosedStatus.ID,
		Percentage:          nil,
		Comment:             MeetingActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:         currentUserID,
	}

	id, err := masus.repo.Create(meetingActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := masus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = masus.taskSvc.UpdateActionPlanTaskStatus("meeting_action_plan_id", actionPlanID, existingMeetingActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = masus.taskSvc.UpdateActionPlanTaskStatus("meeting_action_plan_id", actionPlanID, *existingMeetingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"meeting_register_id":%d, "action_plan_id":%d}`, meetingId, actionPlanID)
		queryParam := ""
		entity := "meeting-action-plan"
		action := "close-meeting-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingMeetingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "meeting_action_plan",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Meeting Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل الاجتماع",
			TitleEn:              existingMeetingActionPlan.Title,
			TitleAr:              existingMeetingActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := masus.userSvc.GetByID(existingMeetingActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := masus.userSvc.GetByID(*existingMeetingActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Meeting Action Plan %s: Resolved", existingMeetingActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"MeetingReferenceNumber": meeting.ReferenceNumber,
			"MeetingTitle":           meeting.Title,
			"ActionPlanTitle":        existingMeetingActionPlan.Title,
			"ActionPlanDescription":  existingMeetingActionPlan.Description,
			"ResolvedByFirstNameEn":  resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":   resolvedBy.LastNameEn,
			"ResolvedOn":             time.Now().Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%s/review-meeting/%d/action-plans?action_plan_id=%d", appURL, meetingId, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "meeting/templates/email/meeting_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = masus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingMeetingActionPlan, err := masus.meetingActionPlanSvc.GetByID(meetingId, actionPlanID)
		if err != nil {
			return
		}

		err = masus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.MeetingActionPlanStatusUpdate{}.TableName(),
			id,
			existingMeetingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.MeetingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (masus *MeetingActionPlanStatusUpdateService) Revert(meetingId uint, actionPlanID uint, currentUserID uint, MeetingActionPlanReviewUserStatusUpdateFormDTO *meetingDTOs.MeetingActionPlanReviewUserStatusUpdateFormDTO) error {
	meeting, err := masus.meetingSvc.GetByID(meetingId)
	if err != nil {
		return err
	}
	existingMeetingActionPlan, err := masus.meetingActionPlanSvc.GetByID(meetingId, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingMeetingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	meetingActionPlanStatusResolved, err := masus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	meetingActionPlanReopenStatus, err := masus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := masus.repo.BeginTransaction()

	if existingMeetingActionPlan.ActionPlanStatusID != meetingActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The Meeting action plan status must be resolved.")
	}

	existingMeetingActionPlan.ActionPlanStatusID = meetingActionPlanReopenStatus.ID

	if err := masus.repo.MeetingActionPlanUpdateStatus(existingMeetingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	meetingActionPlanUpdate := &models.MeetingActionPlanStatusUpdate{
		MeetingActionPlanID: actionPlanID,
		ActionPlanStatusID:  meetingActionPlanReopenStatus.ID,
		Comment:             MeetingActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:         currentUserID,
	}

	id, err := masus.repo.Create(meetingActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := masus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = masus.taskSvc.UpdateActionPlanTaskStatus("meeting_action_plan_id", actionPlanID, *existingMeetingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"meeting_register_id":%d, "action_plan_id":%d}`, meetingId, actionPlanID)
		queryParam := ""
		entity := "meeting-action-plan"
		action := "reopen-meeting-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingMeetingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "meeting_action_plan",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Meeting Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل الاجتماع",
			TitleEn:              existingMeetingActionPlan.Title,
			TitleAr:              existingMeetingActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := masus.userSvc.GetByID(existingMeetingActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := masus.userSvc.GetByID(*existingMeetingActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Meeting Action Plan %s: Reopened", existingMeetingActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"MeetingReferenceNumber": meeting.ReferenceNumber,
			"MeetingTitle":           meeting.Title,
			"ActionPlanTitle":        existingMeetingActionPlan.Title,
			"ActionPlanDescription":  existingMeetingActionPlan.Description,
			"ResolvedByFirstNameEn":  resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":   resolvedBy.LastNameEn,
			"ResolvedOn":             time.Now().Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%s/review-meeting/%d/action-plans?action_plan_id=%d", appURL, meetingId, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "meeting/templates/email/meeting_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = masus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "meeting_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingMeetingActionPlan, err := masus.meetingActionPlanSvc.GetByID(meetingId, actionPlanID)
		if err != nil {
			return
		}

		err = masus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.MeetingActionPlanStatusUpdate{}.TableName(),
			id,
			existingMeetingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.MeetingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
