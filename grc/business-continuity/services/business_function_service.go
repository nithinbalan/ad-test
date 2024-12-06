package services

import (
	"errors"
	"fmt"
	businessFunctionDTOs "jethur/grc/business-continuity/dtos"
	businessFunctionRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BusinessFunctionServiceInterface interface {
	GetAll(params *businessFunctionDTOs.BusinessFunctionListQueryParamsDTO) (*businessFunctionDTOs.BusinessFunctionListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessFunctionDTOs.BusinessFunctionDetailDTO, error)
	GetByID(id uint) (*models.BusinessFunction, error)
	Create(currentUserID uint, businessFunctionDTOs *businessFunctionDTOs.BusinessFunctionFormDTO) (uint, error)
	Update(id uint, currentUserID uint, businesFunction *businessFunctionDTOs.BusinessFunctionFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BusinessFunctionService struct {
	repo           businessFunctionRepositories.BusinessFunctionRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewBusinessFunctionService(
	repo businessFunctionRepositories.BusinessFunctionRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) BusinessFunctionServiceInterface {
	return &BusinessFunctionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (bfs *BusinessFunctionService) GetAll(params *businessFunctionDTOs.BusinessFunctionListQueryParamsDTO) (*businessFunctionDTOs.BusinessFunctionListResponseDTO, error) {
	return bfs.repo.GetAll(params)
}

func (bfs BusinessFunctionService) GetByID(id uint) (*models.BusinessFunction, error) {
	businesFunction, err := bfs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return businesFunction, nil
}

func (bfs BusinessFunctionService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessFunctionDTOs.BusinessFunctionDetailDTO, error) {
	// Initialize your DTO
	var businessFunctionDTO businessFunctionDTOs.BusinessFunctionDetailDTO

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

	businessFunctionDTO.CreatedAt = businessFunction.CreatedAt
	businessFunctionDTO.UpdatedAt = businessFunction.UpdatedAt

	return &businessFunctionDTO, nil
}

func (bfs *BusinessFunctionService) Create(currentUserID uint, businessFunctionFormDTO *businessFunctionDTOs.BusinessFunctionFormDTO) (uint, error) {
	tx := bfs.repo.BeginTransaction()

	referenceNumber, err := bfs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	businessFunction := &models.BusinessFunction{
		ReferenceNumber:   referenceNumber,
		SeverityLevelID:   businessFunctionFormDTO.SeverityLevelID,
		ResponsibleUserID: businessFunctionFormDTO.ResponsibleUserID,
		Title:             businessFunctionFormDTO.Title,
		Description:       businessFunctionFormDTO.Description,
		CreatedByID:       currentUserID,
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
			CreatedByID:        currentUserID,
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
func (bfs *BusinessFunctionService) Update(id uint, currentUserID uint, businessFunctionFormDTO *businessFunctionDTOs.BusinessFunctionFormDTO) error {
	tx := bfs.repo.BeginTransaction()

	existingBusinessFunction, err := bfs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}
	// oldResponsibleUserID := existingBusinessFunction.ResponsibleUserID

	oldBusinessFunction, err := bfs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingBusinessFunction.SeverityLevelID = businessFunctionFormDTO.SeverityLevelID
	existingBusinessFunction.ResponsibleUserID = businessFunctionFormDTO.ResponsibleUserID
	existingBusinessFunction.Title = businessFunctionFormDTO.Title
	existingBusinessFunction.Description = businessFunctionFormDTO.Description
	existingBusinessFunction.UpdatedByID = &currentUserID

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
				CreatedByID:        currentUserID,
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
		}

		newData := map[string]interface{}{
			"severity_level":   existingBusinessFunction.SeverityLevel.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s", existingBusinessFunction.ResponsibleUser.FirstNameEN, existingBusinessFunction.ResponsibleUser.LastNameEN),
			"title":            existingBusinessFunction.Title,
			"description":      existingBusinessFunction.Description,
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
		businesFunction, err := bfs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBusinessFunction = append(deletedBusinessFunction, map[string]interface{}{
			"id":    id,
			"title": businesFunction.Title,
		})

		if len(businesFunction.Processes) > 0 {
			var deleteControlIDs []*uint
			for _, process := range businesFunction.Processes {
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
