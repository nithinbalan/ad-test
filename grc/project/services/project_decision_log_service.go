package services

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	dtos "jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectDecisionLogServiceInterface interface {
	SetWorkflowService(workflowSvc ProjectDecisionLogWorkflowServiceInterface)
	GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error)
	GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.ProjectDecisionLogDetailDTO, error)
	Create(projectID uint, currentUserID uint, form *dtos.ProjectDecisionLogFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectDecisionLogFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(projectID uint, decisionLogID uint, documentID uint, token string) (string, error)
	GetSummary(projectID uint) (*dtos.ProjectDecisionLogSummaryDTO, error)
}

type ProjectDecisionLogService struct {
	repo        repositories.ProjectDecisionLogRepositoryInterface
	activityLog logServices.ActivityLogServiceInterface
	userSvc     organizationServices.UserServiceInterface
	projectSvc  ProjectServiceInterface
	workflowSvc ProjectDecisionLogWorkflowServiceInterface
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

func (pdls *ProjectDecisionLogService) SetWorkflowService(workflowSvc ProjectDecisionLogWorkflowServiceInterface) {
	pdls.workflowSvc = workflowSvc
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

func (pdls *ProjectDecisionLogService) GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.ProjectDecisionLogDetailDTO, error) {
	if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	decisionLog, err := pdls.repo.GetFullDetailsByID(projectID, id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	var detailDTO dtos.ProjectDecisionLogDetailDTO
	if err := mapstructure.Decode(decisionLog, &detailDTO); err != nil {
		return nil, err
	}

	decisionDateStr := decisionLog.DecisionDate.Format("2006-01-02")
	detailDTO.DecisionDateString = &decisionDateStr

	if decisionLog.SubmittedAt != nil {
		str := decisionLog.SubmittedAt.Format("2006-01-02 15:04:05")
		detailDTO.SubmittedAtString = &str
	}
	if decisionLog.ApprovedAt != nil {
		str := decisionLog.ApprovedAt.Format("2006-01-02 15:04:05")
		detailDTO.ApprovedAtString = &str
	}
	if decisionLog.RejectedAt != nil {
		str := decisionLog.RejectedAt.Format("2006-01-02 15:04:05")
		detailDTO.RejectedAtString = &str
	}

	detailDTO.CreatedAt = decisionLog.CreatedAt
	detailDTO.UpdatedAt = decisionLog.UpdatedAt

	return &detailDTO, nil
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
	var newDocsToCreate []*models.ProjectDecisionLogDocument
	var docsToDelete []*uint
	var documentsToDelete []*models.ProjectDecisionLogDocument

	for _, doc := range form.Documents {
		if (doc.IsNewFile != nil && *doc.IsNewFile) || (doc.ID == nil && doc.FileName != "") {
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
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := pdls.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return utility.NewResourceNotFoundError(err, id)
					}
					tx.Rollback()
					return err
				}
				documentsToDelete = append(documentsToDelete, existingDocument)
				docsToDelete = append(docsToDelete, &existingDocument.ID)
			}
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

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to remove document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
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

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		decisionLog, err := pdls.repo.GetByID(projectID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// 1. Delete review comments and review comment documents
		tx := pdls.repo.BeginTransaction()

		if err := pdls.repo.DeleteReviewComments(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 2. Delete workflow histories
		if err := pdls.repo.DeleteWorkflowHistories(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 3. Delete workflows, reviews, and participants
		if pdls.workflowSvc != nil {
			if err := pdls.workflowSvc.DeleteAllWorkflowsForProjectDecisionLog(id, tx); err != nil {
				tx.Rollback()
				failedIds = append(failedIds, id)
				continue
			}
		}

		// 4. Delete documents
		if err := pdls.repo.DeleteDocumentsByDecisionLogID(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// 5. Delete parent decision log
		if err := pdls.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := tx.Commit(); err != nil {
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
		return failedIds, nil
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

	go func() {
		for _, item := range deletedItems {
			_ = pdls.activityLog.LogActivity(
				models.DeleteAction,
				models.ProjectDecisionLog{}.TableName(),
				item["id"].(uint),
				item,
				currentUserID,
				"project", "project_register", "", nil,
			)
		}
	}()

	return nil, nil
}

func (pdls *ProjectDecisionLogService) GetSummary(projectID uint) (*dtos.ProjectDecisionLogSummaryDTO, error) {
	if projectID != 0 {
		if _, err := pdls.projectSvc.GetByID(projectID); err != nil {
			return nil, err
		}
	}

	return pdls.repo.GetSummary(projectID)
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

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "decision-logs", strconv.FormatUint(uint64(decisionLogID), 10), "documents")
	fullPath := filepath.Join(baseStoragePath, storagePath, document.Name)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	}

	return fullPath, nil
}
