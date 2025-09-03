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
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentActionPlanStatusUpdateServiceInterface interface {
	GetAll(IncidentID uint, IncidentActionPlanID uint, sort string, sortBy string, keyword string) ([]*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error)
	Create(IncidentID uint, IncidentActionPlanID uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) (uint, error)
	Update(IncidentID uint, IncidentActionPlanID uint, id uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) error
	GetLatestIncidentActionPlanStatusUpdate(incidentID uint, actionPlanID uint) (*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(incidentID uint, actionPlanID uint, id uint, token string) (string, error)
	Approve(incidentID uint, actionPlanID uint, currentUserID uint, IncidenttatusReviewUserUpdateFormDTO *incidentDTOs.IncidentActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(incidentID uint, actionPlanID uint, currentUserID uint, IncidenttatusReviewUserUpdateFormDTO *incidentDTOs.IncidentActionPlanReviewUserStatusUpdateFormDTO) error
}

type IncidentActionPlanStatusUpdateService struct {
	repo                  repositories.IncidentActionPlanStatusUpdateRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	incidentActionPlanSvc IncidentActionPlanServiceInterface
	actionPlanStatusSvc   masterServices.ActionPlanStatusServiceInterface
	incidentRegisterSvc   IncidentRegisterServiceInterface
	taskSvc               taskServices.TaskServiceInterface
	taskStatusSvc         taskServices.TaskStatusServiceInterface
	userSvc               organizationServices.UserServiceInterface
	notificationSvc       commonServices.NotificationServiceInterface
}

func NewIncidentActionPlanStatusUpdateService(
	repo repositories.IncidentActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentActionPlanSvc IncidentActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) IncidentActionPlanStatusUpdateServiceInterface {
	return &IncidentActionPlanStatusUpdateService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		incidentActionPlanSvc: incidentActionPlanSvc,
		actionPlanStatusSvc:   actionPlanStatusSvc,
		incidentRegisterSvc:   incidentRegisterSvc,
		taskSvc:               taskSvc,
		taskStatusSvc:         taskStatusSvc,
		userSvc:               userSvc,
		notificationSvc:       notificationSvc,
	}
}

func (iasus *IncidentActionPlanStatusUpdateService) GetAll(IncidentID uint, IncidentActionPlanID uint, sort string, sortBy string, keyword string) ([]*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error) {
	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return nil, err
	}

	return iasus.repo.GetAll(IncidentActionPlanID, sort, sortBy, keyword)
}

func (iasus IncidentActionPlanStatusUpdateService) GetByID(IncidentID uint, IncidentActionPlanID uint, id uint) (*models.IncidentActionPlanStatusUpdate, error) {
	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return nil, err
	}

	incidentActionPlanStatusUpdate, err := iasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentActionPlanStatusUpdate, nil
}

