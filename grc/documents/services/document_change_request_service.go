package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepositories "jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestServiceInterface interface {
	SetDocumentChangeRequestWorkflowService(svc DocumentChangeRequestWorkflowServiceInterface)

	GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestDetailsDTO, error)
	Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	UpdateDocumentChangeRequest(documentChangeRequest *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) error

	CheckoutDocument(id uint, currentUserID uint, token string) (string, error)
	CheckinDocument(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error
	RevertDocumentCheckout(id uint, currentUserID uint) error
	SubmitDocumentChangeRequest(id uint, currentUserID uint) error
	PublishDocumentChangeRequest(id uint, currentUserID uint) error

	GetDocumentPath(id uint, token string, fileType string, tx ...utility.DatabaseTransaction) (string, error)
	GetDocumentCRLastEditSummary(id uint) (*documentDTOs.DocumentLastEditSummaryDTO, error)
}

type DocumentChangeRequestService struct {
	repo                             documentRepositories.DocumentChangeRequestRepositoryInterface
	activityLogSvc                   services.ActivityLogServiceInterface
	documentChangeRequestTypeSvc     DocumentChangeRequestTypeServiceInterface
	documentChangeRequestStatusSvc   DocumentChangeRequestStatusServiceInterface
	documentChangeRequestWorkflowSvc DocumentChangeRequestWorkflowServiceInterface
	documentCRWorkflowStatusSvc      DocumentChangeRequestWorkflowStatusServiceInterface
	documentSvc                      DocumentServiceInterface
	taskSvc                          taskServices.TaskServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
}

func NewDocumentChangeRequestService(
	repo documentRepositories.DocumentChangeRequestRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentChangeRequestTypeSvc DocumentChangeRequestTypeServiceInterface,
	documentChangeRequestStatusSvc DocumentChangeRequestStatusServiceInterface,
	documentCRWorkflowStatusSvc DocumentChangeRequestWorkflowStatusServiceInterface,
	documentSvc DocumentServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) DocumentChangeRequestServiceInterface {
	return &DocumentChangeRequestService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		documentChangeRequestTypeSvc:   documentChangeRequestTypeSvc,
		documentChangeRequestStatusSvc: documentChangeRequestStatusSvc,
		documentCRWorkflowStatusSvc:    documentCRWorkflowStatusSvc,
		documentSvc:                    documentSvc,
		taskSvc:                        taskSvc,
		notificationSvc:                notificationSvc,
	}
}

