package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessFunctionRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BusinessFunctionServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BusinessFunctionListQueryParamsDTO) (*businessContinuityDTOs.BusinessFunctionListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessFunctionDetailDTO, error)
	GetByID(id uint) (*models.BusinessFunction, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BusinessFunctionFormDTO) (uint, error)
	Update(id uint, currentUserID uint, businessFunction *businessContinuityDTOs.BusinessFunctionFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	UpdateBusinessFunction(businessFunction *models.BusinessFunction, tx utility.DatabaseTransaction) error
	ExportBusinessFunction(data []*businessContinuityDTOs.BusinessFunctionListDTO, currentUserID uint) (*os.File, error)

	GetProcessAssociatedRisks(businessFunctionID uint, processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemRiskListResponseDTO, error)
	GetProcessAssociatedIncidents(businessFunctionID uint, processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO, error)
	GetProcessAssociatedControls(businessFunctionID uint, processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemControlListResponseDTO, error)
}

// Service and implementations
type BusinessFunctionService struct {
	repo                     businessFunctionRepositories.BusinessFunctionRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	userSettingSvc           organizationServices.UserSettingServiceInterface
	processAssociatedItemSvc organizationServices.ProcessAssociatedItemServiceInterface
}

func NewBusinessFunctionService(
	repo businessFunctionRepositories.BusinessFunctionRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	processAssociatedItemSvc organizationServices.ProcessAssociatedItemServiceInterface,
) BusinessFunctionServiceInterface {
	return &BusinessFunctionService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		userSvc:                  userSvc,
		userSettingSvc:           userSettingSvc,
		processAssociatedItemSvc: processAssociatedItemSvc,
	}
}

func (bfs *BusinessFunctionService) GetAll(params *businessContinuityDTOs.BusinessFunctionListQueryParamsDTO) (*businessContinuityDTOs.BusinessFunctionListResponseDTO, error) {
	return bfs.repo.GetAll(params)
}

func (bfs BusinessFunctionService) GetByID(id uint) (*models.BusinessFunction, error) {
	businessFunction, err := bfs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return businessFunction, nil
}

func (bfs BusinessFunctionService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessFunctionDetailDTO, error) {
	// Initialize your DTO
	var businessFunctionDTO businessContinuityDTOs.BusinessFunctionDetailDTO

	businessFunction, err := bfs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(businessFunction, &businessFunctionDTO)
	if err != nil {
		return nil, err
	}

	if businessFunction.NextReviewDate != nil {
		nextReviewDate := businessFunction.NextReviewDate.Format("2006-01-02")
		businessFunctionDTO.NextReviewDateString = &nextReviewDate
	}
	if businessFunction.LastReviewDate != nil {
		lastReviewDate := businessFunction.LastReviewDate.Format("2006-01-02")
		businessFunctionDTO.LastReviewDateString = &lastReviewDate
	}

	businessFunctionDTO.CreatedAt = businessFunction.CreatedAt
	businessFunctionDTO.UpdatedAt = businessFunction.UpdatedAt

	return &businessFunctionDTO, nil
}

