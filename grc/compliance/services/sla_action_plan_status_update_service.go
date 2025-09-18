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
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLAActionPlanStatusUpdateServiceInterface interface {
	GetAll(SLAID uint, SLAActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error)
	Create(SLAID uint, SLAActionPlanID uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) (uint, error)
	Update(SLAID uint, SLAActionPlanID uint, id uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) error
	GetLatestSLAActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(slaID uint, actionPlanID uint, id uint, token string) (string, error)
	Approve(slaID uint, actionPlanID uint, currentUserID uint, SLAAStatusReviewUserUpdateFormDTO *complianceDTOs.SLAActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(slaID uint, actionPlanID uint, currentUserID uint, SLAActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.SLAActionPlanReviewUserStatusUpdateFormDTO) error
}

type SLAActionPlanStatusUpdateService struct {
	repo                repositories.SLAActionPlanStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	slaActionPlanSvc    SLAActionPlanServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	slaRegisterSvc      SLARegisterServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	taskStatusSvc       taskServices.TaskStatusServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
}

func NewSLAActionPlanStatusUpdateService(
	repo repositories.SLAActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaActionPlanSvc SLAActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	slaRegisterSvc SLARegisterServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) SLAActionPlanStatusUpdateServiceInterface {
	return &SLAActionPlanStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		slaActionPlanSvc:    slaActionPlanSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		slaRegisterSvc:      slaRegisterSvc,
		taskSvc:             taskSvc,
		taskStatusSvc:       taskStatusSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
	}
}

func (sasus *SLAActionPlanStatusUpdateService) GetAll(SLAID uint, SLAActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetAll(SLAActionPlanID, sort, sortBy, keyword)
}

func (sasus SLAActionPlanStatusUpdateService) GetByID(SLAID uint, SLAActionPlanID uint, id uint) (*models.SLAActionPlanStatusUpdate, error) {
	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return nil, err
	}

	SLAActionPlanStatusUpdate, err := sasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return SLAActionPlanStatusUpdate, nil
}

