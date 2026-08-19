package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/grc/security-operations/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type CSLogSocRoomAccessLogServiceInterface interface {
	GetAll(params *socDtos.CSLogSocRoomAccessLogListQueryParamsDTO) (*socDtos.CSLogSocRoomAccessLogListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*socDtos.CSLogSocRoomAccessLogDetailDTO, error)
	GetByID(id uint) (*models.CSLogSocRoomAccessLog, error)
	Create(currentUserID uint, form *socDtos.CSLogSocRoomAccessLogFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *socDtos.CSLogSocRoomAccessLogFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(id uint, token string) (string, error)
	ExportCSLogSocRoomAccessLog(data []*socDtos.CSLogSocRoomAccessLogListDTO, currentUserID uint) (*os.File, error)
}

type CSLogSocRoomAccessLogService struct {
	repo                repositories.CSLogSocRoomAccessLogRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface
	userSvc             organizationServices.UserServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewCSLogSocRoomAccessLogService(
	repo repositories.CSLogSocRoomAccessLogRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) CSLogSocRoomAccessLogServiceInterface {
	return &CSLogSocRoomAccessLogService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		organizationUnitSvc: organizationUnitSvc,
		userSvc:             userSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (csrals *CSLogSocRoomAccessLogService) GetAll(params *socDtos.CSLogSocRoomAccessLogListQueryParamsDTO) (*socDtos.CSLogSocRoomAccessLogListResponseDTO, error) {
	return csrals.repo.GetAll(params)
}

func (csrals *CSLogSocRoomAccessLogService) GetByID(id uint) (*models.CSLogSocRoomAccessLog, error) {
	entry, err := csrals.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return entry, nil
}

func (csrals *CSLogSocRoomAccessLogService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*socDtos.CSLogSocRoomAccessLogDetailDTO, error) {
	var entryDTO socDtos.CSLogSocRoomAccessLogDetailDTO

	entry, err := csrals.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(entry, &entryDTO)
	if err != nil {
		return nil, err
	}

	entryDTO.CreatedAt = entry.CreatedAt
	entryDTO.UpdatedAt = entry.UpdatedAt

	return &entryDTO, nil
}

func (csrals *CSLogSocRoomAccessLogService) Create(currentUserID uint, form *socDtos.CSLogSocRoomAccessLogFormDTO) (uint, error) {
	tx := csrals.repo.BeginTransaction()
	if err := csrals.validateCSLogSocRoomAccessLogFormDTO(form); err != nil {
		tx.Rollback()
		return 0, err
	}

	parsedDates, err := csrals.parseCSLogSocRoomAccessLogDates(form)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	referenceNumber, err := csrals.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	entry := &models.CSLogSocRoomAccessLog{
		ReferenceNumber:      referenceNumber,
		VisitorFullName:      form.VisitorFullName,
		EntryTime:            parsedDates.EntryTime,
		ExitTime:             parsedDates.ExitTime,
		VisitorIDNumber:      form.VisitorIDNumber,
		DepartmentID:         form.DepartmentID,
		ExternalPartyName:    form.ExternalPartyName,
		PurposeOfEntry:       form.PurposeOfEntry,
		AuthorizingOfficerID: form.AuthorizingOfficerID,
		Notes:                form.Notes,
		CreatedBYID:          currentUserID,
	}

	id, err := csrals.repo.Create(entry, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.CSLogSocRoomAccessLogDocument
	for _, doc := range form.Documents {
		// Encrypt the filename
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("security-operations", "soc-room-access-log", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.CSLogSocRoomAccessLogDocument{
			SocRoomAccessLogID: id,
			URL:                &url,
			Name:               &doc.FileName,
			Size:               &doc.Size,
			Token:              &encryptedFilename,
			OriginalName:       &doc.OriginalName,
			Type:               &doc.Type,
			CreatedBYID:        &currentUserID,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := csrals.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		storagePath := filepath.Join("security-operations", "soc-room-access-log", strconv.FormatUint(uint64(id), 10), "documents")
		for _, doc := range documents {
			if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}

	go func() {
		existingEntry, err := csrals.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch soc room access log", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = csrals.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogSocRoomAccessLog{}.TableName(),
			id,
			existingEntry,
			currentUserID,
			"security_operations", "security_operations_soc_room_access_log", "create soc room access log", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.CSLogSocRoomAccessLog{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (csrals *CSLogSocRoomAccessLogService) Update(id uint, currentUserID uint, form *socDtos.CSLogSocRoomAccessLogFormDTO) error {
	tx := csrals.repo.BeginTransaction()
	existingEntry, err := csrals.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	oldEntry, err := csrals.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := csrals.validateCSLogSocRoomAccessLogFormDTO(form); err != nil {
		tx.Rollback()
		return err
	}
	parsedDates, err := csrals.parseCSLogSocRoomAccessLogDates(form)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingEntry.VisitorFullName = form.VisitorFullName
	existingEntry.EntryTime = parsedDates.EntryTime
	existingEntry.ExitTime = parsedDates.ExitTime
	existingEntry.VisitorIDNumber = form.VisitorIDNumber
	existingEntry.DepartmentID = form.DepartmentID
	existingEntry.ExternalPartyName = form.ExternalPartyName
	existingEntry.PurposeOfEntry = form.PurposeOfEntry
	existingEntry.AuthorizingOfficerID = form.AuthorizingOfficerID
	existingEntry.Notes = form.Notes
	now := time.Now()
	existingEntry.UpdatedAt = &now
	existingEntry.UpdatedBYID = &currentUserID

	if err := csrals.repo.Update(existingEntry, tx); err != nil {
		tx.Rollback()
		return err
	}

	// documents update
	var documents []*models.CSLogSocRoomAccessLogDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.CSLogSocRoomAccessLogDocument

	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("security-operations", "soc-room-access-log", strconv.FormatUint(uint64(id), 10), "documents")

			document := &models.CSLogSocRoomAccessLogDocument{
				SocRoomAccessLogID: id,
				URL:                &url,
				Name:               &doc.FileName,
				Size:               &doc.Size,
				Token:              &encryptedFilename,
				OriginalName:       &doc.OriginalName,
				Type:               &doc.Type,
				CreatedBYID:        &currentUserID,
			}
			documents = append(documents, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := csrals.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					tx.Rollback()
					return err
				}
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
				documentsToDelete = append(documentsToDelete, existingDocument)
			}
		}
	}

	if len(documents) > 0 {
		if err := csrals.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := csrals.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	storagePath := filepath.Join("security-operations", "soc-room-access-log", strconv.FormatUint(uint64(id), 10), "documents")
	if len(documents) > 0 {
		for _, doc := range documents {
			if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(documentsToDelete) > 0 {
		for _, doc := range documentsToDelete {
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}

	go func() {
		updatedEntry, err := csrals.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch soc room access log", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		changedData := make(map[string]interface{})

		if len(documents) > 0 || len(documentsToDelete) > 0 {
			documentsAdded := []string{}
			for _, doc := range documents {
				documentsAdded = append(documentsAdded, *doc.Name)
			}
			documentsRemoved := []string{}
			for _, doc := range documentsToDelete {
				documentsRemoved = append(documentsRemoved, *doc.Name)
			}
			changedData["documents"] = map[string]interface{}{
				"added":   documentsAdded,
				"removed": documentsRemoved,
			}
		}

		oldData := map[string]interface{}{
			"visitor_full_name":   oldEntry.VisitorFullName,
			"entry_time":          oldEntry.EntryTime.Format("2006-01-02 15:04:05"),
			"exit_time":           oldEntry.ExitTime.Format("2006-01-02 15:04:05"),
			"visitor_id_number":   oldEntry.VisitorIDNumber,
			"department":          utility.GetFieldValue(oldEntry.Department),
			"external_party_name": oldEntry.ExternalPartyName,
			"purpose_of_entry":    oldEntry.PurposeOfEntry,
			"authorizing_officer": utility.GetFieldValue(oldEntry.AuthorizingOfficer),
			"notes":               oldEntry.Notes,
		}
		newData := map[string]interface{}{
			"visitor_full_name":   updatedEntry.VisitorFullName,
			"entry_time":          updatedEntry.EntryTime.Format("2006-01-02 15:04:05"),
			"exit_time":           updatedEntry.ExitTime.Format("2006-01-02 15:04:05"),
			"visitor_id_number":   updatedEntry.VisitorIDNumber,
			"department":          utility.GetFieldValue(updatedEntry.Department),
			"external_party_name": updatedEntry.ExternalPartyName,
			"purpose_of_entry":    updatedEntry.PurposeOfEntry,
			"authorizing_officer": utility.GetFieldValue(updatedEntry.AuthorizingOfficer),
			"notes":               updatedEntry.Notes,
		}
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := csrals.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.CSLogSocRoomAccessLog{}.TableName(),
				id,
				changedData,
				currentUserID,
				"security_operations", "security_operations_soc_room_access_log", "update soc room access log", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.CSLogSocRoomAccessLog{}.TableName(),
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

func (csrals *CSLogSocRoomAccessLogService) validateCSLogSocRoomAccessLogFormDTO(form *socDtos.CSLogSocRoomAccessLogFormDTO) error {
	if _, err := csrals.organizationUnitSvc.GetByID(form.DepartmentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("department does not exist"), "not_exist", "department")
	}
	if _, err := csrals.userSvc.GetByID(form.AuthorizingOfficerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("authorizing officer does not exist"), "not_exist", "authorizing officer")
	}
	return nil
}

func (csrals *CSLogSocRoomAccessLogService) parseCSLogSocRoomAccessLogDates(form *socDtos.CSLogSocRoomAccessLogFormDTO) (*socDtos.CSLogSocRoomAccessLogParsedDates, error) {
	parsed := &socDtos.CSLogSocRoomAccessLogParsedDates{}

	entryTime, err := utility.ParseDateTime(form.EntryTime)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid entry time"), "invalid_date", "entry time")
	}
	parsed.EntryTime = entryTime

	exitTime, err := utility.ParseDateTime(form.ExitTime)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid exit time"), "invalid_date", "exit time")
	}
	parsed.ExitTime = exitTime

	if entryTime != nil && exitTime != nil && !exitTime.After(*entryTime) {
		return nil, utility.NewValidationError(fmt.Errorf("exit time must be after entry time"), "invalid_date", "exit time")
	}

	return parsed, nil
}

func (csrals *CSLogSocRoomAccessLogService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedEntries []map[string]interface{}

	for _, id := range ids {
		tx := csrals.repo.BeginTransaction()
		entry, err := csrals.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if len(entry.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range entry.Documents {
				if doc != nil && doc.ID != nil {
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := csrals.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				failedIds = append(failedIds, id)
				continue
			}
		}

		if err := csrals.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if len(entry.Documents) > 0 {
			storagePath := filepath.Join("security-operations", "soc-room-access-log", strconv.FormatUint(uint64(id), 10), "documents")
			for _, doc := range entry.Documents {
				if doc != nil {
					if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
						utility.LogError("Failed to delete document file", map[string]interface{}{
							"recordID": id,
							"error":    err.Error(),
						})
					}
				}
			}
		}

		successfulIds = append(successfulIds, id)
		deletedEntries = append(deletedEntries, map[string]interface{}{
			"id":                id,
			"reference_number":  entry.ReferenceNumber,
			"visitor_full_name": entry.VisitorFullName,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := csrals.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CSLogSocRoomAccessLog{}.TableName(),
				successfulIds,
				deletedEntries,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (csrals *CSLogSocRoomAccessLogService) GetDocumentPath(id uint, token string) (string, error) {
	existingEntry, err := csrals.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	if len(existingEntry.Documents) == 0 {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("no documents found for the soc room access log"), id)
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("resource not found"), token)
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("resource not found"), token)
	}

	storagePath := filepath.Join("security-operations", "soc-room-access-log", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("file not found"), token)
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (csrals *CSLogSocRoomAccessLogService) ExportCSLogSocRoomAccessLog(data []*socDtos.CSLogSocRoomAccessLogListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := csrals.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Visitor Name",
			"Entry Time",
			"Exit Time",
			"Visitor ID Number",
			"Department",
			"External Party",
			"Authorizing Officer",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Visitor Name",
			"Entry Time",
			"Exit Time",
			"Visitor ID Number",
			"Department",
			"External Party",
			"Authorizing Officer",
		}
	}
	// Convert the soc room access log data to [][]interface{} based on language setting
	rows := make([][]interface{}, len(data))
	for i, entry := range data {
		if userSetting.Language == "ar" {
			// Populate row data in Arabic
			rows[i] = []interface{}{
				i + 1,
				entry.ReferenceNumber,
				entry.VisitorFullName,
				entry.EntryTime.Format("2006-01-02 15:04:05"),
				entry.ExitTime.Format("2006-01-02 15:04:05"),
				entry.VisitorIDNumber,
				entry.DepartmentTitleAr,
				entry.ExternalPartyName,
				entry.AuthorizingOfficerFirstNameAr + " " + entry.AuthorizingOfficerLastNameAr,
			}
		} else {
			// Populate row data in English
			rows[i] = []interface{}{
				i + 1,
				entry.ReferenceNumber,
				entry.VisitorFullName,
				entry.EntryTime.Format("2006-01-02 15:04:05"),
				entry.ExitTime.Format("2006-01-02 15:04:05"),
				entry.VisitorIDNumber,
				entry.DepartmentTitleEn,
				entry.ExternalPartyName,
				entry.AuthorizingOfficerFirstNameEn + " " + entry.AuthorizingOfficerLastNameEn,
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
