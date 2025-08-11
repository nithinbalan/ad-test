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
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	settingServices "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskActionPlanStatusUpdateServiceInterface interface {
	GetAll(RiskID uint, RiskActionPlanID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.RiskActionPlanStatusUpdateListDTO, error)
	Create(RiskID uint, RiskActionPlanID uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) (uint, error)
	Update(RiskID uint, RiskActionPlanID uint, id uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) error
	GetLatestRiskActionPlanStatusUpdate(riskID uint, actionPlanID uint) (*riskDTOs.RiskActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(riskID, actionPlanID, statusUpdateID uint, token string) (string, error)

	Approve(riskID uint, RiskActionPlanID uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(riskID uint, RiskActionPlanID uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanReviewUserStatusUpdateFormDTO) error
}

type RiskActionPlanStatusUpdateService struct {
	repo                repositories.RiskActionPlanStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	riskRegisterSvc     RiskRegisterServiceInterface
	riskActionPlanSvc   RiskActionPlanServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	moduleSvc           settingServices.ModuleServiceInterface
	userSvc             organizationServices.UserServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	taskSvc             taskServices.TaskServiceInterface
	taskStatusSvc       taskServices.TaskStatusServiceInterface
}

func NewRiskActionPlanStatusUpdateService(
	repo repositories.RiskActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskActionPlanSvc RiskActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	moduleSvc settingServices.ModuleServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) RiskActionPlanStatusUpdateServiceInterface {
	return &RiskActionPlanStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		riskRegisterSvc:     riskRegisterSvc,
		riskActionPlanSvc:   riskActionPlanSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
		moduleSvc:           moduleSvc,
		userSvc:             userSvc,
		notificationSvc:     notificationSvc,
		taskSvc:             taskSvc,
		taskStatusSvc:       taskStatusSvc,
	}
}

func (rasus *RiskActionPlanStatusUpdateService) GetAll(RiskID uint, RiskActionPlanID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.RiskActionPlanStatusUpdateListDTO, error) {
	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return nil, err
	}

	return rasus.repo.GetAll(RiskActionPlanID, sort, sortBy, keyword)
}

func (rasus RiskActionPlanStatusUpdateService) GetByID(RiskID uint, RiskActionPlanID uint, id uint) (*models.RiskActionPlanStatusUpdate, error) {
	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return nil, err
	}

	riskActionPlanStatusUpdate, err := rasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlanStatusUpdate, nil
}

