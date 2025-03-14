package services

import (
	"errors"
	"fmt"
	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
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

type TaeRecoveryActionPlanStatusUpdateServiceInterface interface {
	GetAll(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error)
	Create(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) (uint, error)
	Update(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, id uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) error
	GetLatestTaeRecoveryActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(slaID uint, actionPlanID uint, id uint, token string) (string, error)

	Approve(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, currentUserID uint, recoveryReviewUserStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, currentUserID uint, recoveryReviewUserStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanReviewUserStatusUpdateFormDTO) error
}

type TaeRecoveryActionPlanStatusUpdateService struct {
	repo                                 repositories.TaeRecoveryActionPlanStatusUpdateRepositoryInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	testAndExerciseRecoveryActionPlanSvc TestAndExerciseRecoveryActionPlanServiceInterface
	actionPlanStatusSvc                  masterServices.ActionPlanStatusServiceInterface
	testAndExerciseSvc                   TestAndExerciseServiceInterface
	userSvc                              organizationServices.UserServiceInterface
	notificationSvc                      commonServices.NotificationServiceInterface
	taskSvc                              taskServices.TaskServiceInterface
	taskStatusSvc                        taskServices.TaskStatusServiceInterface
}

func NewTaeRecoveryActionPlanStatusUpdateService(
	repo repositories.TaeRecoveryActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	testAndExerciseRecoveryActionPlanSvc TestAndExerciseRecoveryActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	testAndExerciseSvc TestAndExerciseServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) TaeRecoveryActionPlanStatusUpdateServiceInterface {
	return &TaeRecoveryActionPlanStatusUpdateService{
		repo:                                 repo,
		activityLogSvc:                       activityLogSvc,
		testAndExerciseRecoveryActionPlanSvc: testAndExerciseRecoveryActionPlanSvc,
		actionPlanStatusSvc:                  actionPlanStatusSvc,
		testAndExerciseSvc:                   testAndExerciseSvc,
		userSvc:                              userSvc,
		notificationSvc:                      notificationSvc,
		taskSvc:                              taskSvc,
		taskStatusSvc:                        taskStatusSvc,
	}
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) GetAll(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetAll(TaeRecoveryActionPlanID, sort, sortBy, keyword)
}

func (sasus TaeRecoveryActionPlanStatusUpdateService) GetByID(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, id uint) (*models.TaeRecoveryActionPlanStatusUpdate, error) {
	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID); err != nil {
		return nil, err
	}

	TaeRecoveryActionPlanStatusUpdate, err := sasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return TaeRecoveryActionPlanStatusUpdate, nil
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) Create(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := sasus.repo.BeginTransaction()

	actionPlan, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID)
	if err != nil {
		return 0, err
	}

	if err := sasus.validateTaeRecoveryActionPlanStatusUpdateFormDTO(TaeRecoveryActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	TaeRecoveryActionPlanUpdate := &models.TaeRecoveryActionPlanStatusUpdate{
		TestAndExerciseRecoveryActionPlanID: TaeRecoveryActionPlanID,
		ActionPlanStatusID:                  TaeRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                             TaeRecoveryActionPlanStatusUpdateFormDTO.Comment,
		Percentage:                          TaeRecoveryActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:                         currentUserID,
	}

	id, err := sasus.repo.Create(TaeRecoveryActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = sasus.testAndExerciseRecoveryActionPlanSvc.StatusUpdate(TestAndExerciseID, TaeRecoveryActionPlanID, currentUserID, TaeRecoveryActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.TaeRecoveryActionPlanStatusUpdateDocument

	for _, doc := range TaeRecoveryActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TestAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.TaeRecoveryActionPlanStatusUpdateDocument{
			TaeRecoveryActionPlanStatusUpdateID: id,
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
		if err := sasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TestAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	testAndExercise, err := sasus.testAndExerciseSvc.GetFullDetailsByID(TestAndExerciseID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	actionPlanStatus, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if TaeRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && actionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, TestAndExerciseID, TaeRecoveryActionPlanID)
		if actionPlan.IsRecovery {
			urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, TestAndExerciseID, TaeRecoveryActionPlanID)
		}
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:     "business_continuity",
			RecoveryActionPlanID: &TaeRecoveryActionPlanID,
			Title:                actionPlan.Title,
			URLParam:             &urlParam,
			QueryParam:           nil,
			SeverityLevelID:      actionPlan.SeverityLevelID,
			ResponsibleUserID:    *actionPlan.ReviewUserID,
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
			models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
			id,
			TaeRecoveryActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if TaeRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && actionPlan.ReviewUserID != nil {
			// prepare notification
			appURL := utility.GetAppURL()
			urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, TestAndExerciseID, TaeRecoveryActionPlanID)
			redirectUrl := fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, TestAndExerciseID, TaeRecoveryActionPlanID)
			if actionPlan.IsRecovery {
				urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, TestAndExerciseID, TaeRecoveryActionPlanID)
				redirectUrl = fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/recovery-action-plans?action_plan_id=%d", appURL, TestAndExerciseID, TaeRecoveryActionPlanID)
			}
			queryParam := ""
			entity := "test-and-exercise-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*actionPlan.ReviewUserID},
				ModuleIdentifier:     "business_continuity",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Action Plan Status Update",
				HeadingAr:            "تحديث حالة خطة العمل",
				TitleEn:              actionPlan.Title,
				TitleAr:              actionPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			reviewUser, err := sasus.userSvc.GetByID(*actionPlan.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			resolvedBy, err := sasus.userSvc.GetByID(actionPlan.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch updated user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Action Plan %s: Resolved", actionPlan.Title)

			emailData := map[string]interface{}{
				"FirstNameEn":           reviewUser.FirstNameEn,
				"LastNameEn":            reviewUser.LastNameEn,
				"ReferenceNumber":       testAndExercise.ReferenceNumber,
				"TestAndExerciseTitle":  testAndExercise.BcpTestAndExercise.Title,
				"ActionPlanTitle":       actionPlan.Title,
				"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
				"ResolvedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           redirectUrl,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "business-continuity/templates/email/action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) Update(TestAndExerciseID uint, TaeRecoveryActionPlanID uint, id uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) error {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID); err != nil {
		return err
	}

	oldTaeRecoveryActionPlanStatusUpdate, err := sasus.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID, id)

	if err != nil {
		return err
	}

	existingTaeRecoveryActionPlanStatusUpdate, err := sasus.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID, id)

	if err != nil {
		return err
	}

	if err := sasus.validateTaeRecoveryActionPlanStatusUpdateFormDTO(TaeRecoveryActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	// Update TaeRecovery StatusUpdate-- only comments and percentage
	existingTaeRecoveryActionPlanStatusUpdate.Comment = TaeRecoveryActionPlanStatusUpdateFormDTO.Comment
	existingTaeRecoveryActionPlanStatusUpdate.Percentage = TaeRecoveryActionPlanStatusUpdateFormDTO.Percentage
	existingTaeRecoveryActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = sasus.repo.Update(existingTaeRecoveryActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.TaeRecoveryActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.TaeRecoveryActionPlanStatusUpdateDocument

	for _, doc := range TaeRecoveryActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TestAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.TaeRecoveryActionPlanStatusUpdateDocument{
				TaeRecoveryActionPlanStatusUpdateID: id,
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
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TestAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TestAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := sasus.GetByID(TestAndExerciseID, TaeRecoveryActionPlanID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldTaeRecoveryActionPlanStatusUpdate.Comment,
			"percentage": oldTaeRecoveryActionPlanStatusUpdate.Percentage,
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
				models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
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

func (sasus *TaeRecoveryActionPlanStatusUpdateService) validateTaeRecoveryActionPlanStatusUpdateFormDTO(TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) error {

	if _, err := sasus.actionPlanStatusSvc.GetByID(TaeRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (sasus TaeRecoveryActionPlanStatusUpdateService) GetLatestTaeRecoveryActionPlanStatusUpdate(testAndExerciseID uint, actionPlanID uint) (*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetLatestTaeRecoveryActionPlanStatusUpdate(actionPlanID)
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) GetDocumentPath(testAndExerciseID, ActionPlanID, id uint, token string) (string, error) {

	existingTaeRecoveryActionplanStatusUpdate, err := sasus.GetByID(testAndExerciseID, ActionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingTaeRecoveryActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the recovery action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (sasus *TaeRecoveryActionPlanStatusUpdateService) Approve(testAndExerciseID uint, actionPlanID uint, currentUserID uint, recoveryActionPlanReviewUserStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanReviewUserStatusUpdateFormDTO) error {
	existingActionPlan, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	// if currentUserID != *existingActionPlan.ReviewUserID {
	// 	return utility.NewForbiddenError("Permission denied.")
	// }

	actionPlanStatusResolved, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	recoveryActionPlanClosedStatus, err := sasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := sasus.repo.BeginTransaction()

	if *existingActionPlan.ActionPlanStatusId != actionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The action plan status must be resolved.")
	}

	existingActionPlan.ActionPlanStatusId = &recoveryActionPlanClosedStatus.ID

	if err := sasus.repo.RecoveryActionPlanUpdateStatus(existingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	recoveryActionPlanUpdate := &models.TaeRecoveryActionPlanStatusUpdate{
		TestAndExerciseRecoveryActionPlanID: actionPlanID,
		ActionPlanStatusID:                  recoveryActionPlanClosedStatus.ID,
		Percentage:                          nil,
		Comment:                             recoveryActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:                         currentUserID,
	}

	id, err := sasus.repo.Create(recoveryActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := sasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("recovery_action_plan_id", actionPlanID, *existingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Start a new transaction for the second update
	tx2 := sasus.repo.BeginTransaction()
	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("recovery_action_plan_id", actionPlanID, existingActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx2)
	if err != nil {
		tx2.Rollback()
		return err
	}

	tx2.Commit()

	go func() {
		var notificationHeadingEn, notificationHeadingAr string
		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, actionPlanID)
		queryParam := ""
		entity := "test-and-exercise-action-plan"
		action := "close-action-plan"
		notificationHeadingEn = "Test And Exercise Action Plan Status Update"
		notificationHeadingAr = "تحديث حالة خطة عمل الاختبار والتمرين"
		redirectUrl := fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)

		if existingActionPlan.IsRecovery {
			urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, testAndExerciseID, actionPlanID)
			notificationHeadingEn = "Test And Exercise Recovery Action Plan Status Update"
			notificationHeadingAr = "تحديث حالة خطة عمل التعافي من الاختبار والتمرين"
			redirectUrl = fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/recovery-action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "test-and-exercise",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            notificationHeadingEn,
			HeadingAr:            notificationHeadingAr,
			TitleEn:              existingActionPlan.Title,
			TitleAr:              existingActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := sasus.userSvc.GetByID(existingActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := sasus.userSvc.GetByID(*existingActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Action Plan %s: Closed", existingActionPlan.Title)
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ActionPlanTitle":       existingActionPlan.Title,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectUrl,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "business-continuity/templates/email/action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingActionPlan, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID)
		if err != nil {
			return
		}

		err = sasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) Revert(testAndExerciseID uint, actionPlanID uint, currentUserID uint, actionPlanReviewUserStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanReviewUserStatusUpdateFormDTO) error {
	existingActionPlan, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if currentUserID != *existingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	actionPlanStatusResolved, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	actionPlanReopenStatus, err := sasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := sasus.repo.BeginTransaction()

	if *existingActionPlan.ActionPlanStatusId != actionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The action plan status must be resolved.")
	}

	existingActionPlan.ActionPlanStatusId = &actionPlanReopenStatus.ID

	if err := sasus.repo.RecoveryActionPlanUpdateStatus(existingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	actionPlanUpdate := &models.TaeRecoveryActionPlanStatusUpdate{
		TestAndExerciseRecoveryActionPlanID: actionPlanID,
		ActionPlanStatusID:                  actionPlanReopenStatus.ID,
		Comment:                             actionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:                         currentUserID,
	}

	id, err := sasus.repo.Create(actionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := sasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("recovery_action_plan_id", actionPlanID, *existingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		var notificationHeadingEn, notificationHeadingAr string
		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "action_plan_id":%d}`, testAndExerciseID, actionPlanID)
		queryParam := ""
		entity := "test-and-exercise-action-plan"
		action := "reopen-action-plan"
		notificationHeadingEn = "Test And Exercise Action Plan Status Update"
		notificationHeadingAr = "تحديث حالة خطة عمل الاختبار والتمرين"
		redirectUrl := fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)

		if existingActionPlan.IsRecovery {
			urlParam = fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, testAndExerciseID, actionPlanID)
			notificationHeadingEn = "Test And Exercise Recovery Action Plan Status Update"
			notificationHeadingAr = "تحديث حالة خطة عمل التعافي من الاختبار والتمرين"
			redirectUrl = fmt.Sprintf("%sbusiness-continuity/test-and-exercise/%d/recovery-action-plans?action_plan_id=%d", appURL, testAndExerciseID, id)
		}

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "test-and-exercise",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            notificationHeadingEn,
			HeadingAr:            notificationHeadingAr,
			TitleEn:              existingActionPlan.Title,
			TitleAr:              existingActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := sasus.userSvc.GetByID(existingActionPlan.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		resolvedBy, err := sasus.userSvc.GetByID(*existingActionPlan.ReviewUserID)
		if err != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Action Plan %s: Reopened", existingActionPlan.Title)
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"ActionPlanTitle":       existingActionPlan.Title,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           redirectUrl,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "business-continuity/templates/email/action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "test_and_exercise")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingActionPlan, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID)
		if err != nil {
			return
		}

		err = sasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
