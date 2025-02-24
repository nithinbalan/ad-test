package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepositories "jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationSvc "jethur/grc/organization/services"
	settingsSvc "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestServiceInterface interface {
	SetDocumentChangeRequestWorkflowService(svc DocumentChangeRequestWorkflowServiceInterface)

	GetAll(params *documentDTOs.DocumentChangeRequestQueryParamsDTO) (*documentDTOs.DocumentChangeRequestListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestDetailsDTO, error)
	Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
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

	ExportDocumentChangeRequest(data []*documentDTOs.DocumentChangeRequestListDTO, currentUserID uint) (*os.File, error)
}

type DocumentChangeRequestService struct {
	repo                             documentRepositories.DocumentChangeRequestRepositoryInterface
	activityLogSvc                   services.ActivityLogServiceInterface
	documentChangeRequestTypeSvc     DocumentChangeRequestTypeServiceInterface
	documentChangeRequestStatusSvc   DocumentChangeRequestStatusServiceInterface
	documentChangeRequestWorkflowSvc DocumentChangeRequestWorkflowServiceInterface
	documentCRWorkflowStatusSvc      DocumentChangeRequestWorkflowStatusServiceInterface
	documentSvc                      DocumentServiceInterface
	documentStatusSvc                DocumentStatusServiceInterface
	taskSvc                          taskServices.TaskServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
	userSvc                          organizationSvc.UserServiceInterface
	documentAccessSvc                DocumentAccessServiceInterface
	orgLevelSvc                      settingsSvc.OrganizationLevelServiceInterface
	userSettingSvc                   organizationSvc.UserSettingServiceInterface
}

