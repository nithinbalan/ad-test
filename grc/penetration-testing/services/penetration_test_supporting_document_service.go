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
	"jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestSupportingDocumentServiceInterface interface {
	GetAll(penetrationTestID uint, params *dtos.PenetrationTestSupportingDocumentListQueryParamsDTO) (*dtos.PenetrationTestSupportingDocumentListResponseDTO, error)
	GetByID(penetrationTestID uint, id uint) (*models.PenetrationTestSupportingDocument, error)
	GetFullDetailsByID(penetrationTestID uint, id uint) (*dtos.PenetrationTestSupportingDocumentDetailDTO, error)
	Create(penetrationTestID uint, currentUserID uint, form *dtos.PenetrationTestSupportingDocumentFormDTO) (uint, error)
	Update(penetrationTestID uint, id uint, currentUserID uint, form *dtos.PenetrationTestSupportingDocumentFormDTO) error
	Delete(penetrationTestID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(penetrationTestID uint, supportingDocID uint, documentID uint, token string) (string, error)
}

type PenetrationTestSupportingDocumentService struct {
	repo               repositories.PenetrationTestSupportingDocumentRepositoryInterface
	activityLog        logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	penetrationTestSvc PenetrationTestServiceInterface
}

func NewPenetrationTestSupportingDocumentService(
	repo repositories.PenetrationTestSupportingDocumentRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
) PenetrationTestSupportingDocumentServiceInterface {
	return &PenetrationTestSupportingDocumentService{
		repo:               repo,
		activityLog:        activityLog,
		userSvc:            userSvc,
		penetrationTestSvc: penetrationTestSvc,
	}
}

func (ptsds *PenetrationTestSupportingDocumentService) GetAll(penetrationTestID uint, params *dtos.PenetrationTestSupportingDocumentListQueryParamsDTO) (*dtos.PenetrationTestSupportingDocumentListResponseDTO, error) {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	return ptsds.repo.GetAll(penetrationTestID, params)
}

func (ptsds *PenetrationTestSupportingDocumentService) GetByID(penetrationTestID uint, id uint) (*models.PenetrationTestSupportingDocument, error) {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	supportingDoc, err := ptsds.repo.GetByID(penetrationTestID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}
	return supportingDoc, nil
}

