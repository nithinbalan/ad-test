package services

import (
	"errors"
	"fmt"
	"jethur/config"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SLARegisterServiceInterface interface {
	GetAll(params *complianceDTOs.SLARegisterListQueryParamsDTO) (*complianceDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.SLARegisterDetailDTO, error)
	Create(currentUserID uint, slaRegisterFormDTO *complianceDTOs.SLARegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, sla *complianceDTOs.SLARegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(id uint, token string) (string, error)
}

// Service and implementations
type SLARegisterService struct {
	repo               complianceRepositories.SLARegisterRepositoryInterface
	slaCategorySvc     SLACategoryServiceInterface
	severityLevelSvc   masterServices.SeverityLevelServiceInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	processSvc         organizationServices.ProcessServiceInterface
	userSvc            organizationServices.UserServiceInterface
}

func NewSLARegisterService(
	repo complianceRepositories.SLARegisterRepositoryInterface,
	slaCategorySvc SLACategoryServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	userSvc organizationServices.UserServiceInterface) SLARegisterServiceInterface {
	return &SLARegisterService{
		repo:               repo,
		slaCategorySvc:     slaCategorySvc,
		severityLevelSvc:   severityLevelSvc,
		lifecycleStatusSvc: lifecycleStatusSvc,
		activityLogSvc:     activityLogSvc,
		processSvc:         processSvc,
		userSvc:            userSvc,
	}
}

func (srs *SLARegisterService) GetAll(params *complianceDTOs.SLARegisterListQueryParamsDTO) (*complianceDTOs.SLARegisterListResponseDTO, error) {
	return srs.repo.GetAll(params)
}

func (srs SLARegisterService) GetByID(id uint) (*models.SLARegister, error) {
	sla, err := srs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return sla, nil
}

func (srs SLARegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.SLARegisterDetailDTO, error) {
	// Initialize your DTO
	var slaRegisterDTO complianceDTOs.SLARegisterDetailDTO

	slaRegister, err := srs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(slaRegister, &slaRegisterDTO)
	if err != nil {
		return nil, err
	}

	for _, document := range slaRegisterDTO.Documents {
		if document.IssueDate != nil {
			issueDate := document.IssueDate.Format("2006-01-02")
			document.IssueDateString = &issueDate
		}
		if document.ExpiryDate != nil {
			expiryDate := document.ExpiryDate.Format("2006-01-02")
			document.ExpiryDateString = &expiryDate
		}
		if document.SA1Date != nil {
			sa1Date := document.SA1Date.Format("2006-01-02")
			document.SA1DateString = &sa1Date
		}
		if document.SA2Date != nil {
			sa2Date := document.SA2Date.Format("2006-01-02")
			document.SA2DateString = &sa2Date
		}
	}

	slaRegisterDTO.CreatedAt = slaRegister.CreatedAt
	slaRegisterDTO.UpdatedAt = slaRegister.UpdatedAt

	return &slaRegisterDTO, nil
}

func (srs *SLARegisterService) Create(currentUserID uint, slaRegisterFormDTO *complianceDTOs.SLARegisterFormDTO) (uint, error) {
	tx := srs.repo.BeginTransaction()

	if err := srs.validateSLARegisterFormDTO(slaRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := srs.parseSLARegisterDates(slaRegisterFormDTO)
	if err != nil {
		return 0, err
	}
	referenceNumber, err := srs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	responsibleUnitID, err := srs.userSvc.GetUserPrimaryUnitID(uint(slaRegisterFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
	}
	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	if err := srs.validateSLADocumentFormDTO(&slaRegisterFormDTO.Document); err != nil {
		return 0, err
	}

	slaRegister := &models.SLARegister{
		SLACategoryID:     slaRegisterFormDTO.SLACategoryID,
		SeverityLevelID:   slaRegisterFormDTO.SeverityLevelID,
		ResponsibleUnitID: responsibleUnitID,
		ResponsibleUserID: slaRegisterFormDTO.ResponsibleUserID,
		LifecycleStatusID: lifecycleStatus.ID,
		ReferenceNumber:   referenceNumber,
		TitleEn:           slaRegisterFormDTO.TitleEn,
		TitleAr:           slaRegisterFormDTO.TitleAr,
		Description:       slaRegisterFormDTO.Description,
		Client:            slaRegisterFormDTO.Client,
		CreatedByID:       currentUserID,
	}

	id, err := srs.repo.Create(slaRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save Processes
	var processes []*models.SLAProcess
	for _, processID := range slaRegisterFormDTO.ProcessIDs {
		process := &models.SLAProcess{
			SLARegisterID: id,
			ProcessID:     *processID,
			CreatedByID:   currentUserID,
		}
		processes = append(processes, process)
	}

	if len(processes) > 0 {
		if err := srs.repo.CreateSLAProcesses(processes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	slaCheck := &models.SLACheck{
		SLARegisterID: id,
		CreatedByID:   currentUserID,
	}

	if err := srs.repo.CreateSLACheck(slaCheck, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Encrypt the filename
	dataToEncrypt := []byte(slaRegisterFormDTO.Document.FileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return 0, err
	}
	url := utility.GenerateDocumentURL("compliance", "sla-register", strconv.FormatUint(uint64(slaRegister.ID), 10), "agreements")

	slaDocument := &models.SLADocument{
		SLARegisterID: id,
		Version:       slaRegisterFormDTO.Version,
		IssueDate:     dates.IssueDate,
		ExpiryDate:    dates.ExpiryDate,
		SA1Date:       dates.SA1Date,
		SA2Date:       dates.SA2Date,
		URL:           url,
		Name:          slaRegisterFormDTO.Document.FileName,
		Size:          slaRegisterFormDTO.Document.Size,
		Token:         encryptedFilename,
		OriginalName:  slaRegisterFormDTO.Document.OriginalName,
		Type:          slaRegisterFormDTO.Document.Type,
		IsLatest:      true,
		CreatedByID:   currentUserID,
	}

	// Save SLA document
	if err := srs.repo.CreateDocument(slaDocument, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Store SLA Document File using utility.StoreFile
	storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(id), 10), "agreements")

	if err := utility.StoreFile(slaRegisterFormDTO.Document.FileName, storagePath); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := srs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SLARegister{}.TableName(),
			id,
			slaRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SLARegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (srs *SLARegisterService) Update(id uint, currentUserID uint, slaRegisterFormDTO *complianceDTOs.SLARegisterFormDTO) error {
	tx := srs.repo.BeginTransaction()
	existingSLARegister, err := srs.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := srs.validateSLARegisterFormDTO(slaRegisterFormDTO); err != nil {
		return err
	}

	if err := srs.validateSLADocumentFormDTO(&slaRegisterFormDTO.Document); err != nil {
		return err
	}

	dates, err := srs.parseSLARegisterDates(slaRegisterFormDTO)
	if err != nil {
		return err
	}

	if existingSLARegister.ResponsibleUserID != slaRegisterFormDTO.ResponsibleUserID {
		responsibleUnitID, err := srs.userSvc.GetUserPrimaryUnitID(uint(existingSLARegister.ResponsibleUserID))
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingSLARegister.ResponsibleUnitID = responsibleUnitID
	}

	existingSLARegister.SLACategoryID = slaRegisterFormDTO.SLACategoryID
	existingSLARegister.ResponsibleUserID = slaRegisterFormDTO.ResponsibleUserID
	existingSLARegister.SeverityLevelID = slaRegisterFormDTO.SeverityLevelID
	existingSLARegister.LifecycleStatusID = slaRegisterFormDTO.LifecycleStatusID
	existingSLARegister.TitleEn = slaRegisterFormDTO.TitleEn
	existingSLARegister.TitleAr = slaRegisterFormDTO.TitleAr
	existingSLARegister.Description = slaRegisterFormDTO.Description
	existingSLARegister.Client = slaRegisterFormDTO.Client
	existingSLARegister.UpdatedByID = currentUserID

	if err := srs.repo.Update(existingSLARegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingSLA, err := srs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingProcessIDs := make(map[uint]bool)
	for _, proc := range existingSLA.Processes {
		existingProcessIDs[proc.ProcessID] = true
	}

	// Collect new process IDs
	newProcessIDs := make(map[uint]bool)
	for _, proc := range slaRegisterFormDTO.ProcessIDs {
		newProcessIDs[*proc] = true
	}

	// Add new processes
	var processesToAdd []*models.SLAProcess
	for _, newProcessID := range slaRegisterFormDTO.ProcessIDs {
		if _, exists := existingProcessIDs[*newProcessID]; !exists {
			// Create a new sla process record
			process := &models.SLAProcess{
				SLARegisterID: id,
				ProcessID:     *newProcessID,
				CreatedByID:   currentUserID,
			}
			processesToAdd = append(processesToAdd, process)
		}
	}

	// Save the new processes
	if len(processesToAdd) > 0 {
		if err := srs.repo.CreateSLAProcesses(processesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove processes not in the new list
	var deleteProcessIDs []*uint

	for _, existingProcess := range existingSLA.Processes {
		if _, exists := newProcessIDs[existingProcess.ProcessID]; !exists {
			// Append process to delete list
			deleteProcessIDs = append(deleteProcessIDs, existingProcess.ID)
		}
	}

	// Delete processes not in the new list
	if len(deleteProcessIDs) > 0 {
		if err := srs.repo.DeleteSLAProcesses(deleteProcessIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingSLADocument, err := srs.repo.GetLatestDocumentBySLAID(id)
	if err != nil {
		return err
	}
	existingSLADocument.Version = slaRegisterFormDTO.Version
	existingSLADocument.IssueDate = dates.IssueDate
	existingSLADocument.ExpiryDate = dates.ExpiryDate
	existingSLADocument.SA1Date = dates.SA1Date
	existingSLADocument.SA2Date = dates.SA2Date
	existingSLADocument.UpdatedByID = currentUserID
	existingSLADocument.IsLatest = true

	if slaRegisterFormDTO.Document.IsNewFile != nil && *slaRegisterFormDTO.Document.IsNewFile {
		existingSLADocumentFileName := existingSLADocument.Name
		//Encrypt the filename
		dataToEncrypt := []byte(slaRegisterFormDTO.Document.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}
		url := utility.GenerateDocumentURL("compliance", "sla-register", strconv.FormatUint(uint64(id), 10), "agreements")

		existingSLADocument.URL = url
		existingSLADocument.Name = slaRegisterFormDTO.Document.FileName
		existingSLADocument.Size = slaRegisterFormDTO.Document.Size
		existingSLADocument.Token = encryptedFilename
		existingSLADocument.OriginalName = slaRegisterFormDTO.Document.OriginalName
		existingSLADocument.Type = slaRegisterFormDTO.Document.Type
		if err := srs.repo.CreateDocument(existingSLADocument, tx); err != nil {
			tx.Rollback()
			return err
		}
		storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(existingSLADocument.SLARegisterID), 10), "agreements")
		if err := utility.RemoveFile(existingSLADocumentFileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}
		// Store SLA Document File using utility.StoreFile
		if err := utility.StoreFile(slaRegisterFormDTO.Document.FileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

	} else {
		if err := srs.repo.CreateDocument(existingSLADocument, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := srs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SLARegister{}.TableName(),
			id,
			slaRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.SLARegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (srs *SLARegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := srs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*complianceDTOs.SLADocumentDTO)

	for _, id := range ids {
		slaRegister, err := srs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Check if controls exist
		if len(slaRegister.Processes) > 0 {
			var deleteProcessIDs []*uint
			for _, process := range slaRegister.Processes {
				deleteProcessIDs = append(deleteProcessIDs, process.ID)
			}
			if err := srs.repo.DeleteSLAProcesses(deleteProcessIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		// Delete each SLA document if documents exist
		if len(slaRegister.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range slaRegister.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := srs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		// Delete the SLA Checks
		if err := srs.repo.DeleteSLAChecksBySLAID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		// Delete the SLA Register
		if err := srs.repo.Delete(id, tx); err != nil {
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
		for slaID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(slaID), 10), "agreements")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SLARegister{}.TableName(),
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

func (srs *SLARegisterService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return nil, err
	}

	idsToUpdate, err := srs.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := srs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.SLARegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (srs *SLARegisterService) Restore(ids []uint, currentUserID uint) ([]uint, error) {

	lifecycleStatus, err := srs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := srs.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := srs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				models.RestoreAction,
				models.SLARegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (srs *SLARegisterService) validateSLARegisterFormDTO(slaRegisterFormDTO *complianceDTOs.SLARegisterFormDTO) error {
	if _, err := srs.slaCategorySvc.GetByID(slaRegisterFormDTO.SLACategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("sla category does not exist"), "not_exist", "sla category")
	}

	if slaRegisterFormDTO.LifecycleStatusID != 0 {
		if _, err := srs.lifecycleStatusSvc.GetByID(slaRegisterFormDTO.LifecycleStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("status does not exist"), "not_exist", "lifecycle_status")
		}
	}
	if slaRegisterFormDTO.SeverityLevelID != 0 {
		if _, err := srs.severityLevelSvc.GetByID(slaRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}
	processIds := slaRegisterFormDTO.ProcessIDs

	processIDs := make(map[uint]bool)

	for _, processId := range processIds {
		if processIDs[*processId] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", processId), "duplicate", "process")
		}

		if _, err := srs.processSvc.GetByID(*processId); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processId), "not_exist", "process")
		}

		processIDs[*processId] = true
	}

	return nil
}

func (srs *SLARegisterService) parseSLARegisterDates(slaRegisterForm *complianceDTOs.SLARegisterFormDTO) (*complianceDTOs.ParsedSLADates, error) {
	parsed := &complianceDTOs.ParsedSLADates{}
	var err error

	// Validate issue_date
	if slaRegisterForm.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(slaRegisterForm.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	// Validate expiry_date
	if slaRegisterForm.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(slaRegisterForm.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	// Validate SA1Date
	if slaRegisterForm.SA1Date != "" {
		parsed.SA1Date, err = utility.ParseDate(slaRegisterForm.SA1Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa1 date"), "invalid_date", "sa1 date")
		}
	}

	// Validate SA2Date
	if slaRegisterForm.SA2Date != "" {
		parsed.SA2Date, err = utility.ParseDate(slaRegisterForm.SA2Date)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid sa2 date"), "invalid_date", "sa2 date")
		}
	}

	return parsed, nil
}

func (srs *SLARegisterService) validateSLADocumentFormDTO(document *complianceDTOs.SLADocumentFormDTO) error {

	//@Review
	//All items are mandatory during creation

	// Validate file type (PDF)
	if document.Type != ".pdf" {
		return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", "type")
	}

	// @todo
	// Validate file size (if needed)
	// For example:
	// if document.Size > MaxFileSize {
	//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
	// }

	return nil
}

func (srs *SLARegisterService) GetDocumentPath(id uint, token string) (string, error) {
	existingSLARegister, err := srs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	if len(existingSLARegister.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the sla register")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}
	storagePath := filepath.Join("compliance", "sla-register", strconv.FormatUint(uint64(id), 10), "agreements")

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
