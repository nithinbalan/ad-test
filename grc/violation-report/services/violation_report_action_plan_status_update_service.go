package services

import (
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/grc/violation-report/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ViolationReportActionPlanStatusUpdateServiceInterface interface {
	GetAll(ViolationReportID uint, ViolationReportActionPlanID uint, sort string, sortBy string, keyword string) ([]*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error)
	Create(ViolationReportID uint, ViolationReportActionPlanID uint, currentUserID uint, violationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ViolationReportID uint, ViolationReportActionPlanID uint, id uint, currentUserID uint, violationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO) error
	GetLatestViolationReportActionPlanStatusUpdate(violationReportID uint, actionPlanID uint) (*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(violationReportID uint, actionPlanID uint, id uint, token string) (string, error)
	Approve(violationReportID uint, actionPlanID uint, currentUserID uint, ViolationReportActionPlanReviewUserStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(violationReportID uint, actionPlanID uint, currentUserID uint, ViolationReportActionPlanReviewUserStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanReviewUserStatusUpdateFormDTO) error
}

type ViolationReportActionPlanStatusUpdateService struct {
	repo                         repositories.ViolationReportActionPlanStatusUpdateRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	violationReportActionPlanSvc ViolationReportActionPlanServiceInterface
	actionPlanStatusSvc          masterServices.ActionPlanStatusServiceInterface
	violationReportSvc           ViolationReportServiceInterface
	taskSvc                      taskServices.TaskServiceInterface
	taskStatusSvc                taskServices.TaskStatusServiceInterface
	userSvc                      organizationServices.UserServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
}

func NewViolationReportActionPlanStatusUpdateService(
	repo repositories.ViolationReportActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	violationReportActionPlanSvc ViolationReportActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	violationReportSvc ViolationReportServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ViolationReportActionPlanStatusUpdateServiceInterface {
	return &ViolationReportActionPlanStatusUpdateService{
		repo:                         repo,
		activityLogSvc:               activityLogSvc,
		violationReportActionPlanSvc: violationReportActionPlanSvc,
		actionPlanStatusSvc:          actionPlanStatusSvc,
		violationReportSvc:           violationReportSvc,
		taskSvc:                      taskSvc,
		taskStatusSvc:                taskStatusSvc,
		userSvc:                      userSvc,
		notificationSvc:              notificationSvc,
	}
}

func (vrapsus *ViolationReportActionPlanStatusUpdateService) GetAll(ViolationReportID uint, ViolationReportActionPlanID uint, sort string, sortBy string, keyword string) ([]*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error) {
	if _, err := vrapsus.violationReportActionPlanSvc.GetByID(ViolationReportID, ViolationReportActionPlanID); err != nil {
		return nil, err
	}

	return vrapsus.repo.GetAll(ViolationReportActionPlanID, sort, sortBy, keyword)
}

func (vrapsus ViolationReportActionPlanStatusUpdateService) GetByID(ViolationReportID uint, ViolationReportActionPlanID uint, id uint) (*models.ViolationReportActionPlanStatusUpdate, error) {
	if _, err := vrapsus.violationReportActionPlanSvc.GetByID(ViolationReportID, ViolationReportActionPlanID); err != nil {
		return nil, err
	}

	violationReportActionPlanStatusUpdate, err := vrapsus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return violationReportActionPlanStatusUpdate, nil
}

func (vrapsus *ViolationReportActionPlanStatusUpdateService) Create(ViolationReportID uint, ViolationReportActionPlanID uint, currentUserID uint, violationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO) (uint, error) {
	if err := vrapsus.validateViolationReportActionPlanStatusUpdateFormDTO(violationReportActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	violationReport, err := vrapsus.violationReportSvc.GetByID(ViolationReportID)
	if err != nil {
		return 0, err
	}

	violationReportActionPlan, err := vrapsus.violationReportActionPlanSvc.GetByID(ViolationReportID, ViolationReportActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := vrapsus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	tx := vrapsus.repo.BeginTransaction()

	violationReportActionPlanUpdate := &models.ViolationReportActionPlanStatusUpdate{
		ViolationReportActionPlanID: ViolationReportActionPlanID,
		ActionPlanStatusID:          violationReportActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                     violationReportActionPlanStatusUpdateFormDTO.Comment,
		Percentage:                  violationReportActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:                 currentUserID,
	}

	id, err := vrapsus.repo.Create(violationReportActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = vrapsus.violationReportActionPlanSvc.StatusUpdate(ViolationReportID, ViolationReportActionPlanID, currentUserID, violationReportActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ViolationReportActionPlanStatusUpdateDocument

	for _, doc := range violationReportActionPlanStatusUpdateFormDTO.Documents {
		// Skip empty or invalid documents
		if doc.FileName == "" || doc.Size <= 0 {
			continue
		}

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(ViolationReportID), 10), "action-plans", strconv.FormatUint(uint64(ViolationReportActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ViolationReportActionPlanStatusUpdateDocument{
			ViolationReportActionPlanStatusUpdateID: id,
			URL:                                     url,
			Name:                                    doc.FileName,
			Size:                                    doc.Size,
			Token:                                   encryptedFilename,
			OriginalName:                            doc.OriginalName,
			Type:                                    doc.Type,
			CreatedByID:                             currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := vrapsus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(ViolationReportID), 10), "action-plans", strconv.FormatUint(uint64(ViolationReportActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if violationReportActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && violationReportActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, ViolationReportID, ViolationReportActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:     "violation_report_register",
			IncidentActionPlanID: &ViolationReportActionPlanID,
			Title:                violationReportActionPlan.Title,
			Description:          violationReportActionPlan.Description,
			URLParam:             &urlParam,
			QueryParam:           nil,
			SeverityLevelID:      violationReportActionPlan.SeverityLevelID,
			ResponsibleUserID:    *violationReportActionPlan.ReviewUserID,
		}

		if _, err := vrapsus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := vrapsus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportActionPlanStatusUpdate{}.TableName(),
			id,
			violationReportActionPlanUpdate,
			currentUserID,
			"violation_report", "violation_report_register", "action_plan_status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ViolationReportActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if violationReportActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && violationReportActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, ViolationReportID, ViolationReportActionPlanID)
			queryParam := ""
			entity := "violation-report-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*violationReportActionPlan.ReviewUserID},
				ModuleIdentifier:     "violation_report",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Violation Report Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة عمل تقرير المخالفة",
				TitleEn:              violationReportActionPlan.Title,
				TitleAr:              violationReportActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := vrapsus.userSvc.GetByID(*violationReportActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := vrapsus.userSvc.GetByID(violationReportActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Violation Report Action Plan %s: Resolved", violationReportActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":                    reviewUser.FirstNameEn,
				"LastNameEn":                     reviewUser.LastNameEn,
				"ViolationReportReferenceNumber": violationReport.ReferenceNumber,
				"ViolationReportTitle":           violationReport.TitleEn,
				"ActionPlanTitle":                violationReportActionPlan.Title,
				"ActionPlanDescription":          violationReportActionPlan.Description,
				"ResolvedByFirstNameEn":          resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":           resolvedBy.LastNameEn,
				"ResolvedOn":                     time.Now().Format("2006-01-02"),
				"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/action-plans?action_plan_id=%d", appURL, ViolationReportID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = vrapsus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (vrapsus *ViolationReportActionPlanStatusUpdateService) validateViolationReportActionPlanStatusUpdateFormDTO(violationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO) error {
	if _, err := vrapsus.actionPlanStatusSvc.GetByID(violationReportActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (vrapsus ViolationReportActionPlanStatusUpdateService) GetLatestViolationReportActionPlanStatusUpdate(violationReportID uint, actionPlanID uint) (*violationReportDTOs.ViolationReportActionPlanStatusUpdateListDTO, error) {
	if _, err := vrapsus.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID); err != nil {
		return nil, err
	}

	return vrapsus.repo.GetLatestViolationReportActionPlanStatusUpdate(actionPlanID)
}

func (vrapsus *ViolationReportActionPlanStatusUpdateService) GetDocumentPath(ViolationReportID uint, ViolationReportActionPlanID uint, id uint, token string) (string, error) {
	existingViolationReportActionPlanStatusUpdate, err := vrapsus.GetByID(ViolationReportID, ViolationReportActionPlanID, id)
	if err != nil {
		return "", err
	}

	if len(existingViolationReportActionPlanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the violation report action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(ViolationReportID), 10), "action-plans", strconv.FormatUint(uint64(ViolationReportActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (vrapsus *ViolationReportActionPlanStatusUpdateService) Update(ViolationReportID uint, ViolationReportActionPlanID uint, id uint, currentUserID uint, ViolationReportActionPlanStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanStatusUpdateFormDTO) error {
	tx := vrapsus.repo.BeginTransaction()

	if _, err := vrapsus.violationReportActionPlanSvc.GetByID(ViolationReportID, ViolationReportActionPlanID); err != nil {
		return err
	}

	oldViolationReportActionPlanStatusUpdate, err := vrapsus.GetByID(ViolationReportID, ViolationReportActionPlanID, id)
	if err != nil {
		return err
	}

	existingViolationReportActionPlanStatusUpdate, err := vrapsus.GetByID(ViolationReportID, ViolationReportActionPlanID, id)
	if err != nil {
		return err
	}

	if err := vrapsus.validateViolationReportActionPlanStatusUpdateFormDTO(ViolationReportActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Violation Report StatusUpdate-- only comments and percentage
	existingViolationReportActionPlanStatusUpdate.Comment = ViolationReportActionPlanStatusUpdateFormDTO.Comment
	existingViolationReportActionPlanStatusUpdate.Percentage = ViolationReportActionPlanStatusUpdateFormDTO.Percentage
	existingViolationReportActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = vrapsus.repo.Update(existingViolationReportActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ViolationReportActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ViolationReportActionPlanStatusUpdateDocument

	for _, doc := range ViolationReportActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Skip empty or invalid documents
			if doc.FileName == "" || doc.Size <= 0 {
				continue
			}

			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// generate URL
			url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(ViolationReportID), 10), "action-plans", strconv.FormatUint(uint64(ViolationReportActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new document instance
			document := &models.ViolationReportActionPlanStatusUpdateDocument{
				ViolationReportActionPlanStatusUpdateID: id,
				URL:                                     url,
				Name:                                    doc.FileName,
				Size:                                    doc.Size,
				Token:                                   encryptedFileName,
				OriginalName:                            doc.OriginalName,
				Type:                                    doc.Type,
				CreatedByID:                             currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := vrapsus.repo.GetDocumentByID(*doc.ID)
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
		if err := vrapsus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := vrapsus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(ViolationReportID), 10), "action-plans", strconv.FormatUint(uint64(ViolationReportActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(ViolationReportID), 10), "action-plans", strconv.FormatUint(uint64(ViolationReportActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlanStatusUpdate, err := vrapsus.GetByID(ViolationReportID, ViolationReportActionPlanID, id)
		if err != nil {
			tx.Rollback()
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldViolationReportActionPlanStatusUpdate.Comment,
			"percentage": oldViolationReportActionPlanStatusUpdate.Percentage,
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
			err := vrapsus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ViolationReportActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				"violation_report", "violation_report_register", "action_plan_status_update", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ViolationReportActionPlan{}.TableName(),
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

func (vrapsus *ViolationReportActionPlanStatusUpdateService) Approve(violationReportID uint, actionPlanID uint, currentUserID uint, ViolationReportActionPlanReviewUserStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanReviewUserStatusUpdateFormDTO) error {
	violationReport, err := vrapsus.violationReportSvc.GetByID(violationReportID)
	if err != nil {
		return err
	}
	existingViolationReportActionPlan, err := vrapsus.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingViolationReportActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this violation report action plan.")
	}

	if currentUserID != *existingViolationReportActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	violationReportActionPlanStatusResolved, err := vrapsus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	violationReportActionPlanClosedStatus, err := vrapsus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := vrapsus.repo.BeginTransaction()

	if existingViolationReportActionPlan.ActionPlanStatusID != violationReportActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The violation report action plan status must be resolved.")
	}

	existingViolationReportActionPlan.ActionPlanStatusID = violationReportActionPlanClosedStatus.ID

	if err := vrapsus.repo.ViolationReportActionPlanUpdateStatus(existingViolationReportActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	violationReportActionPlanUpdate := &models.ViolationReportActionPlanStatusUpdate{
		ViolationReportActionPlanID: actionPlanID,
		ActionPlanStatusID:          violationReportActionPlanClosedStatus.ID,
		Percentage:                  nil,
		Comment:                     ViolationReportActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:                 currentUserID,
	}

	id, err := vrapsus.repo.Create(violationReportActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := vrapsus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = vrapsus.taskSvc.UpdateActionPlanTaskStatus("violation_report_action_plan_id", actionPlanID, existingViolationReportActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = vrapsus.taskSvc.UpdateActionPlanTaskStatus("violation_report_action_plan_id", actionPlanID, *existingViolationReportActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, actionPlanID)
		queryParam := ""
		entity := "violation-report-action-plan"
		action := "close-violation-report-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingViolationReportActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "violation_report",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Violation Report Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل تقرير المخالفة",
			TitleEn:              existingViolationReportActionPlan.Title,
			TitleAr:              existingViolationReportActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := vrapsus.userSvc.GetByID(existingViolationReportActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := vrapsus.userSvc.GetByID(*existingViolationReportActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Violation Report Action Plan %s: Resolved", existingViolationReportActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                    responsibleUser.FirstNameEn,
			"LastNameEn":                     responsibleUser.LastNameEn,
			"ViolationReportReferenceNumber": violationReport.ReferenceNumber,
			"ViolationReportTitle":           violationReport.TitleEn,
			"ActionPlanTitle":                existingViolationReportActionPlan.Title,
			"ActionPlanDescription":          existingViolationReportActionPlan.Description,
			"ResolvedByFirstNameEn":          resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":           resolvedBy.LastNameEn,
			"ResolvedOn":                     time.Now().Format("2006-01-02"),
			"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/action-plan?action_plan_id=%d", appURL, violationReportID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = vrapsus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingViolationReportActionPlan, err := vrapsus.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID)
		if err != nil {
			return
		}

		err = vrapsus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.ViolationReportActionPlanStatusUpdate{}.TableName(),
			id,
			existingViolationReportActionPlan,
			currentUserID,
			"violation_report", "violation_report_register", "action_plan_status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ViolationReportActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (vrapsus *ViolationReportActionPlanStatusUpdateService) Revert(violationReportID uint, actionPlanID uint, currentUserID uint, ViolationReportActionPlanReviewUserStatusUpdateFormDTO *violationReportDTOs.ViolationReportActionPlanReviewUserStatusUpdateFormDTO) error {
	violationReport, err := vrapsus.violationReportSvc.GetByID(violationReportID)
	if err != nil {
		return err
	}
	existingViolationReportActionPlan, err := vrapsus.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingViolationReportActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	violationReportActionPlanStatusResolved, err := vrapsus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	violationReportActionPlanClosedStatus, err := vrapsus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := vrapsus.repo.BeginTransaction()

	if existingViolationReportActionPlan.ActionPlanStatusID != violationReportActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The violation report action plan status must be resolved.")
	}

	existingViolationReportActionPlan.ActionPlanStatusID = violationReportActionPlanClosedStatus.ID

	if err := vrapsus.repo.ViolationReportActionPlanUpdateStatus(existingViolationReportActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	violationReportActionPlanUpdate := &models.ViolationReportActionPlanStatusUpdate{
		ViolationReportActionPlanID: actionPlanID,
		ActionPlanStatusID:          violationReportActionPlanClosedStatus.ID,
		Comment:                     ViolationReportActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:                 currentUserID,
	}

	id, err := vrapsus.repo.Create(violationReportActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := vrapsus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = vrapsus.taskSvc.UpdateActionPlanTaskStatus("violation_report_action_plan_id", actionPlanID, *existingViolationReportActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"violation_report_id":%d, "action_plan_id":%d}`, violationReportID, actionPlanID)
		queryParam := ""
		entity := "violation-report-action-plan"
		action := "reopen-violation-report-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingViolationReportActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "violation_report",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Violation Report Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل تقرير المخالفة",
			TitleEn:              existingViolationReportActionPlan.Title,
			TitleAr:              existingViolationReportActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := vrapsus.userSvc.GetByID(existingViolationReportActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := vrapsus.userSvc.GetByID(*existingViolationReportActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Violation Report Action Plan %s: Reopened", existingViolationReportActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":                    responsibleUser.FirstNameEn,
			"LastNameEn":                     responsibleUser.LastNameEn,
			"ViolationReportReferenceNumber": violationReport.ReferenceNumber,
			"ViolationReportTitle":           violationReport.TitleEn,
			"ActionPlanTitle":                existingViolationReportActionPlan.Title,
			"ActionPlanDescription":          existingViolationReportActionPlan.Description,
			"ResolvedByFirstNameEn":          resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":           resolvedBy.LastNameEn,
			"ResolvedOn":                     time.Now().Format("2006-01-02"),
			"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/action-plan?action_plan_id=%d", appURL, violationReportID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "violation-report/templates/email/violation_report_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = vrapsus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingViolationReportActionPlan, err := vrapsus.violationReportActionPlanSvc.GetByID(violationReportID, actionPlanID)
		if err != nil {
			return
		}

		err = vrapsus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.ViolationReportActionPlanStatusUpdate{}.TableName(),
			id,
			existingViolationReportActionPlan,
			currentUserID,
			"violation_report", "violation_report_register", "action_plan_status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ViolationReportActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
