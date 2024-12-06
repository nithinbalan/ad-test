package services

import (
	"errors"
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	documentRepositories "jethur/grc/documents/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentChangeRequestInitiateServiceInterface interface {
	GetAll(params *documentDTOs.DocumentChangeRequestInitiateQueryParamsDTO) (*documentDTOs.DocumentChangeRequestInitiateListResponseDTO, error)
	GetByID(id uint) (*models.DocumentChangeRequestInitiate, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestInitiateDetailsDTO, error)
	Create(currentUserID uint, documentID uint, formDTO *documentDTOs.DocumentChangeRequestInitiateFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestInitiateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type DocumentChangeRequestInitiateService struct {
	repo                                   documentRepositories.DocumentChangeRequestInitiateRepositoryInterface
	activityLogSvc                         services.ActivityLogServiceInterface
	documentChangeRequestInitiateTypeSvc   DocumentChangeRequestInitiateTypeServiceInterface
	documentChangeRequestInitiateStatusSvc DocumentChangeRequestInitiateStatusServiceInterface
}

func NewDocumentChangeRequestInitiateService(
	repo documentRepositories.DocumentChangeRequestInitiateRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	documentChangeRequestInitiateTypeSvc DocumentChangeRequestInitiateTypeServiceInterface,
	documentChangeRequestInitiateStatusSvc DocumentChangeRequestInitiateStatusServiceInterface,
) DocumentChangeRequestInitiateServiceInterface {
	return &DocumentChangeRequestInitiateService{
		repo:                                   repo,
		activityLogSvc:                         activityLogSvc,
		documentChangeRequestInitiateTypeSvc:   documentChangeRequestInitiateTypeSvc,
		documentChangeRequestInitiateStatusSvc: documentChangeRequestInitiateStatusSvc,
	}
}

func (dcris *DocumentChangeRequestInitiateService) GetAll(params *documentDTOs.DocumentChangeRequestInitiateQueryParamsDTO) (*documentDTOs.DocumentChangeRequestInitiateListResponseDTO, error) {
	return dcris.repo.GetAll(params)
}

func (dcris *DocumentChangeRequestInitiateService) GetByID(id uint) (*models.DocumentChangeRequestInitiate, error) {
	initiate, err := dcris.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return initiate, nil
}

func (dcris *DocumentChangeRequestInitiateService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentChangeRequestInitiateDetailsDTO, error) {
	var detailsDTO documentDTOs.DocumentChangeRequestInitiateDetailsDTO
	initiate, err := dcris.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(initiate, &detailsDTO)
	if err != nil {
		return nil, err
	}
	detailsDTO.CreatedAt = initiate.CreatedAt
	detailsDTO.UpdatedAt = initiate.UpdatedAt
	return &detailsDTO, nil
}

func (dcris *DocumentChangeRequestInitiateService) Create(currentUserID uint, documentID uint, formDTO *documentDTOs.DocumentChangeRequestInitiateFormDTO) (uint, error) {
	var tx = dcris.repo.BeginTransaction()

	if err := dcris.validateFormDTO(formDTO); err != nil {
		return 0, err
	}

	changeRequestInitiateStatus, err := dcris.documentChangeRequestInitiateStatusSvc.GetByType("initiated")
	if err != nil {
		return 0, err
	}

	initiate := &models.DocumentChangeRequestInitiate{
		DocumentsID:                           documentID,
		DocumentChangeRequestInitiateTypeID:   formDTO.DocumentChangeRequestInitiateTypeID,
		ReasonForChange:                       formDTO.ReasonForChange,
		DocumentChangeRequestInitiateStatusID: changeRequestInitiateStatus.ID,
		CreatedBYID:                           currentUserID,
	}

	id, err := dcris.repo.Create(initiate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.DocumentChangeRequestInitiateDocument
	for _, doc := range formDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		url := utility.GenerateDocumentURL("document", "change_request_initiates", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.DocumentChangeRequestInitiateDocument{
			DocumentChangeRequestInitiateID: id,
			URL:                             url,
			Name:                            doc.FileName,
			Size:                            doc.Size,
			Token:                           encryptedFileName,
			OriginalName:                    &doc.OriginalName,
			Type:                            doc.Type,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := dcris.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("document", "change_request_initiates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (dcris *DocumentChangeRequestInitiateService) Update(id uint, currentUserID uint, formDTO *documentDTOs.DocumentChangeRequestInitiateFormDTO) error {
	tx := dcris.repo.BeginTransaction()

	existing, err := dcris.GetByID(id)
	if err != nil {
		return err
	}

	if err := dcris.validateFormDTO(formDTO); err != nil {
		return err
	}

	existing.DocumentChangeRequestInitiateTypeID = formDTO.DocumentChangeRequestInitiateTypeID
	existing.ReasonForChange = formDTO.ReasonForChange
	existing.UpdatedBYID = &currentUserID

	if err := dcris.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.DocumentChangeRequestInitiateDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.DocumentChangeRequestInitiateDocument

	for _, doc := range formDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("document", "change_request_initiates", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
			document := &models.DocumentChangeRequestInitiateDocument{
				DocumentChangeRequestInitiateID: id,
				URL:                             url,
				Name:                            doc.FileName,
				Size:                            doc.Size,
				Token:                           encryptedFileName,
				OriginalName:                    &doc.OriginalName,
				Type:                            doc.Type,
			}
			documents = append(documents, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := dcris.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}
				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := dcris.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := dcris.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("document", "change_request_initiates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to store file %s: %v\n", doc.Name, err)
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("document", "change_request_initiates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", doc.Name, err)
		}
	}

	return nil
}

func (dcris *DocumentChangeRequestInitiateService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := dcris.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*documentDTOs.DocumentChangeRequestInitiateDocumentDTO)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		initiate, err := dcris.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if initiate.Documents != nil && len(initiate.Documents) > 0 {
			for _, document := range initiate.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], document)
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}

			if len(deleteDocumentIds) > 0 {
				if err := dcris.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}
		}

		if err := dcris.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		for initiateId, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("document", "change_request_initiates", strconv.FormatUint(uint64(initiateId), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}

	tx.Commit()

	go func() {
		err := dcris.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.DocumentChangeRequestInitiate{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return failedIds, nil
}

func (dcris *DocumentChangeRequestInitiateService) validateFormDTO(dto *documentDTOs.DocumentChangeRequestInitiateFormDTO) error {
	if _, err := dcris.documentChangeRequestInitiateTypeSvc.GetByID(dto.DocumentChangeRequestInitiateTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("document change request initiate type does not exist"), "not_exist", "change request initiate type")
	}
	return nil
}
