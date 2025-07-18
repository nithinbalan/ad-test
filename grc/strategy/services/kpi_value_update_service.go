package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KPIValueUpdateServiceInterface interface {
	GetAll(kpiRegisterID uint, params *strategyDTOs.KPIValueUpdateListQueryParamsDTO) (*strategyDTOs.KPIValueUpdateListResponseDTO, error)
	GetFullDetailsByID(id uint) (*strategyDTOs.KPIValueUpdateDetailsDTO, error)
	Create(kpiRegisterID, currentUserID uint, kpiValueUpdateFormDTO *strategyDTOs.KPIValueUpdateFormDTO) (uint, error)
	Update(kpiRegisterID, updateId, currentUserID uint, kpiValueUpdateFormDTO *strategyDTOs.KPIValueUpdateFormDTO) error
	Delete(kpiRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(id uint, kpiValueUpdateId uint, token string) (string, error)
}

// Service and implementations
type KPIValueUpdateService struct {
	repo               strategyRepositories.KPIValueUpdateRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	kpiRegisterSvc     KPIRegisterServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	userSvc            organizationServices.UserServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
}

func NewKPIValueUpdateService(
	repo strategyRepositories.KPIValueUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kpiRegisterSvc KPIRegisterServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) KPIValueUpdateServiceInterface {
	return &KPIValueUpdateService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		kpiRegisterSvc:     kpiRegisterSvc,
		reviewFrequencySvc: reviewFrequencySvc,
		userSvc:            userSvc,
		notificationSvc:    notificationSvc,
	}
}

func (kvus *KPIValueUpdateService) GetAll(kpiRegisterID uint, params *strategyDTOs.KPIValueUpdateListQueryParamsDTO) (*strategyDTOs.KPIValueUpdateListResponseDTO, error) {
	_, err := kvus.kpiRegisterSvc.GetByID(kpiRegisterID)
	if err != nil {
		return nil, err
	}

	return kvus.repo.GetAll(kpiRegisterID, params)
}

func (kvus *KPIValueUpdateService) Create(kpiRegisterID, currentUserID uint, kpiValueUpdateFormDTO *strategyDTOs.KPIValueUpdateFormDTO) (uint, error) {
	tx := kvus.repo.BeginTransaction()
	existingKPIRegister, err := kvus.kpiRegisterSvc.GetByID(kpiRegisterID)
	if err != nil {
		return 0, err
	}

	reviewFrequency, err := kvus.reviewFrequencySvc.GetByID(*existingKPIRegister.ReviewFrequencyID)
	if err != nil {
		return 0, err
	}

	nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)

	dates, err := kvus.parseKPIValueUpdateDates(kpiValueUpdateFormDTO)
	if err != nil {
		return 0, err
	}
	// Set all previous value updates for this KPIRegister to IsLatest = false
if err := kvus.repo.SetAllIsLatestFalse(kpiRegisterID, tx); err != nil {
	tx.Rollback()
	return 0, err
}

