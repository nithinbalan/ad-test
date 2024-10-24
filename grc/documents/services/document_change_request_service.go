package services

import (
	"errors"
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestServiceInterface interface {
	GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestDetailsDTO, error)
	Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) error
	//Delete(id uint, currentUserID uint) error
}

type DocumentChangeRequestService struct {
	repo repositories.DocumentChangeRequestRepositoryInterface
}

func NewDocumentChangeRequestService(repo repositories.DocumentChangeRequestRepositoryInterface) DocumentChangeRequestServiceInterface {
	return &DocumentChangeRequestService{
		repo: repo,
	}
}

func (dcrs *DocumentChangeRequestService) GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error) {
	return dcrs.repo.GetAll(params)
}

func (dcrs *DocumentChangeRequestService) GetByID(id uint) (*models.DocumentChangeRequest, error) {
	documentChangeRequest, err := dcrs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return documentChangeRequest, nil
}

func (dcrs *DocumentChangeRequestService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestDetailsDTO, error) {
	var documentChangeRequestDetailsDTO documentDTOs.DocumentChangeRequestDetailsDTO
	documentChangeRequest, err := dcrs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Use mapstructure to map the model into DTO
	err = mapstructure.Decode(documentChangeRequest, &documentChangeRequestDetailsDTO)
	if err != nil {
		return nil, err
	}

	// Manually map additional fields if necessary
	documentChangeRequestDetailsDTO.CreatedAt = documentChangeRequest.CreatedAt
	documentChangeRequestDetailsDTO.UpdatedAt = documentChangeRequest.UpdatedAt

	return &documentChangeRequestDetailsDTO, nil
}

func (dcrs *DocumentChangeRequestService) Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) (uint, error) {
	// Begin transaction
	tx := dcrs.repo.BeginTransaction()
	defer tx.Rollback()

	documentChangeRequest := &models.DocumentChangeRequest{
		DocumentsID:                   formDTO.DocumentsID,
		DocumentChangeRequestTypeID:   formDTO.DocumentChangeRequestTypeID,
		ReasonForChange:               formDTO.ReasonForChange,
		DocumentChangeRequestStatusID: 1, // Assuming 'initiated' is the default status
		CreatedBYID:                   currentUserID,
	}

	id, err := dcrs.repo.Create(documentChangeRequest, tx)
	if err != nil {
		return 0, err
	}

	// Process documents
	var documents []*models.DocumentChangeRequestSupportingDocument

	for _, doc := range formDTO.Documents {
		// Encrypt filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}

		// Generate URL
		url := utility.GenerateDocumentURL("document", "change-request", strconv.FormatUint(uint64(id), 10), "supporting-documents", encryptedFileName)

		// Create document instance
		document := &models.DocumentChangeRequestSupportingDocument{
			DocumentChangeRequestID: id,
			URL:                     url,
			Name:                    doc.FileName,
			Size:                    doc.Size,
			Token:                   encryptedFileName,
			OriginalName:            &doc.OriginalName,
			Type:                    doc.Type,
			CreatedBYID:             currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := dcrs.repo.CreateSupportingDocuments(documents, tx); err != nil {
			return 0, err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("document", "change-request", strconv.FormatUint(uint64(id), 10), "supporting-documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	// Commit transaction
	tx.Commit()

	return id, nil
}

func (dcrs *DocumentChangeRequestService) Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) error {
	tx := dcrs.repo.BeginTransaction()

	existingRequest, err := dcrs.repo.GetByID(id)
	if err != nil {
		return err
	}

	existingRequest.DocumentChangeRequestTypeID = formDTO.DocumentChangeRequestTypeID
	existingRequest.ReasonForChange = formDTO.ReasonForChange
	existingRequest.UpdatedBYID = &currentUserID

	err = dcrs.repo.Update(existingRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Process documents
	var documents []*models.DocumentChangeRequestSupportingDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.DocumentChangeRequestSupportingDocument

	for _, doc := range formDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate URL
			url := utility.GenerateDocumentURL("document", "change-request", strconv.FormatUint(uint64(existingRequest.DocumentsID), 10), "supporting-documents", encryptedFileName)

			document := &models.DocumentChangeRequestSupportingDocument{
				DocumentChangeRequestID: id,
				URL:                     url,
				Name:                    doc.FileName,
				Size:                    doc.Size,
				Token:                   encryptedFileName,
				OriginalName:            &doc.OriginalName,
				Type:                    doc.Type,
				CreatedBYID:             currentUserID,
			}

			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := dcrs.repo.GetDocumentByID(*doc.ID)
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
		if err := dcrs.repo.CreateSupportingDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := dcrs.repo.DeleteSupportingDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("document", "change-request", strconv.FormatUint(uint64(existingRequest.DocumentsID), 10), "supporting-documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("document", "change-request", strconv.FormatUint(uint64(existingRequest.DocumentsID), 10), "supporting-documents")
		fmt.Println(storagePath)
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

// func (dcrs *DocumentChangeRequestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
// 	tx := dcrs.repo.BeginTransaction()

// 	var failedIds []uint
// 	var successfulIds []uint
// 	documentsToDelete := make(map[uint][]*documentDTOs.DocumentChangeRequestSupportingDocumentDTO)

// 	for _, id := range ids {
// 		var deleteDocumentIds []*uint
// 		changeRequest, err := dcrs.GetFullDetailsByID(id, tx)
// 		if err != nil {
// 			failedIds = append(failedIds, id)
// 			continue
// 		}

// 		if changeRequest.Documents != nil && len(changeRequest.Documents) > 0 {
// 			for _, document := range changeRequest.Documents {
// 				documentsToDelete[id] = append(documentsToDelete[id], document)
// 				deleteDocumentIds = append(deleteDocumentIds, document.ID)
// 			}

// 			if len(deleteDocumentIds) > 0 {
// 				if err := dcrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
// 					failedIds = append(failedIds, id)
// 					continue
// 				}
// 			}
// 		}

// 		if err := dcrs.repo.Delete(id, tx); err != nil {
// 			failedIds = append(failedIds, id)
// 			continue
// 		} else {
// 			successfulIds = append(successfulIds, id)
// 		}
// 	}

// 	if len(failedIds) > 0 {
// 		tx.Rollback()
// 		return failedIds, nil
// 	}

// 	// Deleting files should not affect the main operation result
// 	if len(successfulIds) > 0 {
// 		for changeRequestId, documents := range documentsToDelete {
// 			for _, doc := range documents {
// 				storagePath := filepath.Join("document", "change-request", strconv.FormatUint(uint64(changeRequestId), 10), "supporting-documents")
// 				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
// 					tx.Rollback()
// 					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
// 				}
// 			}
// 		}
// 	}

// 	tx.Commit()

// 	go func() {
// 		err := dcrs.activityLogSvc.LogBulkActivities(
// 			models.DeleteAction,
// 			models.DocumentChangeRequest{}.TableName(),
// 			successfulIds,
// 			nil,
// 			currentUserID,
// 			nil,
// 		)
// 		if err != nil {
// 			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
// 		}
// 	}()

// 	return failedIds, nil
// }
