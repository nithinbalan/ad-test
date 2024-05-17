package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type ControlActionPlanStatusUpdateServiceInterface interface {
	GetAll(controlID uint, controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error)
	Create(controlID uint, controlActionPlanID uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) (uint, error)
	Update(controlID uint, controlActionPlanID uint, id uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error
}

type ControlActionPlanStatusUpdateService struct {
	repo                 repositories.ControlActionPlanStatusUpdateRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	controlActionPlanSvc ControlActionPlanServiceInterface
	actionPlanStatusSvc  masterServices.ActionPlanStatusServiceInterface
}

func NewControlActionPlanStatusUpdateService(
	repo repositories.ControlActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlActionPlanSvc ControlActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) ControlActionPlanStatusUpdateServiceInterface {
	return &ControlActionPlanStatusUpdateService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		controlActionPlanSvc: controlActionPlanSvc,
		actionPlanStatusSvc:  actionPlanStatusSvc,
	}
}

func (casus *ControlActionPlanStatusUpdateService) GetAll(controlID uint, controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(controlActionPlanID, sort, sortBy, keyword)
}

func (casus ControlActionPlanStatusUpdateService) GetByID(controlID uint, controlActionPlanID uint, id uint) (*models.ControlActionPlanStatusUpdate, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return nil, err
	}

	controlActionPlanStatusUpdate, err := casus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return controlActionPlanStatusUpdate, nil
}

func (casus *ControlActionPlanStatusUpdateService) Create(controlID uint, controlActionPlanID uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return 0, err
	}

	if err := casus.validateControlActionPlanStatusUpdateFormDTO(controlActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	controlActionPlanUpdate := &models.ControlActionPlanStatusUpdate{
		ControlActionPlanID: controlActionPlanID,
		ActionPlanStatusID:  controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:             controlActionPlanStatusUpdateFormDTO.Comment,
		Percentage:          controlActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:         currentUserID,
	}

	id, err := casus.repo.Create(controlActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = casus.controlActionPlanSvc.StatusUpdate(controlID, controlActionPlanID, currentUserID, controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID, controlActionPlanStatusUpdateFormDTO.Percentage, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ComplianceActionPlanStatusUpdateDocument

	for _, doc := range complianceActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "status-updates", encryptedFilename)

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

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := casus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

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

func (casus *ControlActionPlanStatusUpdateService) Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, ComplianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {
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

	//@todo validate docuemnt form

	// Update Compliance StatusUpdate-- only comments and percentage
	existingComplianceActionPlanStatusUpdate.Comment = ComplianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlanStatusUpdate.Percentage = ComplianceActionPlanStatusUpdateFormDTO.Percentage
	existingComplianceActionPlanStatusUpdate.UpdatedBy = &currentUserID

	err = casus.repo.Update(existingComplianceActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ComplianceActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ComplianceActionPlanStatusUpdateDocument

	for _, doc := range ComplianceActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "status-updates", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.ComplianceActionPlanStatusUpdateDocument{
				ComplianceActionPlanStatusUpdateID: id,
				URL:                                url,
				Icon:                               doc.Icon,
				Name:                               doc.FileName,
				Size:                               doc.Size,
				Token:                              encryptedFileName,
				Type:                               doc.Type,
				CreatedBy:                          currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := casus.repo.GetDocumentByID(*doc.ID)
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
		if err := casus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := casus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(ComplianceID), 10), "action-plans", strconv.FormatUint(uint64(ComplianceActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (casus *ControlActionPlanStatusUpdateService) validateControlActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {

	if _, err := casus.actionPlanStatusSvc.GetByID(complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}
