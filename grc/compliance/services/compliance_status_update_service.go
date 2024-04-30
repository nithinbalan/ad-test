package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceStatusUpdateServiceInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error)
	Create(complianceID uint, currentUserID uint, complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, ComplianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error
}

// Service and implementations
type ComplianceStatusUpdateService struct {
	repo                  repositories.ComplianceStatusUpdateRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	complianceRegisterSvc ComplianceRegisterServiceInterface
	complianceStatusSvc   ComplianceStatusServiceInterface
}

func NewComplianceStatusUpdateService(
	repo repositories.ComplianceStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceStatusUpdateServiceInterface {
	return &ComplianceStatusUpdateService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		complianceRegisterSvc: complianceRegisterSvc,
		complianceStatusSvc:   complianceStatusSvc,
	}
}

func (csus *ComplianceStatusUpdateService) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error) {
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	return csus.repo.GetAll(complianceID, sort, sortBy, keyword)
}

func (csus ComplianceStatusUpdateService) GetByID(complianceID uint, id uint) (*models.ComplianceStatusUpdate, error) {
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceStatusUpdate, err := csus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceStatusUpdate, nil
}

func (csus *ComplianceStatusUpdateService) Create(complianceID uint, currentUserID uint, complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) (uint, error) {
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return 0, err
	}

	if err := csus.validateComplianceStatusUpdateFormDTO(complianceStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	complianceUpdate := &models.ComplianceStatusUpdate{
		ComplianceRegisterID: complianceID,
		ComplianceStatusID:   complianceStatusUpdateFormDTO.ComplianceStatusID,
		Comment:              complianceStatusUpdateFormDTO.Comment,
		CreatedBy:            currentUserID,
	}

	id, err := csus.repo.Create(complianceUpdate)
	if err != nil {
		return 0, err
	}

	err = csus.complianceRegisterSvc.StatusUpdate(complianceID, currentUserID, complianceStatusUpdateFormDTO.ComplianceStatusID)
	if err != nil {
		return 0, err
	}

	go func() {
		err := csus.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceStatusUpdate{}.TableName(),
			id,
			complianceUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.ComplianceStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (csus *ComplianceStatusUpdateService) Update(complianceID uint, id uint, currentUserID uint, ComplianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error {
	if _, err := csus.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceStatusUpdate, err := csus.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	// Update Compliance StatusUpdate -- update only comments
	existingComplianceStatusUpdate.Comment = ComplianceStatusUpdateFormDTO.Comment
	existingComplianceStatusUpdate.UpdatedBy = &currentUserID

	err = csus.repo.Update(existingComplianceStatusUpdate)

	if err != nil {
		return err
	}

	return nil
}

func (csus *ComplianceStatusUpdateService) validateComplianceStatusUpdateFormDTO(complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error {

	if _, err := csus.complianceStatusSvc.GetByID(complianceStatusUpdateFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}
