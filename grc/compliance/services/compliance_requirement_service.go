package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	"jethur/utility"
	"log"
)

type ComplianceRequirementServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceRequirement, error)
	Create(currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRequirementFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceRequirementService struct {
	repo           repositories.ComplianceRequirementRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewComplianceRequirementService(
	repo repositories.ComplianceRequirementRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface) ComplianceRequirementServiceInterface {
	return &ComplianceRequirementService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (crs *ComplianceRequirementService) GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs ComplianceRequirementService) GetByID(id uint) (*models.ComplianceRequirement, error) {
	return crs.repo.GetByID(id)
}

func (crs *ComplianceRequirementService) Create(currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error) {

	// Proceed with creating the compliance requirement record
	id, err := crs.repo.Create(currentUserID, complianceRequirementFormDTO)
	if err != nil {
		return 0, err
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceRequirement{}.TableName(),
			id,
			complianceRequirementFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (crs *ComplianceRequirementService) Update(id uint, currentUserID uint, complianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {

	existingComplianceRequirement, err := crs.repo.GetByID(id)

	if err != nil {
		return err
	}

	err = crs.repo.Update(id, currentUserID, complianceRequirementFormDTO)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{"ComplianceRequirementCategoryID",
		"ComplianceStatusID",
		"ResponsibleUserID",
		"SeverityLevelID",
		"Title",
		"Deadline"}
	changes, err := utility.CompareFields(existingComplianceRequirement, complianceRequirementFormDTO, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.ComplianceRequirement{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return nil
}

func (crs *ComplianceRequirementService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := crs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = crs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRequirement{}.TableName(),
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
