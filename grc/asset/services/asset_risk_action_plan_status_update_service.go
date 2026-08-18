package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskActionPlanStatusUpdateServiceInterface interface {
	GetAll(riskID uint, riskActionPlanID uint, sort string, sortBy string, keyword string) ([]*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error)
	Create(riskID uint, riskActionPlanID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO) (uint, error)
	Update(riskID uint, riskActionPlanID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO) error
	GetLatestRiskActionPlanStatusUpdate(riskID uint, riskActionPlanID uint) (*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(riskID, riskActionPlanID, statusUpdateID uint, token string) (string, error)

	Approve(riskID uint, riskActionPlanID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(riskID uint, riskActionPlanID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanReviewUserStatusUpdateFormDTO) error
}

type AssetRiskActionPlanStatusUpdateService struct {
	repo                   repositories.AssetRiskActionPlanStatusUpdateRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	assetRiskRegisterSvc   AssetRiskRegisterServiceInterface
	assetRiskActionPlanSvc AssetRiskActionPlanServiceInterface
	actionPlanStatusSvc    masterServices.ActionPlanStatusServiceInterface
	userSvc                organizationServices.UserServiceInterface
	notificationSvc        commonServices.NotificationServiceInterface
	taskSvc                taskServices.TaskServiceInterface
	taskStatusSvc          taskServices.TaskStatusServiceInterface
}

func NewAssetRiskActionPlanStatusUpdateService(
	repo repositories.AssetRiskActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRiskRegisterSvc AssetRiskRegisterServiceInterface,
	assetRiskActionPlanSvc AssetRiskActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) AssetRiskActionPlanStatusUpdateServiceInterface {
	return &AssetRiskActionPlanStatusUpdateService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		assetRiskRegisterSvc:   assetRiskRegisterSvc,
		assetRiskActionPlanSvc: assetRiskActionPlanSvc,
		actionPlanStatusSvc:    actionPlanStatusSvc,
		userSvc:                userSvc,
		notificationSvc:        notificationSvc,
		taskSvc:                taskSvc,
		taskStatusSvc:          taskStatusSvc,
	}
}

func (arasus *AssetRiskActionPlanStatusUpdateService) GetAll(riskID uint, riskActionPlanID uint, sort string, sortBy string, keyword string) ([]*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error) {
	if _, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID); err != nil {
		return nil, err
	}

	return arasus.repo.GetAll(riskActionPlanID, sort, sortBy, keyword)
}

func (arasus AssetRiskActionPlanStatusUpdateService) GetByID(riskID uint, riskActionPlanID uint, id uint) (*models.RiskActionPlanStatusUpdate, error) {
	if _, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID); err != nil {
		return nil, err
	}

	riskActionPlanStatusUpdate, err := arasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlanStatusUpdate, nil
}

