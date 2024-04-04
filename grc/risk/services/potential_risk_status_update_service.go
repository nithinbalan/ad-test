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
	Create(potentialRiskID uint, currentUserID uint, potentialRiskStatusUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) (uint, error)
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
		return nil, utility.NewValidationError(fmt.Errorf("potential risk does not exist"), "not_exist", "potential risk")
	}

	return prsus.repo.GetAll(potentialRiskID, sort, sortBy, keyword)
}

func (prsus *PotentialRiskStatusUpdateService) Create(potentialRiskID uint, currentUserID uint, potentialRiskUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) (uint, error) {

	if _, err := prsus.potentialRiskSvc.GetByID(potentialRiskID); err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("potential risk does not exist"), "not_exist", "potential risk")
	}

	if err := prsus.validatePotentialRiskUpdateFormDTO(potentialRiskID, potentialRiskUpdateFormDTO); err != nil {
		return 0, err
	}

	potentialRiskUpdate := &models.PotentialRiskStatusUpdate{
		PotentialRiskID:       potentialRiskID,
		PotentialRiskStatusID: potentialRiskUpdateFormDTO.PotentialRiskStatusID,
		Comment:               potentialRiskUpdateFormDTO.Comment,
		CreatedBy:             currentUserID,
	}

	id, err := prsus.repo.Create(potentialRiskUpdate)
	if err != nil {
		return 0, err
	}

	err = prsus.potentialRiskSvc.StatusUpdate(potentialRiskID, currentUserID, potentialRiskUpdateFormDTO.PotentialRiskStatusID)
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

func (prsus *PotentialRiskStatusUpdateService) validatePotentialRiskUpdateFormDTO(potentialRiskID uint, potentialRiskUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) error {

	if _, err := prsus.potentialRiskStatusSvc.GetByID(potentialRiskUpdateFormDTO.PotentialRiskStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("potential risk status does not exist"), "not_exist", "potential risk status")
	}

	return nil
}