func (sasus *SLAActionPlanStatusUpdateService) Create(SLAID uint, SLAActionPlanID uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return 0, err
	}

	if err := sasus.validateSLAActionPlanStatusUpdateFormDTO(SLAActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	SLAActionPlanUpdate := &models.SLAActionPlanStatusUpdate{
		SLAActionPlanID:    SLAActionPlanID,
		ActionPlanStatusID: SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:            SLAActionPlanStatusUpdateFormDTO.Comment,
		Percentage:         SLAActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err := sasus.repo.Create(SLAActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = sasus.slaActionPlanSvc.StatusUpdate(SLAID, SLAActionPlanID, currentUserID, SLAActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.SLAActionPlanStatusUpdateDocument

	for _, doc := range SLAActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.SLAActionPlanStatusUpdateDocument{
			SLAActionPlanStatusUpdateID: id,
			URL:                         url,
			Name:                        doc.FileName,
			Size:                        doc.Size,
			Token:                       encryptedFilename,
			OriginalName:                doc.OriginalName,
			Type:                        doc.Type,
			CreatedByID:                 currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := sasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	sla, err := sasus.slaRegisterSvc.GetByID(SLAID)
	if err != nil {
		return 0, err
	}

	slaActionPlan, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && slaActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"sla_contracts_id":%d, "action_plan_id":%d}`, SLAID, SLAActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:  "sla_contracts",
			SLAActionPlanID:   &SLAActionPlanID,
			Title:             slaActionPlan.Title,
			Description:       slaActionPlan.Description,
			URLParam:          &urlParam,
			QueryParam:        nil,
			SeverityLevelID:   slaActionPlan.SeverityLevelID,
			ResponsibleUserID: *slaActionPlan.ReviewUserID,
		}

		if _, err := sasus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := sasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLAActionPlanStatusUpdate{}.TableName(),
			id,
			SLAActionPlanUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLAActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && slaActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"sla_contracts_id":%d, "action_plan_id":%d}`, SLAID, SLAActionPlanID)
			queryParam := ""
			entity := "sla-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*slaActionPlan.ReviewUserID},
				ModuleIdentifier:     "sla_contracts",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Sla Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة عمل SLA",
				TitleEn:              slaActionPlan.Title,
				TitleAr:              slaActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := sasus.userSvc.GetByID(*slaActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := sasus.userSvc.GetByID(slaActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Sla Action Plan %s: Resolved", slaActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           reviewUser.FirstNameEn,
				"LastNameEn":            reviewUser.LastNameEn,
				"SlaReferenceNumber":    sla.ReferenceNumber,
				"SlaTitle":              sla.TitleEn,
				"ActionPlanTitle":       slaActionPlan.Title,
				"ActionPlanDescription": slaActionPlan.Description,
				"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
				"ResolvedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/compliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, SLAID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/sla_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla_register_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (sasus *SLAActionPlanStatusUpdateService) Update(SLAID uint, SLAActionPlanID uint, id uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) error {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return err
	}

	oldSLAActionPlanStatusUpdate, err := sasus.GetByID(SLAID, SLAActionPlanID, id)
	if err != nil {
		return err
	}

	existingSLAActionPlanStatusUpdate, err := sasus.GetByID(SLAID, SLAActionPlanID, id)
	if err != nil {
		return err
	}

	if err := sasus.validateSLAActionPlanStatusUpdateFormDTO(SLAActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	// Update SLA StatusUpdate-- only comments and percentage
	existingSLAActionPlanStatusUpdate.Comment = SLAActionPlanStatusUpdateFormDTO.Comment
	existingSLAActionPlanStatusUpdate.Percentage = SLAActionPlanStatusUpdateFormDTO.Percentage
	existingSLAActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = sasus.repo.Update(existingSLAActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.SLAActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.SLAActionPlanStatusUpdateDocument

	for _, doc := range SLAActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.SLAActionPlanStatusUpdateDocument{
				SLAActionPlanStatusUpdateID: id,
				URL:                         url,
				Name:                        doc.FileName,
				Size:                        doc.Size,
				Token:                       encryptedFileName,
				OriginalName:                doc.OriginalName,
				Type:                        doc.Type,
				CreatedByID:                 currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := sasus.repo.GetDocumentByID(*doc.ID)
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
		if err := sasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := sasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := sasus.GetByID(SLAID, SLAActionPlanID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldSLAActionPlanStatusUpdate.Comment,
			"percentage": oldSLAActionPlanStatusUpdate.Percentage,
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
			err := sasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SLAActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.SLAActionPlanStatusUpdate{}.TableName(),
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

func (sasus *SLAActionPlanStatusUpdateService) validateSLAActionPlanStatusUpdateFormDTO(SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) error {
	if _, err := sasus.actionPlanStatusSvc.GetByID(SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (sasus SLAActionPlanStatusUpdateService) GetLatestSLAActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetLatestSLAActionPlanStatusUpdate(actionPlanID)
}

func (sasus *SLAActionPlanStatusUpdateService) GetDocumentPath(slaID, ActionPlanID, id uint, token string) (string, error) {
	existingSLAActionplanStatusUpdate, err := sasus.GetByID(slaID, ActionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingSLAActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the sla action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (sasus *SLAActionPlanStatusUpdateService) Approve(slaID uint, actionPlanID uint, currentUserID uint, SLAAStatusReviewUserUpdateFormDTO *complianceDTOs.SLAActionPlanReviewUserStatusUpdateFormDTO) error {
	sla, err := sasus.slaRegisterSvc.GetByID(slaID)
	if err != nil {
		return err
	}
	existingSLAActionPlan, err := sasus.slaActionPlanSvc.GetByID(slaID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingSLAActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this SLA action plan.")
	}

	if currentUserID != *existingSLAActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	SLAActionPlanStatusResolved, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	slaActionPlanClosedStatus, err := sasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := sasus.repo.BeginTransaction()

	if existingSLAActionPlan.ActionPlanStatusID != SLAActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The sla action plan status must be resolved.")
	}

	existingSLAActionPlan.ActionPlanStatusID = slaActionPlanClosedStatus.ID

	if err := sasus.repo.SLAActionPlanUpdateStatus(existingSLAActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	SLAActionPlanUpdate := &models.SLAActionPlanStatusUpdate{
		SLAActionPlanID:    actionPlanID,
		ActionPlanStatusID: slaActionPlanClosedStatus.ID,
		Percentage:         nil,
		Comment:            SLAAStatusReviewUserUpdateFormDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := sasus.repo.Create(SLAActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := sasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("sla_action_plan_id", actionPlanID, existingSLAActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("sla_action_plan_id", actionPlanID, *existingSLAActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"sla_register_id":%d, "action_plan_id":%d}`, slaID, actionPlanID)
		queryParam := ""
		entity := "sla-action-plan"
		action := "close-sla-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingSLAActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "sla_contracts",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "SLA Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل SLA",
			TitleEn:              existingSLAActionPlan.Title,
			TitleAr:              existingSLAActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := sasus.userSvc.GetByID(existingSLAActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := sasus.userSvc.GetByID(*existingSLAActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("SLA Action Plan %s: Resolved", existingSLAActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"SLAReferenceNumber":    sla.ReferenceNumber,
			"SLATitle":              sla.TitleEn,
			"ActionPlanTitle":       existingSLAActionPlan.Title,
			"ActionPlanDescription": existingSLAActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/compliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/sla_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla_register_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingSLAActionPlan, err := sasus.slaActionPlanSvc.GetByID(slaID, actionPlanID)
		if err != nil {
			return
		}

		err = sasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.SLAActionPlanStatusUpdate{}.TableName(),
			id,
			existingSLAActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.SLAActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sasus *SLAActionPlanStatusUpdateService) Revert(slaID uint, actionPlanID uint, currentUserID uint, SLAActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.SLAActionPlanReviewUserStatusUpdateFormDTO) error {
	sla, err := sasus.slaRegisterSvc.GetByID(slaID)
	if err != nil {
		return err
	}
	existingSLAActionPlan, err := sasus.slaActionPlanSvc.GetByID(slaID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingSLAActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	slaActionPlanStatusResolved, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	slaActionPlanReopenStatus, err := sasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := sasus.repo.BeginTransaction()

	if existingSLAActionPlan.ActionPlanStatusID != slaActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The sla action plan status must be resolved.")
	}

	existingSLAActionPlan.ActionPlanStatusID = slaActionPlanReopenStatus.ID

	if err := sasus.repo.SLAActionPlanUpdateStatus(existingSLAActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	SLAActionPlanUpdate := &models.SLAActionPlanStatusUpdate{
		SLAActionPlanID:    actionPlanID,
		ActionPlanStatusID: slaActionPlanReopenStatus.ID,
		Comment:            SLAActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := sasus.repo.Create(SLAActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := sasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("sla_action_plan_id", actionPlanID, *existingSLAActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"sla_register_id":%d, "action_plan_id":%d}`, slaID, actionPlanID)
		queryParam := ""
		entity := "sla-action-plan"
		action := "reopen-sla-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingSLAActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "sla_contracts",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "SLA Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل SLA",
			TitleEn:              existingSLAActionPlan.Title,
			TitleAr:              existingSLAActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := sasus.userSvc.GetByID(existingSLAActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := sasus.userSvc.GetByID(*existingSLAActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("SLA Action Plan %s: Reopened", existingSLAActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"SLAReferenceNumber":    sla.ReferenceNumber,
			"SLATitle":              sla.TitleEn,
			"ActionPlanTitle":       existingSLAActionPlan.Title,
			"ActionPlanDescription": existingSLAActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/compliance/sla-contracts/%d/action-plans?action_plan_id=%d", appURL, slaID, actionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/sla_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "sla_register_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingComplianceActionPlan, err := sasus.slaActionPlanSvc.GetByID(slaID, actionPlanID)
		if err != nil {
			return
		}

		err = sasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.SLAActionPlanStatusUpdate{}.TableName(),
			id,
			existingComplianceActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.SLAActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
