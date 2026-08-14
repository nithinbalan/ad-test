package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditFindingReassessmentServiceInterface interface {
	GetAll(findingID uint, params *auditDTOs.AuditFindingReassessmentListQueryParamsDTO) (*auditDTOs.AuditFindingReassessmentListResponseDTO, error)
	GetByID(findingID uint, id uint) (*models.AuditFindingReassessment, error)
	GetFullDetailsByID(findingID uint, id uint) (*auditDTOs.AuditFindingReassessmentDetailDTO, error)
	Create(findingID uint, currentUserID uint, formDTO *auditDTOs.AuditFindingReassessmentFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, formDTO *auditDTOs.AuditFindingReassessmentFormDTO) error
	Delete(findingID uint, id uint, currentUserID uint) error
	GetDocumentPath(findingID uint, id uint, token string) (string, error)

	// Action plan
	GetActionPlans(findingID uint, reassessmentID uint, params *auditDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditDTOs.AuditFindingActionPlanListResponseDTO, error)
	CreateActionPlan(findingID uint, reassessmentID uint, currentUserID uint, formDTO *auditDTOs.AuditFindingActionPlanFormDTO) (uint, error)
}

type AuditFindingReassessmentService struct {
	repo                      repositories.AuditFindingReassessmentRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	auditFindingSvc           AuditFindingServiceInterface
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface
}

func NewAuditFindingReassessmentService(
	repo repositories.AuditFindingReassessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditFindingSvc AuditFindingServiceInterface,
	auditFindingActionPlanSvc AuditFindingActionPlanServiceInterface,
) AuditFindingReassessmentServiceInterface {
	return &AuditFindingReassessmentService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		auditFindingSvc:           auditFindingSvc,
		auditFindingActionPlanSvc: auditFindingActionPlanSvc,
	}
}

func (afrs *AuditFindingReassessmentService) GetAll(findingID uint, params *auditDTOs.AuditFindingReassessmentListQueryParamsDTO) (*auditDTOs.AuditFindingReassessmentListResponseDTO, error) {
	if _, err := afrs.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	return afrs.repo.GetAll(findingID, params)
}

func (afrs *AuditFindingReassessmentService) GetActionPlans(findingID uint, reassessmentID uint, params *auditDTOs.AuditFindingActionPlanListQueryParamsDTO) (*auditDTOs.AuditFindingActionPlanListResponseDTO, error) {
	if _, err := afrs.GetByID(findingID, reassessmentID); err != nil {
		return nil, err
	}

	if params.Filters == nil {
		params.Filters = make(map[string][]string)
	}
	params.Filters["audit_finding_action_plans.audit_finding_reassessment_id"] = []string{strconv.FormatUint(uint64(reassessmentID), 10)}

	return afrs.auditFindingActionPlanSvc.GetAll(findingID, params)
}

func (afrs *AuditFindingReassessmentService) CreateActionPlan(findingID uint, reassessmentID uint, currentUserID uint, formDTO *auditDTOs.AuditFindingActionPlanFormDTO) (uint, error) {
	if _, err := afrs.GetByID(findingID, reassessmentID); err != nil {
		return 0, err
	}

	formDTO.AuditFindingReassessmentID = &reassessmentID

	return afrs.auditFindingActionPlanSvc.Create(findingID, currentUserID, formDTO)
}

func (afrs *AuditFindingReassessmentService) GetByID(findingID uint, id uint) (*models.AuditFindingReassessment, error) {
	if _, err := afrs.auditFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditFindingReassessment, err := afrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if auditFindingReassessment.AuditFindingID != findingID {
		return nil, utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, id)
	}

	return auditFindingReassessment, nil
}

func (afrs *AuditFindingReassessmentService) GetFullDetailsByID(findingID uint, id uint) (*auditDTOs.AuditFindingReassessmentDetailDTO, error) {
	if _, err := afrs.GetByID(findingID, id); err != nil {
		return nil, err
	}

	auditFindingReassessment, err := afrs.repo.GetDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var details auditDTOs.AuditFindingReassessmentDetailDTO
	if err := mapstructure.Decode(auditFindingReassessment, &details); err != nil {
		return nil, err
	}

	details.CreatedAt = auditFindingReassessment.CreatedAt

	latestID, err := afrs.repo.GetLatestID(findingID)
	if err != nil {
		return nil, err
	}
	details.IsLatest = latestID != 0 && latestID == id

	return &details, nil
}

