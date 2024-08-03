package services

import (
	"errors"
	"fmt"
	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KRIValueUpdateServiceInterface interface {
	GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.KRIValueUpdateDetailsDTO, error)
	Create(kriRegisterID, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) (uint, error)
	Update(kriRegisterID, updateId, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) error
	Delete(kriRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type KRIValueUpdateService struct {
	repo               riskRepositories.KRIValueUpdateRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	kriRegisterSvc     KRIRegisterServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
}

func NewKRIValueUpdateService(
	repo riskRepositories.KRIValueUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriRegisterSvc KRIRegisterServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
) KRIValueUpdateServiceInterface {
	return &KRIValueUpdateService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		kriRegisterSvc:     kriRegisterSvc,
		reviewFrequencySvc: reviewFrequencySvc,
	}
}

func (kvus *KRIValueUpdateService) GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error) {
	_, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return nil, err
	}

	return kvus.repo.GetAll(kriRegisterID, params)
}

func (kvus *KRIValueUpdateService) Create(kriRegisterID, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) (uint, error) {
	tx := kvus.repo.BeginTransaction()
	existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return 0, err
	}

	reviewFrequency, err := kvus.reviewFrequencySvc.GetByID(existingKRIRegister.ReviewFrequencyID)
	if err != nil {
		return 0, err
	}

	nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)

	dates, err := kvus.parseKRIValueUpdateDates(kriValueUpdateFormDTO)
	if err != nil {
		return 0, err
	}
	kriValueUpdate := &models.KRIValueUpdate{
		KRIRegisterID:   kriRegisterID,
		ReviewedDate:    *dates.ReviewedDate,
		Value:           kriValueUpdateFormDTO.Value,
		LowThreshold:    *existingKRIRegister.LowThreshold,
		MediumThreshold: *existingKRIRegister.MediumThreshold,
		HighThreshold:   *existingKRIRegister.HighThreshold,
		Remarks:         kriValueUpdateFormDTO.Remarks,
		CreatedByID:     currentUserID,
	}

	id, err := kvus.repo.Create(kriValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	existingKRIRegister.CurrentValue = &kriValueUpdateFormDTO.Value
	existingKRIRegister.LastReviewNote = kriValueUpdateFormDTO.Remarks
	existingKRIRegister.LastReviewedDate = dates.ReviewedDate
	existingKRIRegister.NextReviewDate = nextReviewDate

	err = kvus.kriRegisterSvc.UpdateValue(existingKRIRegister, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.KRIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kriValueUpdateFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := fmt.Sprintf("%s/kri-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kriRegisterID, encryptedFilename)

		// Create a new KRIValueUpdateDocument instance
		document := &models.KRIValueUpdateDocument{
			KRIValueUpdateID: id,
			URL:              url,
			Name:             doc.FileName,
			Size:             doc.Size,
			Token:            encryptedFilename,
			OriginalName:     doc.OriginalName,
			Type:             doc.Type,
		}
		// Append the document to the slice
		documents = append(documents, document)

	}
	if len(documents) > 0 {
		if err := kvus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range kriValueUpdateFormDTO.Documents {
			storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := kvus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRIValueUpdate{}.TableName(),
			id,
			kriValueUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KRIValueUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (kvus KRIValueUpdateService) GetFullDetailsByID(id uint) (*riskDTOs.KRIValueUpdateDetailsDTO, error) {
	var valueUpdate riskDTOs.KRIValueUpdateDetailsDTO

	existingKRIUpdate, err := kvus.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(existingKRIUpdate, &valueUpdate)
	if err != nil {
		return nil, err
	}

	
	reviewedDate := existingKRIUpdate.ReviewedDate.Format("2006-01-02")
	valueUpdate.ReviewedDateString = &reviewedDate

	valueUpdate.CreatedAt = existingKRIUpdate.CreatedAt
	valueUpdate.UpdatedAt = existingKRIUpdate.UpdatedAt

	return &valueUpdate, nil
}

func (kvus *KRIValueUpdateService) Update(kriRegisterID, updateId uint, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) error {
	tx := kvus.repo.BeginTransaction()

	existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	dates, err := kvus.parseKRIValueUpdateDates(kriValueUpdateFormDTO)
	if err != nil {
		return err
	}

	existingKRIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}

	// Update KRI Value Update
	existingKRIUpdate.KRIRegisterID = kriRegisterID
	existingKRIUpdate.ReviewedDate = *dates.ReviewedDate
	existingKRIUpdate.Value = kriValueUpdateFormDTO.Value
	existingKRIUpdate.Remarks = kriValueUpdateFormDTO.Remarks

	err = kvus.repo.Update(existingKRIUpdate, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingKRIRegister.CurrentValue = &kriValueUpdateFormDTO.Value
	existingKRIRegister.LastReviewNote = kriValueUpdateFormDTO.Remarks

	err = kvus.kriRegisterSvc.UpdateValue(existingKRIRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.KRIValueUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.KRIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kriValueUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := fmt.Sprintf("%s/kri-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kriRegisterID, encryptedFilename)

			// Create a new KRIValueUpdateDocument instance
			document := &models.KRIValueUpdateDocument{
				KRIValueUpdateID: updateId,
				URL:              url,
				Name:             doc.FileName,
				Size:             doc.Size,
				Token:            encryptedFilename,
				OriginalName:     doc.OriginalName,
				Type:             doc.Type,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := kvus.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, doc.ID)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}
	if len(documents) > 0 {
		if err := kvus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := kvus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	return nil
}

func (kvus *KRIValueUpdateService) Delete(kriRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := kvus.repo.BeginTransaction()

	_, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	documentsToDelete := make(map[uint][]*models.KRIValueUpdateDocument)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		valueUpdate, err := kvus.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if valueUpdate.Documents != nil && len(valueUpdate.Documents) > 0 {
			for _, updateDocument := range valueUpdate.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], &updateDocument)
				deleteDocumentIds = append(deleteDocumentIds, &updateDocument.ID)

			}
			if len(deleteDocumentIds) > 0 {
				if err := kvus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}

		}

		if err := kvus.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for updateId, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}
	}

	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := kvus.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KRIValueUpdate{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (kvus *KRIValueUpdateService) parseKRIValueUpdateDates(dto *riskDTOs.KRIValueUpdateFormDTO) (*riskDTOs.KRIValueUpdateParsedDates, error) {
	parsed := &riskDTOs.KRIValueUpdateParsedDates{}

	var err error
	if dto.ReviewedDate != "" {
		parsed.ReviewedDate, err = utility.ParseDate(dto.ReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid review date"), "invalid_date", "review date")
		}
	}

	return parsed, nil
}
