package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRequirementStatusUpdateServiceInterface interface {
	GetAll(ComplianceID uint, ComplianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error)
	Create(ComplianceID uint, ComplianceRequirementID uint, currentUserID uint, complianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) (uint, error)
	Update(ComplianceID uint, ComplianceRequirementID uint, id uint, currentUserID uint, ComplianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) error
}

// Service and implementations
type ComplianceRequirementStatusUpdateService struct {
	repo                     repositories.ComplianceRequirementStatusUpdateRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	complianceRequirementSvc ComplianceRequirementServiceInterface
	complianceStatusSvc      ComplianceStatusServiceInterface
}

func NewComplianceRequirementStatusUpdateService(
	repo repositories.ComplianceRequirementStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRequirementSvc ComplianceRequirementServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceRequirementStatusUpdateServiceInterface {
	return &ComplianceRequirementStatusUpdateService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		complianceStatusSvc:      complianceStatusSvc,
	}
}

func (crsus *ComplianceRequirementStatusUpdateService) GetAll(ComplianceID uint, ComplianceRequirementID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceRequirementStatusUpdateListDTO, error) {
	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return nil, err
	}

	return crsus.repo.GetAll(ComplianceRequirementID, sort, sortBy, keyword)
}

func (crsus ComplianceRequirementStatusUpdateService) GetByID(ComplianceID uint, ComplianceRequirementID uint, id uint) (*models.ComplianceRequirementStatusUpdate, error) {
	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return nil, err
	}

	complianceRequirementStatusUpdate, err := crsus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirementStatusUpdate, nil
}

func (crsus *ComplianceRequirementStatusUpdateService) Create(ComplianceID uint, ComplianceRequirementID uint, currentUserID uint, complianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) (uint, error) {
	// @Review
	// please move custom error return from repo to here
	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return 0, err
	}

	if err := crsus.validateComplianceRequirementStatusUpdateFormDTO(complianceRequirementStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	//@review
	//Please ensure updated by  and udpated at will be null initially
	//please add transaction
	complianceStatusID := complianceRequirementStatusUpdateFormDTO.ComplianceStatusID
	complianceRequirementUpdate := &models.ComplianceRequirementStatusUpdate{
		ComplianceRequirementID: ComplianceRequirementID,
		ComplianceStatusID:      complianceStatusID,
		Comment:                 complianceRequirementStatusUpdateFormDTO.Comment,
		CreatedBy:               currentUserID,
	}

	id, err := crsus.repo.Create(complianceRequirementUpdate)
	if err != nil {
		return 0, err
	}
	err = crsus.complianceRequirementSvc.StatusUpdate(ComplianceID, ComplianceRequirementID, currentUserID, complianceStatusID)
	if err != nil {
		return 0, err
	}

	go func() {
		err := crsus.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceRequirementStatusUpdate{}.TableName(),
			id,
			complianceRequirementUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.ComplianceRequirementStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crsus *ComplianceRequirementStatusUpdateService) Update(ComplianceID uint, ComplianceRequirementID uint, id uint, currentUserID uint, ComplianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) error {
	if _, err := crsus.complianceRequirementSvc.GetByID(ComplianceID, ComplianceRequirementID); err != nil {
		return err
	}

	existingComplianceRequirementStatusUpdate, err := crsus.GetByID(ComplianceID, ComplianceRequirementID, id)

	if err != nil {
		return err
	}

	if err := crsus.validateComplianceRequirementStatusUpdateFormDTO(ComplianceRequirementStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Compliance StatusUpdate
	complianceStatusID := ComplianceRequirementStatusUpdateFormDTO.ComplianceStatusID
	existingComplianceRequirementStatusUpdate.ComplianceStatusID = complianceStatusID
	existingComplianceRequirementStatusUpdate.Comment = ComplianceRequirementStatusUpdateFormDTO.Comment
	existingComplianceRequirementStatusUpdate.UpdatedBy = &currentUserID

	err = crsus.repo.Update(existingComplianceRequirementStatusUpdate)

	if err != nil {
		return err
	}

	err = crsus.complianceRequirementSvc.StatusUpdate(ComplianceID, ComplianceRequirementID, currentUserID, complianceStatusID)
	if err != nil {
		return err
	}

	return nil
}

func (crsus *ComplianceRequirementStatusUpdateService) validateComplianceRequirementStatusUpdateFormDTO(complianceRequirementStatusUpdateFormDTO *complianceDTOs.ComplianceRequirementStatusUpdateFormDTO) error {

	if _, err := crsus.complianceStatusSvc.GetByID(complianceRequirementStatusUpdateFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}