func (ptsds *PenetrationTestSupportingDocumentService) GetFullDetailsByID(penetrationTestID uint, id uint) (*dtos.PenetrationTestSupportingDocumentDetailDTO, error) {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	supportingDoc, err := ptsds.repo.GetFullDetailsByID(penetrationTestID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	var detailDTO dtos.PenetrationTestSupportingDocumentDetailDTO
	if err := mapstructure.Decode(supportingDoc, &detailDTO); err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = supportingDoc.CreatedAt
	detailDTO.UpdatedAt = supportingDoc.UpdatedAt

	if detailDTO.Documents == nil {
		detailDTO.Documents = []*dtos.PenetrationTestSupportingDocumentDocumentDTO{}
	}

	return &detailDTO, nil
}

func (ptsds *PenetrationTestSupportingDocumentService) Create(penetrationTestID uint, currentUserID uint, form *dtos.PenetrationTestSupportingDocumentFormDTO) (uint, error) {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return 0, err
	}

	supportingDoc := &models.PenetrationTestSupportingDocument{
		PenetrationTestID: penetrationTestID,
		TitleEn:           form.TitleEn,
		TitleAr:           form.TitleAr,
		CreatedBYID:       &currentUserID,
		CreatedAt:         time.Now(),
	}

	tx := ptsds.repo.BeginTransaction()

	id, err := ptsds.repo.Create(supportingDoc, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.PenetrationTestSupportingDocumentDocument
	for _, doc := range form.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("penetration-testing", strconv.FormatUint(uint64(penetrationTestID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

		document := &models.PenetrationTestSupportingDocumentDocument{
			PenetrationTestSupportingDocumentID: id,
			URL:                                 url,
			Name:                                doc.FileName,
			Size:                                doc.Size,
			Token:                               encryptedFileName,
			OriginalName:                        &doc.OriginalName,
			Type:                                doc.Type,
			CreatedByID:                         &currentUserID,
			CreatedAt:                           time.Now(),
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := ptsds.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	for _, doc := range documents {
		storagePath := filepath.Join("penetration-testing", strconv.FormatUint(uint64(penetrationTestID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		_ = ptsds.activityLog.LogActivity(
			models.CreateAction,
			models.PenetrationTestSupportingDocument{}.TableName(),
			id,
			supportingDoc,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "", nil,
		)
	}()

	return id, nil
}

func (ptsds *PenetrationTestSupportingDocumentService) Update(penetrationTestID uint, id uint, currentUserID uint, form *dtos.PenetrationTestSupportingDocumentFormDTO) error {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return err
	}

	supportingDoc, err := ptsds.repo.GetFullDetailsByID(penetrationTestID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	oldData := map[string]interface{}{
		"title_en": supportingDoc.TitleEn,
		"title_ar": supportingDoc.TitleAr,
	}

	supportingDoc.TitleEn = form.TitleEn
	supportingDoc.TitleAr = form.TitleAr
	supportingDoc.UpdatedBYID = &currentUserID
	updatedAt := time.Now()
	supportingDoc.UpdatedAt = &updatedAt

	tx := ptsds.repo.BeginTransaction()

	if err := ptsds.repo.Update(supportingDoc, tx); err != nil {
		tx.Rollback()
		return err
	}

	var existingDocsToKeep []*uint
	var newDocsToCreate []*models.PenetrationTestSupportingDocumentDocument
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
			url := utility.GenerateDocumentURL("penetration-testing", strconv.FormatUint(uint64(penetrationTestID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			newDoc := &models.PenetrationTestSupportingDocumentDocument{
				PenetrationTestSupportingDocumentID: id,
				URL:                                 url,
				Name:                                doc.FileName,
				Size:                                doc.Size,
				Token:                               encryptedFileName,
				OriginalName:                        &doc.OriginalName,
				Type:                                doc.Type,
				CreatedByID:                         &currentUserID,
				CreatedAt:                           time.Now(),
			}
			newDocsToCreate = append(newDocsToCreate, newDoc)
		}
	}

	var docsToDelete []*uint
	for _, doc := range supportingDoc.Documents {
		keep := false
		for _, keepID := range existingDocsToKeep {
			if keepID != nil && *keepID == doc.ID {
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
		if err := ptsds.repo.DeleteDocuments(docsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(newDocsToCreate) > 0 {
		if err := ptsds.repo.CreateDocuments(newDocsToCreate, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, doc := range newDocsToCreate {
		storagePath := filepath.Join("penetration-testing", strconv.FormatUint(uint64(penetrationTestID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		newData := map[string]interface{}{
			"title_en": form.TitleEn,
			"title_ar": form.TitleAr,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			_ = ptsds.activityLog.LogActivity(
				models.UpdateAction,
				models.PenetrationTestSupportingDocument{}.TableName(),
				id,
				changes,
				currentUserID,
				"penetration_testing", "penetration_testing_pentest_log", "", nil,
			)
		}
	}()

	return nil
}

func (ptsds *PenetrationTestSupportingDocumentService) Delete(penetrationTestID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return ids, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	tx := ptsds.repo.BeginTransaction()

	for _, id := range ids {
		supportingDoc, err := ptsds.repo.GetByID(penetrationTestID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := ptsds.repo.DeleteDocumentsBySupportingDocumentID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := ptsds.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedItems = append(deletedItems, map[string]interface{}{
			"id":       id,
			"title_en": supportingDoc.TitleEn,
			"title_ar": supportingDoc.TitleAr,
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
		folderPath := filepath.Join("penetration-testing", strconv.FormatUint(uint64(penetrationTestID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10))
		if err := utility.RemoveFolder(folderPath); err != nil {
			utility.LogError("Failed to remove supporting document storage folder", map[string]interface{}{
				"supportingDocumentID": id,
				"error":                err.Error(),
			})
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ptsds.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.PenetrationTestSupportingDocument{}.TableName(),
				successfulIds,
				deletedItems,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete penetration test supporting documents activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (ptsds *PenetrationTestSupportingDocumentService) GetDocumentPath(penetrationTestID uint, supportingDocID uint, documentID uint, token string) (string, error) {
	if _, err := ptsds.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return "", err
	}

	if _, err := ptsds.repo.GetByID(penetrationTestID, supportingDocID); err != nil {
		return "", err
	}

	document, err := ptsds.repo.GetDocumentByID(documentID)
	if err != nil {
		return "", err
	}

	if document.Token != token {
		return "", &utility.ResourceNotFoundError{}
	}

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	storagePath := filepath.Join("penetration-testing", strconv.FormatUint(uint64(penetrationTestID), 10), "supporting-documents", strconv.FormatUint(uint64(supportingDocID), 10), "documents")
	fullPath := filepath.Join(baseStoragePath, storagePath, document.Name)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	}

	return fullPath, nil
}
