package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type IncidentActionPlanStatusUpdateServiceInterface interface {
	GetAll(IncidentID uint, IncidentActionPlanID uint, sort string, sortBy string, keyword string) ([]*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error)
	Create(IncidentID uint, IncidentActionPlanID uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) (uint, error)
	Update(IncidentID uint, IncidentActionPlanID uint, id uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) error
}

type IncidentActionPlanStatusUpdateService struct {
	repo                  repositories.IncidentActionPlanStatusUpdateRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	incidentActionPlanSvc IncidentActionPlanServiceInterface
	actionPlanStatusSvc   masterServices.ActionPlanStatusServiceInterface
}

func NewIncidentActionPlanStatusUpdateService(
	repo repositories.IncidentActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentActionPlanSvc IncidentActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) IncidentActionPlanStatusUpdateServiceInterface {
	return &IncidentActionPlanStatusUpdateService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		incidentActionPlanSvc: incidentActionPlanSvc,
		actionPlanStatusSvc:   actionPlanStatusSvc,
	}
}

func (iasus *IncidentActionPlanStatusUpdateService) GetAll(IncidentID uint, IncidentActionPlanID uint, sort string, sortBy string, keyword string) ([]*incidentDTOs.IncidentActionPlanStatusUpdateListDTO, error) {
	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return nil, err
	}

	return iasus.repo.GetAll(IncidentActionPlanID, sort, sortBy, keyword)
}

func (iasus IncidentActionPlanStatusUpdateService) GetByID(IncidentID uint, IncidentActionPlanID uint, id uint) (*models.IncidentActionPlanStatusUpdate, error) {
	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return nil, err
	}

	incidentActionPlanStatusUpdate, err := iasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incidentActionPlanStatusUpdate, nil
}

func (iasus *IncidentActionPlanStatusUpdateService) Create(IncidentID uint, IncidentActionPlanID uint, currentUserID uint, incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := iasus.repo.BeginTransaction()

	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return 0, err
	}

	if err := iasus.validateIncidentActionPlanStatusUpdateFormDTO(incidentActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	incidentActionPlanUpdate := &models.IncidentActionPlanStatusUpdate{
		IncidentActionPlanID: IncidentActionPlanID,
		ActionPlanStatusID:   incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:              incidentActionPlanStatusUpdateFormDTO.Comment,
		Percentage:           incidentActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:          currentUserID,
	}

	id, err := iasus.repo.Create(incidentActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = iasus.incidentActionPlanSvc.StatusUpdate(IncidentID, IncidentActionPlanID, currentUserID, incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IncidentActionPlanStatusUpdateDocument

	for _, doc := range incidentActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("incidents", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "status-updates", encryptedFilename)

		document := &models.IncidentActionPlanStatusUpdateDocument{
			IncidentActionPlanStatusUpdateID: id,
			URL:                              url,
			Name:                             doc.FileName,
			Size:                             doc.Size,
			Token:                            encryptedFilename,
			OriginalName:                     doc.OriginalName,
			Type:                             doc.Type,
			CreatedByID:                      currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := iasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		err := iasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentActionPlanStatusUpdate{}.TableName(),
			id,
			incidentActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (iasus *IncidentActionPlanStatusUpdateService) Update(IncidentID uint, IncidentActionPlanID uint, id uint, currentUserID uint, IncidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) error {
	tx := iasus.repo.BeginTransaction()

	if _, err := iasus.incidentActionPlanSvc.GetByID(IncidentID, IncidentActionPlanID); err != nil {
		return err
	}

	existingIncidentActionPlanStatusUpdate, err := iasus.GetByID(IncidentID, IncidentActionPlanID, id)

	if err != nil {
		return err
	}

	if err := iasus.validateIncidentActionPlanStatusUpdateFormDTO(IncidentActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Incident StatusUpdate-- only comments and percentage
	existingIncidentActionPlanStatusUpdate.Comment = IncidentActionPlanStatusUpdateFormDTO.Comment
	existingIncidentActionPlanStatusUpdate.Percentage = IncidentActionPlanStatusUpdateFormDTO.Percentage
	existingIncidentActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = iasus.repo.Update(existingIncidentActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.IncidentActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.IncidentActionPlanStatusUpdateDocument

	for _, doc := range IncidentActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("incidents", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "status-updates", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.IncidentActionPlanStatusUpdateDocument{
				IncidentActionPlanStatusUpdateID: id,
				URL:                              url,
				Name:                             doc.FileName,
				Size:                             doc.Size,
				Token:                            encryptedFileName,
				OriginalName:                     doc.OriginalName,
				Type:                             doc.Type,
				CreatedByID:                      currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := iasus.repo.GetDocumentByID(*doc.ID)
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
		if err := iasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := iasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(IncidentID), 10), "action-plans", strconv.FormatUint(uint64(IncidentActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (iasus *IncidentActionPlanStatusUpdateService) validateIncidentActionPlanStatusUpdateFormDTO(incidentActionPlanStatusUpdateFormDTO *incidentDTOs.IncidentActionPlanStatusUpdateFormDTO) error {

	if _, err := iasus.actionPlanStatusSvc.GetByID(incidentActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}
