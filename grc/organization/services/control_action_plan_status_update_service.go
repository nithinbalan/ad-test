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
	GetByID(controlID uint, controlActionPlanID uint, id uint) (*models.ControlActionPlanStatusUpdate, error)
	GetAll(controlID uint, controlActionPlanID uint, sort string, sortBy string, keyword string) ([]*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error)
	Create(controlID uint, controlActionPlanID uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) (uint, error)
	Update(controlID uint, controlActionPlanID uint, id uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error
	GetLatestControlActionPlanStatusUpdate(controlID uint, actionPlanID uint) (*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error)
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

	err = casus.controlActionPlanSvc.StatusUpdate(controlID, controlActionPlanID, currentUserID, controlActionPlanStatusUpdateFormDTO, controlActionPlanStatusUpdateFormDTO.Percentage, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ControlActionPlanStatusUpdateDocument

	for _, doc := range controlActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), encryptedFileName)

		document := &models.ControlActionPlanStatusUpdateDocument{
			ControlActionPlanStatusUpdateID: id,
			URL:                             url,
			Name:                            doc.FileName,
			Size:                            doc.Size,
			Token:                           encryptedFileName,
			OriginalName:                    doc.OriginalName,
			Type:                            doc.Type,
			CreatedByID:                     currentUserID,
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
			storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")
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
			models.ControlActionPlanStatusUpdate{}.TableName(),
			id,
			controlActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (casus *ControlActionPlanStatusUpdateService) Update(controlID uint, controlActionPlanID uint, id uint, currentUserID uint, controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.controlActionPlanSvc.GetByID(controlID, controlActionPlanID); err != nil {
		return err
	}

	existingControlActionPlanStatusUpdate, err := casus.GetByID(controlID, controlActionPlanID, id)

	if err != nil {
		return err
	}

	if err := casus.validateControlActionPlanStatusUpdateFormDTO(controlActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingControlActionPlanStatusUpdate.Comment = controlActionPlanStatusUpdateFormDTO.Comment
	existingControlActionPlanStatusUpdate.Percentage = controlActionPlanStatusUpdateFormDTO.Percentage
	existingControlActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = casus.repo.Update(existingControlActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ControlActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ControlActionPlanStatusUpdateDocument

	for _, doc := range controlActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), encryptedFileName)

			document := &models.ControlActionPlanStatusUpdateDocument{
				ControlActionPlanStatusUpdateID: id,
				URL:                             url,
				Name:                            doc.FileName,
				Size:                            doc.Size,
				Token:                           encryptedFileName,
				OriginalName:                    doc.OriginalName,
				Type:                            doc.Type,
				CreatedByID:                     currentUserID,
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
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "action-plans", strconv.FormatUint(uint64(controlActionPlanID), 10), "action-plan-updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (casus *ControlActionPlanStatusUpdateService) validateControlActionPlanStatusUpdateFormDTO(controlActionPlanStatusUpdateFormDTO *organizationDTOs.ControlActionPlanStatusUpdateFormDTO) error {

	if _, err := casus.actionPlanStatusSvc.GetByID(controlActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (casus ControlActionPlanStatusUpdateService) GetLatestControlActionPlanStatusUpdate(controlID uint, actionPlanID uint) (*organizationDTOs.ControlActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.controlActionPlanSvc.GetByID(controlID, actionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetLatestControlActionPlanStatusUpdate(actionPlanID)
}