func (rasus *RiskActionPlanStatusUpdateService) Create(RiskID uint, RiskActionPlanID uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := rasus.repo.BeginTransaction()

	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return 0, err
	}

	if err := rasus.validateRiskActionPlanStatusUpdateFormDTO(riskActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   RiskActionPlanID,
		ActionPlanStatusID: riskActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:            riskActionPlanStatusUpdateFormDTO.Comment,
		Percentage:         riskActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err := rasus.repo.Create(riskActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = rasus.riskActionPlanSvc.StatusUpdate(RiskID, RiskActionPlanID, currentUserID, riskActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskActionPlanStatusUpdateDocument

	for _, doc := range riskActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.RiskActionPlanStatusUpdateDocument{
			RiskActionPlanStatusUpdateID: id,
			URL:                          url,
			Name:                         doc.FileName,
			Size:                         doc.Size,
			Token:                        encryptedFilename,
			OriginalName:                 doc.OriginalName,
			Type:                         doc.Type,
			CreatedByID:                  currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := rasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		err := rasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			riskActionPlanUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (rasus *RiskActionPlanStatusUpdateService) Update(RiskID uint, RiskActionPlanID uint, id uint, currentUserID uint, RiskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) error {
	tx := rasus.repo.BeginTransaction()

	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return err
	}

	existingRiskActionPlanStatusUpdate, err := rasus.GetByID(RiskID, RiskActionPlanID, id)
	if err != nil {
		return err
	}

	oldActionPlanStatusUpdate, err := rasus.GetByID(RiskID, RiskActionPlanID, id)
	if err != nil {
		return err
	}

	if err := rasus.validateRiskActionPlanStatusUpdateFormDTO(RiskActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Risk StatusUpdate-- only comments and percentage
	existingRiskActionPlanStatusUpdate.Comment = RiskActionPlanStatusUpdateFormDTO.Comment
	existingRiskActionPlanStatusUpdate.Percentage = RiskActionPlanStatusUpdateFormDTO.Percentage
	existingRiskActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = rasus.repo.Update(existingRiskActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.RiskActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskActionPlanStatusUpdateDocument

	for _, doc := range RiskActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.RiskActionPlanStatusUpdateDocument{
				RiskActionPlanStatusUpdateID: id,
				URL:                          url,
				Name:                         doc.FileName,
				Size:                         doc.Size,
				Token:                        encryptedFileName,
				OriginalName:                 doc.OriginalName,
				Type:                         doc.Type,
				CreatedByID:                  currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := rasus.repo.GetDocumentByID(*doc.ID)
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
		if err := rasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlanStatusUpdate, err := rasus.GetByID(RiskID, RiskActionPlanID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
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
			err := rasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}
	}()

	return nil
}

func (rasus *RiskActionPlanStatusUpdateService) validateRiskActionPlanStatusUpdateFormDTO(riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) error {
	if _, err := rasus.actionPlanStatusSvc.GetByID(riskActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (rasus RiskActionPlanStatusUpdateService) GetLatestRiskActionPlanStatusUpdate(riskID uint, actionPlanID uint) (*riskDTOs.RiskActionPlanStatusUpdateListDTO, error) {
	if _, err := rasus.riskActionPlanSvc.GetByID(riskID, actionPlanID); err != nil {
		return nil, err
	}

	return rasus.repo.GetLatestRiskActionPlanStatusUpdate(actionPlanID)
}

func (rasus *RiskActionPlanStatusUpdateService) GetDocumentPath(riskID, actionPlanID, statusUpdateID uint, token string) (string, error) {
	// Retrieve the specific status update to check if documents exist
	statusUpdate, err := rasus.GetByID(riskID, actionPlanID, statusUpdateID)
	if err != nil {
		fmt.Printf("Error fetching status update by ID: %v\n", err)
		return "", err
	}

	// Check if documents are associated with the status update
	if len(statusUpdate.Documents) == 0 {
		fmt.Println("No documents associated with the status update")
		return "", fmt.Errorf("no documents found for the risk action plan update")
	}

	// Decrypt the token to get the filename
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		fmt.Printf("Token decryption failed: %v\n", err)
		return "", fmt.Errorf("resource not found")
	}
	decryptedFilename := string(decryptedFilenameBytes)
	fmt.Printf("Decrypted filename: %s\n", decryptedFilename)

	// Construct the storage path
	storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(statusUpdateID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	fmt.Printf("Constructed file path: %s\n", fullPath)

	// Verify if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("File does not exist at path: %s\n", fullPath)
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		fmt.Printf("Error accessing file at path %s: %v\n", fullPath, err)
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (rasus *RiskActionPlanStatusUpdateService) Approve(riskID uint, riskActionPlanID uint, currentUserID uint, riskActionPlanReviewUserStatusUpdateFormDTO *riskDTOs.RiskActionPlanReviewUserStatusUpdateFormDTO) error {
	existingRiskActionPlan, err := rasus.riskActionPlanSvc.GetByID(riskID, riskActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, riskActionPlanID)
	}

	if currentUserID != *existingRiskActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	riskActionPlanStatusResolved, err := rasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	riskActionPlanClosedStatus, err := rasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	riskRegister, err := rasus.riskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	tx := rasus.repo.BeginTransaction()

	if existingRiskActionPlan.ActionPlanStatusID != riskActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The risk action plan status must be resolved.")
	}

	existingRiskActionPlan.ActionPlanStatusID = riskActionPlanClosedStatus.ID

	if err := rasus.repo.RiskActionPlanUpdateStatus(existingRiskActionPlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   riskActionPlanID,
		ActionPlanStatusID: riskActionPlanClosedStatus.ID,
		Percentage:         nil,
		Comment:            riskActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := rasus.repo.Create(riskActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := rasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = rasus.taskSvc.UpdateActionPlanTaskStatus("risk_action_plan_id", riskActionPlanID, existingRiskActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = rasus.taskSvc.UpdateActionPlanTaskStatus("risk_action_plan_id", riskActionPlanID, *existingRiskActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, riskActionPlanID)
		queryParam := ""
		entity := "risk-action-plan"
		action := "close-risk-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingRiskActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Action Plan Status Update",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
			TitleEn:              existingRiskActionPlan.Title,
			TitleAr:              existingRiskActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := rasus.userSvc.GetByID(existingRiskActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := rasus.userSvc.GetByID(*existingRiskActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Risk Action Plan %s: Closed", existingRiskActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"RiskRegister":          riskRegister.Title,
			"ActionPlanTitle":       existingRiskActionPlan.Title,
			"ActionPlanDescription": existingRiskActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/risk-register/%d/risk-response?action_plan_id=%d", appURL, riskID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = rasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingRiskActionPlan, err := rasus.riskActionPlanSvc.GetByID(riskID, riskActionPlanID)
		if err != nil {
			return
		}

		err = rasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			existingRiskActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (rasus *RiskActionPlanStatusUpdateService) Revert(riskID uint, riskActionPlanID uint, currentUserID uint, riskActionPlanReviewUserStatusUpdateFormDTO *riskDTOs.RiskActionPlanReviewUserStatusUpdateFormDTO) error {
	existingRiskActionPlan, err := rasus.riskActionPlanSvc.GetByID(riskID, riskActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, riskActionPlanID)
	}

	if currentUserID != *existingRiskActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	riskActionPlanStatusResolved, err := rasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	riskActionPlanReopenStatus, err := rasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	riskRegister, err := rasus.riskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	tx := rasus.repo.BeginTransaction()

	if existingRiskActionPlan.ActionPlanStatusID != riskActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The risk action plan status must be resolved.")
	}

	existingRiskActionPlan.ActionPlanStatusID = riskActionPlanReopenStatus.ID

	if err := rasus.repo.RiskActionPlanUpdateStatus(existingRiskActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   riskActionPlanID,
		ActionPlanStatusID: riskActionPlanReopenStatus.ID,
		Comment:            riskActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := rasus.repo.Create(riskActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := rasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = rasus.taskSvc.UpdateActionPlanTaskStatus("risk_action_plan_id", riskActionPlanID, *existingRiskActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, riskActionPlanID)
		queryParam := ""
		entity := "risk-action-plan"
		action := "reopen-risk-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingRiskActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Action Plan Status Update",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
			TitleEn:              existingRiskActionPlan.Title,
			TitleAr:              existingRiskActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := rasus.userSvc.GetByID(existingRiskActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := rasus.userSvc.GetByID(*existingRiskActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Risk Action Plan %s: Reopened", existingRiskActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"RiskRegister":          riskRegister.Title,
			"ActionPlanTitle":       existingRiskActionPlan.Title,
			"ActionPlanDescription": existingRiskActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/risk-register/%d/risk-response?action_plan_id=%d", appURL, riskID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = rasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingRiskActionPlan, err := rasus.riskActionPlanSvc.GetByID(riskID, riskActionPlanID)
		if err != nil {
			return
		}

		err = rasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			existingRiskActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
