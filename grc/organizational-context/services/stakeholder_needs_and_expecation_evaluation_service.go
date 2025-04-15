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

type StakeholderNeedsAndExpectationEvaluationServiceInterface interface {
	GetAll(StakeholderID uint, NeedsAndExpectationID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error)
	Create(StakeholderID uint, NeedsAndExpectationID uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) (uint, error)
	Update(StakeholderID uint, NeedsAndExpectationID uint, id uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) error
	GetLatestStakeholderNeedsAndExpectationEvaluation(StakeholderID uint, actionPlanID uint) (*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error)
}

type StakeholderNeedsAndExpectationEvaluationService struct {
	repo                              repositories.StakeholderNeedsAndExpectationEvaluationRepositoryInterface
	activityLogSvc                    logServices.ActivityLogServiceInterface
	stakeholderSvc                    StakeholderServiceInterface
	stakeholderNeedsAndExpectationSvc StakeholderNeedsAndExpectationServiceInterface
	evaluationStatusSvc               EvaluationStatusServiceInterface
}

func NewStakeholderNeedsAndExpectationEvaluationService(
	repo repositories.StakeholderNeedsAndExpectationEvaluationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	stakeholderSvc StakeholderServiceInterface,
	stakeholderNeedsAndExpectationSvc StakeholderNeedsAndExpectationServiceInterface,
	evaluationStatusSvc EvaluationStatusServiceInterface,
) StakeholderNeedsAndExpectationEvaluationServiceInterface {
	return &StakeholderNeedsAndExpectationEvaluationService{
		repo:                              repo,
		activityLogSvc:                    activityLogSvc,
		stakeholderSvc:                    stakeholderSvc,
		stakeholderNeedsAndExpectationSvc: stakeholderNeedsAndExpectationSvc,
		evaluationStatusSvc:               evaluationStatusSvc,
	}
}

func (snees *StakeholderNeedsAndExpectationEvaluationService) GetAll(StakeholderID uint, NeedsAndExpectationID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error) {
	if _, err := snees.stakeholderNeedsAndExpectationSvc.GetByID(StakeholderID, NeedsAndExpectationID); err != nil {
		return nil, err
	}

	return snees.repo.GetAll(StakeholderID, NeedsAndExpectationID, sort, sortBy, keyword)
}

func (snees StakeholderNeedsAndExpectationEvaluationService) GetByID(StakeholderID uint, NeedsAndExpectationID uint, id uint) (*models.StakeholderNeedsAndExpectationEvaluation, error) {
	if _, err := snees.stakeholderNeedsAndExpectationSvc.GetByID(StakeholderID, NeedsAndExpectationID); err != nil {
		return nil, err
	}

	evaluation, err := snees.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return evaluation, nil
}

func (snees *StakeholderNeedsAndExpectationEvaluationService) Create(StakeholderID uint, NeedsAndExpectationID uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) (uint, error) {
	tx := snees.repo.BeginTransaction()

	if _, err := snees.stakeholderNeedsAndExpectationSvc.GetByID(StakeholderID, NeedsAndExpectationID); err != nil {
		return 0, err
	}

	if err := snees.validateStakeholderNeedsAndExpectationEvaluationFormDTO(evaluationFormDTO); err != nil {
		return 0, err
	}

	dates, err := snees.parseEvaluationDates(evaluationFormDTO)
	if err != nil {
		return 0, err
	}

	evaluation := &models.StakeholderNeedsAndExpectationEvaluation{
		StakeholderNeedsAndExpectationID: NeedsAndExpectationID,
		EvaluationStatusID:               evaluationFormDTO.EvaluationStatusID,
		Comment:                          evaluationFormDTO.Comment,
		Score:                            evaluationFormDTO.Score,
		EvaluatedBYID:                    evaluationFormDTO.EvaluatedByID,
		EvaluatedAt:                      dates.EvaluatedAt,
		CreatedBYID:                      currentUserID,
	}

	id, err := snees.repo.Create(evaluation)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = snees.stakeholderNeedsAndExpectationSvc.StatusUpdate(StakeholderID, NeedsAndExpectationID, currentUserID, evaluationFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := snees.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderNeedsAndExpectationEvaluation{}.TableName(),
			id,
			evaluation,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderNeedsAndExpectationEvaluation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (snees *StakeholderNeedsAndExpectationEvaluationService) Update(StakeholderID uint, NeedsAndExpectationID uint, id uint, currentUserID uint, evaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) error {
	tx := snees.repo.BeginTransaction()

	if _, err := snees.stakeholderNeedsAndExpectationSvc.GetByID(StakeholderID, NeedsAndExpectationID); err != nil {
		return err
	}

	oldStakeholderNeedsAndExpectationEvaluation, err := snees.GetByID(StakeholderID, NeedsAndExpectationID, id)
	if err != nil {
		return err
	}

	existingStakeholderNeedsAndExpectationEvaluation, err := snees.GetByID(StakeholderID, NeedsAndExpectationID, id)
	if err != nil {
		return err
	}

	if err := snees.validateStakeholderNeedsAndExpectationEvaluationFormDTO(evaluationFormDTO); err != nil {
		return err
	}

	//@todo validate docuemnt form

	existingStakeholderNeedsAndExpectationEvaluation.EvaluationStatusID = evaluationFormDTO.EvaluationStatusID
	existingStakeholderNeedsAndExpectationEvaluation.Comment = evaluationFormDTO.Comment
	existingStakeholderNeedsAndExpectationEvaluation.Score = evaluationFormDTO.Score
	existingStakeholderNeedsAndExpectationEvaluation.UpdatedBYID = &currentUserID

	err = snees.repo.Update(existingStakeholderNeedsAndExpectationEvaluation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = snees.stakeholderNeedsAndExpectationSvc.StatusUpdate(StakeholderID, NeedsAndExpectationID, currentUserID, evaluationFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlanstatusUpdate, err := snees.GetByID(StakeholderID, NeedsAndExpectationID, id)
		if err != nil {
			return
		}
		// log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"comment": oldStakeholderNeedsAndExpectationEvaluation.Comment,
			"score":   oldStakeholderNeedsAndExpectationEvaluation.Score,
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
			err := snees.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderNeedsAndExpectationEvaluation{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderNeedsAndExpectationEvaluation{}.TableName(),
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

func (snees *StakeholderNeedsAndExpectationEvaluationService) validateStakeholderNeedsAndExpectationEvaluationFormDTO(evaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) error {
	if _, err := snees.evaluationStatusSvc.GetByID(evaluationFormDTO.EvaluationStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "evaluation status")
	}

	return nil
}

func (capsus StakeholderNeedsAndExpectationEvaluationService) GetLatestStakeholderNeedsAndExpectationEvaluation(StakeholderID uint, actionPlanID uint) (*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error) {
	if _, err := capsus.stakeholderNeedsAndExpectationSvc.GetByID(StakeholderID, actionPlanID); err != nil {
		return nil, err
	}

	return capsus.repo.GetLatestStakeholderNeedsAndExpectationEvaluation(actionPlanID)
}

func (taps *StakeholderNeedsAndExpectationEvaluationService) parseEvaluationDates(dto *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) (*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationParsedDates, error) {
	parsed := &organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationParsedDates{}

	var err error
	if *dto.EvaluatedAt != "" {
		parsed.EvaluatedAt, err = utility.ParseDate(*dto.EvaluatedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid evaluated date"), "invalid_date", "evaluated at")
		}
	}

	return parsed, nil
}
