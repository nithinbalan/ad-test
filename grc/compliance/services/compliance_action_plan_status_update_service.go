package services

import (
	"errors"
	"fmt"
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
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ComplianceActionPlanStatusUpdateServiceInterface interface {
	GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error
	GetLatestComplianceActionPlanStatusUpdate(complianceID uint, actionPlanID uint) (*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(complianceID uint, actionPlanID uint, id uint, token string) (string, error)

	Approve(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error
}

type ComplianceActionPlanStatusUpdateService struct {
	repo                    repositories.ComplianceActionPlanStatusUpdateRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceRegisterSvc   ComplianceRegisterServiceInterface
	complianceActionPlanSvc ComplianceActionPlanServiceInterface
	actionPlanStatusSvc     masterServices.ActionPlanStatusServiceInterface
	userSvc                 organizationServices.UserServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
	taskSvc                 taskServices.TaskServiceInterface
	taskStatusSvc           taskServices.TaskStatusServiceInterface
}

func NewComplianceActionPlanStatusUpdateService(
	repo repositories.ComplianceActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceActionPlanSvc ComplianceActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) ComplianceActionPlanStatusUpdateServiceInterface {
	return &ComplianceActionPlanStatusUpdateService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		complianceRegisterSvc:   complianceRegisterSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
		actionPlanStatusSvc:     actionPlanStatusSvc,
		userSvc:                 userSvc,
		notificationSvc:         notificationSvc,
		taskSvc:                 taskSvc,
		taskStatusSvc:           taskStatusSvc,
	}
}

func (casus *ComplianceActionPlanStatusUpdateService) GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(ComplianceActionPlanID, sort, sortBy, keyword)
}

func (casus ComplianceActionPlanStatusUpdateService) GetByID(ComplianceID uint, ComplianceActionPlanID uint, id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	complianceActionPlanStatusUpdate, err := casus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceActionPlanStatusUpdate, nil
}

func (casus *ComplianceActionPlanStatusUpdateService) Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return 0, err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: ComplianceActionPlanID,
		ActionPlanStatusID:     complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                complianceActionPlanStatusUpdateFormDTO.Comment,
		Percentage:             complianceActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:            currentUserID,
	}

	id, err := casus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ComplianceActionPlanStatusUpdateDocument

	for _, doc := range complianceActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ComplianceActionPlanStatusUpdateDocument{
			ComplianceActionPlanStatusUpdateID: id,
			URL:                                url,
			Name:                               doc.FileName,
			Size:                               doc.Size,
			Token:                              encryptedFilename,
			OriginalName:                       doc.OriginalName,
			Type:                               doc.Type,
			CreatedByID:                        currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := casus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	compliance, err := casus.complianceRegisterSvc.GetByID(ComplianceID)
	if err != nil {
		return 0, err
	}

	complianceActionPlan, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := casus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && complianceActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, ComplianceID, ComplianceActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:       "compliance",
			ComplianceActionPlanID: &ComplianceActionPlanID,
			Title:                  complianceActionPlan.Title,
			Description:            complianceActionPlan.Description,
			URLParam:               &urlParam,
			QueryParam:             nil,
			SeverityLevelID:        complianceActionPlan.SeverityLevelID,
			ResponsibleUserID:      *complianceActionPlan.ReviewUserID,
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
			models.ComplianceActionPlanStatusUpdate{}.TableName(),
			id,
			complianceActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && complianceActionPlan.ReviewUserID != nil {
			// prepare notification
			urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, ComplianceID, ComplianceActionPlanID)
			queryParam := ""
			entity := "compliance-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               *complianceActionPlan.ReviewUserID,
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Action Plan Status Update",
				HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
				TitleEn:              complianceActionPlan.Title,
				TitleAr:              complianceActionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := casus.userSvc.GetByID(*complianceActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := casus.userSvc.GetByID(complianceActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Compliance Action Plan %s: Resolved", complianceActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":               reviewUser.FirstNameEn,
				"LastNameEn":                reviewUser.LastNameEn,
				"ComplianceReferenceNumber": compliance.ReferenceNumber,
				"ComplianceTitle":           compliance.TitleEn,
				"ActionPlanTitle":           complianceActionPlan.Title,
				"ActionPlanDescription":     complianceActionPlan.Description,
				"ResolvedByFirstNameEn":     resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":      resolvedBy.LastNameEn,
				"ResolvedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, ComplianceID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = casus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (casus *ComplianceActionPlanStatusUpdateService) Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return err
	}

	oldComplianceActionPlanStatusUpdate, err := casus.GetByID(ComplianceID, ComplianceActionPlanID, id)

	if err != nil {
		return err
	}

	existingComplianceActionPlanStatusUpdate, err := casus.GetByID(ComplianceID, ComplianceActionPlanID, id)

	if err != nil {
		return err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingComplianceActionPlanStatusUpdate.ActionPlanStatusID = complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingComplianceActionPlanStatusUpdate.Comment = complianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlanStatusUpdate.Percentage = complianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = casus.repo.Update(existingComplianceActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ComplianceActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ComplianceActionPlanStatusUpdateDocument

	for _, doc := range complianceActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new Control test Document instance
			document := &models.ComplianceActionPlanStatusUpdateDocument{
				ComplianceActionPlanStatusUpdateID: id,
				URL:                                url,
				Name:                               doc.FileName,
				Size:                               doc.Size,
				Token:                              encryptedFileName,
				OriginalName:                       doc.OriginalName,
				Type:                               doc.Type,
				CreatedByID:                        currentUserID,
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
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := casus.GetByID(ComplianceID, ComplianceActionPlanID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldComplianceActionPlanStatusUpdate.Comment,
			"percentage": oldComplianceActionPlanStatusUpdate.Percentage,
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
				models.ComplianceActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceActionPlanStatusUpdate{}.TableName(),
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

func (casus *ComplianceActionPlanStatusUpdateService) validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {

	if _, err := casus.actionPlanStatusSvc.GetByID(complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (capsus ComplianceActionPlanStatusUpdateService) GetLatestComplianceActionPlanStatusUpdate(complianceID uint, actionPlanID uint) (*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {
	if _, err := capsus.complianceActionPlanSvc.GetByID(complianceID, actionPlanID); err != nil {
		return nil, err
	}

	return capsus.repo.GetLatestComplianceActionPlanStatusUpdate(actionPlanID)
}

func (capsus *ComplianceActionPlanStatusUpdateService) GetDocumentPath(ComplianceID, ActionPlanID, id uint, token string) (string, error) {

	existingComplianceActionplanStatusUpdate, err := capsus.GetByID(ComplianceID, ActionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingComplianceActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the compliance action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (capus *ComplianceActionPlanStatusUpdateService) Approve(complianceID uint, compliaceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error {
	compliance, err := capus.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return err
	}
	existingComplianceActionPlan, err := capus.complianceActionPlanSvc.GetByID(complianceID, compliaceActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, compliaceActionPlanID)
	}

	if currentUserID != *existingComplianceActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	complianceActionPlanStatusResolved, err := capus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	complianceActionPlanClosedStatus, err := capus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := capus.repo.BeginTransaction()

	if existingComplianceActionPlan.ActionPlanStatusID != complianceActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The compliance action plan status must be resolved.")
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanClosedStatus.ID

	if err := capus.repo.ComplianceActionPlanUpdateStatus(existingComplianceActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: compliaceActionPlanID,
		ActionPlanStatusID:     complianceActionPlanClosedStatus.ID,
		Percentage:             nil,
		Comment:                complianceActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:            currentUserID,
	}

	id, err := capus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := capus.taskStatusSvc.GetByType("closed")
	if err != nil {
		return err
	}

	_, err = capus.taskSvc.UpdateActionPlanTaskStatus("compliance_action_plan_id", compliaceActionPlanID, *existingComplianceActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = capus.taskSvc.UpdateActionPlanTaskStatus("compliance_action_plan_id", compliaceActionPlanID, existingComplianceActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {

		urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, compliaceActionPlanID)
		queryParam := ""
		entity := "compliance-action-plan"
		action := "close-compliance-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               existingComplianceActionPlan.ResponsibleUserID,
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Action Plan Status Update",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
			TitleEn:              existingComplianceActionPlan.Title,
			TitleAr:              existingComplianceActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := capus.userSvc.GetByID(existingComplianceActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := capus.userSvc.GetByID(*existingComplianceActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Action Plan %s: Resolved", existingComplianceActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":               responsibleUser.FirstNameEn,
			"LastNameEn":                responsibleUser.LastNameEn,
			"ComplianceReferenceNumber": compliance.ReferenceNumber,
			"ComplianceTitle":           compliance.TitleEn,
			"ActionPlanTitle":           existingComplianceActionPlan.Title,
			"ActionPlanDescription":     existingComplianceActionPlan.Description,
			"ResolvedByFirstNameEn":     resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":      resolvedBy.LastNameEn,
			"ResolvedOn":                time.Now().Format("2006-01-02"),
			"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = capus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingComplianceActionPlan, err := capus.complianceActionPlanSvc.GetByID(complianceID, compliaceActionPlanID)
		if err != nil {
			return
		}

		err = capus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.ComplianceActionPlanStatusUpdate{}.TableName(),
			id,
			existingComplianceActionPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ComplianceActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (capus *ComplianceActionPlanStatusUpdateService) Revert(complianceID uint, compliaceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error {
	compliance, err := capus.complianceRegisterSvc.GetByID(complianceID)
	if err != nil {
		return err
	}
	existingComplianceActionPlan, err := capus.complianceActionPlanSvc.GetByID(complianceID, compliaceActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, compliaceActionPlanID)
	}

	if currentUserID != *existingComplianceActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	complianceActionPlanStatusResolved, err := capus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	complianceActionPlanReopenStatus, err := capus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := capus.repo.BeginTransaction()

	if existingComplianceActionPlan.ActionPlanStatusID != complianceActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The compliance action plan status must be resolved.")
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanReopenStatus.ID

	if err := capus.repo.ComplianceActionPlanUpdateStatus(existingComplianceActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: compliaceActionPlanID,
		ActionPlanStatusID:     complianceActionPlanReopenStatus.ID,
		Comment:                complianceActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:            currentUserID,
	}

	id, err := capus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := capus.taskStatusSvc.GetByType("closed")
	if err != nil {
		return err
	}

	_, err = capus.taskSvc.UpdateActionPlanTaskStatus("compliance_action_plan_id", compliaceActionPlanID, *existingComplianceActionPlan.ReviewUserID, uint(taskStatus.ID))
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {

		urlParam := fmt.Sprintf(`{"compliance_register_id":%d, "action_plan_id":%d}`, complianceID, compliaceActionPlanID)
		queryParam := ""
		entity := "compliance-action-plan"
		action := "reopen-compliance-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               existingComplianceActionPlan.ResponsibleUserID,
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Action Plan Status Update",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل الامتثال",
			TitleEn:              existingComplianceActionPlan.Title,
			TitleAr:              existingComplianceActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := capus.userSvc.GetByID(existingComplianceActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := capus.userSvc.GetByID(*existingComplianceActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Action Plan %s: Reopened", existingComplianceActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":               responsibleUser.FirstNameEn,
			"LastNameEn":                responsibleUser.LastNameEn,
			"ComplianceReferenceNumber": compliance.ReferenceNumber,
			"ComplianceTitle":           compliance.TitleEn,
			"ActionPlanTitle":           existingComplianceActionPlan.Title,
			"ActionPlanDescription":     existingComplianceActionPlan.Description,
			"ResolvedByFirstNameEn":     resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":      resolvedBy.LastNameEn,
			"ResolvedOn":                time.Now().Format("2006-01-02"),
			"RedirectURL":               fmt.Sprintf("%scompliance/register/%d/action-plans?action_plan_id=%d", appURL, complianceID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = capus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingComplianceActionPlan, err := capus.complianceActionPlanSvc.GetByID(complianceID, compliaceActionPlanID)
		if err != nil {
			return
		}

		err = capus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.ComplianceActionPlanStatusUpdate{}.TableName(),
			id,
			existingComplianceActionPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ComplianceActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
