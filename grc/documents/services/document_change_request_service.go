package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepositories "jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestServiceInterface interface {
	GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestDetailsDTO, error)
	Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type DocumentChangeRequestService struct {
	repo                           documentRepositories.DocumentChangeRequestRepositoryInterface
	activityLogSvc                 services.ActivityLogServiceInterface
	documentChangeRequestTypeSvc   DocumentChangeRequestTypeServiceInterface
	documentChangeRequestStatusSvc DocumentChangeRequestStatusServiceInterface
}

func NewDocumentChangeRequestService(
	repo documentRepositories.DocumentChangeRequestRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentChangeRequestTypeSvc DocumentChangeRequestTypeServiceInterface,
	documentChangeRequestStatusSvc DocumentChangeRequestStatusServiceInterface,
) DocumentChangeRequestServiceInterface {
	return &DocumentChangeRequestService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		documentChangeRequestTypeSvc:   documentChangeRequestTypeSvc,
		documentChangeRequestStatusSvc: documentChangeRequestStatusSvc,
	}
}

func (dcrs *DocumentChangeRequestService) GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error) {
	return dcrs.repo.GetAll(params)
}

func (dcrs *DocumentChangeRequestService) GetByID(id uint) (*models.DocumentChangeRequest, error) {
	changeRequest, err := dcrs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return changeRequest, nil
}

func (dcrs *DocumentChangeRequestService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestDetailsDTO, error) {
	var detailsDTO documentDTOs.DocumentChangeRequestDetailsDTO
	changeRequest, err := dcrs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(changeRequest, &detailsDTO)
	if err != nil {
		return nil, err
	}
	detailsDTO.CreatedAt = changeRequest.CreatedAt
	detailsDTO.UpdatedAt = changeRequest.UpdatedAt
	return &detailsDTO, nil
}

func (dcrs *DocumentChangeRequestService) Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) (uint, error) {
	tx := dcrs.repo.BeginTransaction()

	if err := dcrs.validateFormDTO(formDTO); err != nil {
		return 0, err
	}

	changeRequestStatus, err := dcrs.documentChangeRequestStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	changeRequest := &models.DocumentChangeRequest{
		DocumentID:                      formDTO.DocumentID,
		DocumentChangeRequestInitiateID: formDTO.DocumentChangeRequestInitiateID,
		DocumentChangeRequestTypeID:     formDTO.DocumentChangeRequestTypeID,
		ReasonForChange:                 formDTO.ReasonForChange,
		DocumentChangeRequestStatusID:   changeRequestStatus.ID,
		CreatedBYID:                     currentUserID,
	}

	id, err := dcrs.repo.Create(changeRequest, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	dataToEncrypt := []byte(formDTO.Document.FileName)
	encryptedFileName, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	url := utility.GenerateDocumentURL("document", "change_requests", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
	document := &models.DocumentChangeRequestDocument{
		DocumentChangeRequestID: id,
		URL:                     url,
		Name:                    formDTO.Document.FileName,
		Size:                    formDTO.Document.Size,
		Token:                   encryptedFileName,
		OriginalName:            &formDTO.Document.OriginalName,
		Type:                    formDTO.Document.Type,
	}

	if err := dcrs.repo.CreateDocument(document, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	storagePath := filepath.Join("document", "change_requests", strconv.FormatUint(uint64(id), 10), "documents")
	if err := utility.StoreFile(formDTO.Document.FileName, storagePath); err != nil {
		return 0, err
	}

	tx.Commit()

	go func() {
		if err := utility.ConvertDocxToPDF(storagePath, formDTO.Document.FileName, storagePath); err != nil {
			log.Printf("Failed to convert document to PDF: %v", err)
		}
	}()

	return id, nil
}

func (dcrs *DocumentChangeRequestService) Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) error {
	tx := dcrs.repo.BeginTransaction()

	existingChangeRequest, err := dcrs.GetByID(id)
	if err != nil {
		return err
	}

	if err := dcrs.validateFormDTO(formDTO); err != nil {
		return err
	}

	// Update the existing change request
	existingChangeRequest.DocumentID = formDTO.DocumentID
	existingChangeRequest.DocumentChangeRequestInitiateID = formDTO.DocumentChangeRequestInitiateID
	existingChangeRequest.DocumentChangeRequestTypeID = formDTO.DocumentChangeRequestTypeID
	existingChangeRequest.ReasonForChange = formDTO.ReasonForChange
	existingChangeRequest.UpdatedBYID = &currentUserID

	err = dcrs.repo.Update(existingChangeRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle document changes
	existingDocument, err := dcrs.repo.GetDocumentByChangeRequestID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	var isNewFileProvided bool
	if formDTO.Document.IsNewFile != nil && *formDTO.Document.IsNewFile {
		isNewFileProvided = true
		existingDocumentFileName := existingDocument.Name
		baseFileName := strings.TrimSuffix(existingDocumentFileName, filepath.Ext(existingDocumentFileName))

		// Encrypt file name
		dataToEncrypt := []byte(formDTO.Document.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Generate file URL
		url := utility.GenerateDocumentURL("document", "change_requests", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

		// Update document fields
		existingDocument.URL = url
		existingDocument.Name = formDTO.Document.FileName
		existingDocument.Size = formDTO.Document.Size
		existingDocument.Token = encryptedFileName
		existingDocument.OriginalName = &formDTO.Document.OriginalName
		existingDocument.Type = formDTO.Document.Type

		// Save the document
		err = dcrs.repo.UpdateDocument(existingDocument, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Remove old file and store the new file
		storagePath := filepath.Join("document", "change_requests", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(existingDocument.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
		pdfFileName := baseFileName + ".pdf"
		if err := utility.RemoveFile(pdfFileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

		if err := utility.StoreFile(formDTO.Document.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	// Convert to PDF if a new file was provided
	if isNewFileProvided {
		go func() {
			storagePath := filepath.Join("document", "change_requests", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.ConvertDocxToPDF(storagePath, formDTO.Document.FileName, storagePath); err != nil {
				log.Printf("Failed to convert document to PDF: %v", err)
			}
		}()
	}

	return nil
}

func (dcrs *DocumentChangeRequestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := dcrs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if err := dcrs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()
	return failedIds, nil
}

func (dcrs *DocumentChangeRequestService) validateFormDTO(dto *documentDTOs.DocumentChangeRequestFormDTO) error {
	if _, err := dcrs.documentChangeRequestTypeSvc.GetByID(dto.DocumentChangeRequestTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document change request type does not exist"), "not_exist", "change request type")
	}
	err := dcrs.ValidateDocument(&dto.Document)
	if err != nil {
		return err
	}
	return nil
}

func (dcrs *DocumentChangeRequestService) ValidateDocument(document *commonDTOs.DocumentFormDTO) error {
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}
	// Validate file type (DOC)
	if document.Type != ".docx" {
		return utility.NewValidationError(fmt.Errorf("only word files are allowed"), "invalid_word_file_type", "type")
	}

	return nil
}