func (dcrs *DocumentChangeRequestService) SetDocumentChangeRequestWorkflowService(svc DocumentChangeRequestWorkflowServiceInterface) {
	dcrs.documentChangeRequestWorkflowSvc = svc
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

	referenceNumber, err := dcrs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	if err := dcrs.validateFormDTO(formDTO); err != nil {
		return 0, err
	}

	changeRequestStatus, err := dcrs.documentChangeRequestStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	documentChangeRequestType, err := dcrs.documentChangeRequestTypeSvc.GetByID(formDTO.DocumentChangeRequestTypeID)
	if err != nil {
		return 0, err
	}

	documentFullDetails, err := dcrs.documentSvc.GetFullDetailsByID(formDTO.DocumentID, nil, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var newVersion float64
	currentVersion := documentFullDetails.LatestDocumentVersion.Version
	if documentChangeRequestType.Type == "minor" {
		newVersion = currentVersion + 0.1
	} else if documentChangeRequestType.Type == "major" {
		newVersion = currentVersion + 1.0
	} else {
		tx.Rollback()
		return 0, fmt.Errorf("invalid document change request type")
	}

	changeRequest := &models.DocumentChangeRequest{
		ReferenceNumber:                 referenceNumber,
		DocumentID:                      formDTO.DocumentID,
		DocumentChangeRequestInitiateID: formDTO.DocumentChangeRequestInitiateID,
		DocumentChangeRequestTypeID:     formDTO.DocumentChangeRequestTypeID,
		ReasonForChange:                 formDTO.ReasonForChange,
		DocumentChangeRequestStatusID:   changeRequestStatus.ID,
		Version:                         newVersion,
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
	url := utility.GenerateDocumentURL("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
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

	storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
	if err := utility.StoreFile(formDTO.Document.FileName, storagePath); err != nil {
		return 0, err
	}

	err = dcrs.documentChangeRequestWorkflowSvc.CreateDocumentCRWorkflow(id, formDTO.DocumentID, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
		DocumentChangeRequestID: id,
		ActionType:              "create",
		FileName:                &formDTO.Document.FileName,
		OriginalName:            &formDTO.Document.OriginalName,
		Type:                    &formDTO.Document.Type,
		Size:                    &formDTO.Document.Size,
	}

	err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, &storagePath, currentUserID, tx)
	if err != nil {
		tx.Rollback()
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

	existingCRFullDetails, err := dcrs.repo.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	existingChangeRequestDocumentID := existingChangeRequest.DocumentID

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
		url := utility.GenerateDocumentURL("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

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
		storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
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

		documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
			DocumentChangeRequestID: id,
			ActionType:              "update",
			FileName:                &formDTO.Document.FileName,
			OriginalName:            &formDTO.Document.OriginalName,
			Type:                    &formDTO.Document.Type,
			Size:                    &formDTO.Document.Size,
		}

		err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, &storagePath, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if existingChangeRequestDocumentID != formDTO.DocumentID {

		var deleteDocumentWorkflowIds []*uint
		var deleteDocumentWorkflowParticipantIds []*uint
		var deleteDocumentWorkflowReviewIds []*uint

		if len(existingCRFullDetails.DocumentChangeRequestWorkflows) > 0 {
			for _, docWorkflow := range existingCRFullDetails.DocumentChangeRequestWorkflows {
				deleteDocumentWorkflowIds = append(deleteDocumentWorkflowIds, &docWorkflow.ID)

				if len(docWorkflow.DocumentChangeRequestWorkflowParticipants) > 0 {
					for _, participant := range docWorkflow.DocumentChangeRequestWorkflowParticipants {
						deleteDocumentWorkflowParticipantIds = append(deleteDocumentWorkflowParticipantIds, &participant.ID)
					}
				}

				if len(docWorkflow.DocumentChangeRequestWorkflowReviews) > 0 {
					for _, reviews := range docWorkflow.DocumentChangeRequestWorkflowReviews {
						deleteDocumentWorkflowReviewIds = append(deleteDocumentWorkflowReviewIds, &reviews.ID)
					}
				}
			}
		}

		if err := dcrs.documentChangeRequestWorkflowSvc.DeleteDocumentCRWorkflowParticipants(deleteDocumentWorkflowParticipantIds, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := dcrs.documentChangeRequestWorkflowSvc.DeleteDocumentCRWorkflowReviews(deleteDocumentWorkflowReviewIds, tx); err != nil {
			tx.Rollback()
			return err
		}

		if err := dcrs.documentChangeRequestWorkflowSvc.DeleteDocumentCRWorkflows(deleteDocumentWorkflowIds, tx); err != nil {
			tx.Rollback()
			return err
		}

		err = dcrs.documentChangeRequestWorkflowSvc.CreateDocumentCRWorkflow(id, formDTO.DocumentID, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()

	// Convert to PDF if a new file was provided
	if isNewFileProvided {
		go func() {
			storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
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

	documentsToDelete := make(map[uint][]*documentDTOs.DocumentChangeRequestDocumentDTO)

	for _, id := range ids {

		changeRequest, err := dcrs.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if changeRequest.ChangeRequestDocument != nil {
			documentsToDelete[id] = append(documentsToDelete[id], changeRequest.ChangeRequestDocument)

			// Remove the related document from the database
			if err := dcrs.repo.DeleteDocument(*changeRequest.ChangeRequestDocument.ID, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

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

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for changeRequestID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(changeRequestID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
			destinationPath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(changeRequestID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}
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

func (dcrs *DocumentChangeRequestService) CheckoutDocument(id uint, currentUserID uint, token string) (string, error) {
	tx := dcrs.repo.BeginTransaction()
	documentChangeRequest, err := dcrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	if documentChangeRequest.DocumentChangeRequestStatus.Type != "draft" && documentChangeRequest.DocumentChangeRequestStatus.Type != "in-review" {
		tx.Rollback()
		return "", utility.NewForbiddenError("Checkout only allowed for documents in 'draft' or 'in-review' status.")
	}

	if documentChangeRequest.IsLocked {
		tx.Rollback()
		return "", utility.NewForbiddenError("Document already checkouted.")
	}

	if documentChangeRequest.CreatedByUser.ID != int(currentUserID) && documentChangeRequest.Document.ResponsibleUser.ID != int(currentUserID) {
		tx.Rollback()
		return "", utility.NewForbiddenError("Access denied: user is neither creator nor responsible user.")
	}

	currentDateStr := time.Now().Format("2006-01-02")
	lastCheckoutAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("error parsing current date: %w", err)
	}

	existingChangeRequest, err := dcrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	existingChangeRequest.IsLocked = true
	existingChangeRequest.LastCheckoutBYID = &currentUserID
	existingChangeRequest.LastCheckoutAt = lastCheckoutAt
	err = dcrs.repo.Update(existingChangeRequest, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	var documentPath string
	fileType := "docx"
	documentPath, err = dcrs.GetDocumentPath(id, token, fileType, tx)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("error retrieving document path: %w", err)
	}

	//document cr history
	documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
		DocumentChangeRequestID: id,
		ActionType:              "check-out",
	}

	err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, nil, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return "", err
	}

	tx.Commit()
	return documentPath, nil
}

func (dcrs *DocumentChangeRequestService) CheckinDocument(id uint, currentUserID uint, documentCheckinFormDTO *documentDTOs.DocumentCheckinFormDTO) error {
	tx := dcrs.repo.BeginTransaction()
	documentChangeRequest, err := dcrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = dcrs.ValidateDocument(&documentCheckinFormDTO.Document)
	if err != nil {
		return err
	}

	if documentChangeRequest.DocumentChangeRequestStatus.Type != "draft" && documentChangeRequest.DocumentChangeRequestStatus.Type != "in-review" {
		tx.Rollback()
		return fmt.Errorf("checkin only allowed for documents in 'draft' or 'in-review' status")
	}

	if !documentChangeRequest.IsLocked {
		tx.Rollback()
		return fmt.Errorf("document already checked in")
	}

	if documentChangeRequest.CreatedByUser.ID != int(currentUserID) && documentChangeRequest.Document.ResponsibleUser.ID != int(currentUserID) {
		tx.Rollback()
		return fmt.Errorf("access denied: user is neither creator nor responsible user")
	}

	dataToEncrypt := []byte(documentCheckinFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentchangeRequestDocument, err := dcrs.repo.GetDocumentByID(*documentChangeRequest.ChangeRequestDocument.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentchangeRequestDocument.Name = documentCheckinFormDTO.Document.FileName
	documentchangeRequestDocument.Size = documentCheckinFormDTO.Document.Size
	documentchangeRequestDocument.Token = encryptedFilename
	documentchangeRequestDocument.OriginalName = &documentCheckinFormDTO.Document.OriginalName
	documentchangeRequestDocument.Type = documentCheckinFormDTO.Document.Type

	err = dcrs.repo.UpdateDocument(documentchangeRequestDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	lastCheckoutIn, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	existingDocumentChangeRequest, err := dcrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingDocumentChangeRequest.IsLocked = false
	existingDocumentChangeRequest.LastCheckinBYID = &currentUserID
	existingDocumentChangeRequest.LastCheckinAt = lastCheckoutIn
	err = dcrs.repo.Update(existingDocumentChangeRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	destinationDir := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")

	deleteFilePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
	if err := utility.RemoveFolder(deleteFilePath); err != nil {
		tx.Rollback()
	}

	if err := utility.StoreFile(documentCheckinFormDTO.Document.FileName, destinationDir); err != nil {
		tx.Rollback()
		return err
	}

	//add document cr history
	documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
		DocumentChangeRequestID: id,
		ActionType:              "check-in",
		FileName:                &documentCheckinFormDTO.Document.FileName,
		OriginalName:            &documentCheckinFormDTO.Document.OriginalName,
		Type:                    &documentCheckinFormDTO.Document.Type,
		Size:                    &documentCheckinFormDTO.Document.Size,
		Comment:                 documentCheckinFormDTO.Comment,
	}

	err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, &destinationDir, currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		if err := utility.ConvertDocxToPDF(destinationDir, documentCheckinFormDTO.Document.FileName, destinationDir); err != nil {
			log.Printf("Failed to convert document to PDF: %v", err)
		}
	}()

	return nil
}

func (dcrs *DocumentChangeRequestService) RevertDocumentCheckout(id uint, currentUserID uint) error {
	tx := dcrs.repo.BeginTransaction()
	documentCRFullDetails, err := dcrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if !documentCRFullDetails.IsLocked {
		tx.Rollback()
		return fmt.Errorf("document already checked in")
	}

	existingDocumentCR, err := dcrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingDocumentCR.IsLocked = false

	err = dcrs.repo.Update(existingDocumentCR, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//add document history
	// documentEditHistory := &documentDTOs.DocumentEditHistoryDTO{
	// 	DocumentID:        id,
	// 	DocumentVersionID: document.LatestDocumentVersion.ID,
	// 	ActionType:        "revert-checkout",
	// }

	// err = ds.CreateDocumentEditHistory(documentEditHistory, nil, currentUserID, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	tx.Commit()

	return nil
}

func (dcrs *DocumentChangeRequestService) SubmitDocumentChangeRequest(id uint, currentUserID uint) error {
	tx := dcrs.repo.BeginTransaction()
	documentCRFullDetails, err := dcrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentCRStatus, err := dcrs.documentChangeRequestStatusSvc.GetByType("in-review")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingCRDocument, err := dcrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingCRDocument.DocumentChangeRequestStatusID = documentCRStatus.ID

	err = dcrs.repo.Update(existingCRDocument, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// documentVersion, err := ds.repo.GetDocumentVersionByID(document.LatestDocumentVersion.ID, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	var level int
	if documentCRFullDetails.DocumentChangeRequestStatus.Type == "reverted" {
		// If status is "reverted", use the NextReviewLevel for resubmission
		if documentCRFullDetails.NextReviewLevel != nil {
			level = *documentCRFullDetails.NextReviewLevel
		} else {
			tx.Rollback()
			return fmt.Errorf("next review level is not set for reverted document")
		}
	} else {
		// For non-reverted documents, start from level 1
		level = 1
	}

	currentDateStr := time.Now().Format("2006-01-02")
	submittedAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	if documentCRFullDetails.DocumentChangeRequestStatus.Type != "reverted" {
		existingCRDocument.SubmittedBYID = &currentUserID
		existingCRDocument.SubmittedAt = submittedAt
		existingCRDocument.NextReviewLevel = &level

		err = dcrs.repo.Update(existingCRDocument, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Fetch the workflow level details, including participants (review users)
	docCRWorkflowLevel, err := dcrs.documentChangeRequestWorkflowSvc.GetDocumentCRWorkflowByLevel(documentCRFullDetails.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if documentCRFullDetails.DocumentChangeRequestStatus.Type == "reverted" {
		documentCRWorkflowStatus, err := dcrs.documentCRWorkflowStatusSvc.GetByType("resubmit")
		if err != nil {
			return err
		}
		documentCRWorkflowReview := &models.DocumentChangeRequestWorkflowReview{
			DocumentChangeRequestWorkflowID:       docCRWorkflowLevel.ID,
			DocumentChangeRequestWorkflowStatusID: documentCRWorkflowStatus.ID,
			Comment:                               "Resubmitted",
			CreatedBYID:                           currentUserID,
		}

		err = dcrs.repo.CreateDocumentCRWorkflowReview(documentCRWorkflowReview, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	var taskIDs []uint

	for _, participant := range docCRWorkflowLevel.DocumentChangeRequestWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:                "document",
			DocumentChangeRequestID:         &id,
			DocumentChangeRequestWorkflowID: &docCRWorkflowLevel.ID,
			// Title:                           existingDocument.TitleEn,
			// Description:                     existingDocument.DescriptionEn,
			URLParam:          &urlParam,
			QueryParam:        nil,
			ResponsibleUserID: participant.User.ID,
		}

		taskID, err := dcrs.taskSvc.Create(taskDTO, currentUserID)
		if err != nil {
			return err
		}
		taskIDs = append(taskIDs, taskID)

	}

	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			return err
		}

		err = ds.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	// to do
	// activity log

	//Prepare notification
	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, id)
		entity := "document"
		action := "review-user-assign"

		for _, participant := range docWorkflowLevel.DocumentWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               participant.User.ID,
				ModuleIdentifier:     "document",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      existingDocument.ReferenceNumber,
				HeadingEn:            "Document Review",
				HeadingAr:            "مراجعة الوثيقة",
				TitleEn:              existingDocument.TitleEn,
				TitleAr:              existingDocument.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := ds.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := ds.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Document/%s: Review Assigned", existingDocument.TitleEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"DocumentTitle":         existingDocument.TitleEn,
				"DocumentType":          document.DocumentType.TitleEn,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/document/register/%d", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "document/templates/email/document_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = ds.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (dcrs *DocumentChangeRequestService) UpdateDocumentChangeRequest(documentChangeRequest *models.DocumentChangeRequest, tx ...utility.DatabaseTransaction) error {
	err := dcrs.repo.Update(documentChangeRequest, tx...)
	if err != nil {
		return err
	}

	return nil
}

func (dcrs *DocumentChangeRequestService) GetDocumentPath(id uint, token string, fileType string, tx ...utility.DatabaseTransaction) (string, error) {
	existingDocument, err := dcrs.GetFullDetailsByID(id, tx...)
	if err != nil {
		return "", err
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	baseFilename := strings.TrimSuffix(decryptedFilename, filepath.Ext(decryptedFilename))

	storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(existingDocument.ID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, baseFilename+"."+fileType)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (dcrs *DocumentChangeRequestService) CreateDocumentCREditHistory(documentCREditHistoryDTO *documentDTOs.DocumentCREditHistoryDTO, sourceDir *string, currentUserId uint, tx ...utility.DatabaseTransaction) error {
	if sourceDir != nil {
		dataToEncrypt := []byte(*documentCREditHistoryDTO.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return err
		}
		url := utility.GenerateDocumentURL("document", "change-requests", strconv.FormatUint(uint64(documentCREditHistoryDTO.DocumentChangeRequestID), 10), "document-history")
		documentCREditHistory := &models.DocumentChangeRequestEditHistory{
			DocumentChangeRequestID: documentCREditHistoryDTO.DocumentChangeRequestID,
			ActionType:              documentCREditHistoryDTO.ActionType,
			Comment:                 documentCREditHistoryDTO.Comment,
			URL:                     &url,
			Name:                    documentCREditHistoryDTO.FileName,
			Size:                    documentCREditHistoryDTO.Size,
			Token:                   &encryptedFilename,
			OriginalName:            documentCREditHistoryDTO.OriginalName,
			Type:                    documentCREditHistoryDTO.Type,
			CreatedBYID:             currentUserId,
		}
		id, err := dcrs.repo.CreateDocumentCREditHistory(documentCREditHistory, tx...)
		if err != nil {
			return err
		}
		destinationDir := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(documentCREditHistoryDTO.DocumentChangeRequestID), 10), "document-history", strconv.FormatUint(uint64(id), 10), "documents")
		err = utility.CopyFileFromPaths(*sourceDir, destinationDir, *documentCREditHistoryDTO.FileName)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}
	} else {
		documentCREditHistory := &models.DocumentChangeRequestEditHistory{
			DocumentChangeRequestID: documentCREditHistoryDTO.DocumentChangeRequestID,
			ActionType:              documentCREditHistoryDTO.ActionType,
			CreatedBYID:             currentUserId,
		}

		_, err := dcrs.repo.CreateDocumentCREditHistory(documentCREditHistory, tx...)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dcrs *DocumentChangeRequestService) GetDocumentCRLastEditSummary(id uint) (*documentDTOs.DocumentLastEditSummaryDTO, error) {
	documentChangeRequest, err := dcrs.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	// Fetch the latest edit activity
	activity, err := dcrs.repo.GetDocumentCRLatestEditHistory(id)
	var latestActivityDTO documentDTOs.ActivityDTO

	// Handle case where activity is not found
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			emptyString := ""
			zeroTime := time.Time{}
			latestActivityDTO = documentDTOs.ActivityDTO{
				Activity:  "No activity found",
				By:        nil,
				Comment:   &emptyString,
				Timestamp: zeroTime,
			}
		} else {
			return nil, err
		}
	} else {
		userDTO := mapUserToDTO(activity.CreatedByUser)
		latestActivityDTO = documentDTOs.ActivityDTO{
			Activity:  activity.ActionType,
			By:        userDTO,
			Comment:   activity.Comment,
			Timestamp: activity.CreatedAt,
		}
	}

	// Initialize variables for the current document workflow
	var currentDocumentCRWorkflow *models.DocumentChangeRequestWorkflow
	var participants []*organizationDTOs.UserMinimalDTO
	totalComments := 0

	// Check if there is a next review level and fetch the corresponding workflow
	if documentChangeRequest.NextReviewLevel != nil {
		currentDocumentCRWorkflow, err = dcrs.documentChangeRequestWorkflowSvc.GetDocumentCRWorkflowByLevel(id, *documentChangeRequest.NextReviewLevel)
		if err != nil {
			return nil, err
		}
	}

	// Process workflow participants and count comments
	if currentDocumentCRWorkflow != nil {
		for _, participant := range currentDocumentCRWorkflow.DocumentChangeRequestWorkflowParticipants {
			userDTO := mapUserToDTO(participant.User)
			if userDTO != nil {
				participants = append(participants, userDTO)
			}
		}

		for _, review := range currentDocumentCRWorkflow.DocumentChangeRequestWorkflowReviews {
			if review.Comment != "" {
				totalComments++
			}
		}
	}

	// Prepare the response DTO
	responseDTO := &documentDTOs.DocumentLastEditSummaryDTO{
		LatestActivity: latestActivityDTO,
		CurrentLevel:   documentChangeRequest.NextReviewLevel,
		Participants:   participants,
		TotalComments:  &totalComments,
	}

	return responseDTO, nil
}

func (dcrs *DocumentChangeRequestService) PublishDocumentChangeRequest(id uint, currentUserID uint) error {
	tx := dcrs.repo.BeginTransaction()

	documentCRFullDetails, err := dcrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentFullDetails, err := dcrs.documentSvc.GetFullDetailsByID(uint(documentCRFullDetails.Document.ID), nil, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	documentCRStatus, err := dcrs.documentChangeRequestStatusSvc.GetByType("published")
	if err != nil {
		tx.Rollback()
		return err
	}

	existingDocumentCR, err := dcrs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	publishedAt, err := utility.ParseDate(currentDateStr)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error parsing current date: %w", err)
	}

	existingDocumentCR.DocumentChangeRequestStatusID = documentCRStatus.ID
	existingDocumentCR.PublishedBYID = &currentUserID
	existingDocumentCR.PublishedAt = publishedAt

	err = dcrs.repo.Update(existingDocumentCR, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingLatestDocumentVersion, err := dcrs.documentSvc.GetDocumentVersionByID(documentFullDetails.LatestDocumentVersion.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingLatestDocumentVersion.IsLatestVersion = false

	err = dcrs.documentSvc.UpdateDocumentVersion(existingLatestDocumentVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//generate token
	dataToEncrypt := []byte(documentCRFullDetails.ChangeRequestDocument.Name)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return err
	}

	url := utility.GenerateDocumentURL("document", "published", strconv.FormatUint(uint64(documentCRFullDetails.Document.ID), 10), "documents")

	documentVersion := &models.DocumentVersion{
		DocumentID:      uint(documentCRFullDetails.Document.ID),
		Version:         documentCRFullDetails.Version,
		IsLatestVersion: true,
		URL:             url,
		Name:            documentCRFullDetails.ChangeRequestDocument.Name,
		Size:            documentCRFullDetails.ChangeRequestDocument.Size,
		Token:           encryptedFilename,
		OriginalName:    &documentCRFullDetails.ChangeRequestDocument.OriginalName,
		Type:            documentCRFullDetails.ChangeRequestDocument.Type,
		CreatedBYID:     currentUserID,
		PublishedBYID:   &currentUserID,
		PublishedAt:     publishedAt,
	}

	_, err = dcrs.documentSvc.CreateDocumentVersion(documentVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	fileName := documentCRFullDetails.ChangeRequestDocument.Name
	baseFileName := strings.TrimSuffix(fileName, filepath.Ext(fileName)) // Remove the extension to get the base name

	sourceDir := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
	destinationDir := filepath.Join("document", "published", strconv.FormatUint(uint64(documentCRFullDetails.Document.ID), 10), "documents")

	//copy docx file
	err = utility.CopyFileFromPaths(sourceDir, destinationDir, fileName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error copying .docx file: %w", err)
	}

	//copy pdf file
	pdfFileName := baseFileName + ".pdf"
	err = utility.CopyFileFromPaths(sourceDir, destinationDir, pdfFileName)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error copying .pdf file: %w", err)
	}

	tx.Commit()

	//ds.notifyUsersOnPublish(document, id, currentUserID)

	//to do
	//activity log

	return nil
}
