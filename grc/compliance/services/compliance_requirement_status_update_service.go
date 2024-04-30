package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type ComplianceRequirementStatusUpdateServiceInterface interface {
	GetAll(ComplianceID uint, ComplianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error)
	Create(ComplianceID uint, ComplianceRequirementID uint, currentUserID uint, complianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) (uint, error)
	Update(ComplianceID uint, ComplianceRequirementID uint, id uint, currentUserID uint, ComplianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) error
}

// Service and implementations
type ComplianceRequirementStatusUpdateService struct {
	repo                     repositories.ComplianceRequirementStatusUpdateRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	complianceRequirementSvc ComplianceRequirementServiceInterface
	complianceStatusSvc      ComplianceStatusServiceInterface
}

func NewComplianceRequirementStatusUpdateService(
	repo repositories.ComplianceRequirementStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRequirementSvc ComplianceRequirementServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceRequirementStatusUpdateServiceInterface {
	return &ComplianceRequirementStatusUpdateService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		complianceStatusSvc:      complianceStatusSvc,
	}
}

func (crsus *ComplianceRequirementStatusUpdateService) GetAll(ComplianceID uint, ComplianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error) {
	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return nil, err
	}

	return crsus.repo.GetAll(ComplianceRequirementID, sort, sortBy, keyword)
}

func (crsus ComplianceRequirementStatusUpdateService) GetByID(ComplianceID uint, ComplianceRequirementID uint, id uint) (*models.ComplianceRequirementStatusUpdate, error) {
	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return nil, err
	}

	complianceRequirementStatusUpdate, err := crsus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirementStatusUpdate, nil
}