func NewDocumentChangeRequestService(
	repo documentRepositories.DocumentChangeRequestRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentChangeRequestTypeSvc DocumentChangeRequestTypeServiceInterface,
	documentChangeRequestStatusSvc DocumentChangeRequestStatusServiceInterface,
	documentCRWorkflowStatusSvc DocumentChangeRequestWorkflowStatusServiceInterface,
	documentSvc DocumentServiceInterface,
	documentStatusSvc DocumentStatusServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	documentAccessSvc DocumentAccessServiceInterface,
	orgLevelSvc settingsSvc.OrganizationLevelServiceInterface,
	userSettingSvc organizationSvc.UserSettingServiceInterface,
) DocumentChangeRequestServiceInterface {
	return &DocumentChangeRequestService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		documentChangeRequestTypeSvc:   documentChangeRequestTypeSvc,
		documentChangeRequestStatusSvc: documentChangeRequestStatusSvc,
		documentCRWorkflowStatusSvc:    documentCRWorkflowStatusSvc,
		documentSvc:                    documentSvc,
		documentStatusSvc:              documentStatusSvc,
		taskSvc:                        taskSvc,
		notificationSvc:                notificationSvc,
		userSvc:                        userSvc,
		documentAccessSvc:              documentAccessSvc,
		orgLevelSvc:                    orgLevelSvc,
		userSettingSvc:                 userSettingSvc,
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

func (dcrs *DocumentChangeRequestService) Create(currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {
	var transaction utility.DatabaseTransaction

	if len(tx) > 0 && tx[0] != nil {
		fmt.Println("123")
		transaction = tx[0]
	} else {
		fmt.Println("124")
		transaction = dcrs.repo.BeginTransaction()
		defer transaction.Rollback() // Ensure rollback if transaction was started here and not committed
	}

	referenceNumber, err := dcrs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	if err := dcrs.validateFormDTO(formDTO); err != nil {
		return 0, err
	}
	fmt.Println("125")
	existingChangeRequests, err := dcrs.repo.GetNonPublishedChangeRequests(formDTO.DocumentID)
	if err != nil {
		return 0, err
	}

	if len(existingChangeRequests) > 0 {
		return 0, utility.NewResourceConflictError(errors.New("cannot create a new change request for this document as there are existing change requests not in a published status"))
	}

	changeRequestStatus, err := dcrs.documentChangeRequestStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}
	fmt.Println("126")
	documentChangeRequestType, err := dcrs.documentChangeRequestTypeSvc.GetByID(formDTO.DocumentChangeRequestTypeID)
	if err != nil {
		return 0, err
	}

	documentFullDetails, err := dcrs.documentSvc.GetFullDetailsByID(formDTO.DocumentID, nil, transaction)
	if err != nil {
		return 0, err
	}
	fmt.Println("127")
	var newVersion float64
	currentVersion := documentFullDetails.LatestDocumentVersion.Version
	if documentChangeRequestType.Type == "minor" {
		newVersion = currentVersion + 0.1
	} else if documentChangeRequestType.Type == "major" {
		newVersion = currentVersion + 1.0
	} else if documentChangeRequestType.Type == "restore-with-minor-update" {
		newVersion = currentVersion + 0.1
	} else if documentChangeRequestType.Type == "restore-with-major-update" {
		newVersion = currentVersion + 1.0
	}
	// else {
	// 	tx.Rollback()
	// 	return 0, fmt.Errorf("invalid document change request type")
	// }

	changeRequest := &models.DocumentChangeRequest{
		ReferenceNumber:                 referenceNumber,
		DocumentID:                      formDTO.DocumentID,
		DocumentChangeRequestInitiateID: formDTO.DocumentChangeRequestInitiateID,
		DocumentChangeRequestTypeID:     formDTO.DocumentChangeRequestTypeID,
		ReasonForChange:                 formDTO.ReasonForChange,
		DocumentChangeRequestStatusID:   changeRequestStatus.ID,
		Version:                         newVersion,
		CreatedBYID:                     currentUserID,
		DocumentReviewID:                formDTO.DocumentReviewID,
	}

	id, err := dcrs.repo.Create(changeRequest, transaction)
	if err != nil {
		return 0, err
	}
	fmt.Println("128")
	var storagePath string
	var is_pdf bool = false
	if documentChangeRequestType.Type != "archive" && documentChangeRequestType.Type != "restore" {
		dataToEncrypt := []byte(formDTO.Document.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
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

		if err := dcrs.repo.CreateDocument(document, transaction); err != nil {
			return 0, err
		}

		storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(formDTO.Document.FileName, storagePath); err != nil {
			return 0, err
		}

		is_pdf = true

		documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
			DocumentChangeRequestID: id,
			ActionType:              "create",
			FileName:                &formDTO.Document.FileName,
			OriginalName:            &formDTO.Document.OriginalName,
			Type:                    &formDTO.Document.Type,
			Size:                    &formDTO.Document.Size,
		}

		err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, &storagePath, currentUserID, transaction)
		if err != nil {
			return 0, err
		}
	} else {
		documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
			DocumentChangeRequestID: id,
			ActionType:              "create",
		}

		err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, nil, currentUserID, transaction)
		if err != nil {
			return 0, err
		}
	}
	fmt.Println("129")
	err = dcrs.documentChangeRequestWorkflowSvc.CreateDocumentCRWorkflow(id, formDTO.DocumentID, currentUserID, transaction)
	if err != nil {
		return 0, err
	}
	fmt.Println("130")
	// Commit if this function started the transaction
	if len(tx) == 0 {
		transaction.Commit()
	}

	go func() {
		if is_pdf {
			if err := utility.ConvertDocxToPDF(storagePath, formDTO.Document.FileName, storagePath); err != nil {
				log.Printf("Failed to convert document to PDF: %v", err)
			}
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

	if existingChangeRequest.DocumentID != formDTO.DocumentID {
		existingChangeRequests, err := dcrs.repo.GetNonPublishedChangeRequests(formDTO.DocumentID)
		if err != nil {
			return err
		}

		if len(existingChangeRequests) > 0 {
			return utility.NewForbiddenError("cannot create a new change request for this document as there are existing change requests not in a published status")
		}
	}

	documentChangeRequestType, err := dcrs.documentChangeRequestTypeSvc.GetByID(formDTO.DocumentChangeRequestTypeID)
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

	// err = dcrs.repo.Update(existingChangeRequest, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// Handle document changes
	var isNewFileProvided bool
	var storagePath string
	if existingCRFullDetails.DocumentChangeRequestType.Type == "minor" || existingCRFullDetails.DocumentChangeRequestType.Type == "major" {
		existingDocument, err := dcrs.repo.GetDocumentByChangeRequestID(id)
		if err != nil {
			tx.Rollback()
			return err
		}
		if documentChangeRequestType.Type == "minor" || documentChangeRequestType.Type == "major" {
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
				storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
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
		} else if documentChangeRequestType.Type == "archive" {
			existingChangeRequest.Version = 0
			isNewFileProvided = true
			existingDocumentFileName := existingDocument.Name
			baseFileName := strings.TrimSuffix(existingDocumentFileName, filepath.Ext(existingDocumentFileName))

			err := dcrs.repo.DeleteDocument(existingDocument.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
				DocumentChangeRequestID: id,
				ActionType:              "update",
			}

			err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, nil, currentUserID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Remove old file and store the new file
			storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(existingDocument.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
			pdfFileName := baseFileName + ".pdf"
			if err := utility.RemoveFile(pdfFileName, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	} else if existingCRFullDetails.DocumentChangeRequestType.Type == "archive" {
		if documentChangeRequestType.Type == "minor" || documentChangeRequestType.Type == "major" {
			documentFullDetails, err := dcrs.documentSvc.GetFullDetailsByID(formDTO.DocumentID, nil, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			var newVersion float64
			currentVersion := documentFullDetails.LatestDocumentVersion.Version
			if documentChangeRequestType.Type == "minor" {
				newVersion = currentVersion + 0.1
			} else if documentChangeRequestType.Type == "major" {
				newVersion = currentVersion + 1.0
			}

			existingChangeRequest.Version = newVersion

			if formDTO.Document.IsNewFile != nil && *formDTO.Document.IsNewFile {
				isNewFileProvided = true
				dataToEncrypt := []byte(formDTO.Document.FileName)
				encryptedFileName, err := utility.EncryptData(dataToEncrypt)
				if err != nil {
					tx.Rollback()
					return err
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
					return err
				}

				storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
				if err := utility.StoreFile(formDTO.Document.FileName, storagePath); err != nil {
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
		}
	} else if existingCRFullDetails.DocumentChangeRequestType.Type == "restore" {
		if documentChangeRequestType.Type == "restore-with-minor-update" || documentChangeRequestType.Type == "restore-with-major-update" {
			documentFullDetails, err := dcrs.documentSvc.GetFullDetailsByID(formDTO.DocumentID, nil, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			var newVersion float64
			currentVersion := documentFullDetails.LatestDocumentVersion.Version
			if documentChangeRequestType.Type == "restore-with-minor-update" {
				newVersion = currentVersion + 0.1
			} else if documentChangeRequestType.Type == "restore-with-major-update" {
				newVersion = currentVersion + 1.0
			}

			existingChangeRequest.Version = newVersion

			if formDTO.Document.IsNewFile != nil && *formDTO.Document.IsNewFile {
				isNewFileProvided = true
				dataToEncrypt := []byte(formDTO.Document.FileName)
				encryptedFileName, err := utility.EncryptData(dataToEncrypt)
				if err != nil {
					tx.Rollback()
					return err
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
					return err
				}

				storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
				if err := utility.StoreFile(formDTO.Document.FileName, storagePath); err != nil {
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
		}
	} else if existingCRFullDetails.DocumentChangeRequestType.Type == "restore-with-minor-update" || existingCRFullDetails.DocumentChangeRequestType.Type == "restore-with-major-update" {
		existingDocument, err := dcrs.repo.GetDocumentByChangeRequestID(id)
		if err != nil {
			tx.Rollback()
			return err
		}
		if documentChangeRequestType.Type == "restore-with-minor-update" || documentChangeRequestType.Type == "restore-with-major-update" {

			documentFullDetails, err := dcrs.documentSvc.GetFullDetailsByID(formDTO.DocumentID, nil, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			var newVersion float64
			currentVersion := documentFullDetails.LatestDocumentVersion.Version
			if documentChangeRequestType.Type == "restore-with-minor-update" {
				newVersion = currentVersion + 0.1
			} else if documentChangeRequestType.Type == "restore-with-major-update" {
				newVersion = currentVersion + 1.0
			}

			existingChangeRequest.Version = newVersion

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
				storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
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
		} else if documentChangeRequestType.Type == "restore" {
			existingChangeRequest.Version = 0
			isNewFileProvided = true
			existingDocumentFileName := existingDocument.Name
			baseFileName := strings.TrimSuffix(existingDocumentFileName, filepath.Ext(existingDocumentFileName))

			err := dcrs.repo.DeleteDocument(existingDocument.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			documentCREditHistory := &documentDTOs.DocumentCREditHistoryDTO{
				DocumentChangeRequestID: id,
				ActionType:              "update",
			}

			err = dcrs.CreateDocumentCREditHistory(documentCREditHistory, nil, currentUserID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Remove old file and store the new file
			storagePath = filepath.Join("document", "change-requests", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(existingDocument.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
			pdfFileName := baseFileName + ".pdf"
			if err := utility.RemoveFile(pdfFileName, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = dcrs.repo.Update(existingChangeRequest, tx)
	if err != nil {
		tx.Rollback()
		return err
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

	documentsToDelete := make(map[uint][]*models.DocumentChangeRequestDocument)
	documentCREditHistoryToDelete := make(map[uint][]*models.DocumentChangeRequestEditHistory)

	for _, id := range ids {
		var deleteDocumentCREditHistoryIds []*uint
		var deleteDocumentWorkflowIds []*uint
		var deleteDocumentWorkflowParticipantIds []*uint
		var deleteDocumentWorkflowReviewIds []*uint
		changeRequest, err := dcrs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if len(changeRequest.DocumentCREditHistory) > 0 {
			for _, history := range changeRequest.DocumentCREditHistory {
				documentCREditHistoryToDelete[id] = append(documentCREditHistoryToDelete[id], history)
				deleteDocumentCREditHistoryIds = append(deleteDocumentCREditHistoryIds, &history.ID)
			}
		}

		if err := dcrs.repo.DeleteDocumentCREditHistory(deleteDocumentCREditHistoryIds, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(changeRequest.DocumentChangeRequestWorkflows) > 0 {
			for _, docWorkflow := range changeRequest.DocumentChangeRequestWorkflows {
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
			return nil, err
		}

		if err := dcrs.documentChangeRequestWorkflowSvc.DeleteDocumentCRWorkflowReviews(deleteDocumentWorkflowReviewIds, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		if err := dcrs.documentChangeRequestWorkflowSvc.DeleteDocumentCRWorkflows(deleteDocumentWorkflowIds, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		if changeRequest.ChangeRequestDocument != nil {
			documentsToDelete[id] = append(documentsToDelete[id], changeRequest.ChangeRequestDocument)

			// Remove the related document from the database
			if err := dcrs.repo.DeleteDocument(changeRequest.ChangeRequestDocument.ID, tx); err != nil {
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
		for documentCRID, history := range documentCREditHistoryToDelete {
			for _, doc := range history {
				if doc.Name != nil {
					storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(documentCRID), 10), "document-history", strconv.FormatUint(uint64(doc.ID), 10), "documents")
					if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
						tx.Rollback()
						fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
					}
				}
			}
			destinationPath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(documentCRID), 10), "document-history")
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}

		for changeRequestID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("document", "change-requests", strconv.FormatUint(uint64(changeRequestID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
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

	documentCRType, err := dcrs.documentChangeRequestTypeSvc.GetByID(dto.DocumentChangeRequestTypeID)
	if err != nil {
		return err
	}
	if documentCRType.Type == "minor" || documentCRType.Type == "major" || documentCRType.Type == "restore-with-minor-update" || documentCRType.Type == "restore-with-major-update" {
		err = dcrs.ValidateDocument(&dto.Document)
		if err != nil {
			return err
		}
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

	if documentChangeRequest.DocumentChangeRequestType.Type == "archive" || documentChangeRequest.DocumentChangeRequestType.Type != "restore" {
		tx.Rollback()
		return "", utility.NewForbiddenError("Checkout not allowed for document change request type archive and restore")
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

	// document cr history
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

	if documentChangeRequest.DocumentChangeRequestType.Type == "archive" || documentChangeRequest.DocumentChangeRequestType.Type != "restore" {
		tx.Rollback()
		return utility.NewForbiddenError("Checkin not allowed for document change request type archive and restore")
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

	// add document cr history
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

	if documentCRFullDetails.DocumentChangeRequestType.Type == "archive" || documentCRFullDetails.DocumentChangeRequestType.Type != "restore" {
		tx.Rollback()
		return utility.NewForbiddenError("Checkout revert not allowed for document change request type archive and restore")
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

	// add document history
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
			// Title:                           documentCRFullDetails.Document.TitleEn,
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

		err = dcrs.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	// to do
	// activity log

	// Prepare notification
	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, id)
		entity := "document-change_request"
		action := "review-user-assign"

		for _, participant := range docCRWorkflowLevel.DocumentChangeRequestWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.User.ID},
				ModuleIdentifier:     "document",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      existingCRDocument.ReferenceNumber,
				HeadingEn:            "Document Change Request Review",
				HeadingAr:            "مراجعة طلب تغيير المستند",
				TitleEn:              documentCRFullDetails.Document.TitleEn,
				TitleAr:              documentCRFullDetails.Document.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := dcrs.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := dcrs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Document Change Request/%s: Review Assigned", documentCRFullDetails.Document.TitleEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":               responsibleUser.FirstNameEn,
				"LastNameEn":                responsibleUser.LastNameEn,
				"DocumentTitle":             documentCRFullDetails.Document.TitleEn,
				"DocumentType":              documentCRFullDetails.Document.DocumentType.TitleEn,
				"DocumentChangeRequestType": documentCRFullDetails.DocumentChangeRequestType.TitleEn,
				"ReasonForChange":           documentCRFullDetails.ReasonForChange,
				"AssignedByFirstNameEn":     assignedBy.FirstNameEn,
				"AssignedByLastNameEn":      assignedBy.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%s/document/change-request/%d", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "document/templates/email/document_change_request_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send Notification and Email
			err = dcrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
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

	if documentCRFullDetails.DocumentChangeRequestType.Type == "archive" {
		documentStatus, err := dcrs.documentStatusSvc.GetByType("archived")
		if err != nil {
			tx.Rollback()
			return err
		}

		existingDocument, err := dcrs.documentSvc.GetByID(uint(documentCRFullDetails.Document.ID), tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingDocument.DocumentStatusID = documentStatus.ID

		err = dcrs.documentSvc.UpdateDocument(existingDocument, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

	} else if documentCRFullDetails.DocumentChangeRequestType.Type == "restore" {
		documentStatus, err := dcrs.documentStatusSvc.GetByType("published")
		if err != nil {
			tx.Rollback()
			return err
		}

		existingDocument, err := dcrs.documentSvc.GetByID(uint(documentCRFullDetails.Document.ID), tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingDocument.DocumentStatusID = documentStatus.ID

		err = dcrs.documentSvc.UpdateDocument(existingDocument, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
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

		// generate token
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

		// copy docx file
		err = utility.CopyFileFromPaths(sourceDir, destinationDir, fileName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error copying .docx file: %w", err)
		}

		// copy pdf file
		pdfFileName := baseFileName + ".pdf"
		err = utility.CopyFileFromPaths(sourceDir, destinationDir, pdfFileName)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error copying .pdf file: %w", err)
		}

	}

	tx.Commit()

	dcrs.notifyUsersOnPublish(documentCRFullDetails, id, currentUserID)

	// to do
	// activity log

	return nil
}

func (dcrs *DocumentChangeRequestService) notifyUsersOnPublish(documentChangeRequest *documentDTOs.DocumentChangeRequestDetailsDTO, documentCRID uint, publishedBy uint) {
	go func() {
		entity := "document-change-request"
		action := "published"
		appURL := utility.GetAppURL()
		subject := fmt.Sprintf("Document Change Request/%s: Published Notification", documentChangeRequest.Document.TitleEn)
		redirectURL := fmt.Sprintf("%s/document/register/%d", appURL, documentCRID)

		publishedUser, _ := dcrs.userSvc.GetByID(publishedBy)

		// Collect users to notify based on access type
		var usersToNotify []*organizationDTOs.UserMinimalDTO

		if documentChangeRequest.Document.DocumentAccessType.Type == "public" {
			params := &organizationDTOs.UserListQueryParamsDTO{
				Sort:        "asc",
				SortBy:      "first_name_en",
				Keyword:     "",
				Status:      "",
				Page:        1,
				PageSize:    0,
				UserIDs:     nil,
				Filters:     nil,
				QuickFilter: nil,
				Export:      true,
			}

			userResponse, err := dcrs.userSvc.GetAll(params)
			if err != nil {
				utility.LogError("Failed to fetch all users for public document notification", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			for _, user := range userResponse.Items {
				minimalUser := &organizationDTOs.UserMinimalDTO{
					ID:          int(user.ID),
					FirstNameEN: user.FirstNameEn,
					LastNameEN:  user.LastNameEn,
					Email:       user.Email,
				}
				usersToNotify = append(usersToNotify, minimalUser)
			}
		} else if documentChangeRequest.Document.DocumentAccessType.Type == "shared" {
			const page = 1
			const pageSize = 10000
			sharedUsers, _ := dcrs.documentAccessSvc.GetSharedUsers(uint(documentChangeRequest.Document.ID), page, pageSize)
			usersToNotify = append(usersToNotify, sharedUsers.Items...)

			reviewUsers, _ := dcrs.documentAccessSvc.GetReviewUsers(uint(documentChangeRequest.Document.ID))
			usersToNotify = append(usersToNotify, reviewUsers...)

			primaryLevel, _ := dcrs.orgLevelSvc.GetPrimaryLevel()
			columnName := fmt.Sprintf("organization_unit_%d_id", primaryLevel.Level)

			sharedUnits, _ := dcrs.documentAccessSvc.GetSharedUnits(uint(documentChangeRequest.Document.ID))
			if len(sharedUnits) > 0 {

				unitIDs := make([]uint, len(sharedUnits))
				for i, unit := range sharedUnits {
					unitIDs[i] = uint(unit.ID)
				}

				// Convert unit IDs to string slice
				unitIDsStr := make([]string, len(unitIDs))
				for i, id := range unitIDs {
					unitIDsStr[i] = strconv.Itoa(int(id))
				}

				params := &organizationDTOs.UserListQueryParamsDTO{
					Sort:        "asc",
					SortBy:      "first_name_en",
					Keyword:     "",
					Status:      "",
					Page:        1,
					PageSize:    0,
					UserIDs:     nil,
					Filters:     map[string][]string{columnName: unitIDsStr},
					QuickFilter: nil,
					Export:      true,
				}

				unitUserResponse, err := dcrs.userSvc.GetAll(params)
				if err != nil {
					utility.LogError("Failed to fetch all users for public document notification", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				for _, user := range unitUserResponse.Items {
					minimalUser := &organizationDTOs.UserMinimalDTO{
						ID:          int(user.ID),
						FirstNameEN: user.FirstNameEn,
						LastNameEN:  user.LastNameEn,
						Email:       user.Email,
					}
					usersToNotify = append(usersToNotify, minimalUser)
				}
			}

		}

		// Send notifications to all users in the usersToNotify list
		uniqueUsers := make(map[uint]bool)
		for _, user := range usersToNotify {
			if _, exists := uniqueUsers[uint(user.ID)]; exists {
				continue
			}

			uniqueUsers[uint(user.ID)] = true

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(user.ID)},
				ModuleIdentifier:     "document",
				NotificationType:     string(commonDTOs.PublishType),
				HeadingEn:            "Document Change Request Published",
				HeadingAr:            "تم نشر طلب تغيير المستند",
				ReferenceNumber:      documentChangeRequest.ReferenceNumber,
				TitleEn:              documentChangeRequest.Document.TitleEn,
				TitleAr:              documentChangeRequest.Document.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             fmt.Sprintf(`{"id":%d}`, documentCRID),
				IsSystemNotification: false,
				CreatedByID:          &publishedBy,
			}

			// Prepare email data
			emailTo := []string{user.Email}
			emailData := map[string]interface{}{
				"FirstNameEn":               user.FirstNameEN,
				"LastNameEn":                user.LastNameEN,
				"DocumentTitle":             documentChangeRequest.Document.TitleEn,
				"DocumentType":              documentChangeRequest.Document.DocumentType.TitleEn,
				"DocumentChangeRequestType": documentChangeRequest.DocumentChangeRequestType.TitleEn,
				"PublishedByFirstNameEn":    publishedUser.FirstNameEn,
				"PublishedByFirstLastEn":    publishedUser.LastNameEn,
				"PublishedAt":               time.Now().Format("2006-01-02"),
				"RedirectURL":               redirectURL,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "document/templates/email/document_change_request_publish_user_mail.html",
				EmailData:            emailData,
			}

			// Send both database and email notifications
			err := dcrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "document_change_request")
			if err != nil {
				utility.LogError("Failed to send document publish notification", map[string]interface{}{
					"error": err.Error(),
				})
			}

		}
	}()
}

func (dcrs *DocumentChangeRequestService) ExportDocumentChangeRequest(data []*documentDTOs.DocumentChangeRequestListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := dcrs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"رقم",               // SL No.
			"رقم الإشارة",       // Ref No.
			"العنوان",           // Title
			"نوع الطلب",         // Change Request Type
			"سبب التغيير",       // Reason for Change
			"الإصدار",           // Version
			"تاريخ الإنشاء",     // Created At
			"أنشئ بواسطة",       // Created By
			"تم التحديث بواسطة", // Updated By
			"الحالة",            // Status
		}
	} else {
		headers = []string{
			"SL No.",              // SL No.
			"Ref No.",             // Ref No.
			"Title",               // Title
			"Change Request Type", // Change Request Type
			"Reason for Change",   // Reason for Change
			"Version",             // Version
			"Created At",          // Created At
			"Created By",          // Created By
			"Updated By",          // Updated By
			"Status",              // Status
		}
	}

	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, documentCR := range data {
		var refNo, title, changeRequestType, reasonForChange, createdAt, createdBy, updatedBy, status string
		var version float64

		// Assign values based on language preference
		if userSetting.Language == "ar" {
			refNo = documentCR.ReferenceNumber
			title = documentCR.DocumentTitleAr
			changeRequestType = documentCR.DocumentChangeRequestTypeAr
			reasonForChange = *documentCR.ReasonForChange
			version = documentCR.Version
			createdAt = documentCR.CreatedAt.Format("02 يناير 2006 03:04 م") // Example Arabic date format
			createdBy = utility.GetStringValue(documentCR.CreatedByUserFirstNameAr) + " " + utility.GetStringValue(documentCR.CreatedByUserLastNameAr)
			updatedBy = utility.GetStringValue(documentCR.UpdatedByUserFirstNameAr) + " " + utility.GetStringValue(documentCR.UpdatedByUserLastNameAr)
			status = documentCR.DocumentChangeRequestStatusTitleAr
		} else {
			refNo = documentCR.ReferenceNumber
			title = documentCR.DocumentTitleEn
			changeRequestType = documentCR.DocumentChangeRequestTypeEn
			reasonForChange = *documentCR.ReasonForChange
			version = documentCR.Version
			createdAt = documentCR.CreatedAt.Format("Jan 02, 2006 03:04 PM") // Example English date format
			createdBy = utility.GetStringValue(documentCR.CreatedByUserFirstNameEn) + " " + utility.GetStringValue(documentCR.CreatedByUserLastNameEn)
			updatedBy = utility.GetStringValue(documentCR.UpdatedByUserFirstNameEn) + " " + utility.GetStringValue(documentCR.UpdatedByUserLastNameEn)
			status = documentCR.DocumentChangeRequestStatusTitleEn
		}

		rows[i] = []interface{}{
			i + 1,             // SL No.
			refNo,             // Ref No.
			title,             // Title
			changeRequestType, // Change Request Type
			reasonForChange,   // Reason for Change
			version,           // Version
			createdAt,         // Created At
			createdBy,         // Created By
			updatedBy,         // Updated By
			status,            // Status
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
