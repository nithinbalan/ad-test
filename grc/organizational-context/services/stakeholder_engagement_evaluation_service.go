package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderEngagementEvaluationServiceInterface interface {
	GetAll(StakeholderID uint, EngagementID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error)
	Create(StakeholderID uint, EngagementID uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO) (uint, error)
	Update(StakeholderID uint, EngagementID uint, id uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO) error
	GetLatestStakeholderEngagementEvaluation(StakeholderID uint, actionPlanID uint) (*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error)
}

type StakeholderEngagementEvaluationService struct {
	repo                     repositories.StakeholderEngagementEvaluationRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	stakeholderSvc           StakeholderServiceInterface
	stakeholderEngagementSvc StakeholderEngagementServiceInterface
	evaluationStatusSvc      EvaluationStatusServiceInterface
}

func NewStakeholderEngagementEvaluationService(
	repo repositories.StakeholderEngagementEvaluationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	stakeholderSvc StakeholderServiceInterface,
	stakeholderEngagementSvc StakeholderEngagementServiceInterface,
	evaluationStatusSvc EvaluationStatusServiceInterface,
) StakeholderEngagementEvaluationServiceInterface {
	return &StakeholderEngagementEvaluationService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		stakeholderSvc:           stakeholderSvc,
		stakeholderEngagementSvc: stakeholderEngagementSvc,
		evaluationStatusSvc:      evaluationStatusSvc,
	}
}

func (casus *StakeholderEngagementEvaluationService) GetAll(StakeholderID uint, EngagementID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error) {
	if _, err := casus.stakeholderEngagementSvc.GetByID(StakeholderID, EngagementID); err != nil {
		return nil, err
	}

	return casus.repo.GetAll(StakeholderID, EngagementID, sort, sortBy, keyword)
}

func (casus StakeholderEngagementEvaluationService) GetByID(StakeholderID uint, EngagementID uint, id uint) (*models.StakeholderEngagementEvaluation, error) {
	if _, err := casus.stakeholderEngagementSvc.GetByID(StakeholderID, EngagementID); err != nil {
		return nil, err
	}

	evaluation, err := casus.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return evaluation, nil
}

func (casus *StakeholderEngagementEvaluationService) Create(StakeholderID uint, EngagementID uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO) (uint, error) {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.stakeholderEngagementSvc.GetByID(StakeholderID, EngagementID); err != nil {
		return 0, err
	}

	dates, err := casus.parseEvaluationDates(evaluationFormDTO)
	if err != nil {
		return 0, err
	}

	evaluation := &models.StakeholderEngagementEvaluation{
		StakeholderEngagementID: EngagementID,
		Comment:                 evaluationFormDTO.Comment,
		Score:                   evaluationFormDTO.Score,
		EvaluatedBYID:           evaluationFormDTO.EvaluatedByID,
		EvaluatedAt:             dates.EvaluatedAt,
		CreatedBYID:             currentUserID,
	}

	id, err := casus.repo.Create(evaluation)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = casus.stakeholderEngagementSvc.StatusUpdate(StakeholderID, EngagementID, currentUserID, evaluationFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := casus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderEngagementEvaluation{}.TableName(),
			id,
			evaluation,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderEngagementEvaluation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (casus *StakeholderEngagementEvaluationService) Update(StakeholderID uint, EngagementID uint, id uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO) error {
	tx := casus.repo.BeginTransaction()

	if _, err := casus.stakeholderEngagementSvc.GetByID(StakeholderID, EngagementID); err != nil {
		return err
	}

	oldStakeholderEngagementEvaluation, err := casus.GetByID(StakeholderID, EngagementID, id)
	if err != nil {
		return err
	}

	existingStakeholderEngagementEvaluation, err := casus.GetByID(StakeholderID, EngagementID, id)
	if err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingStakeholderEngagementEvaluation.Comment = evaluationFormDTO.Comment
	existingStakeholderEngagementEvaluation.Score = evaluationFormDTO.Score
	existingStakeholderEngagementEvaluation.UpdatedBYID = &currentUserID

	err = casus.repo.Update(existingStakeholderEngagementEvaluation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = casus.stakeholderEngagementSvc.StatusUpdate(StakeholderID, EngagementID, currentUserID, evaluationFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := casus.GetByID(StakeholderID, EngagementID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment": oldStakeholderEngagementEvaluation.Comment,
			"score":   oldStakeholderEngagementEvaluation.Score,
		}

		newData := map[string]interface{}{
			"comment": existingActionPlanstatusUpdate.Comment,
			"score":   existingActionPlanstatusUpdate.Score,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := casus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderEngagementEvaluation{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderEngagementEvaluation{}.TableName(),
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

func (capsus StakeholderEngagementEvaluationService) GetLatestStakeholderEngagementEvaluation(StakeholderID uint, actionPlanID uint) (*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error) {
	if _, err := capsus.stakeholderEngagementSvc.GetByID(StakeholderID, actionPlanID); err != nil {
		return nil, err
	}

	return capsus.repo.GetLatestStakeholderEngagementEvaluation(actionPlanID)
}

func (taps *StakeholderEngagementEvaluationService) parseEvaluationDates(dto *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO) (*organizationalContextDTOs.StakeholderEngagementEvaluationParsedDates, error) {
	parsed := &organizationalContextDTOs.StakeholderEngagementEvaluationParsedDates{}

	var err error
	if *dto.EvaluatedAt != "" {
		parsed.EvaluatedAt, err = utility.ParseDate(*dto.EvaluatedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid evaluated date"), "invalid_date", "evaluated at")
		}
	}

	return parsed, nil
}
