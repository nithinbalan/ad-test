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
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestFindingStatusUpdateServiceInterface interface {
	GetAll(findingID uint, sort string, sortBy string, keyword string) ([]*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error)
	GetLatest(findingID uint) (*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error)
	Create(findingID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO) error
	GetDocumentPath(findingID uint, id uint, token string) (string, error)
}

type PenetrationTestFindingStatusUpdateService struct {
	repo                      repositories.PenetrationTestFindingStatusUpdateRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	penetrationTestFindingSvc PenetrationTestFindingServiceInterface
	actionPlanStatusSvc       masterServices.ActionPlanStatusServiceInterface
}

func NewPenetrationTestFindingStatusUpdateService(
	repo repositories.PenetrationTestFindingStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestFindingServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) PenetrationTestFindingStatusUpdateServiceInterface {
	return &PenetrationTestFindingStatusUpdateService{repo: repo, activityLogSvc: activityLogSvc, penetrationTestFindingSvc: penetrationTestSvc, actionPlanStatusSvc: actionPlanStatusSvc}
}

func (ptfsus *PenetrationTestFindingStatusUpdateService) GetAll(findingID uint, sort string, sortBy string, keyword string) ([]*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error) {
	if _, err := ptfsus.penetrationTestFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	return ptfsus.repo.GetAll(findingID, sort, sortBy, keyword)
}

func (ptfsus *PenetrationTestFindingStatusUpdateService) GetLatest(findingID uint) (*penetrationTestDTOs.PenetrationTestFindingStatusUpdateListDTO, error) {
	if _, err := ptfsus.penetrationTestFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}
	value, err := ptfsus.repo.GetLatest(findingID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, findingID)
		}
		return nil, err
	}
	return value, nil
}

func (ptfsus *PenetrationTestFindingStatusUpdateService) validatePenetrationTestFindingStatusUpdateFormDTO(form *penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO) error {
	if _, err := ptfsus.actionPlanStatusSvc.GetByID(form.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}
	return nil
}

func (ptfsus *PenetrationTestFindingStatusUpdateService) Create(findingID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO) (uint, error) {
	if _, err := ptfsus.penetrationTestFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}
	if err := ptfsus.validatePenetrationTestFindingStatusUpdateFormDTO(form); err != nil {
		return 0, err
	}
	tx := ptfsus.repo.BeginTransaction()
	value := &models.PenetrationTestFindingStatusUpdate{PenetrationTestFindingID: findingID, ActionPlanStatusID: form.ActionPlanStatusID, Percentage: form.Percentage, Comment: form.Comment, CreatedBYID: currentUserID}
	id, err := ptfsus.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := ptfsus.penetrationTestFindingSvc.StatusUpdate(findingID, currentUserID, form.ActionPlanStatusID, form.Percentage, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	var documents []*models.PenetrationTestFindingStatusUpdateDocument
	for _, doc := range form.Documents {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		url := utility.GenerateDocumentURL("penetration-testing", "penetration-test-findings", uint64(findingID), "status-updates", uint64(id), "documents")
		documents = append(documents, &models.PenetrationTestFindingStatusUpdateDocument{
			PenetrationTestFindingStatusUpdateID: id, URL: url, Name: doc.FileName, Size: doc.Size,
			Token: encryptedFilename, OriginalName: doc.OriginalName, Type: doc.Type, CreatedBYID: currentUserID,
		})
	}
	if len(documents) > 0 {
		if err := ptfsus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		storagePath := filepath.Join("penetration-testing", "penetration-test-findings", strconv.FormatUint(uint64(findingID), 10), "status-updates", strconv.FormatUint(uint64(id), 10), "documents")
		for _, doc := range documents {
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingEntry, err := ptfsus.repo.GetByID(findingID, id)
		if err != nil {
			utility.LogError("Failed to fetch penetration test finding status update", map[string]interface{}{"error": err.Error()})
			return
		}
		if err := ptfsus.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTestFindingStatusUpdate{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "create penetration test finding status update", nil); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
		}
	}()
	return id, nil
}

func (ptfsus *PenetrationTestFindingStatusUpdateService) GetDocumentPath(findingID uint, id uint, token string) (string, error) {
	value, err := ptfsus.repo.GetByID(findingID, id)
	if err != nil {
		return "", err
	}
	if len(value.Documents) == 0 {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("no documents found"), id)
	}
	filenameBytes, err := utility.DecryptData(token)
	if err != nil || string(filenameBytes) == "" {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("resource not found"), id)
	}
	storagePath := filepath.Join("penetration-testing", "penetration-test-findings", strconv.FormatUint(uint64(findingID), 10), "status-updates", strconv.FormatUint(uint64(id), 10), "documents")
	fullPath := filepath.Join(config.GetEnv("STORAGE_PATH", ""), storagePath, string(filenameBytes))
	if _, err := os.Stat(fullPath); err != nil {
		return "", utility.NewResourceNotFoundError(err, id)
	}
	return fullPath, nil
}

func (ptfsus *PenetrationTestFindingStatusUpdateService) Update(findingID uint, id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFindingStatusUpdateFormDTO) error {
	if err := ptfsus.validatePenetrationTestFindingStatusUpdateFormDTO(form); err != nil {
		return err
	}
	tx := ptfsus.repo.BeginTransaction()
	value, err := ptfsus.repo.GetByID(findingID, id, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	oldData := map[string]interface{}{
		"comment":    value.Comment,
		"percentage": value.Percentage,
	}
	value.Percentage = form.Percentage
	value.Comment = form.Comment
	value.UpdatedBYID = &currentUserID
	if err := ptfsus.repo.Update(value, tx); err != nil {
		tx.Rollback()
		return err
	}
	var documents []*models.PenetrationTestFindingStatusUpdateDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.PenetrationTestFindingStatusUpdateDocument
	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("penetration-testing", "penetration-test-findings", uint64(findingID), "status-updates", uint64(id), "documents")
			documents = append(documents, &models.PenetrationTestFindingStatusUpdateDocument{
				PenetrationTestFindingStatusUpdateID: id,
				URL:                                  url,
				Name:                                 doc.FileName,
				Size:                                 doc.Size,
				Token:                                encryptedFilename,
				OriginalName:                         doc.OriginalName,
				Type:                                 doc.Type,
				CreatedBYID:                          currentUserID,
			})
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := ptfsus.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return utility.NewResourceNotFoundError(err, *doc.ID)
				}
				return err
			}
			documentsToDelete = append(documentsToDelete, existingDocument)
			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
		}
	}
	if len(documents) > 0 {
		if err := ptfsus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIDs) > 0 {
		if err := ptfsus.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	storagePath := filepath.Join("penetration-testing", "penetration-test-findings", strconv.FormatUint(uint64(findingID), 10), "status-updates", strconv.FormatUint(uint64(id), 10), "documents")
	for _, doc := range documents {
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, doc := range documentsToDelete {
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		updatedEntry, err := ptfsus.repo.GetByID(findingID, id)
		if err != nil {
			utility.LogError("Failed to fetch penetration test finding status update", map[string]interface{}{"error": err.Error()})
			return
		}
		newData := map[string]interface{}{
			"comment":    updatedEntry.Comment,
			"percentage": updatedEntry.Percentage,
		}
		changedData := utility.CompareFieldChanges(oldData, newData)
		if len(changedData) > 0 {
			if err := ptfsus.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTestFindingStatusUpdate{}.TableName(), id, changedData, currentUserID, "penetration_testing", "penetration_testing_pentest_log", "update penetration test finding status update", nil); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}
	}()
	return nil
}
