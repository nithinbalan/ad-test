package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessContinuityPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BusinessContinuityPlanServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessContinuityPlanDetailDTO, error)
	GetByID(id uint) (*models.BusinessContinuityPlan, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BusinessContinuityPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, businessContinuityPlan *businessContinuityDTOs.BusinessContinuityPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BusinessContinuityPlanService struct {
	repo                businessContinuityPlanRepositories.BusinessContinuityPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	businessFunctionSvc BusinessFunctionServiceInterface
}

func NewBusinessContinuityPlanService(
	repo businessContinuityPlanRepositories.BusinessContinuityPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	businessFunctionSvc BusinessFunctionServiceInterface,
) BusinessContinuityPlanServiceInterface {
	return &BusinessContinuityPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		businessFunctionSvc: businessFunctionSvc,
	}
}

func (bcps *BusinessContinuityPlanService) GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error) {
	return bcps.repo.GetAll(params)
}

func (bcps BusinessContinuityPlanService) GetByID(id uint) (*models.BusinessContinuityPlan, error) {
	businessContinuityPlan, err := bcps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return businessContinuityPlan, nil
}

func (bcps BusinessContinuityPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BusinessContinuityPlanDetailDTO, error) {
	// Initialize your DTO
	var businessContinuityPlanDTO businessContinuityDTOs.BusinessContinuityPlanDetailDTO

	businessContinuityPlan, err := bcps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(businessContinuityPlan, &businessContinuityPlanDTO)
	if err != nil {
		return nil, err
	}

	businessContinuityPlanDTO.CreatedAt = businessContinuityPlan.CreatedAt
	businessContinuityPlanDTO.UpdatedAt = businessContinuityPlan.UpdatedAt

	return &businessContinuityPlanDTO, nil
}

func (bcps *BusinessContinuityPlanService) Create(currentUserID uint, businessContinuityPlanFormDTO *businessContinuityDTOs.BusinessContinuityPlanFormDTO) (uint, error) {
	// // First, check if an business continuity already exists for the given business function
	// existingAudit, err := bcps.repo.GetBusinessContinuityPlanByBusinessFunctionID(businessContinuityPlanFormDTO.BusinessFunctionID)
	// if err != nil {
	// 	return 0, err // Handle DB query error
	// }
	// // If an business continuity plan already exists, return an error without starting a transaction
	// if existingAudit != nil {
	// 	return 0, utility.NewForbiddenError("Business Continuity Plan has already been created for this business function")
	// }

	tx := bcps.repo.BeginTransaction()

	exists, err := bcps.repo.GetBusinessContinuityPlanByBusinessFunctionID(businessContinuityPlanFormDTO.BusinessFunctionID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if exists != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("the business function already exists"), "unique", "business function")
	}

	referenceNumber, err := bcps.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	businessContinuityPlan := &models.BusinessContinuityPlan{
		ReferenceNumber:    referenceNumber,
		BusinessFunctionID: businessContinuityPlanFormDTO.BusinessFunctionID,
		Version:            businessContinuityPlanFormDTO.Version,
		Purpose:            businessContinuityPlanFormDTO.Purpose,
		Scope:              businessContinuityPlanFormDTO.Scope,
		Objectives:         businessContinuityPlanFormDTO.Objectives,
		Assumptions:        businessContinuityPlanFormDTO.Assumptions,
		CreatedBYID:        currentUserID,
	}
	id, err := bcps.repo.Create(businessContinuityPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	businessFunction, err := bcps.businessFunctionSvc.GetByID(businessContinuityPlanFormDTO.BusinessFunctionID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	businessFunction.BcpStatus = true

	if err := bcps.businessFunctionSvc.UpdateBcpStatus(businessFunction, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingBusinessContinuityPlan, err := bcps.GetFullDetailsByID(businessFunction.ID)
		if err != nil {
			utility.LogError("Failed to fetch business continuity plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = bcps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BusinessContinuityPlan{}.TableName(),
			id,
			existingBusinessContinuityPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BusinessContinuityPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (bcps *BusinessContinuityPlanService) Update(id uint, currentUserID uint, businessContinuityPlanFormDTO *businessContinuityDTOs.BusinessContinuityPlanFormDTO) error {
	tx := bcps.repo.BeginTransaction()

	existingBusinessContinuityPlan, err := bcps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldBusinessContinuityPlan, err := bcps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingBusinessContinuityPlan.BusinessFunctionID = businessContinuityPlanFormDTO.BusinessFunctionID
	existingBusinessContinuityPlan.Version = businessContinuityPlanFormDTO.Version
	existingBusinessContinuityPlan.Purpose = businessContinuityPlanFormDTO.Purpose
	existingBusinessContinuityPlan.Scope = businessContinuityPlanFormDTO.Scope
	existingBusinessContinuityPlan.Objectives = businessContinuityPlanFormDTO.Objectives
	existingBusinessContinuityPlan.Assumptions = businessContinuityPlanFormDTO.Assumptions
	existingBusinessContinuityPlan.UpdatedBYID = &currentUserID

	err = bcps.repo.Update(existingBusinessContinuityPlan, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingBusinessContinuityPlan, err := bcps.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"business_function": oldBusinessContinuityPlan.BusinessFunction.Title,
			"version":           oldBusinessContinuityPlan.Version,
			"scope":             oldBusinessContinuityPlan.Scope,
			"purpose":           oldBusinessContinuityPlan.Purpose,
			"objectives":        oldBusinessContinuityPlan.Objectives,
			"assumptions":       oldBusinessContinuityPlan.Assumptions,
		}

		newData := map[string]interface{}{
			"business_function": existingBusinessContinuityPlan.BusinessFunction.Title,
			"version":           existingBusinessContinuityPlan.Version,
			"scope":             existingBusinessContinuityPlan.Scope,
			"purpose":           existingBusinessContinuityPlan.Purpose,
			"objectives":        existingBusinessContinuityPlan.Objectives,
			"assumptions":       existingBusinessContinuityPlan.Assumptions,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := bcps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BusinessContinuityPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BusinessContinuityPlan{}.TableName(),
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

func (bcps *BusinessContinuityPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := bcps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedBusinessContinuityPlan []map[string]interface{}

	for _, id := range ids {
		businessContinuityPlan, err := bcps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedBusinessContinuityPlan = append(deletedBusinessContinuityPlan, map[string]interface{}{
			"id":                id,
			"business_function": businessContinuityPlan.BusinessFunction.Title,
			"version":           businessContinuityPlan.Version,
		})

		// Delete the business
		if err := bcps.repo.Delete(id, tx); err != nil {
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
			err := bcps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BusinessContinuityPlan{}.TableName(),
				successfulIds,
				deletedBusinessContinuityPlan,
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
