package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	BcpTestAndExerciseRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BcpTestAndExerciseServiceInterface interface {
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpTestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.BcpTestAndExerciseListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpTestAndExerciseDetailDTO, error)
	GetByID(id uint) (*models.BcpTestAndExercise, error)
	Create(bcpId uint, currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BcpTestAndExerciseFormDTO) (uint, error)
	Update(id uint, currentUserID uint, BcpTestAndExercise *businessContinuityDTOs.BcpTestAndExerciseFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BcpTestAndExerciseService struct {
	repo           BcpTestAndExerciseRepositories.BcpTestAndExerciseRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewBcpTestAndExerciseService(
	repo BcpTestAndExerciseRepositories.BcpTestAndExerciseRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) BcpTestAndExerciseServiceInterface {
	return &BcpTestAndExerciseService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (btes *BcpTestAndExerciseService) GetAll(bcpId uint, params *businessContinuityDTOs.BcpTestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.BcpTestAndExerciseListResponseDTO, error) {
	return btes.repo.GetAll(bcpId, params)
}

func (btes BcpTestAndExerciseService) GetByID(id uint) (*models.BcpTestAndExercise, error) {
	BcpTestAndExercise, err := btes.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return BcpTestAndExercise, nil
}

func (btes BcpTestAndExerciseService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpTestAndExerciseDetailDTO, error) {
	// Initialize your DTO
	var BcpTestAndExerciseDTO businessContinuityDTOs.BcpTestAndExerciseDetailDTO

	BcpTestAndExercise, err := btes.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(BcpTestAndExercise, &BcpTestAndExerciseDTO)
	if err != nil {
		return nil, err
	}

	BcpTestAndExerciseDTO.CreatedAt = BcpTestAndExercise.CreatedAt
	BcpTestAndExerciseDTO.UpdatedAt = BcpTestAndExercise.UpdatedAt

	return &BcpTestAndExerciseDTO, nil
}

func (btes *BcpTestAndExerciseService) Create(bcpId uint, currentUserID uint, bcpTestAndExerciseFormDTO *businessContinuityDTOs.BcpTestAndExerciseFormDTO) (uint, error) {
	tx := btes.repo.BeginTransaction()

	responsibleUnitID, err := btes.userSvc.GetUserPrimaryUnitID(*bcpTestAndExerciseFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}

	BcpTestAndExercise := &models.BcpTestAndExercise{
		BusinessContinuityPlanID: bcpId,
		ReviewFrequencyID:        bcpTestAndExerciseFormDTO.ReviewFrequencyID,
		ResponsibleUserID:        bcpTestAndExerciseFormDTO.ResponsibleUserID,
		ResponsibleUnitID:        &responsibleUnitID,
		Title:                    bcpTestAndExerciseFormDTO.Title,
		Objective:                bcpTestAndExerciseFormDTO.Objective,
		Scenario:                 bcpTestAndExerciseFormDTO.Scenario,
		SuccessCriteria:          bcpTestAndExerciseFormDTO.SuccessCriteria,
		CreatedBYID:              currentUserID,
	}

	id, err := btes.repo.Create(BcpTestAndExercise, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save BcpProcesses
	var bcpProcesses []*models.BcpTestAndExerciseBcpProcess
	for _, bcpProcessID := range bcpTestAndExerciseFormDTO.BcpProcessIDs {
		bcpProcess := &models.BcpTestAndExerciseBcpProcess{
			BcpTestAndExerciseID: id,
			BcpProcessID:         *bcpProcessID,
			CreatedBYID:          currentUserID,
		}
		bcpProcesses = append(bcpProcesses, bcpProcess)
	}

	if len(bcpProcesses) > 0 {
		if err := btes.repo.CreateBcpTestAndExerciseBcpProcesses(bcpProcesses, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	go func() {
		existingBcpTestAndExercise, err := btes.GetFullDetailsByID(BcpTestAndExercise.ID)
		if err != nil {
			utility.LogError("Failed to fetch bcp communication plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = btes.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BcpTestAndExercise{}.TableName(),
			id,
			existingBcpTestAndExercise,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BcpTestAndExercise{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (btes *BcpTestAndExerciseService) Update(id uint, currentUserID uint, bcpTestAndExerciseFormDTO *businessContinuityDTOs.BcpTestAndExerciseFormDTO) error {
	tx := btes.repo.BeginTransaction()

	existingBcpTestAndExercise, err := btes.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldBcpTestAndExercise, err := btes.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingBcpTestAndExercise.ResponsibleUserID != bcpTestAndExerciseFormDTO.ResponsibleUserID {
		responsibleUnitID, err := btes.userSvc.GetUserPrimaryUnitID(*existingBcpTestAndExercise.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingBcpTestAndExercise.ResponsibleUnitID = &responsibleUnitID
	}

	existingBcpTestAndExercise.ResponsibleUserID = bcpTestAndExerciseFormDTO.ResponsibleUserID
	existingBcpTestAndExercise.ReviewFrequencyID = bcpTestAndExerciseFormDTO.ReviewFrequencyID
	existingBcpTestAndExercise.Title = bcpTestAndExerciseFormDTO.Title
	existingBcpTestAndExercise.Objective = bcpTestAndExerciseFormDTO.Objective
	existingBcpTestAndExercise.Scenario = bcpTestAndExerciseFormDTO.Scenario
	existingBcpTestAndExercise.SuccessCriteria = bcpTestAndExerciseFormDTO.SuccessCriteria
	existingBcpTestAndExercise.UpdatedBYID = &currentUserID

	err = btes.repo.Update(existingBcpTestAndExercise, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Collect existing bcpProcess IDs
	existingTestAndExercise, err := btes.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingBcpProcessIDs := make(map[uint]bool)
	for _, proc := range existingTestAndExercise.BcpTestAndExerciseBcpProcesses {
		existingBcpProcessIDs[proc.BcpProcessID] = true
	}

	// Collect new bcpProcess IDs
	newBcpProcessIDs := make(map[uint]bool)
	for _, proc := range bcpTestAndExerciseFormDTO.BcpProcessIDs {
		newBcpProcessIDs[*proc] = true
	}

	// Add new bcpProcesses
	var bcpProcessesToAdd []*models.BcpTestAndExerciseBcpProcess
	for _, newBcpProcessID := range bcpTestAndExerciseFormDTO.BcpProcessIDs {
		if _, exists := existingBcpProcessIDs[*newBcpProcessID]; !exists {
			// Create a new TestAndExercise bcpProcess record
			bcpProcess := &models.BcpTestAndExerciseBcpProcess{
				BcpTestAndExerciseID: id,
				BcpProcessID:         *newBcpProcessID,
				CreatedBYID:          currentUserID,
			}
			bcpProcessesToAdd = append(bcpProcessesToAdd, bcpProcess)
		}
	}

	// Save the new bcpProcesses
	if len(bcpProcessesToAdd) > 0 {
		if err := btes.repo.CreateBcpTestAndExerciseBcpProcesses(bcpProcessesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove bcpProcesses not in the new list
	var deleteBcpProcessIDs []*uint

	for _, existingProcess := range existingTestAndExercise.BcpTestAndExerciseBcpProcesses {
		if _, exists := newBcpProcessIDs[existingProcess.BcpProcessID]; !exists {
			// Append bcpProcess to delete list
			deleteBcpProcessIDs = append(deleteBcpProcessIDs, existingProcess.ID)
		}
	}

	// Delete bcpProcesses not in the new list
	if len(deleteBcpProcessIDs) > 0 {
		if err := btes.repo.DeleteBcpTestAndExerciseBcpProcesses(deleteBcpProcessIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingBcpTestAndExercise, err := btes.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Determine bcpProcess changes

		processesAdded := []uint{}
		processesRemoved := []uint{}
		for id := range newBcpProcessIDs {
			if !existingBcpProcessIDs[id] {
				processesAdded = append(processesAdded, id)
			}
		}
		for id := range existingBcpProcessIDs {
			if !newBcpProcessIDs[id] {
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
			"responsible_user": fmt.Sprintf("%s %s", oldBcpTestAndExercise.ResponsibleUser.FirstNameEN, oldBcpTestAndExercise.ResponsibleUser.LastNameEN),
			"title":            oldBcpTestAndExercise.Title,
			"review_frequency": oldBcpTestAndExercise.ReviewFrequency.TitleEn,
			"objective":        oldBcpTestAndExercise.Objective,
			"scenario":         oldBcpTestAndExercise.Scenario,
			"success_criteris": oldBcpTestAndExercise.SuccessCriteria,
		}

		newData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", existingBcpTestAndExercise.ResponsibleUser.FirstNameEN, existingBcpTestAndExercise.ResponsibleUser.LastNameEN),
			"title":            existingBcpTestAndExercise.Title,
			"review_frequency": existingBcpTestAndExercise.ReviewFrequency.TitleEn,
			"objective":        existingBcpTestAndExercise.Objective,
			"scenario":         existingBcpTestAndExercise.Scenario,
			"success_criteris": existingBcpTestAndExercise.SuccessCriteria,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := btes.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BcpTestAndExercise{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BcpTestAndExercise{}.TableName(),
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

func (btes *BcpTestAndExerciseService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := btes.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBcpTestAndExercise []map[string]interface{}

	for _, id := range ids {
		bcpTestAndExercise, err := btes.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBcpTestAndExercise = append(deletedBcpTestAndExercise, map[string]interface{}{
			"id":    id,
			"title": bcpTestAndExercise.Title,
		})

		if len(bcpTestAndExercise.BcpTestAndExerciseBcpProcesses) > 0 {
			var deleteprocessIDs []*uint
			for _, process := range bcpTestAndExercise.BcpTestAndExerciseBcpProcesses {
				deleteprocessIDs = append(deleteprocessIDs, process.ID)
			}
			if err := btes.repo.DeleteBcpTestAndExerciseBcpProcesses(deleteprocessIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the business
		if err := btes.repo.Delete(id, tx); err != nil {
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
			err := btes.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BcpTestAndExercise{}.TableName(),
				successfulIds,
				deletedBcpTestAndExercise,
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
