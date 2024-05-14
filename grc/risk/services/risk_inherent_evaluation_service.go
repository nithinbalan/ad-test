package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskInherentEvaluationServiceInterface interface {
	GetAll(riskID uint) ([]*riskDTOs.RiskInherentEvaluationListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskEvaluation, error)
	Create(riskID uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) error
	Delete(riskID uint, id uint, currentUserID uint) (uint, error)
}

// Service and implementations
type RiskInherentEvaluationService struct {
	repo            riskRepositories.RiskInherentEvaluationRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	riskRegisterSvc RiskRegisterServiceInterface
}

func NewRiskInherentEvaluationService(
	repo riskRepositories.RiskInherentEvaluationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
) RiskInherentEvaluationServiceInterface {
	return &RiskInherentEvaluationService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		riskRegisterSvc: riskRegisterSvc,
	}
}

func (ries *RiskInherentEvaluationService) GetAll(riskID uint) ([]*riskDTOs.RiskInherentEvaluationListDTO, error) {
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}
	return ries.repo.GetAll(riskID)
}

func (ries *RiskInherentEvaluationService) GetByID(riskID uint, id uint) (*models.RiskEvaluation, error) {
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return nil, err
	}

	riskInherentEvaluation, err := ries.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskInherentEvaluation, nil
}

func (ries *RiskInherentEvaluationService) Create(riskID uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) (uint, error) {
	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	// Increment the "Why" field by 1
	latestWhy, err := ries.repo.GetLatestWhyByRiskID(riskID)
	if err != nil {
		return 0, err
	}
	why := latestWhy + 1

	riskInherentEvaluation := &models.RiskEvaluation{
		RiskRegisterID: riskID,
		Why:            why,
		Question:       inherentEvaluationFormDTO.Question,
		Answer:         inherentEvaluationFormDTO.Answer,
		CreatedBy:      currentUserID,
	}

	id, err := ries.repo.Create(riskInherentEvaluation)
	if err != nil {
		return 0, err
	}

	go func() {
		err := ries.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskEvaluation{}.TableName(),
			id,
			inherentEvaluationFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskEvaluation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ries *RiskInherentEvaluationService) Update(riskID uint, id uint, currentUserID uint, inherentEvaluationFormDTO *riskDTOs.RiskInherentEvaluationFormDTO) error {

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return err
	}

	existingRiskInherentEvaluation, err := ries.GetByID(riskID, id)

	if err != nil {
		return err
	}

	// Update Risk Inherent Evaluation
	existingRiskInherentEvaluation.Question = inherentEvaluationFormDTO.Question
	existingRiskInherentEvaluation.Answer = inherentEvaluationFormDTO.Answer
	existingRiskInherentEvaluation.UpdatedBy = &currentUserID

	err = ries.repo.Update(existingRiskInherentEvaluation)

	if err != nil {
		return err
	}

	return nil
}

func (ries *RiskInherentEvaluationService) Delete(riskID uint, id uint, currentUserID uint) (uint, error) {
	tx := ries.repo.BeginTransaction()

	if _, err := ries.riskRegisterSvc.GetByID(riskID); err != nil {
		return 0, err
	}

	var failedId uint

	entryToDelete, err := ries.GetByID(riskID, id)
	if err != nil {
		return 0, err
	}

	// Delete the entry
	err = ries.repo.Delete(id, tx)
	if err != nil {
		tx.Rollback()
		failedId = id
		return failedId, err
	}

	// Reassign "Why" values for subsequent entries
	err = ries.repo.UpdateWhyValuesAfterDeletion(riskID, entryToDelete.Why, tx)
	if err != nil {
		tx.Rollback()
		fmt.Printf("Failed to update 'Why' values after deletion: %v\n", err)
		return 0, err
	}

	// Log successful delete activity
	go func() {
		err := ries.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.RiskEvaluation{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
