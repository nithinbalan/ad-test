package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	"jethur/utility"
)

type ComplianceStatusUpdateServiceInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error)
	Create(complianceID uint, currentUserID uint, complianceStatusUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) (uint, error)
}

// Service and implementations
type ComplianceStatusUpdateService struct {
	repo                repositories.ComplianceStatusUpdateRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	complianceSvc       ComplianceRegisterServiceInterface
	complianceStatusSvc ComplianceStatusServiceInterface
}

func NewComplianceStatusUpdateService(
	repo repositories.ComplianceStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceSvc ComplianceRegisterServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceStatusUpdateServiceInterface {
	return &ComplianceStatusUpdateService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		complianceSvc:       complianceSvc,
		complianceStatusSvc: complianceStatusSvc,
	}
}

func (csus *ComplianceStatusUpdateService) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceStatusUpdateListDTO, error) {
	if _, err := csus.complianceSvc.GetByID(complianceID); err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("compliance does not exist"), "not_exist", "compliance")
	}

	return csus.repo.GetAll(complianceID, sort, sortBy, keyword)
}

func (csus *ComplianceStatusUpdateService) Create(complianceID uint, currentUserID uint, complianceUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) (uint, error) {

	if _, err := csus.complianceSvc.GetByID(complianceID); err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("compliance does not exist"), "not_exist", "compliance")
	}

	if err := csus.validateComplianceUpdateFormDTO(complianceUpdateFormDTO); err != nil {
		return 0, err
	}

	complianceUpdate := &models.ComplianceStatusUpdate{
		ComplianceID:       complianceID,
		ComplianceStatusID: complianceUpdateFormDTO.ComplianceStatusID,
		Comment:            complianceUpdateFormDTO.Comment,
		CreatedBy:          currentUserID,
	}

	id, err := csus.repo.Create(complianceUpdate)
	if err != nil {
		return 0, err
	}

	err = csus.complianceSvc.StatusUpdate(complianceID, currentUserID, complianceUpdateFormDTO.ComplianceStatusID)
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

func (csus *ComplianceStatusUpdateService) validateComplianceUpdateFormDTO(complianceUpdateFormDTO *complianceDTOs.ComplianceStatusUpdateFormDTO) error {

	if _, err := csus.complianceStatusSvc.GetByID(complianceUpdateFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}
