package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/identity-access-review/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewActionPlanStatusUpdateServiceInterface interface {
	GetByID(identityReviewUserID uint, identityReviewActionPlanID uint, id uint) (*models.IdentityReviewActionPlanStatusUpdate, error)
	GetAll(identityReviewUserID uint, identityReviewActionPlanID uint, sort string, sortBy string, keyword string) ([]*dtos.IdentityReviewActionPlanStatusUpdateListDTO, error)
	Create(identityReviewUserID uint, identityReviewActionPlanID uint, currentUserID uint, identityReviewActionPlanStatusUpdateFormDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO) (uint, error)
	Update(identityReviewUserID uint, identityReviewActionPlanID uint, id uint, currentUserID uint, identityReviewActionPlanStatusUpdateFormDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO) error
	GetLatestIdentityReviewActionPlanStatusUpdate(identityReviewUserID uint, identityReviewActionPlanID uint) (*dtos.IdentityReviewActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(identityReviewUserID uint, identityReviewActionPlanID uint, id uint, token string) (string, error)
	Approve(identityReviewUserID uint, identityReviewActionPlanID uint, currentUserID uint, identityReviewActionPlanReviewFormDTO *dtos.IdentityReviewActionPlanReviewFormDTO) error
	Revert(identityReviewUserID uint, identityReviewActionPlanID uint, currentUserID uint, identityReviewActionPlanReviewFormDTO *dtos.IdentityReviewActionPlanReviewFormDTO) error
}

type IdentityReviewActionPlanStatusUpdateService struct {
	repo                        repositories.IdentityReviewActionPlanStatusUpdateRepositoryInterface
	activityLogSvc              logServices.ActivityLogServiceInterface
	identityReviewActionPlanSvc IdentityReviewActionPlanServiceInterface
	actionPlanStatusSvc         masterServices.ActionPlanStatusServiceInterface
	taskSvc                     taskServices.TaskServiceInterface
	taskStatusSvc               taskServices.TaskStatusServiceInterface
}

func NewIdentityReviewActionPlanStatusUpdateService(
	repo repositories.IdentityReviewActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	identityReviewActionPlanSvc IdentityReviewActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	taskStatusSvc taskServices.TaskStatusServiceInterface,
) IdentityReviewActionPlanStatusUpdateServiceInterface {
	return &IdentityReviewActionPlanStatusUpdateService{
		repo:                        repo,
		activityLogSvc:              activityLogSvc,
		identityReviewActionPlanSvc: identityReviewActionPlanSvc,
		actionPlanStatusSvc:         actionPlanStatusSvc,
		taskSvc:                     taskSvc,
		taskStatusSvc:               taskStatusSvc,
	}
}

func (irapsus *IdentityReviewActionPlanStatusUpdateService) GetAll(identityReviewUserID uint, actionPlanID uint, sort string, sortBy string, keyword string) ([]*dtos.IdentityReviewActionPlanStatusUpdateListDTO, error) {
	if _, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	return irapsus.repo.GetAll(actionPlanID, sort, sortBy, keyword)
}

func (irapsus IdentityReviewActionPlanStatusUpdateService) GetByID(identityReviewUserID uint, actionPlanID uint, id uint) (*models.IdentityReviewActionPlanStatusUpdate, error) {
	if _, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	identityReviewActionPlanStatusUpdate, err := irapsus.repo.GetByID(actionPlanID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return identityReviewActionPlanStatusUpdate, nil
}

func (irapsus *IdentityReviewActionPlanStatusUpdateService) Create(identityReviewUserID uint, actionPlanID uint, currentUserID uint, identityReviewActionPlanStatusUpdateFormDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO) (uint, error) {
	if err := irapsus.validateIdentityReviewActionPlanStatusUpdateFormDTO(identityReviewActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	actionPlan, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID)
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := irapsus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return 0, err
	}

	tx := irapsus.repo.BeginTransaction()

	//@todo validate document form

	identityReviewActionPlanUpdate := &models.IdentityReviewActionPlanStatusUpdate{
		IdentityReviewActionPlanID: actionPlanID,
		ActionPlanStatusID:         identityReviewActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                    identityReviewActionPlanStatusUpdateFormDTO.Comment,
		Percentage:                 identityReviewActionPlanStatusUpdateFormDTO.Percentage,
		CreatedBYID:                currentUserID,
	}

	id, err := irapsus.repo.Create(identityReviewActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = irapsus.identityReviewActionPlanSvc.StatusUpdate(identityReviewUserID, actionPlanID, currentUserID, identityReviewActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IdentityReviewActionPlanStatusUpdateDocument

	for _, doc := range identityReviewActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.IdentityReviewActionPlanStatusUpdateDocument{
			IdentityReviewActionPlanStatusUpdateID: id,
			URL:                                    url,
			Name:                                   doc.FileName,
			Size:                                   doc.Size,
			Token:                                  encryptedFilename,
			OriginalName:                           doc.OriginalName,
			Type:                                   doc.Type,
			CreatedByID:                            currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := irapsus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if identityReviewActionPlanStatusUpdateFormDTO.ActionPlanStatusID == actionPlanStatus.ID && actionPlan.ReviewUserID != nil {
		urlParam := fmt.Sprintf(`{"identity_review_user_id":%d, "action_plan_id":%d}`, identityReviewUserID, actionPlanID)
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:           "identity_access_review",
			IdentityReviewActionPlanID: &actionPlanID,
			Title:                      actionPlan.Title,
			Description:                actionPlan.Description,
			TargetDate:                 actionPlan.TargetDate,
			URLParam:                   &urlParam,
			QueryParam:                 nil,
			SeverityLevelID:            actionPlan.SeverityLevelID,
			ResponsibleUserID:          *actionPlan.ReviewUserID,
		}

		if _, err := irapsus.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := irapsus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
			id,
			identityReviewActionPlanUpdate,
			currentUserID,
			"identity_access_review", "identity_review_user", "action_plan_status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (irapsus *IdentityReviewActionPlanStatusUpdateService) Update(identityReviewUserID uint, actionPlanID uint, id uint, currentUserID uint, identityReviewActionPlanStatusUpdateFormDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO) error {
	tx := irapsus.repo.BeginTransaction()

	if _, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return err
	}

	oldIdentityReviewActionPlanStatusUpdate, err := irapsus.GetByID(identityReviewUserID, actionPlanID, id)
	if err != nil {
		return err
	}

	existingIdentityReviewActionPlanStatusUpdate, err := irapsus.GetByID(identityReviewUserID, actionPlanID, id)
	if err != nil {
		return err
	}

	if err := irapsus.validateIdentityReviewActionPlanStatusUpdateFormDTO(identityReviewActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Identity Review StatusUpdate -- comment, percentage and status
	existingIdentityReviewActionPlanStatusUpdate.ActionPlanStatusID = identityReviewActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingIdentityReviewActionPlanStatusUpdate.Comment = identityReviewActionPlanStatusUpdateFormDTO.Comment
	existingIdentityReviewActionPlanStatusUpdate.Percentage = identityReviewActionPlanStatusUpdateFormDTO.Percentage
	existingIdentityReviewActionPlanStatusUpdate.UpdatedBYID = &currentUserID

	err = irapsus.repo.Update(existingIdentityReviewActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := irapsus.identityReviewActionPlanSvc.StatusUpdate(identityReviewUserID, actionPlanID, currentUserID, identityReviewActionPlanStatusUpdateFormDTO, tx); err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.IdentityReviewActionPlanStatusUpdateDocument
	var deleteDocumentIds []uint
	var documentsToDelete []*models.IdentityReviewActionPlanStatusUpdateDocument

	for _, doc := range identityReviewActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate URL
			url := utility.GenerateDocumentURL("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new document instance
			document := &models.IdentityReviewActionPlanStatusUpdateDocument{
				IdentityReviewActionPlanStatusUpdateID: id,
				URL:                                    url,
				Name:                                   doc.FileName,
				Size:                                   doc.Size,
				Token:                                  encryptedFileName,
				OriginalName:                           doc.OriginalName,
				Type:                                   doc.Type,
				CreatedByID:                            currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := irapsus.repo.GetDocumentByID(*doc.ID, tx)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := irapsus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := irapsus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingIdentityReviewActionPlanStatusUpdate, err := irapsus.GetByID(identityReviewUserID, actionPlanID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":               oldIdentityReviewActionPlanStatusUpdate.Comment,
			"percentage":            oldIdentityReviewActionPlanStatusUpdate.Percentage,
			"action_plan_status_id": oldIdentityReviewActionPlanStatusUpdate.ActionPlanStatusID,
		}

		newData := map[string]interface{}{
			"comment":               existingIdentityReviewActionPlanStatusUpdate.Comment,
			"percentage":            existingIdentityReviewActionPlanStatusUpdate.Percentage,
			"action_plan_status_id": existingIdentityReviewActionPlanStatusUpdate.ActionPlanStatusID,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := irapsus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				"identity_access_review", "identity_review_user", "action_plan_status_update", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
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

func (irapsus IdentityReviewActionPlanStatusUpdateService) GetLatestIdentityReviewActionPlanStatusUpdate(identityReviewUserID uint, actionPlanID uint) (*dtos.IdentityReviewActionPlanStatusUpdateListDTO, error) {
	if _, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID); err != nil {
		return nil, err
	}

	identityReviewActionPlanStatusUpdate, err := irapsus.repo.GetLatestIdentityReviewActionPlanStatusUpdate(actionPlanID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, actionPlanID)
		}
		return nil, err
	}
	return identityReviewActionPlanStatusUpdate, nil
}

func (irapsus *IdentityReviewActionPlanStatusUpdateService) validateIdentityReviewActionPlanStatusUpdateFormDTO(identityReviewActionPlanStatusUpdateFormDTO *dtos.IdentityReviewActionPlanStatusUpdateFormDTO) error {
	if _, err := irapsus.actionPlanStatusSvc.GetByID(identityReviewActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action_plan_status_id")
	}

	return nil
}

func (irapsus *IdentityReviewActionPlanStatusUpdateService) GetDocumentPath(identityReviewUserID uint, actionPlanID uint, id uint, token string) (string, error) {
	existingIdentityReviewActionPlanStatusUpdate, err := irapsus.GetByID(identityReviewUserID, actionPlanID, id)
	if err != nil {
		return "", err
	}

	if len(existingIdentityReviewActionPlanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the identity review action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("identity-access-review", "identity-review-users", strconv.FormatUint(uint64(identityReviewUserID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (irapsus *IdentityReviewActionPlanStatusUpdateService) Approve(identityReviewUserID uint, actionPlanID uint, currentUserID uint, identityReviewActionPlanReviewFormDTO *dtos.IdentityReviewActionPlanReviewFormDTO) error {
	existingActionPlan, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this identity review action plan.")
	}

	if currentUserID != *existingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	completedStatus, err := irapsus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	closedStatus, err := irapsus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := irapsus.repo.BeginTransaction()

	if existingActionPlan.ActionPlanStatusID != completedStatus.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The identity review action plan status must be resolved.")
	}

	identityReviewActionPlanUpdate := &models.IdentityReviewActionPlanStatusUpdate{
		IdentityReviewActionPlanID: actionPlanID,
		ActionPlanStatusID:         closedStatus.ID,
		Percentage:                 &existingActionPlan.Percentage,
		Comment:                    identityReviewActionPlanReviewFormDTO.Comment,
		CreatedBYID:                currentUserID,
	}

	id, err := irapsus.repo.Create(identityReviewActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	statusForm := &dtos.IdentityReviewActionPlanStatusUpdateFormDTO{
		ActionPlanStatusID: closedStatus.ID,
		Comment:            identityReviewActionPlanReviewFormDTO.Comment,
		Percentage:         &existingActionPlan.Percentage,
	}

	if err := irapsus.identityReviewActionPlanSvc.StatusUpdate(identityReviewUserID, actionPlanID, currentUserID, statusForm, tx); err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := irapsus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = irapsus.taskSvc.UpdateActionPlanTaskStatus("identity_review_action_plan_id", actionPlanID, existingActionPlan.ResponsibleUserID, taskStatus.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = irapsus.taskSvc.UpdateActionPlanTaskStatus("identity_review_action_plan_id", actionPlanID, *existingActionPlan.ReviewUserID, taskStatus.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		existingActionPlan, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID)
		if err != nil {
			return
		}

		err = irapsus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			"identity_access_review", "identity_review_user", "action_plan_status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (irapsus *IdentityReviewActionPlanStatusUpdateService) Revert(identityReviewUserID uint, actionPlanID uint, currentUserID uint, identityReviewActionPlanReviewFormDTO *dtos.IdentityReviewActionPlanReviewFormDTO) error {
	existingActionPlan, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, actionPlanID)
	}

	if existingActionPlan.ReviewUserID == nil {
		return utility.NewForbiddenError("Review user is not assigned to this identity review action plan.")
	}

	if currentUserID != *existingActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	completedStatus, err := irapsus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	reopenStatus, err := irapsus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := irapsus.repo.BeginTransaction()

	if existingActionPlan.ActionPlanStatusID != completedStatus.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The identity review action plan status must be resolved.")
	}

	identityReviewActionPlanUpdate := &models.IdentityReviewActionPlanStatusUpdate{
		IdentityReviewActionPlanID: actionPlanID,
		ActionPlanStatusID:         reopenStatus.ID,
		Percentage:                 &existingActionPlan.Percentage,
		Comment:                    identityReviewActionPlanReviewFormDTO.Comment,
		CreatedBYID:                currentUserID,
	}

	id, err := irapsus.repo.Create(identityReviewActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	statusForm := &dtos.IdentityReviewActionPlanStatusUpdateFormDTO{
		ActionPlanStatusID: reopenStatus.ID,
		Comment:            identityReviewActionPlanReviewFormDTO.Comment,
		Percentage:         &existingActionPlan.Percentage,
	}

	if err := irapsus.identityReviewActionPlanSvc.StatusUpdate(identityReviewUserID, actionPlanID, currentUserID, statusForm, tx); err != nil {
		tx.Rollback()
		return err
	}

	taskStatus, err := irapsus.taskStatusSvc.GetByType("closed")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = irapsus.taskSvc.UpdateActionPlanTaskStatus("identity_review_action_plan_id", actionPlanID, *existingActionPlan.ReviewUserID, taskStatus.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlan, err := irapsus.identityReviewActionPlanSvc.GetByID(identityReviewUserID, actionPlanID)
		if err != nil {
			return
		}

		err = irapsus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
			id,
			existingActionPlan,
			currentUserID,
			"identity_access_review", "identity_review_user", "action_plan_status_update", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.RejectAction,
				"tableName": models.IdentityReviewActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
