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
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectSupportingDocumentServiceInterface interface {
	GetAll(projectID uint, params *dtos.ProjectSupportingDocumentListQueryParamsDTO) (*dtos.ProjectSupportingDocumentListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectSupportingDocument, error)
	GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectSupportingDocumentDetailDTO, error)
	Create(projectID uint, currentUserID uint, form *dtos.ProjectSupportingDocumentFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectSupportingDocumentFormDTO) error
	Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(projectID uint, supportingDocID uint, documentID uint, token string) (string, error)
}

type ProjectSupportingDocumentService struct {
	repo        repositories.ProjectSupportingDocumentRepositoryInterface
	activityLog logServices.ActivityLogServiceInterface
	userSvc     organizationServices.UserServiceInterface
	projectSvc  ProjectServiceInterface
}

func NewProjectSupportingDocumentService(
	repo repositories.ProjectSupportingDocumentRepositoryInterface,
	activityLog logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	projectSvc ProjectServiceInterface,
) ProjectSupportingDocumentServiceInterface {
	return &ProjectSupportingDocumentService{
		repo:        repo,
		activityLog: activityLog,
		userSvc:     userSvc,
		projectSvc:  projectSvc,
	}
}

func (psds *ProjectSupportingDocumentService) GetAll(projectID uint, params *dtos.ProjectSupportingDocumentListQueryParamsDTO) (*dtos.ProjectSupportingDocumentListResponseDTO, error) {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	return psds.repo.GetAll(projectID, params)
}

func (psds *ProjectSupportingDocumentService) GetByID(projectID uint, id uint) (*models.ProjectSupportingDocument, error) {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	supportingDoc, err := psds.repo.GetByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}
	return supportingDoc, nil
}

func (psds *ProjectSupportingDocumentService) GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectSupportingDocumentDetailDTO, error) {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	supportingDoc, err := psds.repo.GetFullDetailsByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	var detailDTO dtos.ProjectSupportingDocumentDetailDTO
	if err := mapstructure.Decode(supportingDoc, &detailDTO); err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = supportingDoc.CreatedAt
	detailDTO.UpdatedAt = supportingDoc.UpdatedAt

	return &detailDTO, nil
}

func (psds *ProjectSupportingDocumentService) Create(projectID uint, currentUserID uint, form *dtos.ProjectSupportingDocumentFormDTO) (uint, error) {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return 0, err
	}

	supportingDoc := &models.ProjectSupportingDocument{
		ProjectID:   projectID,
		TitleEn:     form.TitleEn,
		TitleAr:     form.TitleAr,
		CreatedByID: &currentUserID,
		CreatedAt:   time.Now(),
	}

	tx := psds.repo.BeginTransaction()

	id, err := psds.repo.Create(supportingDoc, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ProjectSupportingDocumentDocument
	for _, doc := range form.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

		document := &models.ProjectSupportingDocumentDocument{
			ProjectSupportingDocumentID: id,
			URL:                         url,
			Name:                        doc.FileName,
			Size:                        doc.Size,
			Token:                       encryptedFileName,
			OriginalName:                &doc.OriginalName,
			Type:                        doc.Type,
			CreatedByID:                 &currentUserID,
			CreatedAt:                   time.Now(),
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := psds.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	for _, doc := range documents {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		_ = psds.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectSupportingDocument{}.TableName(),
			id,
			supportingDoc,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	return id, nil
}

func (psds *ProjectSupportingDocumentService) Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectSupportingDocumentFormDTO) error {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	supportingDoc, err := psds.repo.GetFullDetailsByID(projectID, id)
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
	supportingDoc.UpdatedByID = &currentUserID
	updatedAt := time.Now()
	supportingDoc.UpdatedAt = &updatedAt

	tx := psds.repo.BeginTransaction()

	if err := psds.repo.Update(supportingDoc, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Document handling
	var newDocsToCreate []*models.ProjectSupportingDocumentDocument
	var docsToDelete []*uint
	var documentsToDelete []*models.ProjectSupportingDocumentDocument

	for _, doc := range form.Documents {
		if (doc.IsNewFile != nil && *doc.IsNewFile) || (doc.ID == nil && doc.FileName != "") {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			newDoc := &models.ProjectSupportingDocumentDocument{
				ProjectSupportingDocumentID: id,
				URL:                         url,
				Name:                        doc.FileName,
				Size:                        doc.Size,
				Token:                       encryptedFileName,
				OriginalName:                &doc.OriginalName,
				Type:                        doc.Type,
				CreatedByID:                 &currentUserID,
				CreatedAt:                   time.Now(),
			}
			newDocsToCreate = append(newDocsToCreate, newDoc)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := psds.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						tx.Rollback()
						return utility.NewResourceNotFoundError(err, *doc.ID)
					}
					tx.Rollback()
					return err
				}
				if existingDocument.ProjectSupportingDocumentID != id {
					tx.Rollback()
					return utility.NewResourceNotFoundError(gorm.ErrRecordNotFound, *doc.ID)
				}
				documentsToDelete = append(documentsToDelete, existingDocument)
				docsToDelete = append(docsToDelete, &existingDocument.ID)
			}
		}
	}

	if len(docsToDelete) > 0 {
		if err := psds.repo.DeleteDocuments(docsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(newDocsToCreate) > 0 {
		if err := psds.repo.CreateDocuments(newDocsToCreate, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	for _, doc := range newDocsToCreate {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to store document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			utility.LogError("Failed to remove document file", map[string]interface{}{"error": err.Error(), "doc": doc.Name})
		}
	}

	go func() {
		newData := map[string]interface{}{
			"title_en": form.TitleEn,
			"title_ar": form.TitleAr,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			_ = psds.activityLog.LogActivity(
				models.UpdateAction,
				models.ProjectSupportingDocument{}.TableName(),
				id,
				changes,
				currentUserID,
				"project", "project_register", "", nil,
			)
		}
	}()

	return nil
}

func (psds *ProjectSupportingDocumentService) Delete(projectID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return ids, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	tx := psds.repo.BeginTransaction()

	for _, id := range ids {
		supportingDoc, err := psds.repo.GetByID(projectID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := psds.repo.DeleteDocumentsBySupportingDocumentID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := psds.repo.Delete(id, tx); err != nil {
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
		folderPath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(id), 10))
		if err := utility.RemoveFolder(folderPath); err != nil {
			utility.LogError("Failed to remove supporting document storage folder", map[string]interface{}{
				"supportingDocumentID": id,
				"error":                err.Error(),
			})
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := psds.activityLog.LogBulkActivities(
				models.DeleteAction,
				models.ProjectSupportingDocument{}.TableName(),
				successfulIds,
				deletedItems,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete project supporting documents activity", map[string]interface{}{
					"recordIDs": successfulIds,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (psds *ProjectSupportingDocumentService) GetDocumentPath(projectID uint, supportingDocID uint, documentID uint, token string) (string, error) {
	if _, err := psds.projectSvc.GetByID(projectID); err != nil {
		return "", err
	}

	if _, err := psds.repo.GetByID(projectID, supportingDocID); err != nil {
		return "", err
	}

	document, err := psds.repo.GetDocumentByID(documentID)
	if err != nil {
		return "", err
	}

	if document.Token != token {
		return "", &utility.ResourceNotFoundError{}
	}

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	storagePath := filepath.Join("project", strconv.FormatUint(uint64(projectID), 10), "supporting-documents", strconv.FormatUint(uint64(supportingDocID), 10), "documents")
	fullPath := filepath.Join(baseStoragePath, storagePath, document.Name)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", &utility.ResourceNotFoundError{}
	}

	return fullPath, nil
}
