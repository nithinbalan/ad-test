package services

import (
	"fmt"
	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
)

type KRIValueUpdateServiceInterface interface {
	GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error)
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
		CreatedBy:       currentUserID,
	}

	id, err := kvus.repo.Create(kriValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	existingKRIRegister.CurrentValue = &kriValueUpdateFormDTO.Value
	existingKRIRegister.LastReviewNote = kriValueUpdateFormDTO.Remarks
	existingKRIRegister.LastReviewedDate = *dates.ReviewedDate
	existingKRIRegister.NextReviewDate = nextReviewDate

	err = kvus.kriRegisterSvc.UpdateValue(existingKRIRegister, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.KRIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kriValueUpdateFormDTO.Documents {

		if err := kvus.repo.StoreKRIValueUpdateDocumentFile(kriRegisterID, id, doc.FileName); err != nil {
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
		url := fmt.Sprintf("%s/kri-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kriRegisterID, encryptedFilename)

		// Create a new KRIValueUpdateDocument instance
		document := &models.KRIValueUpdateDocument{
			KRIValueUpdateID: id,
			URL:              url,
			Icon:             doc.Icon,
			Name:             doc.FileName,
			Size:             doc.Size,
			Token:            encryptedFilename,
			Type:             doc.Type,
		}
		// Append the document to the slice
		documents = append(documents, document)

	}
	if len(documents) > 0 {
		if err := kvus.repo.SaveKRIValueUpdateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
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
	// Iterate over each document in the array
	for _, doc := range kriValueUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Handle new file case
			if err := kvus.repo.StoreKRIValueUpdateDocumentFile(existingKRIRegister.ID, updateId, doc.FileName); err != nil {
				tx.Rollback()
				return err
			}
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
				Icon:             doc.Icon,
				Name:             doc.FileName,
				Size:             doc.Size,
				Token:            encryptedFilename,
				Type:             doc.Type,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {

			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := kvus.repo.GetKRIValueUpdateDocumentByID(*doc.ID)
				if err != nil {
					return err
				}

				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				if err := kvus.repo.RemoveKRIValueUpdateDocumentFile(existingKRIRegister.ID, updateId, existingDocument.Name); err != nil {
					tx.Rollback()
					return err
				}
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}
	if len(documents) > 0 {
		if err := kvus.repo.SaveKRIValueUpdateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := kvus.repo.DeleteKRIValueUpdateDocuments(deleteDocumentIds, tx); err != nil {
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
	var deleteDocumentIds []*uint

	for _, id := range ids {

		valueUpdate, err := kvus.repo.GetByID(id)

		if err != nil {
			continue
		}

		for _, updateDocument := range valueUpdate.KRIValueUpdateDocuments {

			document, err := kvus.repo.GetKRIValueUpdateDocumentByID(updateDocument.ID)
			if err != nil {
				return nil, err
			}
			if document == nil {
				tx.Rollback()
				return nil, fmt.Errorf("existing document not found for ID: %d", updateDocument.ID)
			}
			if err := kvus.repo.RemoveKRIValueUpdateDocumentFile(valueUpdate.KRIRegisterID, valueUpdate.ID, document.Name); err != nil {
				tx.Rollback()
				return nil, err
			}
			deleteDocumentIds = append(deleteDocumentIds, &document.ID)

		}
		if len(deleteDocumentIds) > 0 {
			if err := kvus.repo.DeleteKRIValueUpdateDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		err = kvus.repo.Delete(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

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

	tx.Commit()

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
