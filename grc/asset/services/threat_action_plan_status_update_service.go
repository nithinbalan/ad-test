package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"jethur/config"
	threatDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatActionPlanStatusUpdateServiceInterface interface {
	GetAll(ThreatID uint, ThreatActionPlanID uint, sort string, sortBy string, keyword string) ([]*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error)
	Create(ThreatID uint, ThreatActionPlanID uint, currentUserID uint, threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ThreatID uint, ThreatActionPlanID uint, id uint, currentUserID uint, threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) error
	GetLatestThreatActionPlanStatusUpdate(threatID uint, actionPlanID uint) (*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error)
	GetDocumentPath(threatID, actionPlanID, statusUpdateID uint, token string) (string, error)
}

type ThreatActionPlanStatusUpdateService struct {
	repo                repositories.ThreatActionPlanStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	threatActionPlanSvc ThreatActionPlanServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
}

func NewThreatActionPlanStatusUpdateService(
	repo repositories.ThreatActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatActionPlanSvc ThreatActionPlanServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) ThreatActionPlanStatusUpdateServiceInterface {
	return &ThreatActionPlanStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		threatActionPlanSvc: threatActionPlanSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
	}
}

func (tasus *ThreatActionPlanStatusUpdateService) GetAll(ThreatID uint, ThreatActionPlanID uint, sort string, sortBy string, keyword string) ([]*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return nil, err
	}

	return tasus.repo.GetAll(ThreatActionPlanID, sort, sortBy, keyword)
}

func (tasus ThreatActionPlanStatusUpdateService) GetByID(ThreatID uint, ThreatActionPlanID uint, id uint) (*models.ThreatActionPlanStatusUpdate, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return nil, err
	}

	threatActionPlanStatusUpdate, err := tasus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return threatActionPlanStatusUpdate, nil
}

