package services

import (
	"errors"
	"fmt"
	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"gorm.io/gorm"
)

type TaeRecoveryActionPlanStatusUpdateServiceInterface interface {
	GetAll(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error)
	Create(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) (uint, error)
	Update(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, id uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) error
	GetLatestTaeRecoveryActionPlanStatusUpdate(slaID uint, actionPlanID uint) (*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(slaID uint, actionPlanID uint, id uint, token string) (string, error)
}

type TaeRecoveryActionPlanStatusUpdateService struct {
	repo                                 repositories.TaeRecoveryActionPlanStatusUpdateRepositoryInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	testAndExerciseRecoveryActionPlanSvc TestAndExerciseRecoveryActionPlanServiceInterface
	actionPlanStatusSvc                  masterServices.ActionPlanStatusServiceInterface
}

func NewTaeRecoveryActionPlanStatusUpdateService(
	repo repositories.TaeRecoveryActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	testAndExerciseRecoveryActionPlanSvc TestAndExerciseRecoveryActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) TaeRecoveryActionPlanStatusUpdateServiceInterface {
	return &TaeRecoveryActionPlanStatusUpdateService{
		repo:                                 repo,
		activityLogSvc:                       activityLogSvc,
		testAndExerciseRecoveryActionPlanSvc: testAndExerciseRecoveryActionPlanSvc,
		actionPlanStatusSvc:                  actionPlanStatusSvc,
	}
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) GetAll(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, sort string, sortBy string, keyword string) ([]*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetAll(TaeRecoveryActionPlanID, sort, sortBy, keyword)
}

func (sasus TaeRecoveryActionPlanStatusUpdateService) GetByID(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, id uint) (*models.TaeRecoveryActionPlanStatusUpdate, error) {
	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID); err != nil {
		return nil, err
	}

	TaeRecoveryActionPlanStatusUpdate, err := sasus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return TaeRecoveryActionPlanStatusUpdate, nil
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) Create(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID); err != nil {
		return 0, err
	}

	if err := sasus.validateTaeRecoveryActionPlanStatusUpdateFormDTO(TaeRecoveryActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@todo validate document form

	TaeRecoveryActionPlanUpdate := &models.TaeRecoveryActionPlanStatusUpdate{
		TestAndExerciseRecoveryActionPlanID: TaeRecoveryActionPlanID,
		ActionPlanStatusID:                  TaeRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:                             TaeRecoveryActionPlanStatusUpdateFormDTO.Comment,
		Percentage:                          TaeRecoveryActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:                         currentUserID,
	}

	id, err := sasus.repo.Create(TaeRecoveryActionPlanUpdate)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = sasus.testAndExerciseRecoveryActionPlanSvc.StatusUpdate(TaeRecoveryID, TaeRecoveryActionPlanID, currentUserID, TaeRecoveryActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.TaeRecoveryActionPlanStatusUpdateDocument

	for _, doc := range TaeRecoveryActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TaeRecoveryID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.TaeRecoveryActionPlanStatusUpdateDocument{
			TaeRecoveryActionPlanStatusUpdateID: id,
			URL:                                 url,
			Name:                                doc.FileName,
			Size:                                doc.Size,
			Token:                               encryptedFilename,
			OriginalName:                        doc.OriginalName,
			Type:                                doc.Type,
			CreatedByID:                         currentUserID,
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
			storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TaeRecoveryID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
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
			models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
			id,
			TaeRecoveryActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) Update(TaeRecoveryID uint, TaeRecoveryActionPlanID uint, id uint, currentUserID uint, TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) error {
	tx := sasus.repo.BeginTransaction()

	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID); err != nil {
		return err
	}

	oldTaeRecoveryActionPlanStatusUpdate, err := sasus.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID, id)

	if err != nil {
		return err
	}

	existingTaeRecoveryActionPlanStatusUpdate, err := sasus.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID, id)

	if err != nil {
		return err
	}

	if err := sasus.validateTaeRecoveryActionPlanStatusUpdateFormDTO(TaeRecoveryActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	// Update TaeRecovery StatusUpdate-- only comments and percentage
	existingTaeRecoveryActionPlanStatusUpdate.Comment = TaeRecoveryActionPlanStatusUpdateFormDTO.Comment
	existingTaeRecoveryActionPlanStatusUpdate.Percentage = TaeRecoveryActionPlanStatusUpdateFormDTO.Percentage
	existingTaeRecoveryActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = sasus.repo.Update(existingTaeRecoveryActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.TaeRecoveryActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.TaeRecoveryActionPlanStatusUpdateDocument

	for _, doc := range TaeRecoveryActionPlanStatusUpdateFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TaeRecoveryID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new Control test Document instance
			document := &models.TaeRecoveryActionPlanStatusUpdateDocument{
				TaeRecoveryActionPlanStatusUpdateID: id,
				URL:                                 url,
				Name:                                doc.FileName,
				Size:                                doc.Size,
				Token:                               encryptedFileName,
				OriginalName:                        doc.OriginalName,
				Type:                                doc.Type,
				CreatedByID:                         currentUserID,
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
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TaeRecoveryID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(TaeRecoveryID), 10), "recovery-action-plans", strconv.FormatUint(uint64(TaeRecoveryActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := sasus.GetByID(TaeRecoveryID, TaeRecoveryActionPlanID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldTaeRecoveryActionPlanStatusUpdate.Comment,
			"percentage": oldTaeRecoveryActionPlanStatusUpdate.Percentage,
		}

		newData := map[string]interface{}{
			"comment":    existingActionPlanstatusUpdate.Comment,
			"percentage": existingActionPlanstatusUpdate.Percentage,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := sasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TaeRecoveryActionPlanStatusUpdate{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) validateTaeRecoveryActionPlanStatusUpdateFormDTO(TaeRecoveryActionPlanStatusUpdateFormDTO *businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateFormDTO) error {

	if _, err := sasus.actionPlanStatusSvc.GetByID(TaeRecoveryActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (sasus TaeRecoveryActionPlanStatusUpdateService) GetLatestTaeRecoveryActionPlanStatusUpdate(testAndExerciseID uint, actionPlanID uint) (*businessContinuityDTOs.TaeRecoveryActionPlanStatusUpdateListDTO, error) {
	if _, err := sasus.testAndExerciseRecoveryActionPlanSvc.GetByID(testAndExerciseID, actionPlanID); err != nil {
		return nil, err
	}

	return sasus.repo.GetLatestTaeRecoveryActionPlanStatusUpdate(actionPlanID)
}

func (sasus *TaeRecoveryActionPlanStatusUpdateService) GetDocumentPath(testAndExerciseID, ActionPlanID, id uint, token string) (string, error) {

	existingTaeRecoveryActionplanStatusUpdate, err := sasus.GetByID(testAndExerciseID, ActionPlanID, id)
	if err != nil {
		return "", err
	}
	if len(existingTaeRecoveryActionplanStatusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the recovery action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "recovery-action-plans", strconv.FormatUint(uint64(ActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
