package services

import (
	"errors"
	"fmt"
	"jethur/config"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// @todo - @document support to be added
// @todo - transaction to be added
type ComplianceActionPlanStatusUpdateServiceInterface interface {
	GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error
}

type ComplianceActionPlanStatusUpdateService struct {
	repo                    repositories.ComplianceActionPlanStatusUpdateRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceActionPlanSvc ComplianceActionPlanServiceInterface
	actionPlanStatusSvc     masterServices.ActionPlanStatusServiceInterface
}

func NewComplianceActionPlanStatusUpdateService(
	repo repositories.ComplianceActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceActionPlanSvc ComplianceActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) ComplianceActionPlanStatusUpdateServiceInterface {
	return &ComplianceActionPlanStatusUpdateService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
		actionPlanStatusSvc:     actionPlanStatusSvc,
	}
}

func (casus *ComplianceActionPlanStatusUpdateService) GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(ComplianceActionPlanID, sort, sortBy, keyword)
}

func (casus ComplianceActionPlanStatusUpdateService) GetByID(ComplianceID uint, ComplianceActionPlanID uint, id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	complianceActionPlanStatusUpdate, err := casus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceActionPlanStatusUpdate, nil
}

func (casus *ComplianceActionPlanStatusUpdateService) Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return 0, err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateDocumentsFormDTO(complianceActionPlanStatusUpdateFormDTO.Documents); err != nil {
		return 0, err
	}

	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: ComplianceActionPlanID,
		ActionPlanStatusID:     complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                complianceActionPlanStatusUpdateFormDTO.Comment,
		Percentage:             complianceActionPlanStatusUpdateFormDTO.Percentage,
		CreatedBy:              currentUserID,
	}

	id, err := casus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID)
	if err != nil {
		return 0, err
	}
	// Iterate over each document in the array
	for _, doc := range complianceActionPlanStatusUpdateFormDTO.Documents {

		if err := casus.repo.StoreComplianceActionPlanStatusUpdateDocumentFile(ComplianceID, ComplianceActionPlanID, id, doc.FileName); err != nil {
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
		url := fmt.Sprintf("%s/compliances/%d/action-plans/%d/status-updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), ComplianceID, ComplianceActionPlanID, encryptedFilename)

		// Create a new ComplianceActionPlanStatusUpdateDocument instance
		document := &models.ComplianceActionPlanStatusUpdateDocument{
			ComplianceActionPlanStatusUpdateID: id,
			URL:                                url,
			Icon:                               doc.Icon,
			Name:                               doc.FileName,
			Size:                               doc.Size,
			Token:                              encryptedFilename,
			Type:                               doc.Type,
			CreatedBy:                          currentUserID,
		}

		// Save the ComplianceActionPlanStatusUpdateDocument
		if err := casus.repo.SaveComplianceActionPlanStatusUpdateDocument(document, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := casus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlanStatusUpdate{}.TableName(),
			id,
			complianceActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (casus *ComplianceActionPlanStatusUpdateService) Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, ComplianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return err
	}

	existingComplianceActionPlanStatusUpdate, err := casus.GetByID(ComplianceID, ComplianceActionPlanID, id)

	if err != nil {
		return err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateFormDTO(ComplianceActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateDocumentsFormDTO(ComplianceActionPlanStatusUpdateFormDTO.Documents); err != nil {
		return err
	}

	// Update Compliance StatusUpdate-- only comments and percentage
	existingComplianceActionPlanStatusUpdate.Comment = ComplianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlanStatusUpdate.Percentage = ComplianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlanStatusUpdate.UpdatedBy = &currentUserID

	err = casus.repo.Update(existingComplianceActionPlanStatusUpdate, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	// @review
	//update main action plan staus only if latest status updates is updated
	err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, ComplianceActionPlanStatusUpdateFormDTO.ComplianceStatusID)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Iterate over each document in the array
	for _, doc := range ComplianceActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Handle new file case
			if err := casus.repo.StoreComplianceActionPlanStatusUpdateDocumentFile(ComplianceID, ComplianceActionPlanID, id, doc.FileName); err != nil {
				tx.Rollback()
				return err
			}
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := fmt.Sprintf("%s/compliances/%d/action-plans/%d/status-updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), ComplianceID, ComplianceActionPlanID, encryptedFilename)
			// Create a new ComplianceActionPlanStatusUpdateDocument instance
			document := &models.ComplianceActionPlanStatusUpdateDocument{
				ComplianceActionPlanStatusUpdateID: id,
				URL:                                url,
				Icon:                               doc.Icon,
				Name:                               doc.FileName,
				Size:                               doc.Size,
				Token:                              encryptedFilename,
				Type:                               doc.Type,
				CreatedBy:                          currentUserID,
			}
			if err := casus.repo.SaveComplianceActionPlanStatusUpdateDocument(document, tx); err != nil {
				tx.Rollback()
				return err
			}
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := casus.repo.GetComplianceActionPlanStatusUpdateDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				if err := casus.repo.RemoveComplianceActionPlanStatusUpdateDocumentFile(ComplianceID, ComplianceActionPlanID, id, existingDocument.Name); err != nil {
					tx.Rollback()
					return err
				}
				if err := casus.repo.DeleteComplianceActionPlanStatusUpdateDocumentId(existingDocument.ID, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		} else {
			// Handle existing document update case
			if doc.ID != nil {
				existingDocument, err := casus.repo.GetComplianceActionPlanStatusUpdateDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				existingDocument.UpdatedBy = currentUserID
				if err := casus.repo.SaveComplianceActionPlanStatusUpdateDocument(existingDocument, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}
	tx.Commit()

	return nil
}

func (casus *ComplianceActionPlanStatusUpdateService) validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {

	if _, err := casus.actionPlanStatusSvc.GetByID(complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (casus *ComplianceActionPlanStatusUpdateService) validateComplianceActionPlanStatusUpdateDocumentsFormDTO(documents []*complianceDTOs.ComplianceActionPlanStatusUpdateDocumentFormDTO) error {
	// Check if documents array is empty
	if len(documents) == 0 {
		return utility.NewValidationError(fmt.Errorf("documents array is required"), "required", "documents")
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
