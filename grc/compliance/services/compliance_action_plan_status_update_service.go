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

type ComplianceActionPlanStatusUpdateServiceInterface interface {
	GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error)
	Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error)
	Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error
}

type ComplianceActionPlanStatusUpdateService struct {
	repo                    repositories.ComplianceActionPlanStatusUpdateRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	complianceActionPlanSvc ComplianceActionPlanServiceInterface
	complianceStatusSvc     ComplianceStatusServiceInterface
}

func NewComplianceActionPlanStatusUpdateService(
	repo repositories.ComplianceActionPlanStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceActionPlanSvc ComplianceActionPlanServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceActionPlanStatusUpdateServiceInterface {
	return &ComplianceActionPlanStatusUpdateService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		complianceActionPlanSvc: complianceActionPlanSvc,
		complianceStatusSvc:     complianceStatusSvc,
	}
}

func (casus *ComplianceActionPlanStatusUpdateService) GetAll(ComplianceID uint, ComplianceActionPlanID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceActionPlanStatusUpdateListDTO, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(ComplianceActionPlanID, sort, sortBy, keyword)
}

func (casus ComplianceActionPlanStatusUpdateService) GetByID(ComplianceID uint, ComplianceActionPlanID uint, id uint) (*models.ComplianceActionPlanStatusUpdate, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return nil, err
	}

	complianceActionPlanStatusUpdate, err := casus.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceActionPlanStatusUpdate, nil
}

func (casus *ComplianceActionPlanStatusUpdateService) Create(ComplianceID uint, ComplianceActionPlanID uint, currentUserID uint, complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) (uint, error) {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return 0, err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	complianceActionPlanUpdate := &models.ComplianceActionPlanStatusUpdate{
		ComplianceActionPlanID: ComplianceActionPlanID,
		ComplianceStatusID:     complianceActionPlanStatusUpdateFormDTO.ComplianceStatusID,
		Comment:                complianceActionPlanStatusUpdateFormDTO.Comment,
		CreatedBy:              currentUserID,
	}

	id, err := casus.repo.Create(complianceActionPlanUpdate)
	if err != nil {
		return 0, err
	}

	// err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, complianceActionPlanStatusUpdateFormDTO.ComplianceStatusID)
	// if err != nil {
	// 	return 0, err
	// }

	go func() {
		err := casus.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceActionPlanStatusUpdate{}.TableName(),
			id,
			complianceActionPlanUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.ComplianceActionPlanStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (casus *ComplianceActionPlanStatusUpdateService) Update(ComplianceID uint, ComplianceActionPlanID uint, id uint, currentUserID uint, ComplianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {
	if _, err := casus.complianceActionPlanSvc.GetByID(ComplianceID, ComplianceActionPlanID); err != nil {
		return err
	}

	existingComplianceActionPlanStatusUpdate, err := casus.GetByID(ComplianceID, ComplianceActionPlanID, id)

	if err != nil {
		return err
	}

	if err := casus.validateComplianceActionPlanStatusUpdateFormDTO(ComplianceActionPlanStatusUpdateFormDTO); err != nil {
		return err
	}

	// Update Compliance StatusUpdate
	existingComplianceActionPlanStatusUpdate.ComplianceStatusID = ComplianceActionPlanStatusUpdateFormDTO.ComplianceStatusID
	existingComplianceActionPlanStatusUpdate.Comment = ComplianceActionPlanStatusUpdateFormDTO.Comment
	existingComplianceActionPlanStatusUpdate.UpdatedBy = &currentUserID

	err = casus.repo.Update(existingComplianceActionPlanStatusUpdate)

	if err != nil {
		return err
	}

	// err = casus.complianceActionPlanSvc.StatusUpdate(ComplianceID, ComplianceActionPlanID, currentUserID, ComplianceActionPlanStatusUpdateFormDTO.ComplianceStatusID)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func (casus *ComplianceActionPlanStatusUpdateService) validateComplianceActionPlanStatusUpdateFormDTO(complianceActionPlanStatusUpdateFormDTO *complianceDTOs.ComplianceActionPlanStatusUpdateFormDTO) error {

	if _, err := casus.complianceStatusSvc.GetByID(complianceActionPlanStatusUpdateFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	return nil
}
