package services

import (
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
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
		return nil, err
	}

	return prsus.repo.GetAll(potentialRiskID, sort, sortBy, keyword)
}

func (prsus *PotentialRiskStatusUpdateService) Create(potentialRiskID uint, currentUserID uint, potentialRiskStatusUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) (uint, error) {

	if _, err := prsus.potentialRiskSvc.GetByID(potentialRiskID); err != nil {
		return 0, err
	}

	if err := prsus.validatePotentialRiskUpdateFormDTO(potentialRiskStatusUpdateFormDTO); err != nil {
		return 0, err
	}

	potentialRiskUpdate := &models.PotentialRiskStatusUpdate{
		PotentialRiskID:       potentialRiskID,
		PotentialRiskStatusID: potentialRiskStatusUpdateFormDTO.PotentialRiskStatusID,
		Comment:               potentialRiskStatusUpdateFormDTO.Comment,
		CreatedBy:             currentUserID,
	}

	id, err := prsus.repo.Create(potentialRiskUpdate)
	if err != nil {
		return 0, err
	}

	err = prsus.potentialRiskSvc.StatusUpdate(potentialRiskID, currentUserID, potentialRiskStatusUpdateFormDTO.PotentialRiskStatusID)
	if err != nil {
		return 0, err
	}

	go func() {
		err := prsus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PotentialRiskStatusUpdate{}.TableName(),
			id,
			potentialRiskUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
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

func (prsus *PotentialRiskStatusUpdateService) validatePotentialRiskUpdateFormDTO(potentialRiskStatusUpdateFormDTO *riskDTOs.PotentialRiskStatusUpdateFormDTO) error {

	if _, err := prsus.potentialRiskStatusSvc.GetByID(potentialRiskStatusUpdateFormDTO.PotentialRiskStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("potential risk status does not exist"), "not_exist", "potential risk status")
	}

	return nil
}
