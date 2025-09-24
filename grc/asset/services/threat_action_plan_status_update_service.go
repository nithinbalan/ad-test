package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	threatDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatActionPlanStatusUpdateServiceInterface interface {
	GetAll(ThreatID uint, ThreatActionPlanID uint, sort string, sortBy string, keyword string) ([]*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error)
	Create(ThreatID uint, ThreatActionPlanID uint, currentUserID uint, threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ThreatID uint, ThreatActionPlanID uint, id uint, currentUserID uint, threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) error
	GetLatestThreatActionPlanStatusUpdate(threatID uint, actionPlanID uint) (*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(threatID, actionPlanID, statusUpdateID uint, token string) (string, error)
	Approve(threatID uint, actionPlanID uint, currentUserID uint, ThreatActionPlanReviewUserStatusUpdateFormDTO *threatDTOs.ThreatActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(threatID uint, actionPlanID uint, currentUserID uint, ThreatActionPlanReviewUserStatusUpdateFormDTO *threatDTOs.ThreatActionPlanReviewUserStatusUpdateFormDTO) error
}

type ThreatActionPlanStatusUpdateService struct {
	repo                repositories.ThreatActionPlanStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	threatActionPlanSvc ThreatActionPlanServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	threatSvc           ThreatServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	taskStatusSvc       taskServices.TaskStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewThreatActionPlanStatusUpdateService(
	repo repositories.ThreatActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatActionPlanSvc ThreatActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	threatSvc ThreatServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ThreatActionPlanStatusUpdateServiceInterface {
	return &ThreatActionPlanStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		threatActionPlanSvc: threatActionPlanSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		threatSvc:           threatSvc,
		taskSvc:             taskSvc,
		taskStatusSvc:       taskStatusSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (tasus *ThreatActionPlanStatusUpdateService) GetAll(ThreatID uint, ThreatActionPlanID uint, sort string, sortBy string, keyword string) ([]*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return nil, err
	}

	return tasus.repo.GetAll(ThreatActionPlanID, sort, sortBy, keyword)
}

func (tasus ThreatActionPlanStatusUpdateService) GetByID(ThreatID uint, ThreatActionPlanID uint, id uint) (*models.ThreatActionPlanStatusUpdate, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return nil, err
	}

	threatActionPlanStatusUpdate, err := tasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return threatActionPlanStatusUpdate, nil
}

func (tasus *ThreatActionPlanStatusUpdateService) Create(ThreatID uint, ThreatActionPlanID uint, currentUserID uint, threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) (uint, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return 0, err
	}

	threat, err := tasus.threatSvc.GetByID(ThreatID)
	if err != nil {
		return 0, err
	}

	threatActionPlan, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := tasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	tx := tasus.repo.BeginTransaction()

	if err := tasus.validateThreatActionPlanStatusUpdateFormDTO(threatActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	threatActionPlanUpdate := &models.ThreatActionPlanStatusUpdate{
		ThreatActionPlanID: ThreatActionPlanID,
		ActionPlanStatusID: threatActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:            threatActionPlanStatusUpdateFormDTO.Comment,
		Percentage:         threatActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err := tasus.repo.Create(threatActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tasus.threatActionPlanSvc.StatusUpdate(ThreatID, ThreatActionPlanID, currentUserID, threatActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ThreatActionPlanStatusUpdateDocument

	for _, doc := range threatActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ThreatActionPlanStatusUpdateDocument{
			ThreatActionPlanStatusUpdateID: id,
			URL:                            url,
			Name:                           doc.FileName,
			Size:                           doc.Size,
			Token:                          encryptedFilename,
			OriginalName:                   doc.OriginalName,
			Type:                           doc.Type,
			CreatedByID:                    currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := tasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if threatActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && threatActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, ThreatID, ThreatActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:   "threat",
			ThreatActionPlanID: &ThreatActionPlanID,
			Title:              threatActionPlan.Title,
			Description:        threatActionPlan.Description,
			URLParam:           &urlParam,
			QueryParam:         nil,
			SeverityLevelID:    threatActionPlan.SeverityLevelID,
			ResponsibleUserID:  *threatActionPlan.ReviewUserID,
		}

		if _, err := tasus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := tasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatActionPlanStatusUpdate{}.TableName(),
			id,
			threatActionPlanUpdate,
			currentUserID,
			"asset",
			"threat",
			"threat_action_plan_status_update",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if threatActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && threatActionPlan.ReviewUserID != nil {
			urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, ThreatID, ThreatActionPlanID)
			queryParam := ""
			entity := "threat-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*threatActionPlan.ReviewUserID},
				ModuleIdentifier:     "threat",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Threat Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة العمل الخاصة بالتهديدات",
				TitleEn:              threatActionPlan.Title,
				TitleAr:              threatActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := tasus.userSvc.GetByID(*threatActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := tasus.userSvc.GetByID(threatActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Threat Action Plan %s: Resolved", threatActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn": reviewUser.FirstNameEn,
				"LastNameEn":  reviewUser.LastNameEn,
				// "ThreatReferenceNumber": threat.ReferenceNumber,
				"ThreatTitle":           threat.TitleEn,
				"ActionPlanTitle":       threatActionPlan.Title,
				"ActionPlanDescription": threatActionPlan.Description,
				"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
				"ResolvedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/asset/threats/%d/action-plans?action_plan_id=%d", appURL, ThreatID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "asset/templates/email/threat_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = tasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (tasus *ThreatActionPlanStatusUpdateService) Update(ThreatID uint, ThreatActionPlanID uint, id uint, currentUserID uint, ThreatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) error {
	tx := tasus.repo.BeginTransaction()

	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return err
	}

	oldThreatActionPlanStatusUpdate, err := tasus.GetByID(ThreatID, ThreatActionPlanID, id)
	if err != nil {
		return err
	}

	existingThreatActionPlanStatusUpdate, err := tasus.GetByID(ThreatID, ThreatActionPlanID, id)
	if err != nil {
		return err
	}

	if err := tasus.validateThreatActionPlanStatusUpdateFormDTO(ThreatActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	existingThreatActionPlanStatusUpdate.Comment = ThreatActionPlanStatusUpdateFormDTO.Comment
	existingThreatActionPlanStatusUpdate.Percentage = ThreatActionPlanStatusUpdateFormDTO.Percentage
	existingThreatActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = tasus.repo.Update(existingThreatActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ThreatActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ThreatActionPlanStatusUpdateDocument

	for _, doc := range ThreatActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			document := &models.ThreatActionPlanStatusUpdateDocument{
				ThreatActionPlanStatusUpdateID: id,
				URL:                            url,
				Name:                           doc.FileName,
				Size:                           doc.Size,
				Token:                          encryptedFileName,
				OriginalName:                   doc.OriginalName,
				Type:                           doc.Type,
				CreatedByID:                    currentUserID,
			}

			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := tasus.repo.GetDocumentByID(*doc.ID)
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
		if err := tasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := tasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlanstatusUpdate, err := tasus.GetByID(ThreatID, ThreatActionPlanID, id)
		if err != nil {
			tx.Rollback()
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldThreatActionPlanStatusUpdate.Comment,
			"percentage": oldThreatActionPlanStatusUpdate.Percentage,
		}

		newData := map[string]interface{}{
			"comment":    existingActionPlanstatusUpdate.Comment,
			"percentage": existingActionPlanstatusUpdate.Percentage,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := tasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ThreatActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset",
				"threat",
				"threat_action_plan_status_update",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ThreatActionPlanStatusUpdate{}.TableName(),
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

func (tasus *ThreatActionPlanStatusUpdateService) validateThreatActionPlanStatusUpdateFormDTO(threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) error {
	if _, err := tasus.actionPlanStatusSvc.GetByID(threatActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (tasus ThreatActionPlanStatusUpdateService) GetLatestThreatActionPlanStatusUpdate(threatID uint, actionPlanID uint) (*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return nil, err
	}

	return tasus.repo.GetLatestThreatActionPlanStatusUpdate(actionPlanID)
}

func (tasus *ThreatActionPlanStatusUpdateService) GetDocumentPath(threatID, actionPlanID, statusUpdateID uint, token string) (string, error) {
	statusUpdate, err := tasus.GetByID(threatID, actionPlanID, statusUpdateID)
	if err != nil {
		fmt.Printf("Error fetching status update by ID: %v\n", err)
		return "", err
	}

	if len(statusUpdate.Documents) == 0 {
		fmt.Println("No documents associated with the status update")
		return "", fmt.Errorf("no documents found for the threat action plan update")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		fmt.Printf("Token decryption failed: %v\n", err)
		return "", fmt.Errorf("resource not found")
	}
	decryptedFilename := string(decryptedFilenameBytes)
	fmt.Printf("Decrypted filename: %s\n", decryptedFilename)

	storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(threatID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(statusUpdateID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	fmt.Printf("Constructed file path: %s\n", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("File does not exist at path: %s\n", fullPath)
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		fmt.Printf("Error accessing file at path %s: %v\n", fullPath, err)
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (tasus *ThreatActionPlanStatusUpdateService) Approve(threatID uint, actionPlanID uint, currentUserID uint, ThreatActionPlanReviewUserStatusUpdateFormDTO *threatDTOs.ThreatActionPlanReviewUserStatusUpdateFormDTO) error {
	threat, err := tasus.threatSvc.GetByID(threatID)
	if err != nil {
		return err
	}

	existingThreatActionPlan, err := tasus.threatActionPlanSvc.GetByID(threatID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingThreatActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this Threat action plan.")
	}

	if currentUserID != *existingThreatActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	threatActionPlanStatusResolved, err := tasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	threatActionPlanClosedStatus, err := tasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	taskStatus, err := tasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		return err
	}

	tx := tasus.repo.BeginTransaction()

	if existingThreatActionPlan.ActionPlanStatusID != threatActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The compliance action plan status must be resolved.")
	}

	existingThreatActionPlan.ActionPlanStatusID = threatActionPlanClosedStatus.ID

	if err := tasus.repo.ThreatActionPlanUpdateStatus(existingThreatActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	threatActionPlanUpdate := &models.ThreatActionPlanStatusUpdate{
		ThreatActionPlanID: actionPlanID,
		ActionPlanStatusID: threatActionPlanClosedStatus.ID,
		Percentage:         nil,
		Comment:            ThreatActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := tasus.repo.Create(threatActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err = tasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tasus.taskSvc.UpdateActionPlanTaskStatus("threat_action_plan_id", actionPlanID, existingThreatActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tasus.taskSvc.UpdateActionPlanTaskStatus("threat_action_plan_id", actionPlanID, *existingThreatActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, actionPlanID)
		queryParam := ""
		entity := "threat-action-plan"
		action := "close-threat-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingThreatActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "threat",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Threat Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة العمل الخاصة بالتهديدات",
			TitleEn:              existingThreatActionPlan.Title,
			TitleAr:              existingThreatActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := tasus.userSvc.GetByID(existingThreatActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := tasus.userSvc.GetByID(*existingThreatActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Threat Action Plan %s: Resolved", existingThreatActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ThreatTitle":           threat.TitleEn,
			"ActionPlanTitle":       existingThreatActionPlan.Title,
			"ActionPlanDescription": existingThreatActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/asset/threats/%d/action-plans?action_plan_id=%d", appURL, threatID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "asset/templates/email/threat_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = tasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingThreatActionPlan, err := tasus.threatActionPlanSvc.GetByID(threatID, actionPlanID)
		if err != nil {
			return
		}

		err = tasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.ThreatActionPlanStatusUpdate{}.TableName(),
			id,
			existingThreatActionPlan,
			currentUserID,
			"asset",
			"threat",
			"threat_action_plan_status_update",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ThreatActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (tasus *ThreatActionPlanStatusUpdateService) Revert(threatID uint, actionPlanID uint, currentUserID uint, ThreatActionPlanReviewUserStatusUpdateFormDTO *threatDTOs.ThreatActionPlanReviewUserStatusUpdateFormDTO) error {
	threat, err := tasus.threatSvc.GetByID(threatID)
	if err != nil {
		return err
	}
	existingThreatActionPlan, err := tasus.threatActionPlanSvc.GetByID(threatID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingThreatActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	threatActionPlanStatusResolved, err := tasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	threatActionPlanReopenStatus, err := tasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := tasus.repo.BeginTransaction()

	if existingThreatActionPlan.ActionPlanStatusID != threatActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The threat action plan status must be resolved.")
	}

	existingThreatActionPlan.ActionPlanStatusID = threatActionPlanReopenStatus.ID

	if err := tasus.repo.ThreatActionPlanUpdateStatus(existingThreatActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	threatActionPlanUpdate := &models.ThreatActionPlanStatusUpdate{
		ThreatActionPlanID: actionPlanID,
		ActionPlanStatusID: threatActionPlanReopenStatus.ID,
		Comment:            ThreatActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := tasus.repo.Create(threatActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := tasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tasus.taskSvc.UpdateActionPlanTaskStatus("threat_action_plan_id", actionPlanID, *existingThreatActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"threat_id":%d, "action_plan_id":%d}`, threatID, actionPlanID)
		queryParam := ""
		entity := "threat-action-plan"
		action := "reopen-threat-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingThreatActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "threat",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Threat Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة العمل الخاصة بالتهديدات",
			TitleEn:              existingThreatActionPlan.Title,
			TitleAr:              existingThreatActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := tasus.userSvc.GetByID(existingThreatActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := tasus.userSvc.GetByID(*existingThreatActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Threat Action Plan %s: Reopened", existingThreatActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn": responsibleUser.FirstNameEn,
			"LastNameEn":  responsibleUser.LastNameEn,
			// "ThreatReferenceNumber":    threat.ReferenceNumber,
			"ThreatTitle":           threat.TitleEn,
			"ActionPlanTitle":       existingThreatActionPlan.Title,
			"ActionPlanDescription": existingThreatActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/asset/threats/%d/action-plans?action_plan_id=%d", appURL, threatID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "asset/templates/email/threat_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = tasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "threat")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingThreatActionPlan, err := tasus.threatActionPlanSvc.GetByID(threatID, actionPlanID)
		if err != nil {
			return
		}

		err = tasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.ThreatActionPlanStatusUpdate{}.TableName(),
			id,
			existingThreatActionPlan,
			currentUserID,
			"asset",
			"threat",
			"threat_action_plan_status_update",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ThreatActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
