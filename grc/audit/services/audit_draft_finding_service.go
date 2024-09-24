package services

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditDraftFindingServiceInterface interface {
	GetAll(params *auditDTOs.AuditDraftFindingListQueryParamsDTO) (*auditDTOs.AuditDraftFindingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftFindingDetailDTO, error)
	GetByID(id uint) (*models.AuditDraftFinding, error)
	Create(currentUserID uint, auditDraftFindingFormDTO *auditDTOs.AuditDraftFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditDraftFindingFormDTO *auditDTOs.AuditDraftFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetAuditDraftFindingCountByCategory(auditId uint) ([]auditDTOs.AuditDraftFindingCategoryCountResponseDTO, error)
}

type AuditDraftFindingService struct {
	repo                       repositories.AuditDraftFindingRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	auditSvc                   AuditServiceInterface
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface
	auditTimelineSvc           AuditTimelineServiceInterface
	auditableItemSvc           AuditableItemServiceInterface
	auditFindingCategorySvc    AuditFindingCategoryServiceInterface
	severityLevelSvc           masterServices.SeverityLevelServiceInterface
	auditChecklistSvc          AuditChecklistServiceInterface
}

func NewAuditDraftFindingService(
	repo repositories.AuditDraftFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	auditFindingCategorySvc AuditFindingCategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	auditChecklistSvc AuditChecklistServiceInterface,
) AuditDraftFindingServiceInterface {
	return &AuditDraftFindingService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		auditSvc:                   auditSvc,
		auditDraftFindingStatusSvc: auditDraftFindingStatusSvc,
		auditTimelineSvc:           auditTimelineSvc,
		auditableItemSvc:           auditableItemSvc,
		auditFindingCategorySvc:    auditFindingCategorySvc,
		severityLevelSvc:           severityLevelSvc,
		auditChecklistSvc:          auditChecklistSvc,
	}
}

func (adfs *AuditDraftFindingService) GetAll(params *auditDTOs.AuditDraftFindingListQueryParamsDTO) (*auditDTOs.AuditDraftFindingListResponseDTO, error) {

	return adfs.repo.GetAll(params)
}

func (adfs *AuditDraftFindingService) GetByID(id uint) (*models.AuditDraftFinding, error) {
	auditDraftFinding, err := adfs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftFinding, nil
}

func (adfs *AuditDraftFindingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDraftFindingDetailDTO, error) {
	var auditDraftFindingDTO auditDTOs.AuditDraftFindingDetailDTO

	auditDraftFinding, err := adfs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditDraftFinding, &auditDraftFindingDTO)
	if err != nil {
		return nil, err
	}

	auditDraftFindingDTO.CreatedAt = auditDraftFinding.CreatedAt
	auditDraftFindingDTO.UpdatedAt = auditDraftFinding.UpdatedAt

	return &auditDraftFindingDTO, nil
}

