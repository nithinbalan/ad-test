package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	projectDTOs "jethur/grc/project/dtos"
	projectRepositories "jethur/grc/project/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectActionPlanStatusUpdateServiceInterface interface {
	GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error)
	Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO) error
	GetLatestProjectActionPlanStatusUpdate(projectID uint, actionPlanID uint) (*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(projectID uint, actionPlanID uint, id uint, token string) (string, error)

	Approve(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, projectActionPlanReviewUserStatusUpdateFormDTO *projectDTOs.ProjectActionPlanReviewUserStatusUpdateFormDTO) error
	Revert(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, projectActionPlanReviewUserStatusUpdateFormDTO *projectDTOs.ProjectActionPlanReviewUserStatusUpdateFormDTO) error
}

type ProjectActionPlanStatusUpdateService struct {
	repo                 projectRepositories.ProjectActionPlanStatusUpdateRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	projectSvc           ProjectServiceInterface
	projectActionPlanSvc ProjectActionPlanServiceInterface
	actionPlanStatusSvc  masterServices.ActionPlanStatusServiceInterface
	userSvc              organizationServices.UserServiceInterface
}

func NewProjectActionPlanStatusUpdateService(
	repo projectRepositories.ProjectActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	projectSvc ProjectServiceInterface,
	projectActionPlanSvc ProjectActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) ProjectActionPlanStatusUpdateServiceInterface {
	return &ProjectActionPlanStatusUpdateService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		projectSvc:           projectSvc,
		projectActionPlanSvc: projectActionPlanSvc,
		actionPlanStatusSvc:  actionPlanStatusSvc,
		userSvc:              userSvc,
	}
}

func (pasus *ProjectActionPlanStatusUpdateService) GetAll(ProjectID uint, ProjectActionPlanID uint, sort string, sortBy string, keyword string) ([]*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error) {
	if _, err := pasus.projectActionPlanSvc.GetByID(ProjectID, ProjectActionPlanID); err != nil {
		return nil, err
	}

	return pasus.repo.GetAll(ProjectActionPlanID, sort, sortBy, keyword)
}

func (pasus *ProjectActionPlanStatusUpdateService) GetByID(ProjectID uint, ProjectActionPlanID uint, id uint) (*models.ProjectActionPlanStatusUpdate, error) {
	if _, err := pasus.projectActionPlanSvc.GetByID(ProjectID, ProjectActionPlanID); err != nil {
		return nil, err
	}

	projectActionPlanStatusUpdate, err := pasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return projectActionPlanStatusUpdate, nil
}