func (bfs *BusinessFunctionService) Create(currentUserID uint, businessFunctionFormDTO *businessContinuityDTOs.BusinessFunctionFormDTO) (uint, error) {
	tx := bfs.repo.BeginTransaction()

	referenceNumber, err := bfs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	dates, err := bfs.parseBusinessFunctionDates(businessFunctionFormDTO)
	if err != nil {
		return 0, err
	}

	businessFunction := &models.BusinessFunction{
		ReferenceNumber:   referenceNumber,
		SeverityLevelID:   businessFunctionFormDTO.SeverityLevelID,
		ResponsibleUserID: businessFunctionFormDTO.ResponsibleUserID,
		Title:             businessFunctionFormDTO.Title,
		Description:       businessFunctionFormDTO.Description,
		ReviewFrequencyID: businessFunctionFormDTO.ReviewFrequencyID,
		LastReviewDate:    dates.LastReviewDate,
		NextReviewDate:    dates.NextReviewDate,
		CreatedBYID:       currentUserID,
	}

	id, err := bfs.repo.Create(businessFunction, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Processes
	var processes []*models.BusinessFunctionProcess
	for _, processID := range businessFunctionFormDTO.ProcessIDs {
		process := &models.BusinessFunctionProcess{
			BusinessFunctionID: id,
			ProcessID:          *processID,
			CreatedBYID:        currentUserID,
		}
		processes = append(processes, process)
	}

	if len(processes) > 0 {
		if err := bfs.repo.CreateBusinessFunctionProcesses(processes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	go func() {
		existingBusinessFunction, err := bfs.GetFullDetailsByID(businessFunction.ID)
		if err != nil {
			utility.LogError("Failed to fetch business function", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = bfs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BusinessFunction{}.TableName(),
			id,
			existingBusinessFunction,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BusinessFunction{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (bfs *BusinessFunctionService) Update(id uint, currentUserID uint, businessFunctionFormDTO *businessContinuityDTOs.BusinessFunctionFormDTO) error {
	tx := bfs.repo.BeginTransaction()

	existingBusinessFunction, err := bfs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	dates, err := bfs.parseBusinessFunctionDates(businessFunctionFormDTO)
	if err != nil {
		return err
	}

	oldBusinessFunction, err := bfs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingBusinessFunction.SeverityLevelID = businessFunctionFormDTO.SeverityLevelID
	existingBusinessFunction.ResponsibleUserID = businessFunctionFormDTO.ResponsibleUserID
	existingBusinessFunction.Title = businessFunctionFormDTO.Title
	existingBusinessFunction.Description = businessFunctionFormDTO.Description
	existingBusinessFunction.ReviewFrequencyID = businessFunctionFormDTO.ReviewFrequencyID
	existingBusinessFunction.LastReviewDate = dates.LastReviewDate
	existingBusinessFunction.NextReviewDate = dates.NextReviewDate
	existingBusinessFunction.UpdatedBYID = &currentUserID

	err = bfs.repo.Update(existingBusinessFunction, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Collect existing process IDs
	existingFunction, err := bfs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingProcessIDs := make(map[uint]bool)
	for _, proc := range existingFunction.Processes {
		existingProcessIDs[proc.ProcessID] = true
	}

	// Collect new process IDs
	newProcessIDs := make(map[uint]bool)
	for _, proc := range businessFunctionFormDTO.ProcessIDs {
		newProcessIDs[*proc] = true
	}

	// Add new processes
	var processesToAdd []*models.BusinessFunctionProcess
	for _, newProcessID := range businessFunctionFormDTO.ProcessIDs {
		if _, exists := existingProcessIDs[*newProcessID]; !exists {
			// Create a new business function process record
			process := &models.BusinessFunctionProcess{
				BusinessFunctionID: id,
				ProcessID:          *newProcessID,
				CreatedBYID:        currentUserID,
			}
			processesToAdd = append(processesToAdd, process)
		}
	}

	// Save the new processes
	if len(processesToAdd) > 0 {
		if err := bfs.repo.CreateBusinessFunctionProcesses(processesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove processes not in the new list
	var deleteProcessIDs []*uint

	for _, existingProcess := range existingFunction.Processes {
		if _, exists := newProcessIDs[existingProcess.ProcessID]; !exists {
			// Append process to delete list
			deleteProcessIDs = append(deleteProcessIDs, existingProcess.ID)
		}
	}

	// Delete processes not in the new list
	if len(deleteProcessIDs) > 0 {
		if err := bfs.repo.DeleteBusinessFunctionProcesses(deleteProcessIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingBusinessFunction, err := bfs.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Determine process changes

		processesAdded := []uint{}
		processesRemoved := []uint{}
		for id := range newProcessIDs {
			if !existingProcessIDs[id] {
				processesAdded = append(processesAdded, id)
			}
		}
		for id := range existingProcessIDs {
			if !newProcessIDs[id] {
				processesRemoved = append(processesRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(processesAdded) > 0 || len(processesRemoved) > 0 {
			associatedItems["process_ids"] = map[string]interface{}{
				"added":   processesAdded,
				"removed": processesRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"severity_level":   oldBusinessFunction.SeverityLevel.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s", oldBusinessFunction.ResponsibleUser.FirstNameEN, oldBusinessFunction.ResponsibleUser.LastNameEN),
			"title":            oldBusinessFunction.Title,
			"description":      oldBusinessFunction.Description,
			"review_frequency": utility.GetFieldValue(oldBusinessFunction.ReviewFrequency),
			"last_review_date": oldBusinessFunction.LastReviewDateString,
			"next_review_date": oldBusinessFunction.NextReviewDateString,
		}

		newData := map[string]interface{}{
			"severity_level":   existingBusinessFunction.SeverityLevel.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s", existingBusinessFunction.ResponsibleUser.FirstNameEN, existingBusinessFunction.ResponsibleUser.LastNameEN),
			"title":            existingBusinessFunction.Title,
			"description":      existingBusinessFunction.Description,
			"review_frequency": utility.GetFieldValue(existingBusinessFunction.ReviewFrequency),
			"last_review_date": existingBusinessFunction.LastReviewDateString,
			"next_review_date": existingBusinessFunction.NextReviewDateString,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := bfs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BusinessFunction{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BusinessFunction{}.TableName(),
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

func (bfs *BusinessFunctionService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := bfs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBusinessFunction []map[string]interface{}

	for _, id := range ids {
		businessFunction, err := bfs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBusinessFunction = append(deletedBusinessFunction, map[string]interface{}{
			"id":    id,
			"title": businessFunction.Title,
		})

		if len(businessFunction.Processes) > 0 {
			var deleteControlIDs []*uint
			for _, process := range businessFunction.Processes {
				deleteControlIDs = append(deleteControlIDs, process.ID)
			}
			if err := bfs.repo.DeleteBusinessFunctionProcesses(deleteControlIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the business
		if err := bfs.repo.Delete(id, tx); err != nil {
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

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := bfs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BusinessFunction{}.TableName(),
				successfulIds,
				deletedBusinessFunction,
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

func (bfs *BusinessFunctionService) UpdateBusinessFunction(businessFunction *models.BusinessFunction, tx utility.DatabaseTransaction) error {

	err := bfs.repo.Update(businessFunction, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (bfs *BusinessFunctionService) parseBusinessFunctionDates(dto *businessContinuityDTOs.BusinessFunctionFormDTO) (*businessContinuityDTOs.BusinessFunctionParsedDates, error) {
	parsed := &businessContinuityDTOs.BusinessFunctionParsedDates{}

	var err error

	if dto.NextReviewDate != nil {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.LastReviewDate != nil {
		parsed.LastReviewDate, err = utility.ParseDate(*dto.LastReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
		}
	}

	return parsed, nil
}

func (bfs *BusinessFunctionService) ExportBusinessFunction(data []*businessContinuityDTOs.BusinessFunctionListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := bfs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Severity Level",
			"Responsible User",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Severity Level",
			"Responsible User",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, businessFunction := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				businessFunction.ReferenceNumber,
				businessFunction.Title,
				businessFunction.SeverityLevelTitleAr,
				fmt.Sprintf("%s %s", businessFunction.ResponsibleUserFirstNameAr, businessFunction.ResponsibleUserLastNameAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				businessFunction.ReferenceNumber,
				businessFunction.Title,
				businessFunction.SeverityLevelTitleEn,
				fmt.Sprintf("%s %s", businessFunction.ResponsibleUserFirstNameEn, businessFunction.ResponsibleUserLastNameEn),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (bfs *BusinessFunctionService) GetProcessAssociatedRisks(businessFunctionID uint, processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemRiskListResponseDTO, error) {
	risks, err := bfs.processAssociatedItemSvc.GetAllRisks(processID, params)
	if err != nil {
		return nil, err
	}

	return risks, nil
}

func (bfs *BusinessFunctionService) GetProcessAssociatedIncidents(businessFunctionID uint, processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemIncidentListResponseDTO, error) {
	incidents, err := bfs.processAssociatedItemSvc.GetAllIncidents(processID, params)
	if err != nil {
		return nil, err
	}

	return incidents, nil
}

func (bfs *BusinessFunctionService) GetProcessAssociatedControls(businessFunctionID uint, processID uint, params *organizationDTOs.ProcessAssociatedItemListQueryParamsDTO) (*organizationDTOs.ProcessAssociatedItemControlListResponseDTO, error) {
	controls, err := bfs.processAssociatedItemSvc.GetAllControlsWithPagination(processID, params)
	if err != nil {
		return nil, err
	}

	return controls, nil
}