func (crsus *ComplianceRequirementStatusUpdateService) Create(ComplianceID uint, ComplianceRequirementID uint, currentUserID uint, complianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) (uint, error) {
	tx := crsus.repo.BeginTransaction()

	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return 0, err
	}

	if err := crsus.validateComplianceRequirementStatusUpdateFormDTO(complianceRequirementStatusUpdateFormDTO); err != nil {
		return 0, err
	}
	if err := crsus.validateComplianceRequirementStatusUpdateDocumentsFormDTO(complianceRequirementStatusUpdateFormDTO.Documents); err != nil {
		return 0, err // Returning 0 and the error if validation fails
	}

	complianceStatusID := complianceRequirementStatusUpdateFormDTO.ComplianceStatusID
	complianceRequirementUpdate := &models.ComplianceRequirementStatusUpdate{
		ComplianceRequirementID: ComplianceRequirementID,
		ComplianceStatusID:      complianceStatusID,
		Comment:                 complianceRequirementStatusUpdateFormDTO.Comment,
		CreatedBy:               currentUserID,
	}

	id, err := crsus.repo.Create(complianceRequirementUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	err = crsus.complianceRequirementSvc.StatusUpdate(ComplianceID, ComplianceRequirementID, currentUserID, complianceStatusID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ComplianceRequirementStatusUpdateDocument

	// Iterate over each document in the array
	for _, doc := range complianceRequirementStatusUpdateFormDTO.Documents {
		storagePath := filepath.Join("compliance", "compliance", strconv.FormatUint(uint64(ComplianceID), 10), "requirements", strconv.FormatUint(uint64(ComplianceRequirementID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
			tx.Rollback()
			return 0, err
		}
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("compliances", strconv.FormatUint(uint64(ComplianceID), 10), "requirements", strconv.FormatUint(uint64(ComplianceRequirementID), 10), "status-updates", encryptedFilename)

		// Create a new ComplianceRequirementStatusUpdateDocument instance
		document := &models.ComplianceRequirementStatusUpdateDocument{
			ComplianceRequirementStatusUpdateID: id,
			URL:                                 url,
			Icon:                                doc.Icon,
			Name:                                doc.FileName,
			Size:                                doc.Size,
			Token:                               encryptedFilename,
			Type:                                doc.Type,
			CreatedBy:                           currentUserID,
		}
		// Append the document to the slice
		documents = append(documents, document)

	}
	if len(documents) > 0 {
		// Save the ComplianceRequirementStatusUpdateDocument
		if err := crsus.repo.SaveComplianceRequirementStatusUpdateDocument(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := crsus.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceRequirementStatusUpdate{}.TableName(),
			id,
			complianceRequirementUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.ComplianceRequirementStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crsus *ComplianceRequirementStatusUpdateService) Update(ComplianceID uint, ComplianceRequirementID uint, id uint, currentUserID uint, ComplianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) error {
	tx := crsus.repo.BeginTransaction()

	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return err
	}

	existingComplianceRequirementStatusUpdate, err := crsus.GetByID(ComplianceID, ComplianceRequirementID, id)

	if err != nil {
		return err
	}
	if err := crsus.validateComplianceRequirementStatusUpdateDocumentsFormDTO(ComplianceRequirementStatusUpdateFormDTO.Documents); err != nil {
		return err // Returning 0 and the error if validation fails
	}

	// Update Compliance StatusUpdate -- update comment only
	existingComplianceRequirementStatusUpdate.Comment = ComplianceRequirementStatusUpdateFormDTO.Comment
	existingComplianceRequirementStatusUpdate.UpdatedBy = &currentUserID

	err = crsus.repo.Update(existingComplianceRequirementStatusUpdate, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	var documents []*models.ComplianceRequirementStatusUpdateDocument
	var deleteDocumentIds []*uint
	// Iterate over each document in the array
	for _, doc := range ComplianceRequirementStatusUpdateFormDTO.Documents {
		storagePath := filepath.Join("compliance", "compliance", strconv.FormatUint(uint64(ComplianceID), 10), "requirements", strconv.FormatUint(uint64(ComplianceRequirementID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Handle new file case
			if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
				tx.Rollback()
				return err
			}
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("compliances", strconv.FormatUint(uint64(ComplianceID), 10), "requirements", strconv.FormatUint(uint64(ComplianceRequirementID), 10), "status-updates", encryptedFilename)
			// Create a new ComplianceRequirementStatusUpdateDocument instance
			document := &models.ComplianceRequirementStatusUpdateDocument{
				ComplianceRequirementStatusUpdateID: id,
				URL:                                 url,
				Icon:                                doc.Icon,
				Name:                                doc.FileName,
				Size:                                doc.Size,
				Token:                               encryptedFilename,
				Type:                                doc.Type,
				CreatedBy:                           currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := crsus.repo.GetComplianceRequirementStatusUpdateDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)

				if err := utility.RemoveFile(existingDocument.Name, storagePath); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	if len(documents) > 0 {
		if err := crsus.repo.SaveComplianceRequirementStatusUpdateDocument(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := crsus.repo.DeleteComplianceRequirementStatusUpdateDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	return nil
}

func (crsus *ComplianceRequirementStatusUpdateService) validateComplianceRequirementStatusUpdateFormDTO(complianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) error {

	if _, err := crsus.complianceStatusSvc.GetByID(complianceRequirementStatusUpdateFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}

func (crsus *ComplianceRequirementStatusUpdateService) validateComplianceRequirementStatusUpdateDocumentsFormDTO(documents []*complianceDTOs.ComplianceRequirementStatusUpdateDocumentFormDTO) error {
	// Check if documents array is nil or empty
	if len(documents) == 0 {
		return utility.NewValidationError(fmt.Errorf("documents is required"), "required", "documents")
	}
	for _, document := range documents {
		//@Review
		//All items are mandatory during creation

		// Validate icon
		if document.Icon == "" {
			return utility.NewValidationError(fmt.Errorf("icon is required"), "required", "icon")
		}

		// Validate file type (PDF)
		// if document.Type != ".pdf" {
		//     return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", "type")
		// }

		// @todo
		// Validate file size (if needed)
		// For example:
		// if document.Size > MaxFileSize {
		//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
		// }
	}

	return nil
}