func (afrs *AuditFindingReassessmentService) Create(findingID uint, currentUserID uint, formDTO *auditDTOs.AuditFindingReassessmentFormDTO) (id uint, err error) {
	if _, err = afrs.auditFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}

	if err = afrs.ensureCreatable(findingID); err != nil {
		return 0, err
	}

	tx := afrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	auditFindingReassessment := &models.AuditFindingReassessment{
		AuditFindingID: findingID,
		Outcome:        &formDTO.Outcome,
		Comment:        &formDTO.Comment,
		CreatedBYID:    currentUserID,
	}

	id, err = afrs.repo.Create(auditFindingReassessment, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.AuditFindingReassessmentDocument

	for _, doc := range formDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, encErr := utility.EncryptData(dataToEncrypt)
		if encErr != nil {
			err = encErr
			return 0, err
		}

		url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "reassessments", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.AuditFindingReassessmentDocument{
			AuditFindingReassessmentID: id,
			URL:                        url,
			Name:                       doc.FileName,
			Size:                       doc.Size,
			Token:                      encryptedFilename,
			OriginalName:               doc.OriginalName,
			Type:                       doc.Type,
			CreatedBYID:                currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err = afrs.repo.CreateDocuments(documents, tx); err != nil {
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "reassessments", strconv.FormatUint(uint64(id), 10), "documents")
			if err = utility.StoreFile(doc.Name, storagePath); err != nil {
				return 0, err
			}
		}
	}

	go func() {
		existingAuditFindingReassessment, logErr := afrs.GetFullDetailsByID(findingID, id)
		if logErr != nil {
			utility.LogError("Failed to fetch audit finding reassessment", map[string]interface{}{
				"error": logErr.Error(),
			})
			return
		}

		logErr = afrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFindingReassessment{}.TableName(),
			id,
			existingAuditFindingReassessment,
			currentUserID,
			"audit", "audit_findings", "re-assessment", nil,
		)
		if logErr != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFindingReassessment{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     logErr.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, err
}

