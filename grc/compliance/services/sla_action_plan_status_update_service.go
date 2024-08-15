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

type SLAActionPlanStatusUpdateServiceInterface interface {
	GetAll(SLAID uint, SLAActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error)
	Create(SLAID uint, SLAActionPlanID uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) (uint, error)
	Update(SLAID uint, SLAActionPlanID uint, id uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) error
	GetLatestSLAActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error)
}

type SLAActionPlanStatusUpdateService struct {
	repo                repositories.SLAActionPlanStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	slaActionPlanSvc    SLAActionPlanServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
}

func NewSLAActionPlanStatusUpdateService(
	repo repositories.SLAActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	slaActionPlanSvc SLAActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) SLAActionPlanStatusUpdateServiceInterface {
	return &SLAActionPlanStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		slaActionPlanSvc:    slaActionPlanSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
	}
}

func (sasus *SLAActionPlanStatusUpdateService) GetAll(SLAID uint, SLAActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetAll(SLAActionPlanID, sort, sortBy, keyword)
}

func (sasus SLAActionPlanStatusUpdateService) GetByID(SLAID uint, SLAActionPlanID uint, id uint) (*models.SLAActionPlanStatusUpdate, error) {
	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return nil, err
	}

	SLAActionPlanStatusUpdate, err := sasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return SLAActionPlanStatusUpdate, nil
}

func (sasus *SLAActionPlanStatusUpdateService) Create(SLAID uint, SLAActionPlanID uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return 0, err
	}

	if err := sasus.validateSLAActionPlanStatusUpdateFormDTO(SLAActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	SLAActionPlanUpdate := &models.SLAActionPlanStatusUpdate{
		SLAActionPlanID:    SLAActionPlanID,
		ActionPlanStatusID: SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:            SLAActionPlanStatusUpdateFormDTO.Comment,
		Percentage:         SLAActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err := sasus.repo.Create(SLAActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = sasus.slaActionPlanSvc.StatusUpdate(SLAID, SLAActionPlanID, currentUserID, SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.SLAActionPlanStatusUpdateDocument

	for _, doc := range SLAActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "status-updates", encryptedFilename)

		document := &models.SLAActionPlanStatusUpdateDocument{
			SLAActionPlanStatusUpdateID: id,
			URL:                         url,
			Name:                        doc.FileName,
			Size:                        doc.Size,
			Token:                       encryptedFilename,
			OriginalName:                doc.OriginalName,
			Type:                        doc.Type,
			CreatedByID:                 currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := sasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		err := sasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLAActionPlanStatusUpdate{}.TableName(),
			id,
			SLAActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLAActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (sasus *SLAActionPlanStatusUpdateService) Update(SLAID uint, SLAActionPlanID uint, id uint, currentUserID uint, SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) error {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.slaActionPlanSvc.GetByID(SLAID, SLAActionPlanID); err != nil {
		return err
	}

	existingSLAActionPlanStatusUpdate, err := sasus.GetByID(SLAID, SLAActionPlanID, id)

	if err != nil {
		return err
	}

	if err := sasus.validateSLAActionPlanStatusUpdateFormDTO(SLAActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	// Update SLA StatusUpdate-- only comments and percentage
	existingSLAActionPlanStatusUpdate.Comment = SLAActionPlanStatusUpdateFormDTO.Comment
	existingSLAActionPlanStatusUpdate.Percentage = SLAActionPlanStatusUpdateFormDTO.Percentage
	existingSLAActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = sasus.repo.Update(existingSLAActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.SLAActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.SLAActionPlanStatusUpdateDocument

	for _, doc := range SLAActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "status-updates", encryptedFileName)
			// Create a new Control test Document instance
			document := &models.SLAActionPlanStatusUpdateDocument{
				SLAActionPlanStatusUpdateID: id,
				URL:                         url,
				Name:                        doc.FileName,
				Size:                        doc.Size,
				Token:                       encryptedFileName,
				OriginalName:                doc.OriginalName,
				Type:                        doc.Type,
				CreatedByID:                 currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := sasus.repo.GetDocumentByID(*doc.ID)
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
		if err := sasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := sasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(SLAID), 10), "action-plans", strconv.FormatUint(uint64(SLAActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (sasus *SLAActionPlanStatusUpdateService) validateSLAActionPlanStatusUpdateFormDTO(SLAActionPlanStatusUpdateFormDTO *complianceDTOs.SLAActionPlanStatusUpdateFormDTO) error {

	if _, err := sasus.actionPlanStatusSvc.GetByID(SLAActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (sasus SLAActionPlanStatusUpdateService) GetLatestSLAActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*complianceDTOs.SLAActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.slaActionPlanSvc.GetByID(slaID, actionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetLatestSLAActionPlanStatusUpdate(actionPlanID)
}
