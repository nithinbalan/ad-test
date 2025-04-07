package services

import (
	"errors"
	"fmt"
	"jethur/config"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	StakeholderActionPlanRepositories "jethur/grc/organizational-context/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StakeholderActionPlanServiceInterface interface {
	GetAll(stakeholderId uint, params *organizationalContextDTOs.StakeholderActionPlanQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderActionPlanListDTO, error)
	GetByID(stakeholderID uint, id uint) (*models.StakeholderActionPlan, error)
	GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderActionPlanDetailDTO, error)
	Create(stakeholderID uint, currentUserID uint, dto *organizationalContextDTOs.StakeholderActionPlanFormDTO) (uint, error)
	Update(stakeholderID uint, id uint, currentUserID uint, dto *organizationalContextDTOs.StakeholderActionPlanFormDTO) error
	Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(stakeholderID uint, actionPlanID uint, token string) (string, error)
}

type StakeholderActionPlanService struct {
	repo                StakeholderActionPlanRepositories.StakeholderActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	stakeholderSvc      StakeholderServiceInterface
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
}

func NewStakeholderActionPlanService(
	repo StakeholderActionPlanRepositories.StakeholderActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	stakeholderSvc StakeholderServiceInterface,
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) StakeholderActionPlanServiceInterface {
	return &StakeholderActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		stakeholderSvc:      stakeholderSvc,
		ActionPlanStatusSvc: ActionPlanStatusSvc,
		userSettingSvc:      userSettingSvc,
	}
}

func (saps *StakeholderActionPlanService) GetAll(stakeholderId uint, params *organizationalContextDTOs.StakeholderActionPlanQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderActionPlanListDTO, error) {
	if _, err := saps.stakeholderSvc.GetByID(stakeholderId); err != nil {
		return nil, err
	}
	return saps.repo.GetAll(stakeholderId, params)
}

func (saps *StakeholderActionPlanService) GetByID(stakeholderID uint, id uint) (*models.StakeholderActionPlan, error) {
	actionPlan, err := saps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (saps *StakeholderActionPlanService) GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderActionPlanDetailDTO, error) {
	if _, err := saps.stakeholderSvc.GetByID(stakeholderID); err != nil {
		return nil, err
	}
	var recoveryActionPlan organizationalContextDTOs.StakeholderActionPlanDetailDTO

	actionPlan, err := saps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(actionPlan, &recoveryActionPlan)
	if err != nil {
		return nil, err
	}

	if actionPlan.TargetDate != nil {
		targetDate := actionPlan.TargetDate.Format("2006-01-02")
		recoveryActionPlan.TargetDateString = &targetDate
	}

	recoveryActionPlan.CreatedAt = actionPlan.CreatedAt
	recoveryActionPlan.UpdatedAt = actionPlan.UpdatedAt

	return &recoveryActionPlan, nil
}

