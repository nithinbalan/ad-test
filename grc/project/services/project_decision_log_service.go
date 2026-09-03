package services

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	dtos "jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectDecisionLogServiceInterface interface {
	GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error)
	GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectDecisionLogDetailDTO, error)
	Create(projectID uint, currentUserID uint, form *dtos.ProjectDecisionLogFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectDecisionLogFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(projectID uint, decisionLogID uint, documentID uint, token string) (string, error)
}

type ProjectDecisionLogService struct {
	repo        repositories.ProjectDecisionLogRepositoryInterface
	activityLog logServices.ActivityLogServiceInterface
	userSvc     organizationServices.UserServiceInterface
	projectSvc  ProjectServiceInterface
}

func NewProjectDecisionLogService(
	repo repositories.ProjectDecisionLogRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectSvc ProjectServiceInterface,
) ProjectDecisionLogServiceInterface {
	return &ProjectDecisionLogService{
		repo:        repo,
		activityLog: activityLog,
		userSvc:     userSvc,
		projectSvc:  projectSvc,
	}
}

func (pdls *ProjectDecisionLogService) GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error) {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pdls.repo.GetAll(projectID, params)
}

func (pdls *ProjectDecisionLogService) GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error) {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pdls.repo.GetByID(projectID, id)
}

func (pdls *ProjectDecisionLogService) GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectDecisionLogDetailDTO, error) {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	decisionLog, err := pdls.repo.GetFullDetailsByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	return pdls.mapToDetailDTO(decisionLog)
}