kpiValueUpdate := &models.KPIValueUpdate{
		KPIRegisterID: kpiRegisterID,
		ReviewUserID:  kpiValueUpdateFormDTO.ReviewUserID,
		UnitID:        *existingKPIRegister.UnitID,
		ReviewedDate:  *dates.ReviewedDate,
		Value:         &kpiValueUpdateFormDTO.Value,
		Remarks:       kpiValueUpdateFormDTO.Remarks,
		IsLatest:      true, // Always set to true (1) when creating
		CreatedByID:   currentUserID,
	}

	id, err := kvus.repo.Create(kpiValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	existingKPIRegister.ActualValue = &kpiValueUpdateFormDTO.Value
	existingKPIRegister.LastReviewNote = kpiValueUpdateFormDTO.Remarks
	existingKPIRegister.LastReviewedDate = dates.ReviewedDate
	existingKPIRegister.NextReviewDate = &nextReviewDate

	err = kvus.kpiRegisterSvc.UpdateValue(existingKPIRegister, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.KPIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kpiValueUpdateFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := fmt.Sprintf("%s/kpi-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kpiRegisterID, encryptedFilename)

		// Create a new KPIValueUpdateDocument instance
		document := &models.KPIValueUpdateDocument{
			KPIValueUpdateID: id,
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

		for _, doc := range kpiValueUpdateFormDTO.Documents {
			storagePath := filepath.Join("strategy", "kpi-register", strconv.FormatUint(uint64(kpiRegisterID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	// Prepare input data records
	var inputs []*models.KPIValueUpdateDataInput
	for _, input := range kpiValueUpdateFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIValueUpdateDataInput{
				KPIValueUpdateID: id,
				Input:            input.Input,
				Variable:         input.Variable,
				Value:            input.Value,
				CreatedBYID:      currentUserID,
			}
			inputs = append(inputs, inputData)
		}
	}

	if len(inputs) > 0 {
		if err := kvus.repo.CreateKPIValueUpdateInputData(inputs, tx); err != nil {
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingKpiValueUpdate, err := kvus.GetFullDetailsByID(kpiValueUpdate.ID)
		if err != nil {
			utility.LogError("Failed to fetch kpi value update", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = kvus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KPIValueUpdate{}.TableName(),
			id,
			existingKpiValueUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KPIValueUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (kvus KPIValueUpdateService) GetFullDetailsByID(id uint) (*strategyDTOs.KPIValueUpdateDetailsDTO, error) {
	var valueUpdate strategyDTOs.KPIValueUpdateDetailsDTO

	existingKPIUpdate, err := kvus.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(existingKPIUpdate, &valueUpdate)
	if err != nil {
		return nil, err
	}

	reviewedDate := existingKPIUpdate.ReviewedDate.Format("2006-01-02")
	valueUpdate.ReviewedDateString = &reviewedDate

	valueUpdate.CreatedAt = existingKPIUpdate.CreatedAt
	valueUpdate.UpdatedAt = existingKPIUpdate.UpdatedAt

	return &valueUpdate, nil
}

func (kvus *KPIValueUpdateService) Update(kpiRegisterID, updateId uint, currentUserID uint, kpiValueUpdateFormDTO *strategyDTOs.KPIValueUpdateFormDTO) error {
	tx := kvus.repo.BeginTransaction()

	oldKpiValueUpdate, err := kvus.GetFullDetailsByID(updateId)
	if err != nil {
		return err
	}

	existingKPIRegister, err := kvus.kpiRegisterSvc.GetByID(kpiRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	dates, err := kvus.parseKPIValueUpdateDates(kpiValueUpdateFormDTO)
	if err != nil {
		return err
	}

	existingKPIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}

	// Update KPI Value Update
	existingKPIUpdate.KPIRegisterID = kpiRegisterID
	existingKPIUpdate.ReviewedDate = *dates.ReviewedDate
	existingKPIUpdate.UnitID = *existingKPIRegister.UnitID
	existingKPIUpdate.Value = &kpiValueUpdateFormDTO.Value
	existingKPIUpdate.Remarks = kpiValueUpdateFormDTO.Remarks
	existingKPIUpdate.ReviewUserID = kpiValueUpdateFormDTO.ReviewUserID

	err = kvus.repo.Update(existingKPIUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingKPIRegister.ActualValue = &kpiValueUpdateFormDTO.Value
	existingKPIRegister.LastReviewNote = kpiValueUpdateFormDTO.Remarks

	err = kvus.kpiRegisterSvc.UpdateValue(existingKPIRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.KPIValueUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.KPIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kpiValueUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := fmt.Sprintf("%s/kpi-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kpiRegisterID, encryptedFilename)

			// Create a new KPIValueUpdateDocument instance
			document := &models.KPIValueUpdateDocument{
				KPIValueUpdateID: updateId,
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
		storagePath := filepath.Join("strategy", "kpi-register", strconv.FormatUint(uint64(kpiRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("strategy", "kpi-register", strconv.FormatUint(uint64(kpiRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := kvus.repo.DeleteKPIValueUpdateInputData(updateId, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Prepare input data records
	var inputs []*models.KPIValueUpdateDataInput
	for _, input := range kpiValueUpdateFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIValueUpdateDataInput{
				KPIValueUpdateID: updateId,
				Input:            input.Input,
				Variable:         input.Variable,
				Value:            input.Value,
				CreatedBYID:      currentUserID,
			}
			inputs = append(inputs, inputData)
		}
	}

	if len(inputs) > 0 {
		if err := kvus.repo.CreateKPIValueUpdateInputData(inputs, tx); err != nil {
			return err
		}
	}

	tx.Commit()
	go func() {
		existingKpiValueUpdate, err := kvus.GetFullDetailsByID(updateId)
		if err != nil {
			return
		}
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"review_date": oldKpiValueUpdate.ReviewedDateString,
			"value":       oldKpiValueUpdate.Value,
			"remarks":     oldKpiValueUpdate.Remarks,
		}

		newData := map[string]interface{}{
			"review_date": existingKpiValueUpdate.ReviewedDateString,
			"value":       existingKpiValueUpdate.Value,
			"remarks":     existingKpiValueUpdate.Remarks,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err = kvus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.KPIValueUpdate{}.TableName(),
				updateId,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.KPIValueUpdate{}.TableName(),
					"recordID":  updateId,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (kvus *KPIValueUpdateService) Delete(kpiRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := kvus.repo.BeginTransaction()

	existingKPIRegister, err := kvus.kpiRegisterSvc.GetByID(kpiRegisterID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedKpiValueUpdate []map[string]interface{}

	documentsToDelete := make(map[uint][]*models.KPIValueUpdateDocument)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		valueUpdate, err := kvus.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if len(valueUpdate.Documents) > 0 {
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

		deletedKpiValueUpdate = append(deletedKpiValueUpdate, map[string]interface{}{
			"id":              id,
			"kpi_register_id": kpiRegisterID,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for updateId, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("strategy", "kpi-register", strconv.FormatUint(uint64(kpiRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}

		// Update KPIRegister with the next available KPIValueUpdate value
		nextValueUpdate, err := kvus.repo.GetNextKPIValueUpdate(kpiRegisterID, successfulIds, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		reviewFrequency, err := kvus.reviewFrequencySvc.GetByID(*existingKPIRegister.ReviewFrequencyID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)

		if nextValueUpdate != nil && nextValueUpdate.ID != 0 {
			existingKPIRegister.ActualValue = nextValueUpdate.Value
			existingKPIRegister.LastReviewNote = nextValueUpdate.Remarks
			existingKPIRegister.LastReviewedDate = &nextValueUpdate.ReviewedDate
			existingKPIRegister.NextReviewDate = &nextReviewDate
		} else {
			nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)
			// Set fields to nil if no next value update exists
			existingKPIRegister.ActualValue = nil
			existingKPIRegister.LastReviewNote = nil
			existingKPIRegister.NextReviewDate = &nextReviewDate
		}

		if err := kvus.kpiRegisterSvc.UpdateValue(existingKPIRegister, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := kvus.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KPIValueUpdate{}.TableName(),
				successfulIds,
				deletedKpiValueUpdate,
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

func (kvus *KPIValueUpdateService) GetDocumentPath(id uint, kpiValueUpdateId uint, token string) (string, error) {
	// Get full details of the kpi register by ID
	existingKPIUpdate, err := kvus.GetFullDetailsByID(kpiValueUpdateId)
	if err != nil {
		return "", err
	}

	// Check if there are any associated documents
	if len(existingKPIUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the kpi value update")
	}

	// Decrypt the token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Define the storage path
	storagePath := filepath.Join("strategy", "kpi-register", strconv.FormatUint(uint64(id), 10), "updates", strconv.FormatUint(uint64(kpiValueUpdateId), 10), "documents")

	// Get the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the specified path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (kvus *KPIValueUpdateService) parseKPIValueUpdateDates(dto *strategyDTOs.KPIValueUpdateFormDTO) (*strategyDTOs.KPIValueUpdateParsedDates, error) {
	parsed := &strategyDTOs.KPIValueUpdateParsedDates{}

	var err error
	if dto.ReviewedDate != "" {
		parsed.ReviewedDate, err = utility.ParseDate(dto.ReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid review date"), "invalid_date", "review date")
		}
	}

	return parsed, nil
}
