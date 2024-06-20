package services

import (
	"errors"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	controlDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"path/filepath"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ControlReviewServiceInterface interface {
	GetAll(controlID uint, sort string, sortBy string, keyword string) ([]*controlDTOs.ControlReviewListDTO, error)
	Create(controlID uint, currentUserID uint, controlReviewFormDTO *controlDTOs.ControlReviewFormDTO) (uint, error)
	Update(controlID uint, id uint, currentUserID uint, ControlReviewFormDTO *controlDTOs.ControlReviewFormDTO) error
}

// Service and implementations
type ControlReviewService struct {
	repo                   repositories.ControlReviewRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	controlSvc             ControlServiceInterface
	reviewFrequencySvc     masterServices.ReviewFrequencyServiceInterface
	controlReviewStatusSvc ControlReviewStatusServiceInterface
}

func NewControlReviewService(
	repo repositories.ControlReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlSvc ControlServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	controlReviewStatusSvc ControlReviewStatusServiceInterface,
) ControlReviewServiceInterface {
	return &ControlReviewService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		controlSvc:             controlSvc,
		reviewFrequencySvc:     reviewFrequencySvc,
		controlReviewStatusSvc: controlReviewStatusSvc,
	}
}

func (crs *ControlReviewService) GetAll(controlID uint, sort string, sortBy string, keyword string) ([]*controlDTOs.ControlReviewListDTO, error) {
	if _, err := crs.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	return crs.repo.GetAll(controlID, sort, sortBy, keyword)
}

func (crs ControlReviewService) GetByID(controlID uint, id uint) (*models.ControlReview, error) {
	if _, err := crs.controlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	controlReview, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return controlReview, nil
}

func (crs *ControlReviewService) Create(controlID uint, currentUserID uint, controlReviewFormDTO *controlDTOs.ControlReviewFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	control, err := crs.controlSvc.GetByID(controlID)
	if err != nil {
		return 0, err
	}

	controlReview := &models.ControlReview{
		ControlID:  controlID,
		Comment:    controlReviewFormDTO.Comment,
		ReviewedBy: currentUserID,
		ReviewedAt: time.Now(),
	}

	id, err := crs.repo.Create(controlReview, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if control.ReviewFrequencyID != nil {
		reviewFrequency, err := crs.reviewFrequencySvc.GetByID(*control.ReviewFrequencyID)
		if err != nil {
			return 0, err
		}
		nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)
		control.NextReviewDate = &nextReviewDate
	}

	controlReviewStatus, err := crs.controlReviewStatusSvc.GetByType("reviewed")
	if err != nil {
		return 0, err
	}

	currentTime := time.Now()

	control.LastReviewDate = &currentTime
	control.ControlReviewStatusID = &controlReviewStatus.ID

	if err := crs.controlSvc.ReviewStatusUpdate(control, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ControlReviewDocument

	for _, doc := range controlReviewFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "review-documents")

		document := &models.ControlReviewDocument{
			ControlReviewID: id,
			URL:             url,
			Name:            doc.FileName,
			Size:            doc.Size,
			Token:           encryptedFilename,
			OriginalName:    doc.OriginalName,
			Type:            doc.Type,
			CreatedByID:     currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := crs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "review-documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}

		}
	}
	tx.Commit()

	go func() {
		err := crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlReview{}.TableName(),
			id,
			controlReview,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlReview{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ControlReviewService) Update(controlID uint, id uint, currentUserID uint, controlReviewFormDTO *controlDTOs.ControlReviewFormDTO) error {
	tx := crs.repo.BeginTransaction()

	if _, err := crs.controlSvc.GetByID(controlID); err != nil {
		return err
	}

	existingControlReview, err := crs.GetByID(controlID, id)

	if err != nil {
		return err
	}

	existingControlReview.Comment = controlReviewFormDTO.Comment
	existingControlReview.UpdatedByID = &currentUserID

	err = crs.repo.Update(existingControlReview, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ControlReviewDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ControlReviewDocument

	for _, doc := range controlReviewFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "review-documents")
			// Create a new Control test Document instance
			document := &models.ControlReviewDocument{
				ControlReviewID: id,
				URL:             url,
				Name:            doc.FileName,
				Size:            doc.Size,
				Token:           encryptedFileName,
				OriginalName:    doc.OriginalName,
				Type:            doc.Type,
				CreatedByID:     currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := crs.repo.GetDocumentByID(*doc.ID)
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
		if err := crs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := crs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "review-documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("control", "control-register", strconv.FormatUint(uint64(controlID), 10), "review-documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}
