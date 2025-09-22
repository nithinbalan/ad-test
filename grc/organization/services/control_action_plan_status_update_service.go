package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlActionPlanStatusUpdateServiceInterface interface {
	GetByID(controlID uint, controlActionPlanID uint, id uint) (*models.ControlActionPlanStatusUpdate, error)
	GetAll(controlID uint, controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error)
	Create(controlID uint, controlActionPlanID uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) (uint, error)
	Update(controlID uint, controlActionPlanID uint, id uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error
	GetLatestControlActionPlanStatusUpdate(controlID uint, actionPlanID uint) (*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error)
	Approve(controlID uint, actionPlanID uint, currentUserID uint, controlActionPlanReviewUserStatusUpdateFormDTO *organizationDTOs.ControlActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(controlID uint, actionPlanID uint, currentUserID uint, controlActionPlanReviewUserStatusUpdateFormDTO *organizationDTOs.ControlActionPlanReviewUserStatusUpdateFormDTO) error
}

type ControlActionPlanStatusUpdateService struct {
	repo                 repositories.ControlActionPlanStatusUpdateRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	controlActionPlanSvc ControlActionPlanServiceInterface
	actionPlanStatusSvc  masterServices.ActionPlanStatusServiceInterface
	controlSvc           ControlServiceInterface
	userSvc              UserServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	taskSvc              taskServices.TaskServiceInterface
	taskStatusSvc        taskServices.TaskStatusServiceInterface
}

func NewControlActionPlanStatusUpdateService(
	repo repositories.ControlActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlActionPlanSvc ControlActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	controlSvc ControlServiceInterface,
	userSvc UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) ControlActionPlanStatusUpdateServiceInterface {
	return &ControlActionPlanStatusUpdateService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		controlActionPlanSvc: controlActionPlanSvc,
		actionPlanStatusSvc:  actionPlanStatusSvc,
		controlSvc:           controlSvc,
		userSvc:              userSvc,
		notificationSvc:      notificationSvc,
		taskSvc:              taskSvc,
		taskStatusSvc:        taskStatusSvc,
	}
}

func (casus *ControlActionPlanStatusUpdateService) GetAll(controlID uint, controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(controlActionPlanID, sort, sortBy, keyword)
}

func (casus ControlActionPlanStatusUpdateService) GetByID(controlID uint, controlActionPlanID uint, id uint) (*models.ControlActionPlanStatusUpdate, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return nil, err
	}

	controlActionPlanStatusUpdate, err := casus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return controlActionPlanStatusUpdate, nil
}

