package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type ComplianceStatusUpdateServiceInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error)
	Create(complianceID uint, currentUserID uint, complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, ComplianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error
}

// Service and implementations
type ComplianceStatusUpdateService struct {
	repo                  repositories.ComplianceStatusUpdateRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
	complianceStatusSvc   ComplianceStatusServiceInterface
}

func NewComplianceStatusUpdateService(
	repo repositories.ComplianceStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceStatusUpdateServiceInterface {
	return &ComplianceStatusUpdateService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		complianceRegisterSvc: complianceRegisterSvc,
		complianceStatusSvc:   complianceStatusSvc,
	}
}

func (csus *ComplianceStatusUpdateService) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error) {
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	return csus.repo.GetAll(complianceID, sort, sortBy, keyword)
}

func (csus ComplianceStatusUpdateService) GetByID(complianceID uint, id uint) (*models.ComplianceStatusUpdate, error) {
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceStatusUpdate, err := csus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceStatusUpdate, nil
}

func (csus *ComplianceStatusUpdateService) Create(complianceID uint, currentUserID uint, complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) (uint, error) {
	tx := csus.repo.BeginTransaction()
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return 0, err
	}

	if err := csus.validateComplianceStatusUpdateFormDTO(complianceStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	complianceUpdate := &models.ComplianceStatusUpdate{
		ComplianceRegisterID: complianceID,
		ComplianceStatusID:   complianceStatusUpdateFormDTO.ComplianceStatusID,
		Comment:              complianceStatusUpdateFormDTO.Comment,
		CreatedBy:            currentUserID,
	}

	id, err := csus.repo.Create(complianceUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = csus.complianceRegisterSvc.StatusUpdate(complianceID, currentUserID, complianceStatusUpdateFormDTO.ComplianceStatusID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var documents []*models.ComplianceStatusUpdateDocument

	for _, doc := range complianceStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("compliances", strconv.FormatUint(uint64(complianceID), 10), "status-updates", encryptedFilename)

		document := &models.ComplianceStatusUpdateDocument{
			ComplianceStatusUpdateID: id,
			URL:                      url,
			Name:                     doc.FileName,
			Size:                     doc.Size,
			Token:                    encryptedFilename,
			OriginalName:             doc.OriginalName,
			Type:                     doc.Type,
			CreatedBy:                currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := csus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "status-update-documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}

	tx.Commit()

	go func() {
		err := csus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceStatusUpdate{}.TableName(),
			id,
			complianceUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (csus *ComplianceStatusUpdateService) Update(complianceID uint, id uint, currentUserID uint, complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error {
	tx := csus.repo.BeginTransaction()
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceStatusUpdate, err := csus.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	// Update Compliance StatusUpdate -- update only comments
	existingComplianceStatusUpdate.Comment = complianceStatusUpdateFormDTO.Comment
	existingComplianceStatusUpdate.UpdatedBy = &currentUserID

	err = csus.repo.Update(existingComplianceStatusUpdate)

	if err != nil {
		return err
	}

	var documents []*models.ComplianceStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ComplianceStatusUpdateDocument

	for _, doc := range complianceStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("compliances", strconv.FormatUint(uint64(complianceID), 10), "status-updates", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.ComplianceStatusUpdateDocument{
				ComplianceStatusUpdateID: id,
				URL:                      url,
				Name:                     doc.FileName,
				Size:                     doc.Size,
				Token:                    encryptedFileName,
				OriginalName:             doc.OriginalName,
				Type:                     doc.Type,
				CreatedBy:                currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := csus.repo.GetDocumentByID(*doc.ID)
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
		if err := csus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := csus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "status-update-documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "status-update-documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()


	return nil
}

func (csus *ComplianceStatusUpdateService) validateComplianceStatusUpdateFormDTO(complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error {

	if _, err := csus.complianceStatusSvc.GetByID(complianceStatusUpdateFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}