func (arasus *AssetRiskActionPlanStatusUpdateService) Create(riskID uint, riskActionPlanID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO) (id uint, err error) {
	riskActionPlan, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID)
	if err != nil {
		return 0, err
	}
	if riskID == 0 {
		riskID = riskActionPlan.RiskRegisterID
	}

	actionPlanStatus, err := arasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	if err = arasus.validateAssetRiskActionPlanStatusUpdateFormDTO(formDTO); err != nil {
		return 0, err
	}

	tx := arasus.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   riskActionPlanID,
		ActionPlanStatusID: formDTO.ActionPlanStatusID,
		Comment:            formDTO.Comment,
		Percentage:         formDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err = arasus.repo.Create(riskActionPlanUpdate, tx)
	if err != nil {
		return 0, err
	}

	err = arasus.assetRiskActionPlanSvc.StatusUpdate(riskID, riskActionPlanID, currentUserID, formDTO, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.RiskActionPlanStatusUpdateDocument

	for _, doc := range formDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}

		url := utility.GenerateDocumentURL("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(riskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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
		if err := arasus.repo.CreateDocuments(documents, tx); err != nil {
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(riskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				return 0, err
			}
		}
	}

	if formDTO.ActionPlanStatusID == actionPlanStatus.ID && riskActionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, riskActionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:  "asset_risk",
			RiskActionPlanID:  &riskActionPlanID,
			Title:             riskActionPlan.Title,
			Description:       riskActionPlan.Description,
			URLParam:          &urlParam,
			QueryParam:        nil,
			SeverityLevelID:   riskActionPlan.SeverityLevelID,
			ResponsibleUserID: *riskActionPlan.ReviewUserID,
		}
		if _, err := arasus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			return 0, err
		}
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		logErr := arasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			riskActionPlanUpdate,
			currentUserID,
			"asset", "asset_risk", "action_plan_status_update", nil,
		)
		if logErr != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     logErr.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (arasus *AssetRiskActionPlanStatusUpdateService) Update(riskID uint, riskActionPlanID uint, id uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO) (err error) {
	riskActionPlan, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID)
	if err != nil {
		return err
	}
	if riskID == 0 {
		riskID = riskActionPlan.RiskRegisterID
	}

	existingRiskActionPlanStatusUpdate, err := arasus.GetByID(riskID, riskActionPlanID, id)
	if err != nil {
		return err
	}

	oldActionPlanStatusUpdate, err := arasus.GetByID(riskID, riskActionPlanID, id)
	if err != nil {
		return err
	}

	if err := arasus.validateAssetRiskActionPlanStatusUpdateFormDTO(formDTO); err != nil {
		return err
	}

	tx := arasus.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	// Update Risk StatusUpdate -- only comments and percentage
	existingRiskActionPlanStatusUpdate.Comment = formDTO.Comment
	existingRiskActionPlanStatusUpdate.Percentage = formDTO.Percentage
	existingRiskActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = arasus.repo.Update(existingRiskActionPlanStatusUpdate, tx)
	if err != nil {
		return err
	}

	var documents []*models.RiskActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskActionPlanStatusUpdateDocument

	for _, doc := range formDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				return err
			}

			url := utility.GenerateDocumentURL("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(riskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
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

			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := arasus.repo.GetDocumentByID(*doc.ID)
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
		if err := arasus.repo.CreateDocuments(documents, tx); err != nil {
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := arasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(riskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(riskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		existingActionPlanStatusUpdate, logErr := arasus.GetByID(riskID, riskActionPlanID, id)
		if logErr != nil {
			return
		}

		changedData := make(map[string]interface{})

		oldData := map[string]interface{}{
			"comment":    oldActionPlanStatusUpdate.Comment,
			"percentage": oldActionPlanStatusUpdate.Percentage,
		}

		newData := map[string]interface{}{
			"comment":    existingActionPlanStatusUpdate.Comment,
			"percentage": existingActionPlanStatusUpdate.Percentage,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			logErr := arasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset", "asset_risk", "action_plan_status_update", nil,
			)
			if logErr != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    logErr.Error(),
					"recordID": id,
				})
			}
		}
	}()

	return nil
}

