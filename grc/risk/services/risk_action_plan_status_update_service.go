package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type RiskActionPlanStatusUpdateServiceInterface interface {
	GetAll(RiskID uint, RiskActionPlanID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.RiskActionPlanStatusUpdateListDTO, error)
	Create(RiskID uint, RiskActionPlanID uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) (uint, error)
	Update(RiskID uint, RiskActionPlanID uint, id uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) error
}

type RiskActionPlanStatusUpdateService struct {
	repo                repositories.RiskActionPlanStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	riskActionPlanSvc   RiskActionPlanServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
}

func NewRiskActionPlanStatusUpdateService(
	repo repositories.RiskActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskActionPlanSvc RiskActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) RiskActionPlanStatusUpdateServiceInterface {
	return &RiskActionPlanStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		riskActionPlanSvc:   riskActionPlanSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
	}
}

func (rasus *RiskActionPlanStatusUpdateService) GetAll(RiskID uint, RiskActionPlanID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.RiskActionPlanStatusUpdateListDTO, error) {
	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return nil, err
	}

	return rasus.repo.GetAll(RiskActionPlanID, sort, sortBy, keyword)
}

func (rasus RiskActionPlanStatusUpdateService) GetByID(RiskID uint, RiskActionPlanID uint, id uint) (*models.RiskActionPlanStatusUpdate, error) {
	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return nil, err
	}

	riskActionPlanStatusUpdate, err := rasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskActionPlanStatusUpdate, nil
}

func (rasus *RiskActionPlanStatusUpdateService) Create(RiskID uint, RiskActionPlanID uint, currentUserID uint, riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := rasus.repo.BeginTransaction()

	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return 0, err
	}

	if err := rasus.validateRiskActionPlanStatusUpdateFormDTO(riskActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	riskActionPlanUpdate := &models.RiskActionPlanStatusUpdate{
		RiskActionPlanID:   RiskActionPlanID,
		ActionPlanStatusID: riskActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:            riskActionPlanStatusUpdateFormDTO.Comment,
		Percentage:         riskActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err := rasus.repo.Create(riskActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = rasus.riskActionPlanSvc.StatusUpdate(RiskID, RiskActionPlanID, currentUserID, riskActionPlanStatusUpdateFormDTO.ActionPlanStatusID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskActionPlanStatusUpdateDocument

	for _, doc := range riskActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("risks", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "status-updates", encryptedFilename)

		document := &models.RiskActionPlanStatusUpdateDocument{
			RiskActionPlanStatusUpdateID: id,
			URL:                          url,
			Name:                         doc.FileName,
			Size:                         doc.Size,
			Token:                        encryptedFilename,
			Type:                         doc.Type,
			CreatedByID:                  currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := rasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		err := rasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskActionPlanStatusUpdate{}.TableName(),
			id,
			riskActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (rasus *RiskActionPlanStatusUpdateService) Update(RiskID uint, RiskActionPlanID uint, id uint, currentUserID uint, RiskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) error {
	tx := rasus.repo.BeginTransaction()

	if _, err := rasus.riskActionPlanSvc.GetByID(RiskID, RiskActionPlanID); err != nil {
		return err
	}

	existingRiskActionPlanStatusUpdate, err := rasus.GetByID(RiskID, RiskActionPlanID, id)

	if err != nil {
		return err
	}

	if err := rasus.validateRiskActionPlanStatusUpdateFormDTO(RiskActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Risk StatusUpdate-- only comments and percentage
	existingRiskActionPlanStatusUpdate.Comment = RiskActionPlanStatusUpdateFormDTO.Comment
	existingRiskActionPlanStatusUpdate.Percentage = RiskActionPlanStatusUpdateFormDTO.Percentage
	existingRiskActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = rasus.repo.Update(existingRiskActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.RiskActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskActionPlanStatusUpdateDocument

	for _, doc := range RiskActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("risks", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "status-updates", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.RiskActionPlanStatusUpdateDocument{
				RiskActionPlanStatusUpdateID: id,
				URL:                          url,
				Name:                         doc.FileName,
				Size:                         doc.Size,
				Token:                        encryptedFileName,
				Type:                         doc.Type,
				CreatedByID:                  currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := rasus.repo.GetDocumentByID(*doc.ID)
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
		if err := rasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(RiskID), 10), "action-plans", strconv.FormatUint(uint64(RiskActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (rasus *RiskActionPlanStatusUpdateService) validateRiskActionPlanStatusUpdateFormDTO(riskActionPlanStatusUpdateFormDTO *riskDTOs.RiskActionPlanStatusUpdateFormDTO) error {

	if _, err := rasus.actionPlanStatusSvc.GetByID(riskActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}
