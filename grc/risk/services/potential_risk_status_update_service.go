package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
)

type PotentialRiskStatusUpdateServiceInterface interface {
	GetAll(potentialRiskID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.PotentialRiskStatusUpdateListDTO, error)
	Create(potentialRiskID uint, currentUserID uint, statusUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) (uint, error)
}

// Service and implementations
type PotentialRiskStatusUpdateService struct {
	repo                   repositories.PotentialRiskStatusUpdateRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	potentialRiskSvc       PotentialRiskServiceInterface
	potentialRiskStatusSvc PotentialRiskStatusServiceInterface
}

func NewPotentialRiskStatusUpdateService(
	repo repositories.PotentialRiskStatusUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	potentialRiskSvc PotentialRiskServiceInterface,
	potentialRiskStatusSvc PotentialRiskStatusServiceInterface,
) PotentialRiskStatusUpdateServiceInterface {
	return &PotentialRiskStatusUpdateService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		potentialRiskSvc:       potentialRiskSvc,
		potentialRiskStatusSvc: potentialRiskStatusSvc,
	}
}

func (prsus *PotentialRiskStatusUpdateService) GetAll(potentialRiskID uint, sort string, sortBy string, keyword string) ([]*riskDTOs.PotentialRiskStatusUpdateListDTO, error) {
	if _, err := prsus.potentialRiskSvc.GetByID(potentialRiskID); err != nil {
		return nil, err
	}

	return prsus.repo.GetAll(potentialRiskID, sort, sortBy, keyword)
}

func (prsus *PotentialRiskStatusUpdateService) Create(potentialRiskID uint, currentUserID uint, statusUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) (uint, error) {

	if _, err := prsus.potentialRiskSvc.GetByID(potentialRiskID); err != nil {
		return 0, err
	}

	if err := prsus.validatePotentialRiskUpdateFormDTO(statusUpdateFormDTO); err != nil {
		return 0, err
	}

	potentialRiskUpdate := &models.PotentialRiskStatusUpdate{
		PotentialRiskID:       potentialRiskID,
		PotentialRiskStatusID: statusUpdateFormDTO.PotentialRiskStatusID,
		Comment:               statusUpdateFormDTO.Comment,
		CreatedBy:             currentUserID,
	}

	id, err := prsus.repo.Create(potentialRiskUpdate)
	if err != nil {
		return 0, err
	}

	err = prsus.potentialRiskSvc.StatusUpdate(potentialRiskID, currentUserID, statusUpdateFormDTO.PotentialRiskStatusID)
	if err != nil {
		return 0, err
	}

	go func() {
		err := prsus.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.PotentialRiskStatusUpdate{}.TableName(),
			id,
			potentialRiskUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.PotentialRiskStatusUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (prsus *PotentialRiskStatusUpdateService) validatePotentialRiskUpdateFormDTO(statusUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) error {

	if _, err := prsus.potentialRiskStatusSvc.GetByID(statusUpdateFormDTO.PotentialRiskStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("potential risk status does not exist"), "not_exist", "potential risk status")
	}

	return nil
}