func (arasus *AssetRiskActionPlanStatusUpdateService) validateAssetRiskActionPlanStatusUpdateFormDTO(formDTO *assetDTOs.AssetRiskActionPlanStatusUpdateFormDTO) error {
	if _, err := arasus.actionPlanStatusSvc.GetByID(formDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (arasus AssetRiskActionPlanStatusUpdateService) GetLatestRiskActionPlanStatusUpdate(riskID uint, riskActionPlanID uint) (*assetDTOs.AssetRiskActionPlanStatusUpdateListDTO, error) {
	if _, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID); err != nil {
		return nil, err
	}

	return arasus.repo.GetLatestRiskActionPlanStatusUpdate(riskActionPlanID)
}

func (arasus *AssetRiskActionPlanStatusUpdateService) GetDocumentPath(riskID, riskActionPlanID, statusUpdateID uint, token string) (string, error) {
	statusUpdate, err := arasus.GetByID(riskID, riskActionPlanID, statusUpdateID)
	if err != nil {
		return "", err
	}

	if len(statusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the risk action plan update")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		return "", fmt.Errorf("resource not found")
	}
	decryptedFilename := string(decryptedFilenameBytes)

	storagePath := filepath.Join("asset", "asset-risk-register", strconv.FormatUint(uint64(riskID), 10), "action-plans", strconv.FormatUint(uint64(riskActionPlanID), 10), "updates", strconv.FormatUint(uint64(statusUpdateID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (arasus *AssetRiskActionPlanStatusUpdateService) Approve(riskID uint, riskActionPlanID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanReviewUserStatusUpdateFormDTO) (err error) {
	existingRiskActionPlan, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, riskActionPlanID)
	}
	if riskID == 0 {
		riskID = existingRiskActionPlan.RiskRegisterID
	}

	if existingRiskActionPlan.ReviewUserID == nil || currentUserID != *existingRiskActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	riskActionPlanStatusResolved, err := arasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	riskActionPlanClosedStatus, err := arasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	riskRegister, err := arasus.assetRiskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	if existingRiskActionPlan.ActionPlanStatusID != riskActionPlanStatusResolved.ID {
		return utility.NewForbiddenError("The risk action plan status must be resolved.")
	}

	tx := arasus.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	existingRiskActionPlan.ActionPlanStatusID = riskActionPlanClosedStatus.ID

	if err = arasus.repo.RiskActionPlanUpdateStatus(existingRiskActionPlan, tx); err != nil {
		return err
	}

	closedPercentage := uint(100)
	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   riskActionPlanID,
		ActionPlanStatusID: riskActionPlanClosedStatus.ID,
		Percentage:         &closedPercentage,
		Comment:            formDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := arasus.repo.Create(riskActionPlanUpdate, tx)
	if err != nil {
		return err
	}

	taskStatus, err := arasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		return err
	}

	_, err = arasus.taskSvc.UpdateActionPlanTaskStatus("risk_action_plan_id", riskActionPlanID, *existingRiskActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, riskActionPlanID)
		queryParam := ""
		entity := "asset-risk-action-plan"
		action := "close-risk-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingRiskActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "ISMS Risk Action Plan Status Update",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل امتثال مخاطر ISMS",
			TitleEn:              existingRiskActionPlan.Title,
			TitleAr:              existingRiskActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, logErr := arasus.userSvc.GetByID(existingRiskActionPlan.ResponsibleUserID)
		if logErr != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}
		resolvedBy, logErr := arasus.userSvc.GetByID(*existingRiskActionPlan.ReviewUserID)
		if logErr != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		actionPlanDescription := "-"
		if existingRiskActionPlan.Description != nil && *existingRiskActionPlan.Description != "" {
			actionPlanDescription = *existingRiskActionPlan.Description
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("ISMS Risk Action Plan: %s - Closed", existingRiskActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"RiskRegister":          riskRegister.Title,
			"ActionPlanTitle":       existingRiskActionPlan.Title,
			"ActionPlanDescription": actionPlanDescription,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/asset/risk-mitigations?id=%d", appURL, riskActionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_action_plan_closed_responsible_user_mail.html",
			EmailData:            emailData,
		}

		logErr = arasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk_action_plan")
		if logErr != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": logErr.Error(),
			})
		}

		logErr = arasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			existingRiskActionPlan,
			currentUserID,
			"asset", "asset_risk", "action_plan_status_update", nil,
		)
		if logErr != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     logErr.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (arasus *AssetRiskActionPlanStatusUpdateService) Revert(riskID uint, riskActionPlanID uint, currentUserID uint, formDTO *assetDTOs.AssetRiskActionPlanReviewUserStatusUpdateFormDTO) (err error) {
	existingRiskActionPlan, err := arasus.assetRiskActionPlanSvc.GetByID(riskID, riskActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, riskActionPlanID)
	}
	if riskID == 0 {
		riskID = existingRiskActionPlan.RiskRegisterID
	}

	if existingRiskActionPlan.ReviewUserID == nil || currentUserID != *existingRiskActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	riskActionPlanStatusResolved, err := arasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	riskActionPlanReopenStatus, err := arasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	riskRegister, err := arasus.assetRiskRegisterSvc.GetFullDetailsByID(riskID)
	if err != nil {
		return err
	}

	if existingRiskActionPlan.ActionPlanStatusID != riskActionPlanStatusResolved.ID {
		return utility.NewForbiddenError("The risk action plan status must be resolved.")
	}

	tx := arasus.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}
	}()

	existingRiskActionPlan.ActionPlanStatusID = riskActionPlanReopenStatus.ID

	if err = arasus.repo.RiskActionPlanUpdateStatus(existingRiskActionPlan, tx); err != nil {
		return err
	}

	reopenedPercentage := uint(0)
	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   riskActionPlanID,
		ActionPlanStatusID: riskActionPlanReopenStatus.ID,
		Percentage:         &reopenedPercentage,
		Comment:            formDTO.Comment,
		CreatedByID:        currentUserID,
	}

	id, err := arasus.repo.Create(riskActionPlanUpdate, tx)
	if err != nil {
		return err
	}

	taskStatus, err := arasus.taskStatusSvc.GetByType("closed")
	if err != nil {
		return err
	}

	_, err = arasus.taskSvc.UpdateActionPlanTaskStatus("risk_action_plan_id", riskActionPlanID, *existingRiskActionPlan.ReviewUserID, uint(taskStatus.ID), tx)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"risk_register_id":%d, "action_plan_id":%d}`, riskID, riskActionPlanID)
		queryParam := ""
		entity := "asset-risk-action-plan"
		action := "reopen-risk-action-plan"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{existingRiskActionPlan.ResponsibleUserID},
			ModuleIdentifier:     "asset_risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "ISMS Risk Action Plan Status Update",
			HeadingAr:            "تم تعيين مسؤولية خطة عمل امتثال مخاطر ISMS",
			TitleEn:              existingRiskActionPlan.Title,
			TitleAr:              existingRiskActionPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, logErr := arasus.userSvc.GetByID(existingRiskActionPlan.ResponsibleUserID)
		if logErr != nil {
			utility.LogError("Failed to fetch review user", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}
		resolvedBy, logErr := arasus.userSvc.GetByID(*existingRiskActionPlan.ReviewUserID)
		if logErr != nil {
			utility.LogError("Failed to fetch updated user", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		actionPlanDescription := "-"
		if existingRiskActionPlan.Description != nil && *existingRiskActionPlan.Description != "" {
			actionPlanDescription = *existingRiskActionPlan.Description
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("ISMS Risk Action Plan: %s - Reopened", existingRiskActionPlan.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":           responsibleUser.FirstNameEn,
			"LastNameEn":            responsibleUser.LastNameEn,
			"RiskRegister":          riskRegister.Title,
			"ActionPlanTitle":       existingRiskActionPlan.Title,
			"ActionPlanDescription": actionPlanDescription,
			"ResolvedByFirstNameEn": resolvedBy.FirstNameEn,
			"ResolvedByLastNameEn":  resolvedBy.LastNameEn,
			"ResolvedOn":            time.Now().Format("2006-01-02"),
			"RedirectURL":           fmt.Sprintf("%s/asset/risk-mitigations?id=%d", appURL, riskActionPlanID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/risk_action_plan_reopen_responsible_user_mail.html",
			EmailData:            emailData,
		}

		logErr = arasus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "asset_risk_action_plan")
		if logErr != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": logErr.Error(),
			})
		}

		logErr = arasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			existingRiskActionPlan,
			currentUserID,
			"asset", "asset_risk", "action_plan_status_update", nil,
		)
		if logErr != nil {
			fields := map[string]interface{}{
				"action":    models.RejectAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     logErr.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
