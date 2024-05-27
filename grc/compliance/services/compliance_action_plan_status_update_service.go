package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

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

	//@todo validate document form

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

	err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID, tx)
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
			OriginalName:                       doc.OriginalName,
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
				OriginalName:                       doc.OriginalName,
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

func (casus *ComplianceActionPlanStatusUpdateService) validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {

	if _, err := casus.actionPlanStatusSvc.GetByID(complianceActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}
