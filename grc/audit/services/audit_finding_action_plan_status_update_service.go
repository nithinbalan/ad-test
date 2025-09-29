package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
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

type AuditFindingActionPlanStatusUpdateServiceInterface interface {
	GetAll(auditFindingID uint, auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	Create(auditFindingID uint, auditFindingActionPlanID uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) (uint, error)
	Update(auditFindingID uint, auditFindingActionPlanID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error
	GetLatestAuditFindingActionPlanStatusUpdate(auditFindingID uint, actionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(auditFindingID uint, auditFindingActionPlanID uint, id uint, token string) (string, error)
	Approve(AuditFindingID uint, AuditFindingActionPlanID uint, currentUserID uint, AuditFindingActionPlanReviewUserStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(AuditFindingID uint, AuditFindingActionPlanID uint, currentUserID uint, AuditFindingActionPlanReviewUserStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanReviewUserStatusUpdateFormDTO) error
}

type AuditFindingActionPlanStatusUpdateService struct {
	repo                      repositories.AuditFindingActionPlanStatusUpdateRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface
	actionPlanStatusSvc       masterServices.ActionPlanStatusServiceInterface
	auditFindingSvc           AuditFindingServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	taskSvc                   taskServices.TaskServiceInterface
	taskStatusSvc             taskServices.TaskStatusServiceInterface
}

func NewAuditFindingActionPlanStatusUpdateService(
	repo repositories.AuditFindingActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) AuditFindingActionPlanStatusUpdateServiceInterface {
	return &AuditFindingActionPlanStatusUpdateService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditFindingActionPlanSvc: auditFindingActionPlanSvc,
		actionPlanStatusSvc:       actionPlanStatusSvc,
		auditFindingSvc:           auditFindingSvc,
		userSvc:                   userSvc,
		notificationSvc:           notificationSvc,
		taskSvc:                   taskSvc,
		taskStatusSvc:             taskStatusSvc,
	}
}

func (afs *AuditFindingActionPlanStatusUpdateService) GetAll(auditFindingID uint, auditFindingActionPlanID uint, sort string, sortBy string, keyword string) ([]*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return nil, err
	}

	return afs.repo.GetAll(auditFindingActionPlanID, sort, sortBy, keyword)
}

func (afs AuditFindingActionPlanStatusUpdateService) GetByID(auditFindingID uint, auditFindingActionPlanID uint, id uint) (*models.AuditFindingActionPlanStatusUpdate, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return nil, err
	}

	auditFindingActionPlanStatusUpdate, err := afs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFindingActionPlanStatusUpdate, nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) Create(auditFindingID uint, auditFindingActionPlanID uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) (uint, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return 0, err
	}

	auditFinding, err := afs.auditFindingSvc.GetByID(auditFindingID)
	if err != nil {
		return 0, err
	}

	auditFindingActionPlan, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := afs.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	tx := afs.repo.BeginTransaction()
	if err := afs.validateAuditFindingActionPlanStatusUpdateFormDTO(auditFindingActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	auditFindingActionPlanUpdate := &models.AuditFindingActionPlanStatusUpdate{
		AuditFindingActionPlanID: auditFindingActionPlanID,
		ActionPlanStatusID:       auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                  auditFindingActionPlanStatusUpdateFormDTO.Comment,
		Percentage:               auditFindingActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:              currentUserID,
	}

	id, err := afs.repo.Create(auditFindingActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = afs.auditFindingActionPlanSvc.StatusUpdate(auditFindingID, auditFindingActionPlanID, currentUserID, auditFindingActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.AuditFindingActionPlanStatusUpdateDocument

	for _, doc := range auditFindingActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.AuditFindingActionPlanStatusUpdateDocument{
			AuditFindingActionPlanStatusUpdateID: id,
			URL:                                  url,
			Name:                                 doc.FileName,
			Size:                                 doc.Size,
			Token:                                encryptedFilename,
			OriginalName:                         doc.OriginalName,
			Type:                                 doc.Type,
			CreatedByID:                          currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := afs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	if auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && auditFindingActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"audit_id":%d, "action_plan_id":%d}`, auditFindingID, auditFindingActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:         "audit",
			AuditFindingActionPlanID: &auditFindingActionPlanID,
			Title:                    auditFindingActionPlan.Title,
			Description:              auditFindingActionPlan.Description,
			URLParam:                 &urlParam,
			QueryParam:               nil,
			SeverityLevelID:          auditFindingActionPlan.SeverityLevelID,
			ResponsibleUserID:        *auditFindingActionPlan.ReviewUserID,
		}

		if _, err := afs.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingID, auditFindingActionPlanID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit finding action plan status", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = afs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingActionPlanStatusUpdate{}.TableName(),
			id,
			existingAuditFindingActionPlanStatusUpdate,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && auditFindingActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"audit_id":%d, "action_plan_id":%d}`, auditFindingID, auditFindingActionPlanID)
			queryParam := ""
			entity := "audit-finding-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*auditFindingActionPlan.ReviewUserID},
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Audit Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة عمل التدقيق",
				TitleEn:              auditFindingActionPlan.Title,
				TitleAr:              auditFindingActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := afs.userSvc.GetByID(*auditFindingActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := afs.userSvc.GetByID(auditFindingActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Audit Finding Action Plan %s: Resolved", auditFindingActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":                 reviewUser.FirstNameEn,
				"LastNameEn":                  reviewUser.LastNameEn,
				"AuditFindingReferenceNumber": auditFinding.ReferenceNumber,
				"AuditFindingTitle":           auditFinding.Title,
				"ActionPlanTitle":             auditFindingActionPlan.Title,
				"ActionPlanDescription":       auditFindingActionPlan.Description,
				"ResolvedByFirstNameEn":       resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":        resolvedBy.LastNameEn,
				"ResolvedOn":                  time.Now().Format("2006-01-02"),
				"RedirectURL":                 fmt.Sprintf("%s/audit/audit-findings/%d/action-plans?action_plan_id=%d", appURL, auditFindingID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "audit/templates/email/audit_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = afs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) Update(auditFindingID uint, auditFindingActionPlanID uint, id uint, currentUserID uint, auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error {
	tx := afs.repo.BeginTransaction()

	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, auditFindingActionPlanID); err != nil {
		return err
	}

	existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingID, auditFindingActionPlanID, id)
	if err != nil {
		return err
	}

	oldAuditFindingActionPlanStatusUpdate, err := afs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := afs.validateAuditFindingActionPlanStatusUpdateFormDTO(auditFindingActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Audit Finding StatusUpdate-- only comments and percentage
	existingAuditFindingActionPlanStatusUpdate.Comment = auditFindingActionPlanStatusUpdateFormDTO.Comment
	existingAuditFindingActionPlanStatusUpdate.Percentage = auditFindingActionPlanStatusUpdateFormDTO.Percentage
	existingAuditFindingActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = afs.repo.Update(existingAuditFindingActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.AuditFindingActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditFindingActionPlanStatusUpdateDocument

	for _, doc := range auditFindingActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.AuditFindingActionPlanStatusUpdateDocument{
				AuditFindingActionPlanStatusUpdateID: id,
				URL:                                  url,
				Name:                                 doc.FileName,
				Size:                                 doc.Size,
				Token:                                encryptedFileName,
				OriginalName:                         doc.OriginalName,
				Type:                                 doc.Type,
				CreatedByID:                          currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := afs.repo.GetDocumentByID(*doc.ID)
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
		if err := afs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := afs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "action-plans", strconv.FormatUint(uint64(auditFindingActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingFindingActionPlanStatusUpdate, err := afs.repo.GetByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"percentage": oldAuditFindingActionPlanStatusUpdate.Percentage,
			"comment":    oldAuditFindingActionPlanStatusUpdate.Comment,
		}

		newData := map[string]interface{}{
			"percentage": existingFindingActionPlanStatusUpdate.Percentage,
			"comment":    existingFindingActionPlanStatusUpdate.Comment,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := afs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditFindingActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				"",
				"",
				"",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditFindingActionPlanStatusUpdate{}.TableName(),
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

func (afs AuditFindingActionPlanStatusUpdateService) GetLatestAuditFindingActionPlanStatusUpdate(auditFindingID uint, actionPlanID uint) (*auditDTOs.AuditFindingActionPlanStatusUpdateListDTO, error) {
	if _, err := afs.auditFindingActionPlanSvc.GetByID(auditFindingID, actionPlanID); err != nil {
		return nil, err
	}

	return afs.repo.GetLatestAuditFindingActionPlanStatusUpdate(actionPlanID)
}

func (afs *AuditFindingActionPlanStatusUpdateService) validateAuditFindingActionPlanStatusUpdateFormDTO(auditFindingActionPlanStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanStatusUpdateFormDTO) error {
	if _, err := afs.actionPlanStatusSvc.GetByID(auditFindingActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) GetDocumentPath(auditFindingId uint, actionPlanId uint, id uint, token string) (string, error) {
	existingAuditFindingActionPlanStatusUpdate, err := afs.GetByID(auditFindingId, actionPlanId, id)
	if err != nil {
		return "", err
	}

	if len(existingAuditFindingActionPlanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the audit finding action plan update")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingId), 10), "action-plans", strconv.FormatUint(uint64(actionPlanId), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) Approve(AuditFindingID uint, AuditFindingActionPlanID uint, currentUserID uint, AuditFindingActionPlanReviewUserStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanReviewUserStatusUpdateFormDTO) error {
	auditFinding, err := afs.auditFindingSvc.GetByID(AuditFindingID)
	if err != nil {
		return err
	}
	existingAuditFindingActionPlan, err := afs.auditFindingActionPlanSvc.GetByID(AuditFindingID, AuditFindingActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, AuditFindingActionPlanID)
	}

	if currentUserID != *existingAuditFindingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	auditFindingActionPlanStatusResolved, err := afs.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	auditFindingActionPlanClosedStatus, err := afs.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := afs.repo.BeginTransaction()

	if existingAuditFindingActionPlan.ActionPlanStatusID != auditFindingActionPlanStatusResolved.ID {
	}

	existingAuditFindingActionPlan.ActionPlanStatusID = auditFindingActionPlanClosedStatus.ID

	if err := afs.repo.AuditFindingActionPlanUpdateStatus(existingAuditFindingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	auditFindingActionPlanUpdate := &models.AuditFindingActionPlanStatusUpdate{
		AuditFindingActionPlanID: AuditFindingActionPlanID,
		ActionPlanStatusID:       auditFindingActionPlanClosedStatus.ID,
		Percentage:               nil,
		Comment:                  AuditFindingActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:              currentUserID,
	}

	id, err := afs.repo.Create(auditFindingActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := afs.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = afs.taskSvc.UpdateActionPlanTaskStatus("audit_finding_action_plan_id", AuditFindingActionPlanID, existingAuditFindingActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = afs.taskSvc.UpdateActionPlanTaskStatus("audit_finding_action_plan_id", AuditFindingActionPlanID, *existingAuditFindingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"audit_id":%d, "action_plan_id":%d}`, AuditFindingID, AuditFindingActionPlanID)
		queryParam := ""
		entity := "audit-finding-action-plan"
		action := "close-audit-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingAuditFindingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Audit Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل التدقيق",
			TitleEn:              existingAuditFindingActionPlan.Title,
			TitleAr:              existingAuditFindingActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := afs.userSvc.GetByID(existingAuditFindingActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := afs.userSvc.GetByID(*existingAuditFindingActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Audit Finding Action Plan %s: Resolved", existingAuditFindingActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                 responsibleUser.FirstNameEn,
			"LastNameEn":                  responsibleUser.LastNameEn,
			"AuditFindingReferenceNumber": auditFinding.ReferenceNumber,
			"AuditFindingTitle":           auditFinding.Title,
			"ActionPlanTitle":             existingAuditFindingActionPlan.Title,
			"ActionPlanDescription":       existingAuditFindingActionPlan.Description,
			"ResolvedByFirstNameEn":       resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":        resolvedBy.LastNameEn,
			"ResolvedOn":                  time.Now().Format("2006-01-02"),
			"RedirectURL":                 fmt.Sprintf("%s/audit/audit-findings/%d/action-plans?action_plan_id=%d", appURL, AuditFindingID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "audit/templates/email/audit_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = afs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingAuditFindingActionPlan, err := afs.auditFindingActionPlanSvc.GetByID(AuditFindingID, AuditFindingActionPlanID)
		if err != nil {
			return
		}

		err = afs.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.AuditFindingActionPlanStatusUpdate{}.TableName(),
			id,
			existingAuditFindingActionPlan,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.AuditFindingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (afs *AuditFindingActionPlanStatusUpdateService) Revert(AuditFindingID uint, AuditFindingActionPlanID uint, currentUserID uint, AuditFindingActionPlanReviewUserStatusUpdateFormDTO *auditDTOs.AuditFindingActionPlanReviewUserStatusUpdateFormDTO) error {
	auditFinding, err := afs.auditFindingSvc.GetByID(AuditFindingID)
	if err != nil {
		return err
	}
	existingAuditFindingActionPlan, err := afs.auditFindingActionPlanSvc.GetByID(AuditFindingID, AuditFindingActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, AuditFindingActionPlanID)
	}

	if currentUserID != *existingAuditFindingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	auditFindingActionPlanStatusResolved, err := afs.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	auditFindingActionPlanReopenStatus, err := afs.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := afs.repo.BeginTransaction()

	if existingAuditFindingActionPlan.ActionPlanStatusID != auditFindingActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The audit finding action plan status must be resolved.")
	}

	existingAuditFindingActionPlan.ActionPlanStatusID = auditFindingActionPlanReopenStatus.ID

	if err := afs.repo.AuditFindingActionPlanUpdateStatus(existingAuditFindingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	auditFindingActionPlanUpdate := &models.AuditFindingActionPlanStatusUpdate{
		AuditFindingActionPlanID: AuditFindingActionPlanID,
		ActionPlanStatusID:       auditFindingActionPlanReopenStatus.ID,
		Comment:                  AuditFindingActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:              currentUserID,
	}

	id, err := afs.repo.Create(auditFindingActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := afs.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = afs.taskSvc.UpdateActionPlanTaskStatus("audit_finding_action_plan_id", AuditFindingActionPlanID, *existingAuditFindingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"audit_id":%d, "action_plan_id":%d}`, AuditFindingID, AuditFindingActionPlanID)
		queryParam := ""
		entity := "audit-finding-action-plan"
		action := "reopen-audit-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingAuditFindingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Audit Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل التدقيق",
			TitleEn:              existingAuditFindingActionPlan.Title,
			TitleAr:              existingAuditFindingActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := afs.userSvc.GetByID(existingAuditFindingActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := afs.userSvc.GetByID(*existingAuditFindingActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Audit Finding Action Plan %s: Reopened", existingAuditFindingActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                 responsibleUser.FirstNameEn,
			"LastNameEn":                  responsibleUser.LastNameEn,
			"AuditFindingReferenceNumber": auditFinding.ReferenceNumber,
			"AuditFindingTitle":           auditFinding.Title,
			"ActionPlanTitle":             existingAuditFindingActionPlan.Title,
			"ActionPlanDescription":       existingAuditFindingActionPlan.Description,
			"ResolvedByFirstNameEn":       resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":        resolvedBy.LastNameEn,
			"ResolvedOn":                  time.Now().Format("2006-01-02"),
			"RedirectURL":                 fmt.Sprintf("%s/audit/audit-findings/%d/action-plans?action_plan_id=%d", appURL, AuditFindingID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "audit/templates/email/audit_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = afs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit_finding_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingAuditFindingActionPlan, err := afs.auditFindingActionPlanSvc.GetByID(AuditFindingID, AuditFindingActionPlanID)
		if err != nil {
			return
		}

		err = afs.activityLogSvc.LogActivity(
			models.RejectAction,
			models.AuditFindingActionPlanStatusUpdate{}.TableName(),
			id,
			existingAuditFindingActionPlan,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.AuditFindingActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