func (iasus *IncidentActionPlanStatusUpdateService) Create(IncidentID uint, IncidentActionPlanID uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) (uint, error) {
	if err := iasus.validateIncidentActionPlanStatusUpdateFormDTO(incidentActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	incident, err := iasus.incidentRegisterSvc.GetByID(IncidentID)
	if err != nil {
		return 0, err
	}

	incidentActionPlan, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := iasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	tx := iasus.repo.BeginTransaction()

	//@todo validate document form

	incidentActionPlanUpdate := &models.IncidentActionPlanStatusUpdate{
		IncidentActionPlanID: IncidentActionPlanID,
		ActionPlanStatusID:   incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:              incidentActionPlanStatusUpdateFormDTO.Comment,
		Percentage:           incidentActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:          currentUserID,
	}

	id, err := iasus.repo.Create(incidentActionPlanUpdate,tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = iasus.incidentActionPlanSvc.StatusUpdate(IncidentID, IncidentActionPlanID, currentUserID, incidentActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IncidentActionPlanStatusUpdateDocument

	for _, doc := range incidentActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.IncidentActionPlanStatusUpdateDocument{
			IncidentActionPlanStatusUpdateID: id,
			URL:                              url,
			Name:                             doc.FileName,
			Size:                             doc.Size,
			Token:                            encryptedFilename,
			OriginalName:                     doc.OriginalName,
			Type:                             doc.Type,
			CreatedByID:                      currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := iasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && incidentActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"incident_register_id":%d, "action_plan_id":%d}`, IncidentID, IncidentActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:     "incident_register",
			IncidentActionPlanID: &IncidentActionPlanID,
			Title:                incidentActionPlan.Title,
			Description:          incidentActionPlan.Description,
			URLParam:             &urlParam,
			QueryParam:           nil,
			SeverityLevelID:      incidentActionPlan.SeverityLevelID,
			ResponsibleUserID:    *incidentActionPlan.ReviewUserID,
		}

		if _, err := iasus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := iasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentActionPlanStatusUpdate{}.TableName(),
			id,
			incidentActionPlanUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && incidentActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"incident_register_id":%d, "action_plan_id":%d}`, IncidentID, IncidentActionPlanID)
			queryParam := ""
			entity := "incident-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*incidentActionPlan.ReviewUserID},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Incident Action Plan Status Update",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
				TitleEn:              incidentActionPlan.Title,
				TitleAr:              incidentActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := iasus.userSvc.GetByID(*incidentActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := iasus.userSvc.GetByID(incidentActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Incidenta Action Plan %s: Resolved", incidentActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":             reviewUser.FirstNameEn,
				"LastNameEn":              reviewUser.LastNameEn,
				"IncidentReferenceNumber": incident.ReferenceNumber,
				"IncidentTitle":           incident.Title,
				"ActionPlanTitle":         incidentActionPlan.Title,
				"ActionPlanDescription":   incidentActionPlan.Description,
				"ResolvedByFirstNameEn":   resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":    resolvedBy.LastNameEn,
				"ResolvedOn":              time.Now().Format("2006-01-02"),
				"RedirectURL":             fmt.Sprintf("%s/compliance/register/%d/action-plans?action_plan_id=%d", appURL, IncidentID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "incident/templates/email/incident_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = iasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (iasus *IncidentActionPlanStatusUpdateService) Update(IncidentID uint, IncidentActionPlanID uint, id uint, currentUserID uint, IncidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) error {
	tx := iasus.repo.BeginTransaction()

	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return err
	}

	oldIncidentActionPlanStatusUpdate, err := iasus.GetByID(IncidentID, IncidentActionPlanID, id)
	if err != nil {
		return err
	}

	existingIncidentActionPlanStatusUpdate, err := iasus.GetByID(IncidentID, IncidentActionPlanID, id)
	if err != nil {
		return err
	}

	if err := iasus.validateIncidentActionPlanStatusUpdateFormDTO(IncidentActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Incident StatusUpdate-- only comments and percentage
	existingIncidentActionPlanStatusUpdate.Comment = IncidentActionPlanStatusUpdateFormDTO.Comment
	existingIncidentActionPlanStatusUpdate.Percentage = IncidentActionPlanStatusUpdateFormDTO.Percentage
	existingIncidentActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = iasus.repo.Update(existingIncidentActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.IncidentActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.IncidentActionPlanStatusUpdateDocument

	for _, doc := range IncidentActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.IncidentActionPlanStatusUpdateDocument{
				IncidentActionPlanStatusUpdateID: id,
				URL:                              url,
				Name:                             doc.FileName,
				Size:                             doc.Size,
				Token:                            encryptedFileName,
				OriginalName:                     doc.OriginalName,
				Type:                             doc.Type,
				CreatedByID:                      currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := iasus.repo.GetDocumentByID(*doc.ID)
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
		if err := iasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := iasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlanstatusUpdate, err := iasus.GetByID(IncidentID, IncidentActionPlanID, id)
		if err != nil {
			tx.Rollback()
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldIncidentActionPlanStatusUpdate.Comment,
			"percentage": oldIncidentActionPlanStatusUpdate.Percentage,
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
			err := iasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IncidentActionPlan{}.TableName(),
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

func (iasus IncidentActionPlanStatusUpdateService) GetLatestIncidentActionPlanStatusUpdate(incidentID uint, actionPlanID uint) (*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error) {
	if _, err := iasus.incidentActionPlanSvc.GetByID(incidentID, actionPlanID); err != nil {
		return nil, err
	}

	return iasus.repo.GetLatestIncidentActionPlanStatusUpdate(actionPlanID)
}

func (iasus *IncidentActionPlanStatusUpdateService) validateIncidentActionPlanStatusUpdateFormDTO(incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) error {
	if _, err := iasus.actionPlanStatusSvc.GetByID(incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (iasus *IncidentActionPlanStatusUpdateService) GetDocumentPath(IncidentID uint, IncidentActionPlanID uint, id uint, token string) (string, error) {
	existingIncidentActionplanStatusUpdate, err := iasus.GetByID(IncidentID, IncidentActionPlanID, id)
	if err != nil {
		return "", err
	}

	if len(existingIncidentActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the incident action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (iasus *IncidentActionPlanStatusUpdateService) Approve(incidentID uint, actionPlanID uint, currentUserID uint, IncidentActionPlanReviewUserStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanReviewUserStatusUpdateFormDTO) error {
	incident, err := iasus.incidentRegisterSvc.GetByID(incidentID)
	if err != nil {
		return err
	}
	existingIncidentActionPlan, err := iasus.incidentActionPlanSvc.GetByID(incidentID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingIncidentActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this Incident action plan.")
	}

	if currentUserID != *existingIncidentActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	incidentActionPlanStatusResolved, err := iasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	incidentActionPlanClosedStatus, err := iasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := iasus.repo.BeginTransaction()

	if existingIncidentActionPlan.ActionPlanStatusID != incidentActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The Incident action plan status must be resolved.")
	}

	existingIncidentActionPlan.ActionPlanStatusID = incidentActionPlanClosedStatus.ID

	if err := iasus.repo.IncidentActionPlanUpdateStatus(existingIncidentActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	incidentActionPlanUpdate := &models.IncidentActionPlanStatusUpdate{
		IncidentActionPlanID: actionPlanID,
		ActionPlanStatusID:   incidentActionPlanClosedStatus.ID,
		Percentage:           nil,
		Comment:              IncidentActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:          currentUserID,
	}

	id, err := iasus.repo.Create(incidentActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := iasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = iasus.taskSvc.UpdateActionPlanTaskStatus("incident_action_plan_id", actionPlanID, existingIncidentActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = iasus.taskSvc.UpdateActionPlanTaskStatus("incident_action_plan_id", actionPlanID, *existingIncidentActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"incident_id":%d, "action_plan_id":%d}`, incidentID, actionPlanID)
		queryParam := ""
		entity := "incident-action-plan"
		action := "close-incident-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingIncidentActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Incident Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل الحوادث",
			TitleEn:              existingIncidentActionPlan.Title,
			TitleAr:              existingIncidentActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := iasus.userSvc.GetByID(existingIncidentActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := iasus.userSvc.GetByID(*existingIncidentActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Incident Action Plan %s: Resolved", existingIncidentActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":             responsibleUser.FirstNameEn,
			"LastNameEn":              responsibleUser.LastNameEn,
			"IncidentReferenceNumber": incident.ReferenceNumber,
			"IncidentTitle":           incident.Title,
			"ActionPlanTitle":         existingIncidentActionPlan.Title,
			"ActionPlanDescription":   existingIncidentActionPlan.Description,
			"ResolvedByFirstNameEn":   resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":    resolvedBy.LastNameEn,
			"ResolvedOn":              time.Now().Format("2006-01-02"),
			"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/action-plan?action_plan_id=%d", appURL, incidentID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "incident/templates/email/incident_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = iasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingIncidentActionPlan, err := iasus.incidentActionPlanSvc.GetByID(incidentID, actionPlanID)
		if err != nil {
			return
		}

		err = iasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.IncidentActionPlanStatusUpdate{}.TableName(),
			id,
			existingIncidentActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.IncidentActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (iasus *IncidentActionPlanStatusUpdateService) Revert(incidentID uint, actionPlanID uint, currentUserID uint, IncidentActionPlanReviewUserStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanReviewUserStatusUpdateFormDTO) error {
	incident, err := iasus.incidentRegisterSvc.GetByID(incidentID)
	if err != nil {
		return err
	}
	existingIncidentActionPlan, err := iasus.incidentActionPlanSvc.GetByID(incidentID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingIncidentActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	incidentActionPlanStatusResolved, err := iasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	incidentActionPlanClosedStatus, err := iasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := iasus.repo.BeginTransaction()

	if existingIncidentActionPlan.ActionPlanStatusID != incidentActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The incident action plan status must be resolved.")
	}

	existingIncidentActionPlan.ActionPlanStatusID = incidentActionPlanClosedStatus.ID

	if err := iasus.repo.IncidentActionPlanUpdateStatus(existingIncidentActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	incidentActionPlanUpdate := &models.IncidentActionPlanStatusUpdate{
		IncidentActionPlanID: incidentID,
		ActionPlanStatusID:   incidentActionPlanClosedStatus.ID,
		Comment:              IncidentActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:          currentUserID,
	}

	id, err := iasus.repo.Create(incidentActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := iasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = iasus.taskSvc.UpdateActionPlanTaskStatus("incident_action_plan_id", actionPlanID, *existingIncidentActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"incident_id":%d, "action_plan_id":%d}`, incidentID, actionPlanID)
		queryParam := ""
		entity := "incident-action-plan"
		action := "reopen-incident-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingIncidentActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Incident Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل الحوادث",
			TitleEn:              existingIncidentActionPlan.Title,
			TitleAr:              existingIncidentActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := iasus.userSvc.GetByID(existingIncidentActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := iasus.userSvc.GetByID(*existingIncidentActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("incident Action Plan %s: Reopened", existingIncidentActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":             responsibleUser.FirstNameEn,
			"LastNameEn":              responsibleUser.LastNameEn,
			"IncidentReferenceNumber": incident.ReferenceNumber,
			"IncidentTitle":           incident.Title,
			"ActionPlanTitle":         existingIncidentActionPlan.Title,
			"ActionPlanDescription":   existingIncidentActionPlan.Description,
			"ResolvedByFirstNameEn":   resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":    resolvedBy.LastNameEn,
			"ResolvedOn":              time.Now().Format("2006-01-02"),
			"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/action-plan?action_plan_id=%d", appURL, incidentID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "incident/templates/email/incident_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = iasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingIncidentActionPlan, err := iasus.incidentActionPlanSvc.GetByID(incidentID, actionPlanID)
		if err != nil {
			return
		}

		err = iasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.IncidentActionPlanStatusUpdate{}.TableName(),
			id,
			existingIncidentActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.IncidentActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
