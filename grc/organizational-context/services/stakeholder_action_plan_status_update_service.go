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
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderActionPlanStatusUpdateServiceInterface interface {
	GetAll(StakeholderID uint, StakeholderActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderActionPlanStatusUpdateListDTO, error)
	Create(StakeholderID uint, StakeholderActionPlanID uint, currentUserID uint, stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO) (uint, error)
	Update(StakeholderID uint, StakeholderActionPlanID uint, id uint, currentUserID uint, stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO) error
	GetLatestStakeholderActionPlanStatusUpdate(stakeholderID uint, actionPlanID uint) (*organizationalContextDTOs.StakeholderActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(stakeholderID uint, actionPlanID uint, id uint, token string) (string, error)

	Approve(StakeholderID uint, StakeholderActionPlanID uint, currentUserID uint, stakeholderActionPlanReviewUserStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(StakeholderID uint, StakeholderActionPlanID uint, currentUserID uint, stakeholderActionPlanReviewUserStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanReviewUserStatusUpdateFormDTO) error
}

type StakeholderActionPlanStatusUpdateService struct {
	repo                     repositories.StakeholderActionPlanStatusUpdateRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	stakeholderSvc           StakeholderServiceInterface
	stakeholderActionPlanSvc StakeholderActionPlanServiceInterface
	actionPlanStatusSvc      masterServices.ActionPlanStatusServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	notificationSvc          commonServices.NotificationServiceInterface
	taskSvc                  taskServices.TaskServiceInterface
	taskStatusSvc            taskServices.TaskStatusServiceInterface
}

func NewStakeholderActionPlanStatusUpdateService(
	repo repositories.StakeholderActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	stakeholderSvc StakeholderServiceInterface,
	stakeholderActionPlanSvc StakeholderActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) StakeholderActionPlanStatusUpdateServiceInterface {
	return &StakeholderActionPlanStatusUpdateService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		stakeholderSvc:           stakeholderSvc,
		stakeholderActionPlanSvc: stakeholderActionPlanSvc,
		actionPlanStatusSvc:      actionPlanStatusSvc,
		userSvc:                  userSvc,
		notificationSvc:          notificationSvc,
		taskSvc:                  taskSvc,
		taskStatusSvc:            taskStatusSvc,
	}
}

func (casus *StakeholderActionPlanStatusUpdateService) GetAll(StakeholderID uint, StakeholderActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.stakeholderActionPlanSvc.GetByID(StakeholderID, StakeholderActionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(StakeholderActionPlanID, sort, sortBy, keyword)
}

func (casus StakeholderActionPlanStatusUpdateService) GetByID(StakeholderID uint, StakeholderActionPlanID uint, id uint) (*models.StakeholderActionPlanStatusUpdate, error) {
	if _, err := casus.stakeholderActionPlanSvc.GetByID(StakeholderID, StakeholderActionPlanID); err != nil {
		return nil, err
	}

	stakeholderActionPlanStatusUpdate, err := casus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return stakeholderActionPlanStatusUpdate, nil
}

func (casus *StakeholderActionPlanStatusUpdateService) Create(StakeholderID uint, StakeholderActionPlanID uint, currentUserID uint, stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.stakeholderActionPlanSvc.GetByID(StakeholderID, StakeholderActionPlanID); err != nil {
		return 0, err
	}

	if err := casus.validateStakeholderActionPlanStatusUpdateFormDTO(stakeholderActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	stakeholderActionPlanUpdate := &models.StakeholderActionPlanStatusUpdate{
		StakeholderActionPlanID: StakeholderActionPlanID,
		ActionPlanStatusID:      stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                 stakeholderActionPlanStatusUpdateFormDTO.Comment,
		Percentage:              stakeholderActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:             currentUserID,
	}

	id, err := casus.repo.Create(stakeholderActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = casus.stakeholderActionPlanSvc.StatusUpdate(StakeholderID, StakeholderActionPlanID, currentUserID, stakeholderActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.StakeholderActionPlanStatusUpdateDocument

	for _, doc := range stakeholderActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("organizational-context", "stakeholder", strconv.FormatUint(uint64(StakeholderID), 10), "action-plans", strconv.FormatUint(uint64(StakeholderActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.StakeholderActionPlanStatusUpdateDocument{
			StakeholderActionPlanStatusUpdateID: id,
			URL:                                 url,
			Name:                                doc.FileName,
			Size:                                doc.Size,
			Token:                               encryptedFilename,
			OriginalName:                        doc.OriginalName,
			Type:                                doc.Type,
			CreatedByID:                         currentUserID,
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
			storagePath := filepath.Join("organizational-context", "stakeholder", strconv.FormatUint(uint64(StakeholderID), 10), "action-plans", strconv.FormatUint(uint64(StakeholderActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	stakeholder, err := casus.stakeholderSvc.GetByID(StakeholderID)
	if err != nil {
		return 0, err
	}

	stakeholderActionPlan, err := casus.stakeholderActionPlanSvc.GetByID(StakeholderID, StakeholderActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := casus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && stakeholderActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, StakeholderID, StakeholderActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:        "stakeholder_analysis",
			StakeholderActionPlanID: &StakeholderActionPlanID,
			Title:                   stakeholderActionPlan.Title,
			Description:             stakeholderActionPlan.Description,
			URLParam:                &urlParam,
			QueryParam:              nil,
			SeverityLevelID:         stakeholderActionPlan.SeverityLevelID,
			ResponsibleUserID:       *stakeholderActionPlan.ReviewUserID,
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
			models.StakeholderActionPlanStatusUpdate{}.TableName(),
			id,
			stakeholderActionPlanUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && stakeholderActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, StakeholderID, StakeholderActionPlanID)
			queryParam := ""
			entity := "stakeholder-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*stakeholderActionPlan.ReviewUserID},
				ModuleIdentifier:     "stakeholder_analysis",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Stakeholder Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة عمل أصحاب المصلحة",
				TitleEn:              stakeholderActionPlan.Title,
				TitleAr:              stakeholderActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := casus.userSvc.GetByID(*stakeholderActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := casus.userSvc.GetByID(stakeholderActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Stakeholder Action Plan %s: Resolved", stakeholderActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           reviewUser.FirstNameEn,
				"LastNameEn":            reviewUser.LastNameEn,
				"StakeholderTitle":      stakeholder.TitleEn,
				"ActionPlanTitle":       stakeholderActionPlan.Title,
				"ActionPlanDescription": stakeholderActionPlan.Description,
				"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
				"ResolvedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%sstakeholder/%d/action-plans?action_plan_id=%d", appURL, StakeholderID, StakeholderActionPlanID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organizational-context/templates/email/stakeholder_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = casus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (casus *StakeholderActionPlanStatusUpdateService) Update(StakeholderID uint, StakeholderActionPlanID uint, id uint, currentUserID uint, stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO) error {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.stakeholderActionPlanSvc.GetByID(StakeholderID, StakeholderActionPlanID); err != nil {
		return err
	}

	oldStakeholderActionPlanStatusUpdate, err := casus.GetByID(StakeholderID, StakeholderActionPlanID, id)
	if err != nil {
		return err
	}

	existingStakeholderActionPlanStatusUpdate, err := casus.GetByID(StakeholderID, StakeholderActionPlanID, id)
	if err != nil {
		return err
	}

	if err := casus.validateStakeholderActionPlanStatusUpdateFormDTO(stakeholderActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingStakeholderActionPlanStatusUpdate.ActionPlanStatusID = stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingStakeholderActionPlanStatusUpdate.Comment = stakeholderActionPlanStatusUpdateFormDTO.Comment
	existingStakeholderActionPlanStatusUpdate.Percentage = stakeholderActionPlanStatusUpdateFormDTO.Percentage
	existingStakeholderActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = casus.repo.Update(existingStakeholderActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = casus.stakeholderActionPlanSvc.StatusUpdate(StakeholderID, StakeholderActionPlanID, currentUserID, stakeholderActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.StakeholderActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.StakeholderActionPlanStatusUpdateDocument

	for _, doc := range stakeholderActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("organizational-context", "stakeholder", strconv.FormatUint(uint64(StakeholderID), 10), "action-plans", strconv.FormatUint(uint64(StakeholderActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new Control test Document instance
			document := &models.StakeholderActionPlanStatusUpdateDocument{
				StakeholderActionPlanStatusUpdateID: id,
				URL:                                 url,
				Name:                                doc.FileName,
				Size:                                doc.Size,
				Token:                               encryptedFileName,
				OriginalName:                        doc.OriginalName,
				Type:                                doc.Type,
				CreatedByID:                         currentUserID,
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
		storagePath := filepath.Join("organizational-context", "stakeholder", strconv.FormatUint(uint64(StakeholderID), 10), "action-plans", strconv.FormatUint(uint64(StakeholderActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("organizational-context", "stakeholder", strconv.FormatUint(uint64(StakeholderID), 10), "action-plans", strconv.FormatUint(uint64(StakeholderActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := casus.GetByID(StakeholderID, StakeholderActionPlanID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldStakeholderActionPlanStatusUpdate.Comment,
			"percentage": oldStakeholderActionPlanStatusUpdate.Percentage,
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
			err := casus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderActionPlanStatusUpdate{}.TableName(),
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

func (casus *StakeholderActionPlanStatusUpdateService) validateStakeholderActionPlanStatusUpdateFormDTO(stakeholderActionPlanStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanStatusUpdateFormDTO) error {
	if _, err := casus.actionPlanStatusSvc.GetByID(stakeholderActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (capsus StakeholderActionPlanStatusUpdateService) GetLatestStakeholderActionPlanStatusUpdate(stakeholderID uint, actionPlanID uint) (*organizationalContextDTOs.StakeholderActionPlanStatusUpdateListDTO, error) {
	if _, err := capsus.stakeholderActionPlanSvc.GetByID(stakeholderID, actionPlanID); err != nil {
		return nil, err
	}

	return capsus.repo.GetLatestStakeholderActionPlanStatusUpdate(actionPlanID)
}

func (capsus *StakeholderActionPlanStatusUpdateService) GetDocumentPath(StakeholderID, ActionPlanID, id uint, token string) (string, error) {
	existingStakeholderActionplanStatusUpdate, err := capsus.GetByID(StakeholderID, ActionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingStakeholderActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the stakeholder action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("organizational-context", "stakeholder", strconv.FormatUint(uint64(StakeholderID), 10), "action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (capus *StakeholderActionPlanStatusUpdateService) Approve(stakeholderID uint, stakeholderActionPlanID uint, currentUserID uint, stakeholderActionPlanReviewUserStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanReviewUserStatusUpdateFormDTO) error {
	stakeholder, err := capus.stakeholderSvc.GetByID(stakeholderID)
	if err != nil {
		return err
	}
	existingStakeholderActionPlan, err := capus.stakeholderActionPlanSvc.GetByID(stakeholderID, stakeholderActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, stakeholderActionPlanID)
	}

	// if currentUserID != *existingStakeholderActionPlan.ReviewUserID {
	// 	return utility.NewForbiddenError("Permission denied.")
	// }

	stakeholderActionPlanStatusResolved, err := capus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	stakeholderActionPlanClosedStatus, err := capus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := capus.repo.BeginTransaction()

	if existingStakeholderActionPlan.ActionPlanStatusID != stakeholderActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The stakeholder action plan status must be resolved.")
	}

	existingStakeholderActionPlan.ActionPlanStatusID = stakeholderActionPlanClosedStatus.ID

	if err := capus.repo.StakeholderActionPlanUpdateStatus(existingStakeholderActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	stakeholderActionPlanUpdate := &models.StakeholderActionPlanStatusUpdate{
		StakeholderActionPlanID: stakeholderActionPlanID,
		ActionPlanStatusID:      stakeholderActionPlanClosedStatus.ID,
		Percentage:              nil,
		Comment:                 stakeholderActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:             currentUserID,
	}

	id, err := capus.repo.Create(stakeholderActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := capus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = capus.taskSvc.UpdateActionPlanTaskStatus("stakeholder_action_plan_id", stakeholderActionPlanID, *existingStakeholderActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Start a new transaction for the second update
	tx2 := capus.repo.BeginTransaction()
	_, err = capus.taskSvc.UpdateActionPlanTaskStatus("stakeholder_action_plan_id", stakeholderActionPlanID, existingStakeholderActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx2)
	if err != nil {
		tx2.Rollback()
		return err
	}

	tx2.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, stakeholderID, stakeholderActionPlanID)
		queryParam := ""
		entity := "stakeholder-action-plan"
		action := "close-stakeholder-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingStakeholderActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "stakeholder_analysis",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Stakeholder Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل أصحاب المصلحة",
			TitleEn:              existingStakeholderActionPlan.Title,
			TitleAr:              existingStakeholderActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := capus.userSvc.GetByID(existingStakeholderActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := capus.userSvc.GetByID(*existingStakeholderActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Stakeholder Action Plan %s: Resolved", existingStakeholderActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"StakeholderTitle":      stakeholder.TitleEn,
			"ActionPlanTitle":       existingStakeholderActionPlan.Title,
			"ActionPlanDescription": existingStakeholderActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%sstakeholder/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, stakeholderActionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organizational-context/templates/email/stakeholder_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = capus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingStakeholderActionPlan, err := capus.stakeholderActionPlanSvc.GetByID(stakeholderID, stakeholderActionPlanID)
		if err != nil {
			return
		}

		err = capus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.StakeholderActionPlanStatusUpdate{}.TableName(),
			id,
			existingStakeholderActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.StakeholderActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (capus *StakeholderActionPlanStatusUpdateService) Revert(stakeholderID uint, stakeholderActionPlanID uint, currentUserID uint, stakeholderActionPlanReviewUserStatusUpdateFormDTO *organizationalContextDTOs.StakeholderActionPlanReviewUserStatusUpdateFormDTO) error {
	stakeholder, err := capus.stakeholderSvc.GetByID(stakeholderID)
	if err != nil {
		return err
	}
	existingStakeholderActionPlan, err := capus.stakeholderActionPlanSvc.GetByID(stakeholderID, stakeholderActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, stakeholderActionPlanID)
	}

	if currentUserID != *existingStakeholderActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	stakeholderActionPlanStatusResolved, err := capus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	stakeholderActionPlanReopenStatus, err := capus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := capus.repo.BeginTransaction()

	if existingStakeholderActionPlan.ActionPlanStatusID != stakeholderActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The stakeholder action plan status must be resolved.")
	}

	existingStakeholderActionPlan.ActionPlanStatusID = stakeholderActionPlanReopenStatus.ID

	if err := capus.repo.StakeholderActionPlanUpdateStatus(existingStakeholderActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	stakeholderActionPlanUpdate := &models.StakeholderActionPlanStatusUpdate{
		StakeholderActionPlanID: stakeholderActionPlanID,
		ActionPlanStatusID:      stakeholderActionPlanReopenStatus.ID,
		Comment:                 stakeholderActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:             currentUserID,
	}

	id, err := capus.repo.Create(stakeholderActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := capus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = capus.taskSvc.UpdateActionPlanTaskStatus("stakeholder_action_plan_id", stakeholderActionPlanID, *existingStakeholderActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"stakeholder_id":%d, "action_plan_id":%d}`, stakeholderID, stakeholderActionPlanID)
		queryParam := ""
		entity := "stakeholder-action-plan"
		action := "reopen-stakeholder-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingStakeholderActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "stakeholder_analysis",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Stakeholder Action Plan Status Update",
			HeadingAr:            "تحديث حالة خطة عمل أصحاب المصلحة",
			TitleEn:              existingStakeholderActionPlan.Title,
			TitleAr:              existingStakeholderActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := capus.userSvc.GetByID(existingStakeholderActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := capus.userSvc.GetByID(*existingStakeholderActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Stakeholder Action Plan %s: Reopened", existingStakeholderActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"StakeholderTitle":      stakeholder.TitleEn,
			"ActionPlanTitle":       existingStakeholderActionPlan.Title,
			"ActionPlanDescription": existingStakeholderActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%sstakeholder/%d/action-plans?action_plan_id=%d", appURL, stakeholderID, stakeholderActionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organizational-context/templates/email/stakeholder_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = capus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "stakeholder_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingStakeholderActionPlan, err := capus.stakeholderActionPlanSvc.GetByID(stakeholderID, stakeholderActionPlanID)
		if err != nil {
			return
		}

		err = capus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.StakeholderActionPlanStatusUpdate{}.TableName(),
			id,
			existingStakeholderActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.StakeholderActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