func (adfs *AuditDraftFindingService) Create(currentUserID uint, auditDraftFindingFormDTO *auditDTOs.AuditDraftFindingFormDTO) (uint, error) {
	tx := adfs.repo.BeginTransaction()

	if err := adfs.validateAuditDraftFindingFormDTO(auditDraftFindingFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	auditDraftFindingStatus, err := adfs.auditDraftFindingStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	auditDraftFinding := &models.AuditDraftFinding{
		Title:                     auditDraftFindingFormDTO.Title,
		Description:               auditDraftFindingFormDTO.Description,
		AuditID:                   auditDraftFindingFormDTO.AuditID,
		AuditableItemID:           auditDraftFindingFormDTO.AuditableItemID,
		AuditChecklistID:          auditDraftFindingFormDTO.AuditChecklistID,
		AuditFindingCategoryID:    auditDraftFindingFormDTO.AuditFindingCategoryID,
		SeverityLevelID:           &auditDraftFindingFormDTO.SeverityLevelID,
		Recommendations:           auditDraftFindingFormDTO.Recommendations,
		AuditDraftFindingStatusID: auditDraftFindingStatus.ID,
		CreatedBYID:               currentUserID,
	}

	id, err := adfs.repo.Create(auditDraftFinding, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.AuditDraftFindingDocument

	// Iterate over each document in the array
	for _, doc := range auditDraftFindingFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("audit", "draft-audit-findings", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new IncidentDocument instance
		document := &models.AuditDraftFindingDocument{
			AuditDraftFindingID: id,
			URL:                 url,
			Name:                doc.FileName,
			Size:                doc.Size,
			Token:               encryptedFilename,
			OriginalName:        doc.OriginalName,
			Type:                doc.Type,
			CreatedBYID:         currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the IncidentDocument
		if err := adfs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "draft-audit-findings", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := adfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFinding{}.TableName(),
			id,
			auditDraftFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (adfs *AuditDraftFindingService) Update(id uint, currentUserID uint, auditDraftFindingFormDTO *auditDTOs.AuditDraftFindingFormDTO) error {
	tx := adfs.repo.BeginTransaction()

	existingAuditDraftFinding, err := adfs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := adfs.validateAuditDraftFindingFormDTO(auditDraftFindingFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditDraftFinding.AuditableItemID = auditDraftFindingFormDTO.AuditableItemID
	existingAuditDraftFinding.Title = auditDraftFindingFormDTO.Title
	existingAuditDraftFinding.Description = auditDraftFindingFormDTO.Description
	existingAuditDraftFinding.AuditFindingCategoryID = auditDraftFindingFormDTO.AuditFindingCategoryID
	existingAuditDraftFinding.AuditChecklistID = auditDraftFindingFormDTO.AuditChecklistID
	existingAuditDraftFinding.SeverityLevelID = &auditDraftFindingFormDTO.SeverityLevelID
	existingAuditDraftFinding.Recommendations = auditDraftFindingFormDTO.Recommendations
	existingAuditDraftFinding.UpdatedBYID = &currentUserID

	err = adfs.repo.Update(existingAuditDraftFinding, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//documents update
	var documents []*models.AuditDraftFindingDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.AuditDraftFindingDocument

	// Iterate over each document in the array
	for _, doc := range auditDraftFindingFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("audit", "draft-audit-findings", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new AuditDraftFindingDocument instance
			document := &models.AuditDraftFindingDocument{
				AuditDraftFindingID: id,
				URL:                 url,
				Name:                doc.FileName,
				Size:                doc.Size,
				Token:               encryptedFilename,
				OriginalName:        doc.OriginalName,
				Type:                doc.Type,
				CreatedBYID:         currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := adfs.repo.GetDocumentByID(*doc.ID)
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
		if err := adfs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := adfs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "draft-audit-findings", strconv.FormatUint(uint64(existingAuditDraftFinding.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("audit", "draft-audit-findings", strconv.FormatUint(uint64(existingAuditDraftFinding.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		err := adfs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditDraftFinding{}.TableName(),
			id,
			auditDraftFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditDraftFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (adfs *AuditDraftFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

	tx := adfs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*auditDTOs.AuditDraftFindingDocumentDTO)

	for _, id := range ids {
		// Get existing audit finding details
		auditDraftFinding, err := adfs.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if len(auditDraftFinding.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range auditDraftFinding.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := adfs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the audit finding itself
		if err := adfs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		for auditDraftFindingID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(auditDraftFindingID), 10), "documents")
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
				err := adfs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditDraftFinding{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditDraftFinding{}.TableName(),
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

func (adfs *AuditDraftFindingService) validateAuditDraftFindingFormDTO(dto *auditDTOs.AuditDraftFindingFormDTO) error {

	if _, err := adfs.auditSvc.GetByID(dto.AuditID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit does not exist"), "not_exist", "audit")
	}
	if dto.AuditableItemID != nil {
		if _, err := adfs.auditableItemSvc.GetByID(*dto.AuditableItemID); err != nil {
			return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
		}
	}
	if dto.AuditChecklistID != nil {
		if _, err := adfs.auditChecklistSvc.GetByID(dto.AuditID, *dto.AuditChecklistID); err != nil {
			return utility.NewValidationError(fmt.Errorf("audit checklist does not exist"), "not_exist", "audit checklist")
		}
	}
	if _, err := adfs.auditFindingCategorySvc.GetByID(dto.AuditFindingCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit finding category does not exist"), "not_exist", "audit finding category")
	}

	if _, err := adfs.severityLevelSvc.GetByID(dto.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	return nil
}

func (adfs *AuditDraftFindingService) GetAuditDraftFindingCountByCategory(auditId uint) ([]auditDTOs.AuditDraftFindingCategoryCountResponseDTO, error) {
	if _, err := adfs.auditSvc.GetByID(auditId); err != nil {
		return nil, err
	}

	response, err := adfs.repo.GetAuditDraftFindingCountByCategory(auditId)
	if err != nil {
		return nil, err
	}

	return response, nil
}