func (tasus *ThreatActionPlanStatusUpdateService) Create(ThreatID uint, ThreatActionPlanID uint, currentUserID uint, threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) (uint, error) {
	tx := tasus.repo.BeginTransaction()

	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return 0, err
	}

	if err := tasus.validateThreatActionPlanStatusUpdateFormDTO(threatActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	threatActionPlanUpdate := &models.ThreatActionPlanStatusUpdate{
		ThreatActionPlanID: ThreatActionPlanID,
		ActionPlanStatusID: threatActionPlanStatusUpdateFormDTO.ActionPlanStatusID,
		Comment:            threatActionPlanStatusUpdateFormDTO.Comment,
		Percentage:         threatActionPlanStatusUpdateFormDTO.Percentage,
		CreatedByID:        currentUserID,
	}

	id, err := tasus.repo.Create(threatActionPlanUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tasus.threatActionPlanSvc.StatusUpdate(ThreatID, ThreatActionPlanID, currentUserID, threatActionPlanStatusUpdateFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ThreatActionPlanStatusUpdateDocument

	for _, doc := range threatActionPlanStatusUpdateFormDTO.Documents {

		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ThreatActionPlanStatusUpdateDocument{
			ThreatActionPlanStatusUpdateID: id,
			URL:                            url,
			Name:                           doc.FileName,
			Size:                           doc.Size,
			Token:                          encryptedFilename,
			OriginalName:                   doc.OriginalName,
			Type:                           doc.Type,
			CreatedByID:                    currentUserID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := tasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	tx.Commit()

	go func() {
		err := tasus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatActionPlanStatusUpdate{}.TableName(),
			id,
			threatActionPlanUpdate,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ThreatActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (tasus *ThreatActionPlanStatusUpdateService) Update(ThreatID uint, ThreatActionPlanID uint, id uint, currentUserID uint, ThreatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) error {
	tx := tasus.repo.BeginTransaction()

	if _, err := tasus.threatActionPlanSvc.GetByID(ThreatID, ThreatActionPlanID); err != nil {
		return err
	}

	oldThreatActionPlanStatusUpdate, err := tasus.GetByID(ThreatID, ThreatActionPlanID, id)

	if err != nil {
		return err
	}

	existingThreatActionPlanStatusUpdate, err := tasus.GetByID(ThreatID, ThreatActionPlanID, id)
	if err != nil {
		return err
	}

	if err := tasus.validateThreatActionPlanStatusUpdateFormDTO(ThreatActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	existingThreatActionPlanStatusUpdate.Comment = ThreatActionPlanStatusUpdateFormDTO.Comment
	existingThreatActionPlanStatusUpdate.Percentage = ThreatActionPlanStatusUpdateFormDTO.Percentage
	existingThreatActionPlanStatusUpdate.UpdatedByID = &currentUserID

	err = tasus.repo.Update(existingThreatActionPlanStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ThreatActionPlanStatusUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ThreatActionPlanStatusUpdateDocument

	for _, doc := range ThreatActionPlanStatusUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			document := &models.ThreatActionPlanStatusUpdateDocument{
				ThreatActionPlanStatusUpdateID: id,
				URL:                            url,
				Name:                           doc.FileName,
				Size:                           doc.Size,
				Token:                          encryptedFileName,
				OriginalName:                   doc.OriginalName,
				Type:                           doc.Type,
				CreatedByID:                    currentUserID,
			}

			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := tasus.repo.GetDocumentByID(*doc.ID)
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
		if err := tasus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := tasus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(ThreatID), 10), "action-plans", strconv.FormatUint(uint64(ThreatActionPlanID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActionPlanstatusUpdate, err := tasus.GetByID(ThreatID, ThreatActionPlanID, id)
		if err != nil {
			tx.Rollback()
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment":    oldThreatActionPlanStatusUpdate.Comment,
			"percentage": oldThreatActionPlanStatusUpdate.Percentage,
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
			err := tasus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ThreatActionPlanStatusUpdate{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ThreatActionPlanStatusUpdate{}.TableName(),
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

func (tasus *ThreatActionPlanStatusUpdateService) validateThreatActionPlanStatusUpdateFormDTO(threatActionPlanStatusUpdateFormDTO *threatDTOs.ThreatActionPlanStatusUpdateFormDTO) error {
	if _, err := tasus.actionPlanStatusSvc.GetByID(threatActionPlanStatusUpdateFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (tasus ThreatActionPlanStatusUpdateService) GetLatestThreatActionPlanStatusUpdate(threatID uint, actionPlanID uint) (*threatDTOs.ThreatActionPlanStatusUpdateListDTO, error) {
	if _, err := tasus.threatActionPlanSvc.GetByID(threatID, actionPlanID); err != nil {
		return nil, err
	}

	return tasus.repo.GetLatestThreatActionPlanStatusUpdate(actionPlanID)
}

func (tasus *ThreatActionPlanStatusUpdateService) GetDocumentPath(threatID, actionPlanID, statusUpdateID uint, token string) (string, error) {
	statusUpdate, err := tasus.GetByID(threatID, actionPlanID, statusUpdateID)
	if err != nil {
		fmt.Printf("Error fetching status update by ID: %v\n", err)
		return "", err
	}

	if len(statusUpdate.Documents) == 0 {
		fmt.Println("No documents associated with the status update")
		return "", fmt.Errorf("no documents found for the threat action plan update")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		fmt.Printf("Token decryption failed: %v\n", err)
		return "", fmt.Errorf("resource not found")
	}
	decryptedFilename := string(decryptedFilenameBytes)
	fmt.Printf("Decrypted filename: %s\n", decryptedFilename)

	storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(threatID), 10), "action-plans", strconv.FormatUint(uint64(actionPlanID), 10), "updates", strconv.FormatUint(uint64(statusUpdateID), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)
	fmt.Printf("Constructed file path: %s\n", fullPath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fmt.Printf("File does not exist at path: %s\n", fullPath)
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		fmt.Printf("Error accessing file at path %s: %v\n", fullPath, err)
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