func (afrs *AuditFindingReassessmentService) Update(findingID uint, id uint, currentUserID uint, formDTO *auditDTOs.AuditFindingReassessmentFormDTO) (err error) {
	existingAuditFindingReassessment, err := afrs.GetByID(findingID, id)
	if err != nil {
		return err
	}

	if err := afrs.ensureLatestReassessment(findingID, id); err != nil {
		return err
	}

	oldAuditFindingReassessment, err := afrs.GetFullDetailsByID(findingID, id)
	if err != nil {
		return err
	}

	tx := afrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	existingAuditFindingReassessment.Outcome = &formDTO.Outcome
	existingAuditFindingReassessment.Comment = &formDTO.Comment
	existingAuditFindingReassessment.UpdatedBYID = &currentUserID

	if err = afrs.repo.Update(existingAuditFindingReassessment, tx); err != nil {
		return err
	}

	var documents []*models.AuditFindingReassessmentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditFindingReassessmentDocument

	for _, doc := range formDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, encErr := utility.EncryptData(dataToEncrypt)
			if encErr != nil {
				err = encErr
				return err
			}

			url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "reassessments", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.AuditFindingReassessmentDocument{
				AuditFindingReassessmentID: id,
				URL:                        url,
				Name:                       doc.FileName,
				Size:                       doc.Size,
				Token:                      encryptedFileName,
				OriginalName:               doc.OriginalName,
				Type:                       doc.Type,
				CreatedBYID:                currentUserID,
			}

			documents = append(documents, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				var existingDocument *models.AuditFindingReassessmentDocument
				existingDocument, err = afrs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						err = utility.NewResourceNotFoundError(err, *doc.ID)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err = afrs.repo.CreateDocuments(documents, tx); err != nil {
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err = afrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			return err
		}
	}

	storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "reassessments", strconv.FormatUint(uint64(id), 10), "documents")

	for _, doc := range documents {
		if err = utility.StoreFile(doc.Name, storagePath); err != nil {
			return err
		}
	}

	for _, document := range documentsToDelete {
		if err = utility.RemoveFile(document.Name, storagePath); err != nil {
			return err
		}
	}

	go func() {
		newAuditFindingReassessment, logErr := afrs.GetFullDetailsByID(findingID, id)
		if logErr != nil {
			return
		}

		oldData := map[string]interface{}{
			"outcome": oldAuditFindingReassessment.Outcome,
			"comment": oldAuditFindingReassessment.Comment,
		}
		newData := map[string]interface{}{
			"outcome": newAuditFindingReassessment.Outcome,
			"comment": newAuditFindingReassessment.Comment,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)
		if len(changedData) > 0 {
			logErr := afrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditFindingReassessment{}.TableName(),
				id,
				changedData,
				currentUserID,
				"audit", "audit_findings", "re-assessment", nil,
			)
			if logErr != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditFindingReassessment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     logErr.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (afrs *AuditFindingReassessmentService) Delete(findingID uint, id uint, currentUserID uint) (err error) {
	existingAuditFindingReassessment, err := afrs.GetFullDetailsByID(findingID, id)
	if err != nil {
		return err
	}

	if err := afrs.ensureLatestReassessment(findingID, id); err != nil {
		return err
	}

	tx := afrs.repo.BeginTransaction()

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}

		if err != nil {
			tx.Rollback()
			return
		}

		err = tx.Commit()
	}()

	if len(existingAuditFindingReassessment.Documents) > 0 {
		var deleteDocumentIds []*uint
		for _, doc := range existingAuditFindingReassessment.Documents {
			if doc != nil {
				deleteDocumentIds = append(deleteDocumentIds, doc.ID)
			}
		}
		if err = afrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			return err
		}
	}

	if err = afrs.repo.Delete(id, tx); err != nil {
		return err
	}

	for _, doc := range existingAuditFindingReassessment.Documents {
		if doc == nil {
			continue
		}
		storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "reassessments", strconv.FormatUint(uint64(id), 10), "documents")
		if removeErr := utility.RemoveFile(doc.Name, storagePath); removeErr != nil {
			utility.LogError("Failed to delete reassessment document file", map[string]interface{}{
				"error": removeErr.Error(),
				"file":  doc.Name,
			})
		}
	}

	go func() {
		logErr := afrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditFindingReassessment{}.TableName(),
			id,
			map[string]interface{}{
				"id":               id,
				"audit_finding_id": existingAuditFindingReassessment.AuditFindingID,
				"outcome":          existingAuditFindingReassessment.Outcome,
				"comment":          existingAuditFindingReassessment.Comment,
			},
			currentUserID,
			"audit", "audit_findings", "re-assessment", nil,
		)
		if logErr != nil {
			utility.LogError("Failed to log delete activity", map[string]interface{}{
				"error": logErr.Error(),
			})
		}
	}()

	return err
}

func (afrs *AuditFindingReassessmentService) GetDocumentPath(findingID uint, id uint, token string) (string, error) {
	existingAuditFindingReassessment, err := afrs.GetByID(findingID, id)
	if err != nil {
		return "", err
	}

	if len(existingAuditFindingReassessment.Documents) == 0 {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("no documents found for the audit finding reassessment"), id)
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("resource not found"), id)
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("resource not found"), id)
	}

	storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(findingID), 10), "reassessments", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("file not found"), id)
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

// Private helpers

func (afrs *AuditFindingReassessmentService) ensureLatestReassessment(findingID uint, id uint) error {
	latestID, err := afrs.repo.GetLatestID(findingID)
	if err != nil {
		return err
	}

	if latestID == 0 || latestID != id {
		return utility.NewValidationError(fmt.Errorf("only the latest reassessment can be edited or deleted"), "modify_requires_latest", "reassessment")
	}

	return nil
}

func (afrs *AuditFindingReassessmentService) ensureCreatable(findingID uint) error {
	latestID, err := afrs.repo.GetLatestID(findingID)
	if err != nil {
		return err
	}

	if latestID == 0 {
		return nil
	}

	latestAuditFindingReassessment, err := afrs.repo.GetByID(latestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, latestID)
		}
		return err
	}

	if latestAuditFindingReassessment.Outcome == nil || *latestAuditFindingReassessment.Outcome != "fail" {
		return utility.NewValidationError(fmt.Errorf("a new reassessment can only be created if none exists or the latest has failed"), "latest_must_be_failed", "reassessment")
	}

	return nil
}
