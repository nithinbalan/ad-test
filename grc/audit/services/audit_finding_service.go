package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditFindingServiceInterface interface {
	GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error)
	GetByID(id uint) (*models.AuditFinding, error)
	Create(currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AuditFindingService struct {
	repo                  repositories.AuditFindingRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	auditSvc              AuditServiceInterface
	auditFindingStatusSvc AuditFindingStatusServiceInterface
	auditTimelineSvc      AuditTimelineServiceInterface
	auditableItemSvc      AuditableItemServiceInterface
	auditFindingCategory  AuditFindingCategoryServiceInterface
	severityLevelSvc      masterServices.SeverityLevelServiceInterface
}

func NewAuditFindingService(
	repo repositories.AuditFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditFindingStatusSvc AuditFindingStatusServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	auditFindingCategory AuditFindingCategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
) AuditFindingServiceInterface {
	return &AuditFindingService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		auditSvc:              auditSvc,
		auditFindingStatusSvc: auditFindingStatusSvc,
		auditTimelineSvc:      auditTimelineSvc,
		auditableItemSvc:      auditableItemSvc,
		auditFindingCategory:  auditFindingCategory,
		severityLevelSvc:      severityLevelSvc,
	}
}

func (afs *AuditFindingService) GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO) (*auditDTOs.AuditFindingListResponseDTO, error) {

	return afs.repo.GetAll(params)
}

func (afs *AuditFindingService) GetByID(id uint) (*models.AuditFinding, error) {
	auditFinding, err := afs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditFinding, nil
}

func (afs *AuditFindingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error) {
	var auditFindingDTO auditDTOs.AuditFindingDetailDTO

	auditFinding, err := afs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFinding, &auditFindingDTO)
	if err != nil {
		return nil, err
	}

	auditFindingDTO.CreatedAt = auditFinding.CreatedAt
	auditFindingDTO.UpdatedAt = auditFinding.UpdatedAt

	return &auditFindingDTO, nil
}

func (afs *AuditFindingService) Create(currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) (uint, error) {
	tx := afs.repo.BeginTransaction()

	if err := afs.validateAuditFindingFormDTO(auditFindingFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditFindingStatus, err := afs.auditFindingStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	referenceNumber, err := afs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	auditFinding := &models.AuditFinding{
		ReferenceNumber:        referenceNumber,
		Title:                  auditFindingFormDTO.Title,
		Description:            auditFindingFormDTO.Description,
		AuditID:                auditFindingFormDTO.AuditID,
		AuditableItemID:        auditFindingFormDTO.AuditableItemID,
		AuditFindingCategoryID: auditFindingFormDTO.AuditFindingCategoryID,
		SeverityLevelID:        auditFindingFormDTO.SeverityLevelID,
		Recommendations:        auditFindingFormDTO.Recommendations,
		AuditFindingStatusID:   auditFindingStatus.ID,
		CreatedAt:              time.Now(),
		CreatedBYID:            currentUserID,
	}

	id, err := afs.repo.Create(auditFinding, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.AuditFindingDocument

	// Iterate over each document in the array
	for _, doc := range auditFindingFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new IncidentDocument instance
		document := &models.AuditFindingDocument{
			AuditFindingID: id,
			URL:            url,
			Name:           doc.FileName,
			Size:           doc.Size,
			Token:          encryptedFilename,
			OriginalName:   doc.OriginalName,
			Type:           doc.Type,
			CreatedByID:    currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the IncidentDocument
		if err := afs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := afs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditFinding{}.TableName(),
			id,
			auditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (afs *AuditFindingService) Update(id uint, currentUserID uint, auditFindingFormDTO *auditDTOs.AuditFindingFormDTO) error {
	tx := afs.repo.BeginTransaction()

	existingAuditFinding, err := afs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := afs.validateAuditFindingFormDTO(auditFindingFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditFinding.AuditableItemID = auditFindingFormDTO.AuditableItemID
	existingAuditFinding.AuditFindingCategoryID = auditFindingFormDTO.AuditFindingCategoryID
	existingAuditFinding.SeverityLevelID = auditFindingFormDTO.SeverityLevelID
	existingAuditFinding.Recommendations = auditFindingFormDTO.Recommendations
	existingAuditFinding.UpdatedBYID = &currentUserID

	err = afs.repo.Update(existingAuditFinding, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//documents update
	var documents []*models.AuditFindingDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditFindingDocument

	// Iterate over each document in the array
	for _, doc := range auditFindingFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("audit", "audit-findings", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new AuditFindingDocument instance
			document := &models.AuditFindingDocument{
				AuditFindingID: id,
				URL:            url,
				Name:           doc.FileName,
				Size:           doc.Size,
				Token:          encryptedFilename,
				OriginalName:   doc.OriginalName,
				Type:           doc.Type,
				CreatedByID:    currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := afs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
				documentsToDelete = append(documentsToDelete, existingDocument)
			}
		}
	}
	if len(documents) > 0 {
		if err := afs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := afs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(existingAuditFinding.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(existingAuditFinding.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		err := afs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditFinding{}.TableName(),
			id,
			auditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (afs *AuditFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

	tx := afs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*auditDTOs.AuditFindingDocumentDTO)

	for _, id := range ids {
		// Get existing audit finding details
		auditFinding, err := afs.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if len(auditFinding.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range auditFinding.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := afs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the audit finding itself
		if err := afs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		for auditFindingID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditFindingID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}

	tx.Commit()
	if len(successfulIDs) > 0 {

		go func() {
			for _, id := range successfulIDs {
				err := afs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditFinding{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditFinding{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (afs *AuditFindingService) validateAuditFindingFormDTO(dto *auditDTOs.AuditFindingFormDTO) error {

	if _, err := afs.auditSvc.GetByID(dto.AuditID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit does not exist"), "not_exist", "audit")
	}

	if _, err := afs.auditableItemSvc.GetByID(dto.AuditableItemID); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
	}

	if _, err := afs.auditFindingCategory.GetByID(dto.AuditFindingCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit finding category does not exist"), "not_exist", "audit finding category")
	}

	if _, err := afs.severityLevelSvc.GetByID(dto.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	return nil
}