func (pdls *ProjectDecisionLogService) Create(projectID uint, currentUserID uint, form *dtos.ProjectDecisionLogFormDTO) (uint, error) {
	project, err := pdls.projectSvc.GetByID(projectID)
	if err != nil {
		return 0, err
	}

	if _, err := pdls.userSvc.GetByID(form.RequestedByUserID); err != nil {
		return 0, utility.NewValidationError(errors.New("requested by user not found"), "user_not_found", "project_decision_log")
	}

	decisionDate, err := time.Parse("2006-01-02", form.DecisionDate)
	if err != nil {
		return 0, utility.NewValidationError(errors.New("invalid decision date format"), "invalid_date", "project_decision_log")
	}

	refNumber, err := pdls.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	decisionLog := &models.ProjectDecisionLog{
		ReferenceNumber:   refNumber,
		ProjectID:         project.ID,
		Decision:          form.Decision,
		DecisionNotes:     form.DecisionNotes,
		Reason:            form.Reason,
		RequestedByUserID: form.RequestedByUserID,
		DecisionDate:      decisionDate,
		CreatedByID:       currentUserID,
		CreatedAt:         time.Now(),
	}

	tx := pdls.repo.BeginTransaction()

	id, err := pdls.repo.Create(decisionLog, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectDecisionLogDocument
	for _, doc := range form.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.ProjectDecisionLogDocument{
			ProjectDecisionLogID: id,
			URL:                  url,
			Name:                 doc.FileName,
			Size:                 doc.Size,
			Token:                encryptedFileName,
			OriginalName:         &doc.OriginalName,
			Type:                 doc.Type,
			CreatedByID:          currentUserID,
			CreatedAt:            time.Now(),
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := pdls.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	for _, doc := range documents {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		_ = pdls.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectDecisionLog{}.TableName(),
			id,
			decisionLog,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	return id, nil
}

func (pdls *ProjectDecisionLogService) Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectDecisionLogFormDTO) error {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	decisionLog, err := pdls.repo.GetFullDetailsByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	if _, err := pdls.userSvc.GetByID(form.RequestedByUserID); err != nil {
		return utility.NewValidationError(errors.New("requested by user not found"), "user_not_found", "project_decision_log")
	}

	decisionDate, err := time.Parse("2006-01-02", form.DecisionDate)
	if err != nil {
		return utility.NewValidationError(errors.New("invalid decision date format"), "invalid_date", "project_decision_log")
	}

	oldData := map[string]interface{}{
		"decision":             decisionLog.Decision,
		"decision_notes":       decisionLog.DecisionNotes,
		"reason":               decisionLog.Reason,
		"requested_by_user_id": decisionLog.RequestedByUserID,
		"decision_date":        decisionLog.DecisionDate.Format("2006-01-02"),
	}

	decisionLog.Decision = form.Decision
	decisionLog.DecisionNotes = form.DecisionNotes
	decisionLog.Reason = form.Reason
	decisionLog.RequestedByUserID = form.RequestedByUserID
	decisionLog.DecisionDate = decisionDate
	decisionLog.UpdatedByID = &currentUserID
	updatedAt := time.Now()
	decisionLog.UpdatedAt = &updatedAt

	tx := pdls.repo.BeginTransaction()

	if err := pdls.repo.Update(decisionLog, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Document handling
	var existingDocsToKeep []*uint
	var newDocsToCreate []*models.ProjectDecisionLogDocument
	for _, doc := range form.Documents {
		if doc.ID != nil {
			existingDocsToKeep = append(existingDocsToKeep, doc.ID)
		} else {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			newDoc := &models.ProjectDecisionLogDocument{
				ProjectDecisionLogID: id,
				URL:                  url,
				Name:                 doc.FileName,
				Size:                 doc.Size,
				Token:                encryptedFileName,
				OriginalName:         &doc.OriginalName,
				Type:                 doc.Type,
				CreatedByID:          currentUserID,
				CreatedAt:            time.Now(),
			}
			newDocsToCreate = append(newDocsToCreate, newDoc)
		}
	}

	var docsToDelete []*uint
	for _, doc := range decisionLog.Documents {
		keep := false
		for _, keepID := range existingDocsToKeep {
			if *keepID == doc.ID {
				keep = true
				break
			}
		}
		if !keep {
			docID := doc.ID
			docsToDelete = append(docsToDelete, &docID)
		}
	}

	if len(docsToDelete) > 0 {
		if err := pdls.repo.DeleteDocuments(docsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(newDocsToCreate) > 0 {
		if err := pdls.repo.CreateDocuments(newDocsToCreate, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, doc := range newDocsToCreate {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		newData := map[string]interface{}{
			"decision":             form.Decision,
			"decision_notes":       form.DecisionNotes,
			"reason":               form.Reason,
			"requested_by_user_id": form.RequestedByUserID,
			"decision_date":        form.DecisionDate,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			_ = pdls.activityLog.LogActivity(
				models.UpdateAction,
				models.ProjectDecisionLog{}.TableName(),
				id,
				changes,
				currentUserID,
				"project", "project_register", "", nil,
			)
		}
	}()

	return nil
}

func (pdls *ProjectDecisionLogService) Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	tx := pdls.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		decisionLog, err := pdls.repo.GetByID(projectID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// 1. Delete documents
		if err := pdls.repo.DeleteDocumentsByDecisionLogID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// 2. Delete parent decision log
		if err := pdls.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedItems = append(deletedItems, map[string]interface{}{
			"id":               id,
			"reference_number": decisionLog.ReferenceNumber,
			"decision":         decisionLog.Decision,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for _, id := range successfulIds {
		folderPath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(id), 10))
		if err := utility.RemoveFolder(folderPath); err != nil {
			utility.LogError("Failed to remove decision log storage folder", map[string]interface{}{
				"decisionLogID": id,
				"path":          folderPath,
				"error":         err.Error(),
			})
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := pdls.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectDecisionLog{}.TableName(),
				successfulIds,
				deletedItems,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete project decision logs activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (pdls *ProjectDecisionLogService) GetDocumentPath(projectID uint, decisionLogID uint, documentID uint, token string) (string, error) {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return "", err
	}

	if _, err := pdls.repo.GetByID(projectID, decisionLogID); err != nil {
		return "", err
	}

	document, err := pdls.repo.GetDocumentByID(documentID)
	if err != nil {
		return "", err
	}

	if document.Token != token {
		return "", &utility.ResourceNotFoundError{}
	}

	storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(decisionLogID), 10), "documents", document.Name)
	fullPath := filepath.Join(os.Getenv("STORAGE_ABSOLUTE_PATH"), storagePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	}

	return fullPath, nil
}

func (pdls *ProjectDecisionLogService) mapToDetailDTO(decisionLog *models.ProjectDecisionLog) (*dtos.ProjectDecisionLogDetailDTO, error) {
	var requestedByUser *organizationDTOs.UserMinimalDTO
	if decisionLog.RequestedByUser != nil {
		_ = mapstructure.Decode(decisionLog.RequestedByUser, &requestedByUser)
	}

	var createdByUser *organizationDTOs.UserMinimalDTO
	if decisionLog.CreatedByUser != nil {
		_ = mapstructure.Decode(decisionLog.CreatedByUser, &createdByUser)
	}

	var updatedByUser *organizationDTOs.UserMinimalDTO
	if decisionLog.UpdatedByUser != nil {
		_ = mapstructure.Decode(decisionLog.UpdatedByUser, &updatedByUser)
	}

	var documents []*dtos.ProjectDecisionLogDocumentDTO
	for _, doc := range decisionLog.Documents {
		documents = append(documents, &dtos.ProjectDecisionLogDocumentDTO{
			ID:           doc.ID,
			Name:         doc.Name,
			OriginalName: doc.OriginalName,
			Size:         doc.Size,
			Type:         doc.Type,
			URL:          doc.URL,
			Token:        doc.Token,
		})
	}

	var updatedAtStr *string
	if decisionLog.UpdatedAt != nil {
		str := decisionLog.UpdatedAt.Format("2006-01-02 15:04:05")
		updatedAtStr = &str
	}

	return &dtos.ProjectDecisionLogDetailDTO{
		ID:              decisionLog.ID,
		ReferenceNumber: decisionLog.ReferenceNumber,
		ProjectID:       decisionLog.ProjectID,
		Decision:        decisionLog.Decision,
		DecisionNotes:   decisionLog.DecisionNotes,
		Reason:          decisionLog.Reason,
		DecisionDate:    decisionLog.DecisionDate.Format("2006-01-02"),
		RequestedByUser: requestedByUser,
		CreatedByUser:   createdByUser,
		UpdatedByUser:   updatedByUser,
		CreatedAt:       decisionLog.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       updatedAtStr,
		Documents:       documents,
	}, nil
}