func (saps *StakeholderActionPlanService) Create(stakeholderID uint, currentUserID uint, stakeholderActionPlanFormDTO *organizationalContextDTOs.StakeholderActionPlanFormDTO) (uint, error) {
	tx := saps.repo.BeginTransaction()

	if _, err := saps.stakeholderSvc.GetByID(stakeholderID); err != nil {
		return 0, err
	}

	responsibleUnitID, err := saps.userSvc.GetUserPrimaryUnitID(*stakeholderActionPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(errors.New("organization unit not set for responsible user"), "not_exist", "responsible unit")
	}

	actionPlanStatus, err := saps.ActionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	dates, err := saps.parseStakeholderActionPlanDates(stakeholderActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlan := &models.StakeholderActionPlan{
		StakeholderID:      stakeholderID,
		SeverityLevelID:    stakeholderActionPlanFormDTO.SeverityLevelID,
		Title:              stakeholderActionPlanFormDTO.Title,
		Requirement:        stakeholderActionPlanFormDTO.Requirement,
		Description:        stakeholderActionPlanFormDTO.Description,
		ResponsibleUserID:  stakeholderActionPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:  &responsibleUnitID,
		AssignedOwnerID:    stakeholderActionPlanFormDTO.AssignedOwnerID,
		TargetDate:         dates.TargetDate,
		ActionPlanStatusID: &actionPlanStatus.ID,
		CreatedBYID:        currentUserID,
	}

	id, err := saps.repo.Create(actionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.StakeholderActionPlanDocument
	for _, doc := range stakeholderActionPlanFormDTO.Documents {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}
		url := utility.GenerateDocumentURL("organization-context", "stakeholder", strconv.FormatUint(uint64(stakeholderID), 10), "actionPlan", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.StakeholderActionPlanDocument{
			StakeholderActionPlanID: id,
			URL:                     url,
			Name:                    doc.FileName,
			Size:                    doc.Size,
			Token:                   encryptedFileName,
			OriginalName:            doc.OriginalName,
			Type:                    doc.Type,
			CreatedBYID:             currentUserID,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := saps.repo.CreateDocuments(documents, tx); err != nil {
			return 0, err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("organization-context", "stakeholder", strconv.FormatUint(uint64(stakeholderID), 10), "actionPlan", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	tx.Commit()
	go func() {
		existingStakeholderActionPlan, err := saps.GetFullDetailsByID(stakeholderID, actionPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch stakeholder", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = saps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderActionPlan{}.TableName(),
			id,
			existingStakeholderActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (saps *StakeholderActionPlanService) Update(stakeholderID uint, id uint, currentUserID uint, stakeholderActionPlanFormDTO *organizationalContextDTOs.StakeholderActionPlanFormDTO) error {
	tx := saps.repo.BeginTransaction()

	existingStakeholderActionPlan, err := saps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldActionPlan, err := saps.GetFullDetailsByID(stakeholderID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingStakeholderActionPlan.ResponsibleUserID != stakeholderActionPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := saps.userSvc.GetUserPrimaryUnitID(*existingStakeholderActionPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingStakeholderActionPlan.ResponsibleUnitID = &responsibleUnitID
	}

	dates, err := saps.parseStakeholderActionPlanDates(stakeholderActionPlanFormDTO)
	if err != nil {
		return err
	}

	existingStakeholderActionPlan.SeverityLevelID = stakeholderActionPlanFormDTO.SeverityLevelID
	existingStakeholderActionPlan.Title = stakeholderActionPlanFormDTO.Title
	existingStakeholderActionPlan.ResponsibleUserID = stakeholderActionPlanFormDTO.ResponsibleUserID
	existingStakeholderActionPlan.TargetDate = dates.TargetDate
	existingStakeholderActionPlan.Requirement = stakeholderActionPlanFormDTO.Requirement
	existingStakeholderActionPlan.Description = stakeholderActionPlanFormDTO.Description
	existingStakeholderActionPlan.AssignedOwnerID = stakeholderActionPlanFormDTO.AssignedOwnerID
	existingStakeholderActionPlan.UpdatedBYID = &currentUserID

	err = saps.repo.Update(existingStakeholderActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle documents
	var documents []*models.StakeholderActionPlanDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.StakeholderActionPlanDocument

	for _, doc := range stakeholderActionPlanFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate URL
			url := utility.GenerateDocumentURL("organization-context", "stakeholder", strconv.FormatUint(uint64(stakeholderID), 10), "actionPlan", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

			// Create a new actionPlan Document instance
			document := &models.StakeholderActionPlanDocument{
				StakeholderActionPlanID: id,
				URL:                     url,
				Name:                    doc.FileName,
				Size:                    doc.Size,
				Token:                   encryptedFileName,
				OriginalName:            doc.OriginalName,
				Type:                    doc.Type,
			}

			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := saps.repo.GetDocumentByID(*doc.ID)
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

	// Create new documents if any
	if len(documents) > 0 {
		if err := saps.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete documents if any
	if len(deleteDocumentIds) > 0 {
		if err := saps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	// Store new files to storage
	for _, doc := range documents {
		storagePath := filepath.Join("organization-context", "stakeholder", strconv.FormatUint(uint64(stakeholderID), 10), "actionPlan", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to store file %s: %v\n", doc.Name, err)
		}
	}

	// Remove deleted files from storage
	for _, document := range documentsToDelete {
		storagePath := filepath.Join("organization-context", "stakeholder", strconv.FormatUint(uint64(stakeholderID), 10), "actionPlan", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", document.Name, err)
		}
	}

	go func() {
		existingActionPlan, err := saps.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"title":            oldActionPlan.Title,
			"responsible_user": utility.GetFieldValue(oldActionPlan.ResponsibleUser),
			"assigned_owner":   utility.GetFieldValue(oldActionPlan.AssignedOwner),
			"requirement":      oldActionPlan.Requirement,
			"description":      oldActionPlan.Description,
			"severity_level":   utility.GetFieldValue(oldActionPlan.SeverityLevel),
			"target_date":      oldActionPlan.TargetDateString,
		}

		newData := map[string]interface{}{
			"title":            existingActionPlan.Title,
			"responsible_user": utility.GetFieldValue(existingActionPlan.ResponsibleUser),
			"assigned_owner":   utility.GetFieldValue(existingActionPlan.AssignedOwner),
			"requirement":      existingActionPlan.Requirement,
			"description":      existingActionPlan.Description,
			"severity_level":   utility.GetFieldValue(existingActionPlan.SeverityLevel),
			"target_date":      existingActionPlan.TargetDateString,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := saps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderActionPlan{}.TableName(),
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

func (saps *StakeholderActionPlanService) Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := saps.repo.BeginTransaction()

	if _, err := saps.stakeholderSvc.GetByID(stakeholderID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedActionPlans []map[string]interface{}
	documentsToDelete := make(map[uint][]*organizationalContextDTOs.StakeholderActionPlanDocumentDTO)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		actionPlan, err := saps.GetFullDetailsByID(stakeholderID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedActionPlans = append(deletedActionPlans, map[string]interface{}{
			"id":          id,
			"stakeholder": actionPlan.Title,
		})

		if actionPlan.Documents != nil && len(actionPlan.Documents) > 0 {
			for _, document := range actionPlan.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], document)
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}

			if len(deleteDocumentIds) > 0 {
				if err := saps.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
					failedIDs = append(failedIDs, id)
					continue
				}
			}
		}

		if err := saps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	if len(successfulIDs) > 0 {
		go func() {
			err := saps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.StakeholderActionPlan{}.TableName(),
				successfulIDs,
				deletedActionPlans,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}

func (saps *StakeholderActionPlanService) parseStakeholderActionPlanDates(dto *organizationalContextDTOs.StakeholderActionPlanFormDTO) (*organizationalContextDTOs.StakeholderActionPlanParsedDates, error) {
	parsed := &organizationalContextDTOs.StakeholderActionPlanParsedDates{}

	var err error

	if &dto.TargetDate != nil {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}

func (saps *StakeholderActionPlanService) GetDocumentPath(stakeholderID uint, actionPlanID uint, token string) (string, error) {
	statusUpdate, err := saps.GetFullDetailsByID(stakeholderID, actionPlanID)
	if err != nil {
		fmt.Printf("Error fetching status update by ID: %v\n", err)
		return "", err
	}
	if len(statusUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the stakeholder action plan")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil || len(decryptedFilenameBytes) == 0 {
		fmt.Printf("Token decryption failed: %v\n", err)
		return "", fmt.Errorf("resource not found")
	}
	decryptedFilename := string(decryptedFilenameBytes)
	fmt.Printf("Decrypted filename: %s\n", decryptedFilename)

	storagePath := filepath.Join("organization-context", "stakeholder", strconv.FormatUint(uint64(stakeholderID), 10), "actionPlan", strconv.FormatUint(uint64(actionPlanID), 10), "documents")

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
