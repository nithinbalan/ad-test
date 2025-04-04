package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetOtherDocumentServiceInterface interface {
	GetAll(assetRegisterID uint) ([]*dtos.AssetOtherDocumentListDTO, error)
	GetByID(assetRegisterID, id uint) (*models.AssetOtherDocument, error)
	GetFullDetailsByID(assetRegisterID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetOtherDocumentDetailDTO, error)
	Create(assetRegisterID uint, form *dtos.AssetOtherDocumentFormDTO, currentUserID uint) ([]uint, error)
	Update(assetRegisterID uint, id uint, form *dtos.AssetOtherDocumentFormDTO, currentUserID uint) error
	Delete(assetRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AssetOtherDocumentService struct {
	repo             repositories.AssetOtherDocumentRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	assetRegisterSvc AssetRegisterServiceInterface
	categorySvc      AssetDocumentCategoryServiceInterface
}

func NewAssetOtherDocumentService(
	repo repositories.AssetOtherDocumentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	categorySvc AssetDocumentCategoryServiceInterface,
) AssetOtherDocumentServiceInterface {
	return &AssetOtherDocumentService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		assetRegisterSvc: assetRegisterSvc,
		categorySvc:      categorySvc,
	}
}

func (aods *AssetOtherDocumentService) GetAll(assetRegisterID uint) ([]*dtos.AssetOtherDocumentListDTO, error) {
	if _, err := aods.assetRegisterSvc.GetByID(assetRegisterID); err != nil {
		return nil, err
	}

	docs, err := aods.repo.GetAll(assetRegisterID)
	if err != nil {
		return nil, err
	}

	groupMap := make(map[uint][]*dtos.AssetOtherDocumentListDTO)
	for _, doc := range docs {
		groupMap[doc.AssetDocumentCategoryID] = append(groupMap[doc.AssetDocumentCategoryID], doc)
	}

	var result []*dtos.AssetOtherDocumentListDTO
	for _, group := range groupMap {
		parent := group[0]

		// Build a clean copy of group without recursive Documents field
		var cleanedDocs []*dtos.AssetOtherDocumentListDTO
		for _, doc := range group {
			copied := *doc
			copied.Documents = nil // prevent cycle
			cleanedDocs = append(cleanedDocs, &copied)
		}

		parentCopy := *parent
		parentCopy.Documents = cleanedDocs

		result = append(result, &parentCopy)
	}

	return result, nil
}

func (aods *AssetOtherDocumentService) GetByID(assetRegisterID uint, id uint) (*models.AssetOtherDocument, error) {
	if _, err := aods.assetRegisterSvc.GetByID(assetRegisterID); err != nil {
		return nil, err
	}
	doc, err := aods.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return doc, nil
}
func (aods *AssetOtherDocumentService) GetFullDetailsByID(assetRegisterID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetOtherDocumentDetailDTO, error) {
	if _, err := aods.assetRegisterSvc.GetByID(assetRegisterID); err != nil {
		return nil, err
	}
	doc, err := aods.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	var dto dtos.AssetOtherDocumentDetailDTO
	err = mapstructure.Decode(doc, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (aods *AssetOtherDocumentService) Create(assetRegisterID uint, form *dtos.AssetOtherDocumentFormDTO, currentUserID uint) ([]uint, error) {
	tx := aods.repo.BeginTransaction()

	// Validate category existence
	if _, err := aods.categorySvc.GetByID(form.AssetDocumentCategoryID); err != nil {
		tx.Rollback()
		return nil, utility.NewValidationError(fmt.Errorf("invalid category ID"), "not_exist", "asset_document_category_id")
	}

	var createdIDs []uint
	var documents []*models.AssetOtherDocument

	for _, doc := range form.Documents {
		// Encrypt the filename
		encryptedToken, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Construct storage URL
		url := utility.GenerateDocumentURL("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")

		// Prepare model
		document := &models.AssetOtherDocument{
			AssetRegisterID:         assetRegisterID,
			AssetDocumentCategoryID: form.AssetDocumentCategoryID,
			URL:                     &url,
			Name:                    &doc.FileName,
			Size:                    &doc.Size,
			Token:                   &encryptedToken,
			OriginalName:            &doc.OriginalName,
			Type:                    &doc.Type,
			CreatedBYID:             currentUserID,
		}

		documents = append(documents, document)
	}
	if len(documents) > 0 {
		if err := aods.repo.CreateAssetOtherDocuments(documents, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		// Store files physically
		for _, doc := range documents {
			storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")
			if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}

	tx.Commit()

	return createdIDs, nil
}

func (aods *AssetOtherDocumentService) Update(assetRegisterID uint, id uint, form *dtos.AssetOtherDocumentFormDTO, currentUserID uint) error {
	if _, err := aods.assetRegisterSvc.GetByID(assetRegisterID); err != nil {
		return err
	}

	// Validate category
	if _, err := aods.categorySvc.GetByID(form.AssetDocumentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid category id"), "not_exist", "asset_document_category_id")
	}

	tx := aods.repo.BeginTransaction()

	var assetDocuments []*models.AssetOtherDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AssetOtherDocument

	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			encryptedFileName, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")

			document := &models.AssetOtherDocument{
				AssetRegisterID:         assetRegisterID,
				AssetDocumentCategoryID: form.AssetDocumentCategoryID,
				URL:                     &url,
				Name:                    &doc.FileName,
				Size:                    &doc.Size,
				Token:                   &encryptedFileName,
				OriginalName:            &doc.OriginalName,
				Type:                    &doc.Type,
				CreatedBYID:             currentUserID,
			}
			assetDocuments = append(assetDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {

			if doc.ID != nil {
				existingDoc, err := aods.repo.GetByID(*doc.ID)
				if err != nil {
					tx.Rollback()
					return err
				}
				documentsToDelete = append(documentsToDelete, existingDoc)
				deleteDocumentIds = append(deleteDocumentIds, &existingDoc.ID)
			}
		} else {
			// Check for category ID change
			if doc.ID != nil {
				existingDoc, err := aods.repo.GetByID(*doc.ID)
				if err != nil {
					tx.Rollback()
					return err
				}

				if existingDoc.AssetDocumentCategoryID != form.AssetDocumentCategoryID {
					existingDoc.AssetDocumentCategoryID = form.AssetDocumentCategoryID
					existingDoc.UpdatedBYID = &currentUserID

					if err := aods.repo.Update(existingDoc); err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
	}

	// Save new documents
	if len(assetDocuments) > 0 {
		if err := aods.repo.CreateAssetOtherDocuments(assetDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete old documents
	if len(deleteDocumentIds) > 0 {
		if err := aods.repo.DeleteAssetOtherDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store new files
	for _, doc := range assetDocuments {
		storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove deleted files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")
		if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (aods *AssetOtherDocumentService) Delete(assetRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds, successIds []uint
	var deletedItems []map[string]interface{}
	documentsToDelete := make([]*models.AssetOtherDocument, 0)

	tx := aods.repo.BeginTransaction()

	for _, id := range ids {
		doc, err := aods.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Queue for deletion from storage
		documentsToDelete = append(documentsToDelete, doc)

		// Delete DB record
		if err := aods.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		successIds = append(successIds, id)
		deletedItems = append(deletedItems, map[string]interface{}{
			"id":   doc.ID,
			"name": doc.Name,
			"type": doc.Type,
		})
	}

	// If no successful deletes, rollback and return
	if len(successIds) == 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Delete physical files
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")
		if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
			fmt.Printf("⚠️ Failed to remove file %s: %v\n", doc.Name, err)
		}
	}

	// Clean folder if needed
	if len(documentsToDelete) > 0 {
		folderPath := filepath.Join("assets", "asset-register", strconv.FormatUint(uint64(assetRegisterID), 10), "other-documents")
		if err := utility.RemoveFolder(folderPath); err != nil {
			fmt.Printf("⚠️ Failed to remove folder %s: %v\n", folderPath, err)
		}
	}

	// Log activity
	go func() {
		_ = aods.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.AssetOtherDocument{}.TableName(),
			successIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}()

	return failedIds, nil
}
