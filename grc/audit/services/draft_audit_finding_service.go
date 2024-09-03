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

type DraftAuditFindingServiceInterface interface {
	GetAll(params *auditDTOs.DraftAuditFindingListQueryParamsDTO) (*auditDTOs.DraftAuditFindingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.DraftAuditFindingDetailDTO, error)
	GetByID(id uint) (*models.DraftAuditFinding, error)
	Create(currentUserID uint, draftAuditFindingFormDTO *auditDTOs.DraftAuditFindingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, draftAuditFindingFormDTO *auditDTOs.DraftAuditFindingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type DraftAuditFindingService struct {
	repo                       repositories.DraftAuditFindingRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	auditSvc                   AuditServiceInterface
	draftAuditFindingStatusSvc DraftAuditFindingStatusServiceInterface
	auditTimelineSvc           AuditTimelineServiceInterface
	auditableItemSvc           AuditableItemServiceInterface
	auditFindingCategorySvc    AuditFindingCategoryServiceInterface
	severityLevelSvc           masterServices.SeverityLevelServiceInterface
	auditChecklistSvc          AuditChecklistServiceInterface
}

func NewDraftAuditFindingService(
	repo repositories.DraftAuditFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditSvc AuditServiceInterface,
	draftAuditFindingStatusSvc DraftAuditFindingStatusServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	auditFindingCategorySvc AuditFindingCategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	auditChecklistSvc AuditChecklistServiceInterface,
) DraftAuditFindingServiceInterface {
	return &DraftAuditFindingService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		auditSvc:                   auditSvc,
		draftAuditFindingStatusSvc: draftAuditFindingStatusSvc,
		auditTimelineSvc:           auditTimelineSvc,
		auditableItemSvc:           auditableItemSvc,
		auditFindingCategorySvc:    auditFindingCategorySvc,
		severityLevelSvc:           severityLevelSvc,
		auditChecklistSvc:          auditChecklistSvc,
	}
}

func (dafs *DraftAuditFindingService) GetAll(params *auditDTOs.DraftAuditFindingListQueryParamsDTO) (*auditDTOs.DraftAuditFindingListResponseDTO, error) {

	return dafs.repo.GetAll(params)
}

func (dafs *DraftAuditFindingService) GetByID(id uint) (*models.DraftAuditFinding, error) {
	draftAuditFinding, err := dafs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return draftAuditFinding, nil
}

func (dafs *DraftAuditFindingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.DraftAuditFindingDetailDTO, error) {
	var draftAuditFindingDTO auditDTOs.DraftAuditFindingDetailDTO

	draftAuditFinding, err := dafs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(draftAuditFinding, &draftAuditFindingDTO)
	if err != nil {
		return nil, err
	}

	draftAuditFindingDTO.CreatedAt = draftAuditFinding.CreatedAt
	draftAuditFindingDTO.UpdatedAt = draftAuditFinding.UpdatedAt

	return &draftAuditFindingDTO, nil
}

func (dafs *DraftAuditFindingService) Create(currentUserID uint, draftAuditFindingFormDTO *auditDTOs.DraftAuditFindingFormDTO) (uint, error) {
	tx := dafs.repo.BeginTransaction()

	if err := dafs.validateDraftAuditFindingFormDTO(draftAuditFindingFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	draftAuditFindingStatus, err := dafs.draftAuditFindingStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	draftAuditFinding := &models.DraftAuditFinding{
		Title:                     draftAuditFindingFormDTO.Title,
		Description:               draftAuditFindingFormDTO.Description,
		AuditID:                   draftAuditFindingFormDTO.AuditID,
		AuditableItemID:           draftAuditFindingFormDTO.AuditableItemID,
		AuditChecklistID:          draftAuditFindingFormDTO.AuditChecklistID,
		AuditFindingCategoryID:    draftAuditFindingFormDTO.AuditFindingCategoryID,
		SeverityLevelID:           &draftAuditFindingFormDTO.SeverityLevelID,
		Recommendations:           draftAuditFindingFormDTO.Recommendations,
		DraftAuditFindingStatusID: draftAuditFindingStatus.ID,
		CreatedBYID:               currentUserID,
	}

	id, err := dafs.repo.Create(draftAuditFinding, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.DraftAuditFindingDocument

	// Iterate over each document in the array
	for _, doc := range draftAuditFindingFormDTO.Documents {

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
		document := &models.DraftAuditFindingDocument{
			DraftAuditFindingID: id,
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
		if err := dafs.repo.CreateDocuments(documents, tx); err != nil {
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
		err := dafs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.DraftAuditFinding{}.TableName(),
			id,
			draftAuditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.DraftAuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (dafs *DraftAuditFindingService) Update(id uint, currentUserID uint, draftAuditFindingFormDTO *auditDTOs.DraftAuditFindingFormDTO) error {
	tx := dafs.repo.BeginTransaction()

	existingDraftAuditFinding, err := dafs.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := dafs.validateDraftAuditFindingFormDTO(draftAuditFindingFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingDraftAuditFinding.AuditableItemID = draftAuditFindingFormDTO.AuditableItemID
	existingDraftAuditFinding.Title = draftAuditFindingFormDTO.Title
	existingDraftAuditFinding.Description = draftAuditFindingFormDTO.Description
	existingDraftAuditFinding.AuditFindingCategoryID = draftAuditFindingFormDTO.AuditFindingCategoryID
	existingDraftAuditFinding.AuditChecklistID = draftAuditFindingFormDTO.AuditChecklistID
	existingDraftAuditFinding.SeverityLevelID = &draftAuditFindingFormDTO.SeverityLevelID
	existingDraftAuditFinding.Recommendations = draftAuditFindingFormDTO.Recommendations
	existingDraftAuditFinding.UpdatedBYID = &currentUserID

	err = dafs.repo.Update(existingDraftAuditFinding, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//documents update
	var documents []*models.DraftAuditFindingDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.DraftAuditFindingDocument

	// Iterate over each document in the array
	for _, doc := range draftAuditFindingFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("audit", "draft-audit-findings", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new DraftAuditFindingDocument instance
			document := &models.DraftAuditFindingDocument{
				DraftAuditFindingID: id,
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
				existingDocument, err := dafs.repo.GetDocumentByID(*doc.ID)
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
		if err := dafs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := dafs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("audit", "draft-audit-findings", strconv.FormatUint(uint64(existingDraftAuditFinding.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("audit", "draft-audit-findings", strconv.FormatUint(uint64(existingDraftAuditFinding.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		err := dafs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.DraftAuditFinding{}.TableName(),
			id,
			draftAuditFindingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.DraftAuditFinding{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (dafs *DraftAuditFindingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

	tx := dafs.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	documentsToDelete := make(map[uint][]*auditDTOs.DraftAuditFindingDocumentDTO)

	for _, id := range ids {
		// Get existing audit finding details
		draftAuditFinding, err := dafs.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if len(draftAuditFinding.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range draftAuditFinding.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := dafs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the audit finding itself
		if err := dafs.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		for draftAuditFindingID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("audit", "audit-findings", strconv.FormatUint(uint64(draftAuditFindingID), 10), "documents")
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
				err := dafs.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.DraftAuditFinding{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)

				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.DraftAuditFinding{}.TableName(),
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

func (dafs *DraftAuditFindingService) validateDraftAuditFindingFormDTO(dto *auditDTOs.DraftAuditFindingFormDTO) error {

	if _, err := dafs.auditSvc.GetByID(dto.AuditID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit does not exist"), "not_exist", "audit")
	}

	if _, err := dafs.auditableItemSvc.GetByID(*dto.AuditableItemID); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
	}

	if _, err := dafs.auditChecklistSvc.GetByID(dto.AuditID, *dto.AuditChecklistID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit checklist does not exist"), "not_exist", "audit checklist")
	}

	if _, err := dafs.auditFindingCategorySvc.GetByID(dto.AuditFindingCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit finding category does not exist"), "not_exist", "audit finding category")
	}

	if _, err := dafs.severityLevelSvc.GetByID(dto.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	return nil
}
