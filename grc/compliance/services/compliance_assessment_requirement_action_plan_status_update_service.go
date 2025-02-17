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

type ComplianceAssessmentRequirementActionPlanStatusUpdateServiceInterface interface {
	GetAll(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	Create(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error)
	Update(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error
	GetLatestComplianceAssessmentRequirementActionPlanStatusUpdate(assessmentID uint, requirementID uint, actionPlanID uint) (*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(assessmentID uint, requirementID uint, actionPlanID uint, id uint, token string) (string, error)

	Approve(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error
}

type ComplianceAssessmentRequirementActionPlanStatusUpdateService struct {
	repo                    repositories.ComplianceActionPlanStatusUpdateRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceAssessmentSvc ComplianceAssessmentServiceInterface
	complianceActionPlanSvc ComplianceAssessmentRequirementActionPlanServiceInterface
	actionPlanStatusSvc     masterServices.ActionPlanStatusServiceInterface
	userSvc                 organizationServices.UserServiceInterface
	notificationSvc         commonServices.NotificationServiceInterface
	taskSvc                 taskServices.TaskServiceInterface
	taskStatusSvc           taskServices.TaskStatusServiceInterface
}

func NewComplianceAssessmentRequirementActionPlanStatusUpdateService(
	repo repositories.ComplianceActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
	complianceActionPlanSvc ComplianceAssessmentRequirementActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) ComplianceAssessmentRequirementActionPlanStatusUpdateServiceInterface {
	return &ComplianceAssessmentRequirementActionPlanStatusUpdateService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		complianceAssessmentSvc: complianceAssessmentSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
		actionPlanStatusSvc:     actionPlanStatusSvc,
		userSvc:                 userSvc,
		notificationSvc:         notificationSvc,
		taskSvc:                 taskSvc,
		taskStatusSvc:           taskStatusSvc,
	}
}

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) GetAll(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {
	if _, err := carasus.complianceActionPlanSvc.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	return carasus.repo.GetAll(ComplianceActionPlanID, sort, sortBy, keyword)
}

func (carasus ComplianceAssessmentRequirementActionPlanStatusUpdateService) GetByID(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	if _, err := carasus.complianceActionPlanSvc.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	complianceActionPlanStatusUpdate, err := carasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceActionPlanStatusUpdate, nil
}

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) Create(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := carasus.repo.BeginTransaction()

	if _, err := carasus.complianceActionPlanSvc.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID); err != nil {
		return 0, err
	}

	if err := carasus.validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO); err != nil {
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

	id, err := carasus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = carasus.complianceActionPlanSvc.StatusUpdate(AssessmentID, RequirementID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO, tx)
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

		url := utility.GenerateDocumentURL("compliance", "compliance-assessment", strconv.FormatUint(uint64(AssessmentID), 10), "requirement", strconv.FormatUint(uint64(RequirementID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
		if err := carasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(AssessmentID), 10), "requirement", strconv.FormatUint(uint64(RequirementID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	complianceAssessment, err := carasus.complianceAssessmentSvc.GetFullDetailsByID(AssessmentID)
	if err != nil {
		return 0, err
	}

	complianceActionPlan, err := carasus.complianceActionPlanSvc.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := carasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && complianceActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, AssessmentID, RequirementID, ComplianceActionPlanID)
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

		if _, err := carasus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := carasus.activityLogSvc.LogActivity(
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
			urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, AssessmentID, RequirementID, ComplianceActionPlanID)
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

			reviewUser, err := carasus.userSvc.GetByID(*complianceActionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := carasus.userSvc.GetByID(complianceActionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Compliance Assessment Requirement Action Plan %s: Resolved", complianceActionPlan.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           reviewUser.FirstNameEn,
				"LastNameEn":            reviewUser.LastNameEn,
				"ComplianceAssessment":  complianceAssessment.ComplianceRegister.TitleEn,
				"ActionPlanTitle":       complianceActionPlan.Title,
				"ActionPlanDescription": complianceActionPlan.Description,
				"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
				"ResolvedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%scompliance/assessment/%d/action-plans?action_plan_id=%d", appURL, AssessmentID, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_assessment_action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = carasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) Update(AssessmentID uint, RequirementID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {
	tx := carasus.repo.BeginTransaction()

	if _, err := carasus.complianceActionPlanSvc.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID); err != nil {
		return err
	}

	oldComplianceActionPlanStatusUpdate, err := carasus.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID, id)

	if err != nil {
		return err
	}

	existingComplianceActionPlanStatusUpdate, err := carasus.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID, id)

	if err != nil {
		return err
	}

	if err := carasus.validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingComplianceActionPlanStatusUpdate.ActionPlanStatusID = complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingComplianceActionPlanStatusUpdate.Comment = complianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlanStatusUpdate.Percentage = complianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = carasus.repo.Update(existingComplianceActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = carasus.complianceActionPlanSvc.StatusUpdate(AssessmentID, RequirementID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO, tx)
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
			url := utility.GenerateDocumentURL("compliance", "compliance-assessment", strconv.FormatUint(uint64(AssessmentID), 10), "requirement", strconv.FormatUint(uint64(RequirementID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
				existingDocument, err := carasus.repo.GetDocumentByID(*doc.ID)
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
		if err := carasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := carasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(AssessmentID), 10), "requirement", strconv.FormatUint(uint64(RequirementID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(AssessmentID), 10), "requirement", strconv.FormatUint(uint64(RequirementID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := carasus.GetByID(AssessmentID, RequirementID, ComplianceActionPlanID, id)
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
			err := carasus.activityLogSvc.LogActivity(
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

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {

	if _, err := carasus.actionPlanStatusSvc.GetByID(complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) GetLatestComplianceAssessmentRequirementActionPlanStatusUpdate(assessmentId uint, requirementID uint, actionPlanID uint) (*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {
	if _, err := carasus.complianceActionPlanSvc.GetByID(assessmentId, requirementID, actionPlanID); err != nil {
		return nil, err
	}

	return carasus.repo.GetLatestComplianceActionPlanStatusUpdate(actionPlanID)
}

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) GetDocumentPath(AssessmentID, RequirementID, ActionPlanID, id uint, token string) (string, error) {

	existingComplianceActionplanStatusUpdate, err := carasus.GetByID(AssessmentID, RequirementID, ActionPlanID, id)
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

	storagePath := filepath.Join("compliance", "compliance-assessment", strconv.FormatUint(uint64(AssessmentID), 10), "requirement", strconv.FormatUint(uint64(RequirementID), 10), "action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) Approve(assessmentID uint, requirementID uint, compliaceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error {
	existingComplianceActionPlan, err := carasus.complianceActionPlanSvc.GetByID(assessmentID, requirementID, compliaceActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, compliaceActionPlanID)
	}

	if currentUserID != *existingComplianceActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	complianceActionPlanStatusResolved, err := carasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	complianceActionPlanClosedStatus, err := carasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	complianceAssessment, err := carasus.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}

	tx := carasus.repo.BeginTransaction()

	if existingComplianceActionPlan.ActionPlanStatusID != complianceActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The compliance action plan status must be resolved.")
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanClosedStatus.ID

	if err := carasus.repo.ComplianceActionPlanUpdateStatus(existingComplianceActionPlan); err != nil {
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

	id, err := carasus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := carasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = carasus.taskSvc.UpdateActionPlanTaskStatus("compliance_action_plan_id", compliaceActionPlanID, *existingComplianceActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Start a new transaction for the second update
	tx2 := carasus.repo.BeginTransaction()
	_, err = carasus.taskSvc.UpdateActionPlanTaskStatus("compliance_action_plan_id", compliaceActionPlanID, existingComplianceActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx2)
	if err != nil {
		tx2.Rollback()
		return err
	}

	tx2.Commit()

	go func() {

		urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, assessmentID, requirementID, compliaceActionPlanID)
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

		responsibleUser, err := carasus.userSvc.GetByID(existingComplianceActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := carasus.userSvc.GetByID(*existingComplianceActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Assessment Requirement Action Plan %s: Closed", existingComplianceActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ComplianceAssessment":  complianceAssessment.ComplianceRegister.TitleEn,
			"ActionPlanTitle":       existingComplianceActionPlan.Title,
			"ActionPlanDescription": existingComplianceActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%scompliance/assessment/%d/action-plans?action_plan_id=%d", appURL, assessmentID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = carasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingComplianceActionPlan, err := carasus.complianceActionPlanSvc.GetByID(assessmentID, requirementID, compliaceActionPlanID)
		if err != nil {
			return
		}

		err = carasus.activityLogSvc.LogActivity(
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

func (carasus *ComplianceAssessmentRequirementActionPlanStatusUpdateService) Revert(assessmentID uint, requirementID uint, compliaceActionPlanID uint, currentUserID uint, complianceActionPlanReviewUserStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanReviewUserStatusUpdateFormDTO) error {
	existingComplianceActionPlan, err := carasus.complianceActionPlanSvc.GetByID(assessmentID, requirementID, compliaceActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, compliaceActionPlanID)
	}

	if currentUserID != *existingComplianceActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	complianceActionPlanStatusResolved, err := carasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	complianceActionPlanReopenStatus, err := carasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	complianceAssessment, err := carasus.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}

	tx := carasus.repo.BeginTransaction()

	if existingComplianceActionPlan.ActionPlanStatusID != complianceActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The compliance action plan status must be resolved.")
	}

	existingComplianceActionPlan.ActionPlanStatusID = complianceActionPlanReopenStatus.ID

	if err := carasus.repo.ComplianceActionPlanUpdateStatus(existingComplianceActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: compliaceActionPlanID,
		ActionPlanStatusID:     complianceActionPlanReopenStatus.ID,
		Comment:                complianceActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:            currentUserID,
	}

	id, err := carasus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := carasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = carasus.taskSvc.UpdateActionPlanTaskStatus("compliance_action_plan_id", compliaceActionPlanID, *existingComplianceActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {

		urlParam := fmt.Sprintf(`{"compliance_assessment_id":%d,"requirement_id":%d, "action_plan_id":%d}`, assessmentID, requirementID, compliaceActionPlanID)
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

		responsibleUser, err := carasus.userSvc.GetByID(existingComplianceActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := carasus.userSvc.GetByID(*existingComplianceActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Assessment Requirement Action Plan %s: Reopened", existingComplianceActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ComplianceAssessment":  complianceAssessment.ComplianceRegister.TitleEn,
			"ActionPlanTitle":       existingComplianceActionPlan.Title,
			"ActionPlanDescription": existingComplianceActionPlan.Description,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%scompliance/assessment/%d/action-plans?action_plan_id=%d", appURL, assessmentID, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = carasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingComplianceActionPlan, err := carasus.complianceActionPlanSvc.GetByID(assessmentID, requirementID, compliaceActionPlanID)
		if err != nil {
			return
		}

		err = carasus.activityLogSvc.LogActivity(
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