func (casus *ControlActionPlanStatusUpdateService) Create(controlID uint, controlActionPlanID uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) (uint, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return 0, err
	}

	control, err := casus.controlSvc.GetByID(controlID)
	if err != nil {
		return 0, err
	}

	controlActionPlan, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := casus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	tx := casus.repo.BeginTransaction()
	if err := casus.validateControlActionPlanStatusUpdateFormDTO(controlActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	controlActionPlanUpdate := &models.ControlActionPlanStatusUpdate{
		ControlActionPlanID: controlActionPlanID,
		ActionPlanStatusID:  controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:             controlActionPlanStatusUpdateFormDTO.Comment,
		Percentage:          &controlActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:         currentUserID,
	}

	id, err := casus.repo.Create(controlActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = casus.controlActionPlanSvc.StatusUpdate(controlID, controlActionPlanID, currentUserID, controlActionPlanStatusUpdateFormDTO, controlActionPlanStatusUpdateFormDTO.Percentage, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ControlActionPlanStatusUpdateDocument

	for _, doc := range controlActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), encryptedFileName)

		document := &models.ControlActionPlanStatusUpdateDocument{
			ControlActionPlanStatusUpdateID: id,
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
		if err := casus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	if controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && controlActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"control_id":%d, "action_plan_id":%d}`, controlID, controlActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:    "controls",
			ControlActionPlanID: &controlActionPlanID,
			Title:               controlActionPlan.Title,
			Description:         controlActionPlan.Description,
			URLParam:            &urlParam,
			QueryParam:          nil,
			SeverityLevelID:     controlActionPlan.SeverityLevelID,
			ResponsibleUserID:   *controlActionPlan.ReviewUserID,
		}

		if _, err := casus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := casus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlActionPlanStatusUpdate{}.TableName(),
			id,
			controlActionPlanUpdate,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && controlActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"control_id":%d, "action_plan_id":%d}`, controlID, controlActionPlanID)
			queryParam := ""
			entity := "control-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*controlActionPlan.ReviewUserID},
				ModuleIdentifier:     "controls",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Control Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة عمل الرقابة",
				TitleEn:              controlActionPlan.Title,
				TitleAr:              controlActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := casus.userSvc.GetByID(*controlActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := casus.userSvc.GetByID(controlActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Control Action Plan %s: Resolved", controlActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":            reviewUser.FirstNameEn,
				"LastNameEn":             reviewUser.LastNameEn,
				"ControlReferenceNumber": control.ReferenceNumber,
				"ControlTitle":           control.TitleEn,
				"ActionPlanTitle":        controlActionPlan.Title,
				"ActionPlanDescription":  controlActionPlan.Description,
				"ResolvedByFirstNameEn":  resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":   resolvedBy.LastNameEn,
				"ResolvedOn":             time.Now().Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/control_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = casus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (casus *ControlActionPlanStatusUpdateService) Update(controlID uint, controlActionPlanID uint, id uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return err
	}

	existingControlActionPlanStatusUpdate, err := casus.GetByID(controlID, controlActionPlanID, id)
	if err != nil {
		return err
	}

	if err := casus.validateControlActionPlanStatusUpdateFormDTO(controlActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingControlActionPlanStatusUpdate.Comment = controlActionPlanStatusUpdateFormDTO.Comment
	existingControlActionPlanStatusUpdate.Percentage = &controlActionPlanStatusUpdateFormDTO.Percentage
	existingControlActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = casus.repo.Update(existingControlActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ControlActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ControlActionPlanStatusUpdateDocument

	for _, doc := range controlActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), encryptedFileName)

			document := &models.ControlActionPlanStatusUpdateDocument{
				ControlActionPlanStatusUpdateID: id,
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
				existingDocument, err := casus.repo.GetDocumentByID(*doc.ID)
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
		if err := casus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := casus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (casus *ControlActionPlanStatusUpdateService) validateControlActionPlanStatusUpdateFormDTO(controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error {
	if _, err := casus.actionPlanStatusSvc.GetByID(controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (casus ControlActionPlanStatusUpdateService) GetLatestControlActionPlanStatusUpdate(controlID uint, actionPlanID uint) (*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetLatestControlActionPlanStatusUpdate(actionPlanID)
}

func (casus *ControlActionPlanStatusUpdateService) Approve(controlID uint, actionPlanID uint, currentUserID uint, controlActionPlanReviewUserStatusUpdateFormDTO *organizationDTOs.ControlActionPlanReviewUserStatusUpdateFormDTO) error {
	control, err := casus.controlSvc.GetByID(controlID)
	if err != nil {
		return err
	}
	existingControlActionPlan, err := casus.controlActionPlanSvc.GetByID(controlID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingControlActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	controlActionPlanStatusResolved, err := casus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	controlActionPlanClosedStatus, err := casus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := casus.repo.BeginTransaction()

	if existingControlActionPlan.ActionPlanStatusID != controlActionPlanStatusResolved.ID {
	}

	existingControlActionPlan.ActionPlanStatusID = controlActionPlanClosedStatus.ID

	if err := casus.repo.ControlActionPlanUpdateStatus(existingControlActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	controlActionPlanUpdate := &models.ControlActionPlanStatusUpdate{
		ControlActionPlanID: actionPlanID,
		ActionPlanStatusID:  controlActionPlanClosedStatus.ID,
		Percentage:          nil,
		Comment:             controlActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:         currentUserID,
	}

	id, err := casus.repo.Create(controlActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := casus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = casus.taskSvc.UpdateActionPlanTaskStatus("control_action_plan_id", actionPlanID, existingControlActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = casus.taskSvc.UpdateActionPlanTaskStatus("control_action_plan_id", actionPlanID, *existingControlActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"conrtol_id":%d, "action_plan_id":%d}`, controlID, actionPlanID)
		queryParam := ""
		entity := "control-action-plan"
		action := "close-control-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingControlActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "controls",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Control Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل الرقابة",
			TitleEn:              existingControlActionPlan.Title,
			TitleAr:              existingControlActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := casus.userSvc.GetByID(existingControlActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := casus.userSvc.GetByID(*existingControlActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Control Action Plan %s: Resolved", existingControlActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"ControlReferenceNumber": control.ReferenceNumber,
			"ControlTitle":           control.TitleEn,
			"ActionPlanTitle":        existingControlActionPlan.Title,
			"ActionPlanDescription":  existingControlActionPlan.Description,
			"ResolvedByFirstNameEn":  resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":   resolvedBy.LastNameEn,
			"ResolvedOn":             time.Now().Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/control_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = casus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingControlActionPlan, err := casus.controlActionPlanSvc.GetByID(controlID, actionPlanID)
		if err != nil {
			return
		}

		err = casus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.ControlActionPlanStatusUpdate{}.TableName(),
			id,
			existingControlActionPlan,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ControlActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (casus *ControlActionPlanStatusUpdateService) Revert(controlID uint, actionPlanID uint, currentUserID uint, controlActionPlanReviewUserStatusUpdateFormDTO *organizationDTOs.ControlActionPlanReviewUserStatusUpdateFormDTO) error {
	control, err := casus.controlSvc.GetByID(controlID)
	if err != nil {
		return err
	}
	existingControlActionPlan, err := casus.controlActionPlanSvc.GetByID(controlID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingControlActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	controlActionPlanStatusResolved, err := casus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	controlActionPlanReopenStatus, err := casus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := casus.repo.BeginTransaction()

	if existingControlActionPlan.ActionPlanStatusID != controlActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The controla action plan status must be resolved.")
	}

	existingControlActionPlan.ActionPlanStatusID = controlActionPlanReopenStatus.ID

	if err := casus.repo.ControlActionPlanUpdateStatus(existingControlActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	contrtolActionPlanUpdate := &models.ControlActionPlanStatusUpdate{
		ControlActionPlanID: actionPlanID,
		ActionPlanStatusID:  controlActionPlanReopenStatus.ID,
		Comment:             controlActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:         currentUserID,
	}

	id, err := casus.repo.Create(contrtolActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := casus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = casus.taskSvc.UpdateActionPlanTaskStatus("control_action_plan_id", actionPlanID, *existingControlActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"control_id":%d, "action_plan_id":%d}`, controlID, actionPlanID)
		queryParam := ""
		entity := "control-action-plan"
		action := "reopen-control-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingControlActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "controls",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Control Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل الرقابة",
			TitleEn:              existingControlActionPlan.Title,
			TitleAr:              existingControlActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := casus.userSvc.GetByID(existingControlActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := casus.userSvc.GetByID(*existingControlActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Control Action Plan %s: Reopened", existingControlActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"ControlReferenceNumber": control.ReferenceNumber,
			"ControlTitle":           control.TitleEn,
			"ActionPlanTitle":        existingControlActionPlan.Title,
			"ActionPlanDescription":  existingControlActionPlan.Description,
			"ResolvedByFirstNameEn":  resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":   resolvedBy.LastNameEn,
			"ResolvedOn":             time.Now().Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%s/compliance/controls/%d/action-plans?action_plan_id=%d", appURL, controlID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/control_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = casus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingControlActionPlan, err := casus.controlActionPlanSvc.GetByID(controlID, actionPlanID)
		if err != nil {
			return
		}

		err = casus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.ControlActionPlanStatusUpdate{}.TableName(),
			id,
			existingControlActionPlan,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ControlActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
