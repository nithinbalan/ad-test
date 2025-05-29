package services

import (
	"errors"
	"fmt"
	"jethur/config"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type UserDocumentServiceInterface interface {
	GetAll(userID uint, params *orgDTOs.UserDocumentListQueryParamsDTO) ([]*orgDTOs.UserDocumentListDTO, error)
	GetFullDetailsByID(userID uint, id uint) (*orgDTOs.UserDocumentDetailDTO, error)
	GetByID(userID uint, id uint) (*models.UserDocument, error)
	Create(userID uint, currentUserID uint, userDocumentFormDTO *orgDTOs.UserDocumentFormDTO) (uint, error)
	Update(userID uint, id uint, currentUserID uint, userDocumentFormDTO *orgDTOs.UserDocumentFormDTO) error
	Delete(userID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(userID uint, id uint, token string) (string, error)
}

type UserDocumentService struct {
	repo                   repositories.UserDocumentRepositoryInterface
	userDocumentTypeSvc    UserDocumentTypeServiceInterface
	userDocumentSubTypeSvc UserDocumentSubTypeServiceInterface
	activityLogSvc         activityLogServices.ActivityLogServiceInterface
	userSvc                UserServiceInterface
}

func NewUserDocumentService(
	repo repositories.UserDocumentRepositoryInterface,
	userDocumentTypeSvc UserDocumentTypeServiceInterface,
	userDocumentSubTypeSvc UserDocumentSubTypeServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	userSvc UserServiceInterface,
) UserDocumentServiceInterface {
	return &UserDocumentService{
		repo:                   repo,
		userDocumentTypeSvc:    userDocumentTypeSvc,
		userDocumentSubTypeSvc: userDocumentSubTypeSvc,
		activityLogSvc:         activityLogSvc,
		userSvc:                userSvc,
	}
}

func (uds *UserDocumentService) GetAll(userID uint, params *orgDTOs.UserDocumentListQueryParamsDTO) ([]*orgDTOs.UserDocumentListDTO, error) {
	if _, err := uds.userSvc.GetByID(userID); err != nil {
		return nil, err
	}
	return uds.repo.GetAll(userID, params)
}

func (uds *UserDocumentService) GetByID(userID uint, id uint) (*models.UserDocument, error) {
	if _, err := uds.userSvc.GetByID(userID); err != nil {
		return nil, err
	}
	userDocument, err := uds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return userDocument, nil
}

func (uds *UserDocumentService) GetFullDetailsByID(userID uint, id uint) (*orgDTOs.UserDocumentDetailDTO, error) {
	var userDocumentDetailDTO orgDTOs.UserDocumentDetailDTO
	userDocument, err := uds.repo.GetFullDetailsByID(id)
	if err != nil {
		return nil, err
	}

	err = mapstructure.Decode(userDocument, &userDocumentDetailDTO)
	if err != nil {
		return nil, err
	}

	if userDocument.ExpiryDate != nil {
		expiryDate := userDocument.ExpiryDate.Format("2006-01-02")
		userDocumentDetailDTO.ExpiryDateString = expiryDate
	}

	return &userDocumentDetailDTO, nil
}

func (uds *UserDocumentService) Create(userID uint, currentUserID uint, formDTO *orgDTOs.UserDocumentFormDTO) (uint, error) {
	tx := uds.repo.BeginTransaction()

	if _, err := uds.userSvc.GetByID(userID); err != nil {
		return 0, err
	}

	if err := uds.validateUserDocumentFormDTO(formDTO); err != nil {
		return 0, err
	}

	dates, err := uds.parseUserDocumentDates(formDTO)
	if err != nil {
		return 0, err
	}

	url := utility.GenerateDocumentURL("organization", "users", strconv.FormatUint(uint64(userID), 10), "user-document")

	dataToEncrypt := []byte(formDTO.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	userDocument := &models.UserDocument{
		UserID:                userID,
		UserDocumentTypeID:    formDTO.UserDocumentTypeID,
		UserDocumentSubTypeID: formDTO.UserDocumentSubTypeID,
		IssuedBy:              &formDTO.IssuedBy,
		ExpiryDate:            dates.ExpiryDate,
		Remarks:               formDTO.Remarks,
		URL:                   url,
		Name:                  formDTO.FileName,
		Size:                  formDTO.Size,
		Type:                  formDTO.Type,
		Token:                 encryptedFilename,
		OriginalName:          formDTO.OriginalName,
		CreatedBYID:           currentUserID,
	}

	id, err := uds.repo.Create(userDocument)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	storagePath := filepath.Join("organization", "users", strconv.FormatUint(uint64(userID), 10), "user-document")
	if err := utility.StoreFile(formDTO.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return id, nil
}

func (uds *UserDocumentService) Update(userID uint, id uint, currentUserID uint, formDTO *orgDTOs.UserDocumentFormDTO) error {
	tx := uds.repo.BeginTransaction()

	if _, err := uds.userSvc.GetByID(userID); err != nil {
		tx.Rollback()
		return err
	}

	oldData, err := uds.GetFullDetailsByID(userID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existing, err := uds.GetByID(userID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := uds.validateUserDocumentFormDTO(formDTO); err != nil {
		return err
	}

	dates, err := uds.parseUserDocumentDates(formDTO)
	if err != nil {
		return err
	}

	url := utility.GenerateDocumentURL("organization", "users", strconv.FormatUint(uint64(userID), 10), "user-document")

	dataToEncrypt := []byte(formDTO.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Define the storage path
	storagePath := filepath.Join("organization", "users", strconv.FormatUint(uint64(userID), 10), "user-document")

	// Remove the old document if exists
	if &existing.Name != nil {
		if err := utility.RemoveFile(existing.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Store new document
	if err := utility.StoreFile(formDTO.FileName, storagePath); err != nil {
		tx.Rollback()
		return err
	}

	existing.UserDocumentTypeID = formDTO.UserDocumentTypeID
	existing.UserDocumentSubTypeID = formDTO.UserDocumentSubTypeID
	existing.IssuedBy = &formDTO.IssuedBy
	existing.ExpiryDate = dates.ExpiryDate
	existing.Remarks = formDTO.Remarks
	existing.URL = url
	existing.Name = formDTO.FileName
	existing.Size = formDTO.Size
	existing.Type = formDTO.Type
	existing.Token = encryptedFilename
	existing.OriginalName = formDTO.OriginalName
	existing.UpdatedBYID = &currentUserID

	err = uds.repo.Update(existing)
	if err != nil {
		return err
	}

	go func() {
		// Fetch updated full details after update
		newData, err := uds.GetFullDetailsByID(userID, id)
		if err != nil {
			return
		}

		oldFields := map[string]interface{}{
			"sub_type_title_en": oldData.UserDocumentSubType.TitleEn,
			"type_title_en":     oldData.UserDocumentType.TitleEn,
			"issued_by":         oldData.IssuedBy,
			"expiry_date":       oldData.ExpiryDateString,
			"remarks":           oldData.Remarks,
			"url":               oldData.URL,
			"name":              oldData.Name,
			"original_name":     oldData.OriginalName,
			"type":              oldData.Type,
			"token":             oldData.Token,
			"size":              oldData.Size,
		}

		newFields := map[string]interface{}{
			"sub_type_title_en": newData.UserDocumentSubType.TitleEn,
			"type_title_en":     newData.UserDocumentType.TitleEn,
			"issued_by":         newData.IssuedBy,
			"expiry_date":       newData.ExpiryDateString,
			"remarks":           newData.Remarks,
			"url":               newData.URL,
			"name":              newData.Name,
			"original_name":     newData.OriginalName,
			"type":              newData.Type,
			"token":             newData.Token,
			"size":              newData.Size,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)

		if len(changedData) > 0 {
			err := uds.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.UserDocument{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log process dependency update", map[string]interface{}{
					"action":    models.UpdateAction,
					"table":     models.UserDocument{}.TableName(),
					"record_id": id,
					"user_id":   currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (uds *UserDocumentService) Delete(userID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := uds.userSvc.GetByID(userID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		// Define the storage path
		storagePath := filepath.Join("organization", "users", strconv.FormatUint(uint64(userID), 10), "user-document")

		userDocument, err := uds.GetFullDetailsByID(userID, id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		// Remove the old document if exists
		if &userDocument.Name != nil {
			if err := utility.RemoveFile(userDocument.Name, storagePath); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}
		if err := uds.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedItems = append(deletedItems, map[string]interface{}{
				"id":                     id,
				"user_document_type":     userDocument.UserDocumentType.TitleEn,
				"user_document_sub_type": userDocument.UserDocumentSubType.TitleEn,
			})
		}
	}

	if len(successfulIds) > 0 {
		go uds.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.UserDocument{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (uds *UserDocumentService) validateUserDocumentFormDTO(userDocumentFormDTO *orgDTOs.UserDocumentFormDTO) error {
	// Validate sub type
	if _, err := uds.userDocumentSubTypeSvc.GetByID(userDocumentFormDTO.UserDocumentSubTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("user document sub type does not exist"), "not_exist", "user document sub type")
	}

	if _, err := uds.userDocumentTypeSvc.GetByID(userDocumentFormDTO.UserDocumentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("user document type does not exist"), "not_exist", "user document sub type")
	}

	return nil
}

func (uds *UserDocumentService) parseUserDocumentDates(formDTO *orgDTOs.UserDocumentFormDTO) (*orgDTOs.UserDocumentParsedDates, error) {
	parsed := &orgDTOs.UserDocumentParsedDates{}

	var err error
	if formDTO.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(formDTO.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	return parsed, nil
}

func (uds *UserDocumentService) GetDocumentPath(userID uint, id uint, token string) (string, error) {
	_, err := uds.GetByID(userID, id)
	if err != nil {
		return "", err
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("organization", "users", strconv.FormatUint(uint64(userID), 10), "user-document")

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