func (pasus *ProjectActionPlanStatusUpdateService) Create(ProjectID uint, ProjectActionPlanID uint, currentUserID uint, projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := pasus.repo.BeginTransaction()

	if _, err := pasus.projectActionPlanSvc.GetByID(ProjectID, ProjectActionPlanID); err != nil {
		return 0, err
	}

	if err := pasus.validateProjectActionPlanStatusUpdateFormDTO(projectActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	projectActionPlanUpdate := &models.ProjectActionPlanStatusUpdate{
		ProjectActionPlanID: ProjectActionPlanID,
		ActionPlanStatusID:  projectActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:             projectActionPlanStatusUpdateFormDTO.Comment,
		Percentage:          projectActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:         currentUserID,
	}

	id, err := pasus.repo.Create(projectActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = pasus.projectActionPlanSvc.StatusUpdate(ProjectID, ProjectActionPlanID, currentUserID, projectActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectActionPlanStatusUpdateDocument

	for _, doc := range projectActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(ProjectID), 10), "action-plans", strconv.FormatUint(uint64(ProjectActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ProjectActionPlanStatusUpdateDocument{
			ProjectActionPlanStatusUpdateID: id,
			URL:                             url,
			Name:                            doc.FileName,
			Size:                            doc.Size,
			Token:                           encryptedFilename,
			OriginalName:                    doc.OriginalName,
			Type:                            doc.Type,
			CreatedByID:                     currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := pasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		// move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("project", strconv.FormatUint(uint64(ProjectID), 10), "action-plans", strconv.FormatUint(uint64(ProjectActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	tx.Commit()

	go func() {
		err := pasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProjectActionPlanStatusUpdate{}.TableName(),
			id,
			projectActionPlanUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ProjectActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (pasus *ProjectActionPlanStatusUpdateService) Update(ProjectID uint, ProjectActionPlanID uint, id uint, currentUserID uint, projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO) error {
	tx := pasus.repo.BeginTransaction()

	if _, err := pasus.projectActionPlanSvc.GetByID(ProjectID, ProjectActionPlanID); err != nil {
		return err
	}

	oldProjectActionPlanStatusUpdate, err := pasus.GetByID(ProjectID, ProjectActionPlanID, id)
	if err != nil {
		return err
	}

	existingProjectActionPlanStatusUpdate, err := pasus.GetByID(ProjectID, ProjectActionPlanID, id)
	if err != nil {
		return err
	}

	if err := pasus.validateProjectActionPlanStatusUpdateFormDTO(projectActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingProjectActionPlanStatusUpdate.ActionPlanStatusID = projectActionPlanStatusUpdateFormDTO.ActionPlanStatusID
	existingProjectActionPlanStatusUpdate.Comment = projectActionPlanStatusUpdateFormDTO.Comment
	existingProjectActionPlanStatusUpdate.Percentage = projectActionPlanStatusUpdateFormDTO.Percentage
	existingProjectActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = pasus.repo.Update(existingProjectActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = pasus.projectActionPlanSvc.StatusUpdate(ProjectID, ProjectActionPlanID, currentUserID, projectActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ProjectActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ProjectActionPlanStatusUpdateDocument

	for _, doc := range projectActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// genereate URL
			url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(ProjectID), 10), "action-plans", strconv.FormatUint(uint64(ProjectActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new Control test Document instance
			document := &models.ProjectActionPlanStatusUpdateDocument{
				ProjectActionPlanStatusUpdateID: id,
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
				existingDocument, err := pasus.repo.GetDocumentByID(*doc.ID)
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
		if err := pasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := pasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(ProjectID), 10), "action-plans", strconv.FormatUint(uint64(ProjectActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(ProjectID), 10), "action-plans", strconv.FormatUint(uint64(ProjectActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := pasus.GetByID(ProjectID, ProjectActionPlanID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldProjectActionPlanStatusUpdate.Comment,
			"percentage": oldProjectActionPlanStatusUpdate.Percentage,
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
			err := pasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProjectActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ProjectActionPlanStatusUpdate{}.TableName(),
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

func (pasus *ProjectActionPlanStatusUpdateService) validateProjectActionPlanStatusUpdateFormDTO(projectActionPlanStatusUpdateFormDTO *projectDTOs.ProjectActionPlanStatusUpdateFormDTO) error {
	if _, err := pasus.actionPlanStatusSvc.GetByID(projectActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (pasus *ProjectActionPlanStatusUpdateService) GetLatestProjectActionPlanStatusUpdate(projectID uint, actionPlanID uint) (*projectDTOs.ProjectActionPlanStatusUpdateListDTO, error) {
	if _, err := pasus.projectActionPlanSvc.GetByID(projectID, actionPlanID); err != nil {
		return nil, err
	}

	return pasus.repo.GetLatestProjectActionPlanStatusUpdate(actionPlanID)
}

func (pasus *ProjectActionPlanStatusUpdateService) GetDocumentPath(projectID, actionPlanID, id uint, token string) (string, error) {
	existingProjectActionplanStatusUpdate, err := pasus.GetByID(projectID, actionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingProjectActionplanStatusUpdate.Documents) == 0 {
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

	storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

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

func (pasus *ProjectActionPlanStatusUpdateService) Approve(projectID uint, projectActionPlanID uint, currentUserID uint, projectActionPlanReviewUserStatusUpdateFormDTO *projectDTOs.ProjectActionPlanReviewUserStatusUpdateFormDTO) error {
	existingProjectActionPlan, err := pasus.projectActionPlanSvc.GetByID(projectID, projectActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, projectActionPlanID)
	}

	if currentUserID != *existingProjectActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	projectActionPlanStatusResolved, err := pasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	projectActionPlanClosedStatus, err := pasus.actionPlanStatusSvc.GetByType("close")
	if err != nil {
		return err
	}

	tx := pasus.repo.BeginTransaction()

	if existingProjectActionPlan.ActionPlanStatusID != projectActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The project action plan status must be resolved.")
	}

	existingProjectActionPlan.ActionPlanStatusID = projectActionPlanClosedStatus.ID

	if err := pasus.repo.ProjectActionPlanUpdateStatus(existingProjectActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	projectActionPlanUpdate := &models.ProjectActionPlanStatusUpdate{
		ProjectActionPlanID: projectActionPlanID,
		ActionPlanStatusID:  projectActionPlanClosedStatus.ID,
		Percentage:          nil,
		Comment:             projectActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:         currentUserID,
	}

	id, err := pasus.repo.Create(projectActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		err = pasus.activityLogSvc.LogActivity(
			models.ApproveAction,
			models.ProjectActionPlanStatusUpdate{}.TableName(),
			id,
			existingProjectActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.ApproveAction,
				"tableName": models.ProjectActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (pasus *ProjectActionPlanStatusUpdateService) Revert(projectID uint, projectActionPlanID uint, currentUserID uint, projectActionPlanReviewUserStatusUpdateFormDTO *projectDTOs.ProjectActionPlanReviewUserStatusUpdateFormDTO) error {
	existingProjectActionPlan, err := pasus.projectActionPlanSvc.GetByID(projectID, projectActionPlanID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, projectActionPlanID)
	}

	if currentUserID != *existingProjectActionPlan.ReviewUserID {
		return utility.NewForbiddenError("Permission denied.")
	}

	projectActionPlanStatusResolved, err := pasus.actionPlanStatusSvc.GetByType("completed")
	if err != nil {
		return err
	}

	projectActionPlanReopenStatus, err := pasus.actionPlanStatusSvc.GetByType("reopen")
	if err != nil {
		return err
	}

	tx := pasus.repo.BeginTransaction()

	if existingProjectActionPlan.ActionPlanStatusID != projectActionPlanStatusResolved.ID {
		tx.Rollback()
		return utility.NewForbiddenError("The project action plan status must be resolved.")
	}

	existingProjectActionPlan.ActionPlanStatusID = projectActionPlanReopenStatus.ID

	if err := pasus.repo.ProjectActionPlanUpdateStatus(existingProjectActionPlan); err != nil {
		tx.Rollback()
		return err
	}

	projectActionPlanUpdate := &models.ProjectActionPlanStatusUpdate{
		ProjectActionPlanID: projectActionPlanID,
		ActionPlanStatusID:  projectActionPlanReopenStatus.ID,
		Comment:             projectActionPlanReviewUserStatusUpdateFormDTO.Comment,
		CreatedByID:         currentUserID,
	}

	id, err := pasus.repo.Create(projectActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err = pasus.activityLogSvc.LogActivity(
			models.RejectAction,
			models.ProjectActionPlanStatusUpdate{}.TableName(),
			id,
			existingProjectActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.RejectAction,
				"tableName": models.ProjectActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}
