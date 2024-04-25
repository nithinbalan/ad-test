package services

import (
	"fmt"
	"jethur/config"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
)

type KRIValueUpdateServiceInterface interface {
	GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error)
	Create(kriRegisterID, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) (uint, error)
	Update(kriRegisterID, updateId, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type KRIValueUpdateService struct {
	repo           riskRepositories.KRIValueUpdateRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	kriRegisterSvc KRIRegisterServiceInterface
}

func NewKRIValueUpdateService(
	repo riskRepositories.KRIValueUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriRegisterSvc KRIRegisterServiceInterface,
) KRIValueUpdateServiceInterface {
	return &KRIValueUpdateService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		kriRegisterSvc: kriRegisterSvc,
	}
}

func (kvus *KRIValueUpdateService) GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error) {
	if err := kvus.validateKRIRegister(kriRegisterID); err != nil {
		return nil, err
	}

	return kvus.repo.GetAll(kriRegisterID, params)
}

func (kvus *KRIValueUpdateService) Create(kriRegisterID, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) (uint, error) {
	tx := kvus.repo.BeginTransaction()

	if err := kvus.validateKRIRegister(kriRegisterID); err != nil {
		return 0, err
	}

	dates, err := kvus.parseKRIValueUpdateDates(kriValueUpdateFormDTO)
	if err != nil {
		return 0, err
	}

	kriValueUpdate := &models.KRIValueUpdate{
		KRIRegisterID:   kriRegisterID,
		ReviewDate:      *dates.ReviewDate,
		Value:           kriValueUpdateFormDTO.Value,
		LowThreshold:    kriValueUpdateFormDTO.LowThreshold,
		MediumThreshold: kriValueUpdateFormDTO.MediumThreshold,
		HighThreshold:   kriValueUpdateFormDTO.HighThreshold,
		Remarks:         kriValueUpdateFormDTO.Remarks,
		CreatedBy:       currentUserID,
	}

	id, err := kvus.repo.Create(kriValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	// if err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// kriRegister := &riskDTOs.KRIRegisterFormDTO{
	// 	KRICategoryID:     existingKRIRegister.KRICategoryID,
	// 	OwnerUnitID:       existingKRIRegister.OwnerUnitID,
	// 	OwnerID:           existingKRIRegister.OwnerID,
	// 	ReviewFrequencyID: existingKRIRegister.ReviewFrequencyID,
	// 	RiskLevelID:       existingKRIRegister.RiskLevelID,
	// 	Title:             existingKRIRegister.Title,
	// 	Description:       existingKRIRegister.Description,
	// 	CalculationType:   existingKRIRegister.CalculationType,
	// 	LowThreshold:      &kriValueUpdateFormDTO.LowThreshold,
	// 	MediumThreshold:   &kriValueUpdateFormDTO.MediumThreshold,
	// 	HighThreshold:     &kriValueUpdateFormDTO.HighThreshold,
	// 	CurrentValue:      &kriValueUpdateFormDTO.Value,
	// 	LastReviewDate:    kriValueUpdateFormDTO.ReviewDate,
	// 	NextReviewDate:    kriValueUpdateFormDTO.ReviewDate,
	// 	LastReviewNote:    kriValueUpdateFormDTO.Remarks,
	// }
	// err = kvus.kriRegisterSvc.Update(kriRegisterID, currentUserID, kriRegister)
	// if err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

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
			logModels.CreateAction,
			models.KRIValueUpdate{}.TableName(),
			id,
			kriValueUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
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

	if err := kvus.validateKRIRegister(kriRegisterID); err != nil {
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
	existingKRIUpdate.ReviewDate = *dates.ReviewDate
	existingKRIUpdate.Value = kriValueUpdateFormDTO.Value
	existingKRIUpdate.LowThreshold = kriValueUpdateFormDTO.LowThreshold
	existingKRIUpdate.MediumThreshold = kriValueUpdateFormDTO.MediumThreshold
	existingKRIUpdate.HighThreshold = kriValueUpdateFormDTO.HighThreshold
	existingKRIUpdate.Remarks = kriValueUpdateFormDTO.Remarks

	err = kvus.repo.Update(existingKRIUpdate, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// kriRegister := &riskDTOs.KRIRegisterFormDTO{
	// 	KRICategoryID:     existingKRIRegister.KRICategoryID,
	// 	OwnerUnitID:       existingKRIRegister.OwnerUnitID,
	// 	OwnerID:           existingKRIRegister.OwnerID,
	// 	ReviewFrequencyID: existingKRIRegister.ReviewFrequencyID,
	// 	RiskLevelID:       existingKRIRegister.RiskLevelID,
	// 	Title:             existingKRIRegister.Title,
	// 	Description:       existingKRIRegister.Description,
	// 	CalculationType:   existingKRIRegister.CalculationType,
	// 	LowThreshold:      existingKRIRegister.LowThreshold,
	// 	MediumThreshold:   existingKRIRegister.MediumThreshold,
	// 	HighThreshold:     existingKRIRegister.HighThreshold,
	// 	CurrentValue:      existingKRIRegister.CurrentValue,
	// 	LastReviewDate:    kriValueUpdateFormDTO.ReviewDate,
	// 	NextReviewDate:    kriValueUpdateFormDTO.ReviewDate,
	// 	LastReviewNote:    kriValueUpdateFormDTO.Remarks,
	// }

	// err = kvus.kriRegisterSvc.Update(kriRegisterID, currentUserID, kriRegister)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

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

func (kvus *KRIValueUpdateService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := kvus.repo.BeginTransaction()

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
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
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
	if dto.ReviewDate != "" {
		parsed.ReviewDate, err = utility.ParseDate(dto.ReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid review date"), "invalid_date", "review date")
		}
	}

	return parsed, nil
}

func (kvus *KRIValueUpdateService) validateKRIRegister(kriRegisterID uint) error {

	if _, err := kvus.kriRegisterSvc.GetByID(kriRegisterID); err != nil {
		return utility.NewValidationError(fmt.Errorf("kri register does not exist"), "not_exist", "kri register")
	}

	return nil
}
