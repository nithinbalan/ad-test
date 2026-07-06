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
	GetFullDetailsByID(bcpId uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpTestAndExerciseDetailDTO, error)
	GetByID(bcpId uint, id uint) (*models.BcpTestAndExercise, error)
	Create(bcpId uint, currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BcpTestAndExerciseFormDTO) (uint, error)
	Update(bcpId uint, id uint, currentUserID uint, bcpTestAndExercise *businessContinuityDTOs.BcpTestAndExerciseFormDTO) error
	Delete(bcpId uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BcpTestAndExerciseService struct {
	repo                      BcpTestAndExerciseRepositories.BcpTestAndExerciseRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface
}

func NewBcpTestAndExerciseService(
	repo BcpTestAndExerciseRepositories.BcpTestAndExerciseRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface,
) BcpTestAndExerciseServiceInterface {
	return &BcpTestAndExerciseService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		userSvc:                   userSvc,
		businessContinuityPlanSvc: businessContinuityPlanSvc,
	}
}

func (btes *BcpTestAndExerciseService) GetAll(bcpId uint, params *businessContinuityDTOs.BcpTestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.BcpTestAndExerciseListResponseDTO, error) {
	if bcpId != 0 {
		if _, err := btes.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
			return nil, err
		}
	}
	return btes.repo.GetAll(bcpId, params)
}

func (btes BcpTestAndExerciseService) GetByID(bcpId uint, id uint) (*models.BcpTestAndExercise, error) {
	if _, err := btes.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return nil, err
	}

	BcpTestAndExercise, err := btes.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return BcpTestAndExercise, nil
}

func (btes BcpTestAndExerciseService) GetFullDetailsByID(bcpId uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BcpTestAndExerciseDetailDTO, error) {
	// Initialize your DTO
	var BcpTestAndExerciseDTO businessContinuityDTOs.BcpTestAndExerciseDetailDTO

	if _, err := btes.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return nil, err
	}

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

	if _, err := btes.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return 0, err
	}

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

	// Save SystemRegisters
	var systemRegisters []*models.BcpTestAndExerciseSystemRegister
	for _, systemRegisterID := range bcpTestAndExerciseFormDTO.SystemRegisterIDs {
		systemRegister := &models.BcpTestAndExerciseSystemRegister{
			BcpTestAndExerciseID: id,
			SystemRegisterID:     *systemRegisterID,
			CreatedBYID:          currentUserID,
		}
		systemRegisters = append(systemRegisters, systemRegister)
	}

	if len(systemRegisters) > 0 {
		if err := btes.repo.CreateBcpTestAndExerciseSystemRegisters(systemRegisters, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	go func() {
		existingBcpTestAndExercise, err := btes.GetFullDetailsByID(bcpId, BcpTestAndExercise.ID)
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
			"business_continuity",
			"bcp",
			"test_and_exercise",
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

func (btes *BcpTestAndExerciseService) Update(bcpId uint, id uint, currentUserID uint, bcpTestAndExerciseFormDTO *businessContinuityDTOs.BcpTestAndExerciseFormDTO) error {
	tx := btes.repo.BeginTransaction()

	if _, err := btes.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return err
	}

	existingBcpTestAndExercise, err := btes.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldBcpTestAndExercise, err := btes.GetFullDetailsByID(bcpId, id, tx)
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

	// Collect existing SystemRegister IDs
	existingTestAndExercise, err := btes.GetFullDetailsByID(bcpId, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingSystemRegisterIDs := make(map[uint]bool)
	for _, sr := range existingTestAndExercise.BcpTestAndExerciseSystemRegisters {
		existingSystemRegisterIDs[sr.SystemRegisterID] = true
	}

	// Collect new SystemRegister IDs
	newSystemRegisterIDs := make(map[uint]bool)
	for _, sr := range bcpTestAndExerciseFormDTO.SystemRegisterIDs {
		newSystemRegisterIDs[*sr] = true
	}

	// Add new SystemRegisters
	var systemRegistersToAdd []*models.BcpTestAndExerciseSystemRegister
	for _, newSystemRegisterID := range bcpTestAndExerciseFormDTO.SystemRegisterIDs {
		if _, exists := existingSystemRegisterIDs[*newSystemRegisterID]; !exists {
			// Create a new TestAndExercise system register record
			systemRegister := &models.BcpTestAndExerciseSystemRegister{
				BcpTestAndExerciseID: id,
				SystemRegisterID:     *newSystemRegisterID,
				CreatedBYID:          currentUserID,
			}
			systemRegistersToAdd = append(systemRegistersToAdd, systemRegister)
		}
	}

	// Save the new SystemRegisters
	if len(systemRegistersToAdd) > 0 {
		if err := btes.repo.CreateBcpTestAndExerciseSystemRegisters(systemRegistersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove SystemRegisters not in the new list
	var deleteSystemRegisterIDs []*uint

	for _, existingSystemRegister := range existingTestAndExercise.BcpTestAndExerciseSystemRegisters {
		if _, exists := newSystemRegisterIDs[existingSystemRegister.SystemRegisterID]; !exists {
			// Append system register to delete list
			deleteSystemRegisterIDs = append(deleteSystemRegisterIDs, existingSystemRegister.ID)
		}
	}

	// Delete SystemRegisters not in the new list
	if len(deleteSystemRegisterIDs) > 0 {
		if err := btes.repo.DeleteBcpTestAndExerciseSystemRegisters(deleteSystemRegisterIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingBcpTestAndExercise, err := btes.GetFullDetailsByID(bcpId, id)
		if err != nil {
			return
		}

		// Determine system register changes

		systemRegistersAdded := []uint{}
		systemRegistersRemoved := []uint{}
		for srID := range newSystemRegisterIDs {
			if !existingSystemRegisterIDs[srID] {
				systemRegistersAdded = append(systemRegistersAdded, srID)
			}
		}
		for srID := range existingSystemRegisterIDs {
			if !newSystemRegisterIDs[srID] {
				systemRegistersRemoved = append(systemRegistersRemoved, srID)
			}
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(systemRegistersAdded) > 0 || len(systemRegistersRemoved) > 0 {
			associatedItems["system_register_ids"] = map[string]interface{}{
				"added":   systemRegistersAdded,
				"removed": systemRegistersRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}

		// log changes
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
				"business_continuity",
				"bcp",
				"test_and_exercise",
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

func (btes *BcpTestAndExerciseService) Delete(bcpId uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := btes.businessContinuityPlanSvc.GetByID(bcpId); err != nil {
		return nil, err
	}
	tx := btes.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBcpTestAndExercise []map[string]interface{}

	for _, id := range ids {
		bcpTestAndExercise, err := btes.GetFullDetailsByID(bcpId, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBcpTestAndExercise = append(deletedBcpTestAndExercise, map[string]interface{}{
			"id":    id,
			"title": bcpTestAndExercise.Title,
		})

		if len(bcpTestAndExercise.BcpTestAndExerciseSystemRegisters) > 0 {
			var deleteSystemRegisterIDs []*uint
			for _, systemRegister := range bcpTestAndExercise.BcpTestAndExerciseSystemRegisters {
				deleteSystemRegisterIDs = append(deleteSystemRegisterIDs, systemRegister.ID)
			}
			if err := btes.repo.DeleteBcpTestAndExerciseSystemRegisters(deleteSystemRegisterIDs, tx); err != nil {
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
