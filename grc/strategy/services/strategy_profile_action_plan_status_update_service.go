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
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/grc/strategy/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type StrategyProfileActionPlanStatusUpdateServiceInterface interface {
	GetAll(StrategyProfileID uint, StrategyProfileActionPlanID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error)
	Create(StrategyProfileID uint, StrategyProfileActionPlanID uint, currentUserID uint, StrategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO) (uint, error)
	Update(StrategyProfileID uint, StrategyProfileActionPlanID uint, id uint, currentUserID uint, StrategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO) error
	GetLatestStrategyProfileActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(slaID uint, actionPlanID uint, id uint, token string) (string, error)

	Approve(StrategyProfileID uint, StrategyProfileActionPlanID uint, currentUserID uint, strategy_profileReviewUserStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(StrategyProfileID uint, StrategyProfileActionPlanID uint, currentUserID uint, strategy_profileReviewUserStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanReviewUserStatusUpdateFormDTO) error
}

type StrategyProfileActionPlanStatusUpdateService struct {
	repo                         repositories.StrategyProfileActionPlanStatusUpdateRepositoryInterface
	activityLogSvc               logServices.ActivityLogServiceInterface
	strategyProfileActionPlanSvc StrategyProfileActionPlanServiceInterface
	actionPlanStatusSvc          masterServices.ActionPlanStatusServiceInterface
	strategyProfileSvc           StrategyProfileServiceInterface
	userSvc                      organizationServices.UserServiceInterface
	notificationSvc              commonServices.NotificationServiceInterface
	taskSvc                      taskServices.TaskServiceInterface
	taskStatusSvc                taskServices.TaskStatusServiceInterface
}

func NewStrategyProfileActionPlanStatusUpdateService(
	repo repositories.StrategyProfileActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	strategyProfileActionPlanSvc StrategyProfileActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) StrategyProfileActionPlanStatusUpdateServiceInterface {
	return &StrategyProfileActionPlanStatusUpdateService{
		repo:                         repo,
		activityLogSvc:               activityLogSvc,
		strategyProfileActionPlanSvc: strategyProfileActionPlanSvc,
		actionPlanStatusSvc:          actionPlanStatusSvc,
		strategyProfileSvc:           strategyProfileSvc,
		userSvc:                      userSvc,
		notificationSvc:              notificationSvc,
		taskSvc:                      taskSvc,
		taskStatusSvc:                taskStatusSvc,
	}
}

func (sasus *StrategyProfileActionPlanStatusUpdateService) GetAll(StrategyProfileID uint, StrategyProfileActionPlanID uint, sort string, sortBy string, keyword string) ([]*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.strategyProfileActionPlanSvc.GetByID(StrategyProfileID, StrategyProfileActionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetAll(StrategyProfileActionPlanID, sort, sortBy, keyword)
}

func (sasus StrategyProfileActionPlanStatusUpdateService) GetByID(StrategyProfileID uint, StrategyProfileActionPlanID uint, id uint) (*models.StrategyProfileActionPlanStatusUpdate, error) {
	if _, err := sasus.strategyProfileActionPlanSvc.GetByID(StrategyProfileID, StrategyProfileActionPlanID); err != nil {
		return nil, err
	}

	StrategyProfileActionPlanStatusUpdate, err := sasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return StrategyProfileActionPlanStatusUpdate, nil
}

func (sasus *StrategyProfileActionPlanStatusUpdateService) Create(StrategyProfileID uint, StrategyProfileActionPlanID uint, currentUserID uint, StrategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := sasus.repo.BeginTransaction()

	actionPlan, err := sasus.strategyProfileActionPlanSvc.GetByID(StrategyProfileID, StrategyProfileActionPlanID)
	if err != nil {
		return 0, err
	}

	if err := sasus.validateStrategyProfileActionPlanStatusUpdateFormDTO(StrategyProfileActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	StrategyProfileActionPlanUpdate := &models.StrategyProfileActionPlanStatusUpdate{
		StrategyProfileActionPlanID: StrategyProfileActionPlanID,
		ActionPlanStatusID:          StrategyProfileActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                     StrategyProfileActionPlanStatusUpdateFormDTO.Comment,
		Percentage:                  StrategyProfileActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:                 currentUserID,
	}

	id, err := sasus.repo.Create(StrategyProfileActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = sasus.strategyProfileActionPlanSvc.StatusUpdate(StrategyProfileID, StrategyProfileActionPlanID, currentUserID, StrategyProfileActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.StrategyProfileActionPlanStatusUpdateDocument

	for _, doc := range StrategyProfileActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("strategy", "strategy-profile", strconv.FormatUint(uint64(StrategyProfileID), 10), "strategy_profile-action-plans", strconv.FormatUint(uint64(StrategyProfileActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.StrategyProfileActionPlanStatusUpdateDocument{
			StrategyProfileActionPlanStatusUpdateID: id,
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
		if err := sasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(StrategyProfileID), 10), "strategy_profile-action-plans", strconv.FormatUint(uint64(StrategyProfileActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	strategyProfile, err := sasus.strategyProfileSvc.GetFullDetailsByID(StrategyProfileID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	actionPlanStatus, err := sasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}
	if StrategyProfileActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && actionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, StrategyProfileID, StrategyProfileActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:            "strategy",
			StrategyProfileActionPlanID: &StrategyProfileActionPlanID,
			Title:                       actionPlan.Title,
			URLParam:                    &urlParam,
			QueryParam:                  nil,
			SeverityLevelID:             actionPlan.SeverityLevelID,
			ResponsibleUserID:           *actionPlan.ReviewUserID,
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
			models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
			id,
			StrategyProfileActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		if StrategyProfileActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && actionPlan.ReviewUserID != nil {
			// prepare notification
			appURL := utility.GetAppURL()
			urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, StrategyProfileID, StrategyProfileActionPlanID)
			redirectUrl := fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, StrategyProfileID, StrategyProfileActionPlanID)
			queryParam := ""
			entity := "strategy-profile-action-plan"
			action := "review-user"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*actionPlan.ReviewUserID},
				ModuleIdentifier:     "strategy",
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
			resolvedBy, err := sasus.userSvc.GetByID(*actionPlan.ResponsibleUserID)
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
				"StrategyProfileTitle":  strategyProfile.Title,
				"ActionPlanTitle":       actionPlan.Title,
				"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
				"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
				"ResolvedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           redirectUrl,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "strategy/templates/email/action_plan_review_user_mail.html",
				EmailData:            emailData,
			}

			// send notification and mail
			err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (sasus *StrategyProfileActionPlanStatusUpdateService) Update(StrategyProfileID uint, StrategyProfileActionPlanID uint, id uint, currentUserID uint, StrategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO) error {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.strategyProfileActionPlanSvc.GetByID(StrategyProfileID, StrategyProfileActionPlanID); err != nil {
		return err
	}

	oldStrategyProfileActionPlanStatusUpdate, err := sasus.GetByID(StrategyProfileID, StrategyProfileActionPlanID, id)

	if err != nil {
		return err
	}

	existingStrategyProfileActionPlanStatusUpdate, err := sasus.GetByID(StrategyProfileID, StrategyProfileActionPlanID, id)

	if err != nil {
		return err
	}

	if err := sasus.validateStrategyProfileActionPlanStatusUpdateFormDTO(StrategyProfileActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	// Update StrategyProfile StatusUpdate-- only comments and percentage
	existingStrategyProfileActionPlanStatusUpdate.Comment = StrategyProfileActionPlanStatusUpdateFormDTO.Comment
	existingStrategyProfileActionPlanStatusUpdate.Percentage = StrategyProfileActionPlanStatusUpdateFormDTO.Percentage
	existingStrategyProfileActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = sasus.repo.Update(existingStrategyProfileActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.StrategyProfileActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.StrategyProfileActionPlanStatusUpdateDocument

	for _, doc := range StrategyProfileActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("strategy", "strategy-profile", strconv.FormatUint(uint64(StrategyProfileID), 10), "strategy_profile-action-plans", strconv.FormatUint(uint64(StrategyProfileActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.StrategyProfileActionPlanStatusUpdateDocument{
				StrategyProfileActionPlanStatusUpdateID: id,
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
		storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(StrategyProfileID), 10), "strategy_profile-action-plans", strconv.FormatUint(uint64(StrategyProfileActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(StrategyProfileID), 10), "strategy_profile-action-plans", strconv.FormatUint(uint64(StrategyProfileActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := sasus.GetByID(StrategyProfileID, StrategyProfileActionPlanID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldStrategyProfileActionPlanStatusUpdate.Comment,
			"percentage": oldStrategyProfileActionPlanStatusUpdate.Percentage,
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
				models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
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

func (sasus *StrategyProfileActionPlanStatusUpdateService) validateStrategyProfileActionPlanStatusUpdateFormDTO(StrategyProfileActionPlanStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanStatusUpdateFormDTO) error {

	if _, err := sasus.actionPlanStatusSvc.GetByID(StrategyProfileActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (sasus StrategyProfileActionPlanStatusUpdateService) GetLatestStrategyProfileActionPlanStatusUpdate(strategyProfileID uint, actionPlanID uint) (*strategyDTOs.StrategyProfileActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetLatestStrategyProfileActionPlanStatusUpdate(actionPlanID)
}

func (sasus *StrategyProfileActionPlanStatusUpdateService) GetDocumentPath(strategyProfileID, ActionPlanID, id uint, token string) (string, error) {

	existingStrategyProfileActionplanStatusUpdate, err := sasus.GetByID(strategyProfileID, ActionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingStrategyProfileActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the strategy_profile action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("strategy", "strategy-profile", strconv.FormatUint(uint64(strategyProfileID), 10), "strategy_profile-action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (sasus *StrategyProfileActionPlanStatusUpdateService) Approve(strategyProfileID uint, actionPlanID uint, currentUserID uint, strategy_profileActionPlanReviewUserStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanReviewUserStatusUpdateFormDTO) error {
	existingActionPlan, err := sasus.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID)
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

	strategy_profileActionPlanClosedStatus, err := sasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := sasus.repo.BeginTransaction()

	if *existingActionPlan.ActionPlanStatusId != actionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The action plan status must be resolved.")
	}

	existingActionPlan.ActionPlanStatusId = &strategy_profileActionPlanClosedStatus.ID

	if err := sasus.repo.ActionPlanUpdateStatus(existingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	strategy_profileActionPlanUpdate := &models.StrategyProfileActionPlanStatusUpdate{
		StrategyProfileActionPlanID: actionPlanID,
		ActionPlanStatusID:          strategy_profileActionPlanClosedStatus.ID,
		Percentage:                  nil,
		Comment:                     strategy_profileActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:                 currentUserID,
	}

	id, err := sasus.repo.Create(strategy_profileActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := sasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("strategy_profile_action_plan_id", actionPlanID, *existingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Start a new transaction for the second update
	tx2 := sasus.repo.BeginTransaction()
	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("strategy_profile_action_plan_id", actionPlanID, *existingActionPlan.ResponsibleUserID, uint(taskStatus.ID), tx2)
	if err != nil {
		tx2.Rollback()
		return err
	}

	tx2.Commit()

	go func() {
		var notificationHeadingEn, notificationHeadingAr string
		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, actionPlanID)
		queryParam := ""
		entity := "strategy-profile-action-plan"
		action := "close-action-plan"
		notificationHeadingEn = "Test And Exercise Action Plan Status Update"
		notificationHeadingAr = "تحديث حالة خطة عمل الاختبار والتمرين"
		redirectUrl := fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id)

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*existingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "strategy-profile",
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

		responsibleUser, err := sasus.userSvc.GetByID(*existingActionPlan.ResponsibleUserID)
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
			TemplateRelativePath: "strategy/templates/email/action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingActionPlan, err := sasus.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID)
		if err != nil {
			return
		}

		err = sasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (sasus *StrategyProfileActionPlanStatusUpdateService) Revert(strategyProfileID uint, actionPlanID uint, currentUserID uint, actionPlanReviewUserStatusUpdateFormDTO *strategyDTOs.StrategyProfileActionPlanReviewUserStatusUpdateFormDTO) error {
	existingActionPlan, err := sasus.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID)
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

	if err := sasus.repo.ActionPlanUpdateStatus(existingActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	actionPlanUpdate := &models.StrategyProfileActionPlanStatusUpdate{
		StrategyProfileActionPlanID: actionPlanID,
		ActionPlanStatusID:          actionPlanReopenStatus.ID,
		Comment:                     actionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:                 currentUserID,
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

	_, err = sasus.taskSvc.UpdateActionPlanTaskStatus("strategy_profile_action_plan_id", actionPlanID, *existingActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		var notificationHeadingEn, notificationHeadingAr string
		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"strategy_profile_id":%d, "action_plan_id":%d}`, strategyProfileID, actionPlanID)
		queryParam := ""
		entity := "strategy-profile-action-plan"
		action := "reopen-action-plan"
		notificationHeadingEn = "Test And Exercise Action Plan Status Update"
		notificationHeadingAr = "تحديث حالة خطة عمل الاختبار والتمرين"
		redirectUrl := fmt.Sprintf("%sstrategy/strategy-profile/%d/action-plans?action_plan_id=%d", appURL, strategyProfileID, id)

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*existingActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "strategy-profile",
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

		responsibleUser, err := sasus.userSvc.GetByID(*existingActionPlan.ResponsibleUserID)
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
			TemplateRelativePath: "strategy/templates/email/action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = sasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "strategy_profile_action_plan")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}

		existingActionPlan, err := sasus.strategyProfileActionPlanSvc.GetByID(strategyProfileID, actionPlanID)
		if err != nil {
			return
		}

		err = sasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.StrategyProfileActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
