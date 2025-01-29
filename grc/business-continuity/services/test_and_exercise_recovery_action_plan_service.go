package services

import (
	"errors"
	"fmt"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	TestAndExerciseRecoveryActionPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TestAndExerciseRecoveryActionPlanServiceInterface interface {
	GetAll(testAndExerciseId uint, params *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListDTO, error)
	GetByID(testAndExerciseID uint, id uint) (*models.TestAndExerciseRecoveryActionPlan, error)
	GetFullDetailsByID(testAndExerciseID uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanDetailDTO, error)
	Create(testAndExerciseID uint, currentUserID uint, dto *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) (uint, error)
	Update(testAndExerciseID uint, id uint, currentUserID uint, dto *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) error
	Delete(testAndExerciseID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type TestAndExerciseRecoveryActionPlanService struct {
	repo                TestAndExerciseRecoveryActionPlanRepositories.TestAndExerciseRecoveryActionPlanRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	userSvc             organizationServices.UserServiceInterface
	testAndExerciseSvc  TestAndExerciseServiceInterface
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
}

func NewTestAndExerciseRecoveryActionPlanService(
	repo TestAndExerciseRecoveryActionPlanRepositories.TestAndExerciseRecoveryActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	testAndExerciseSvc TestAndExerciseServiceInterface,
	ActionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) TestAndExerciseRecoveryActionPlanServiceInterface {
	return &TestAndExerciseRecoveryActionPlanService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		userSvc:             userSvc,
		testAndExerciseSvc:  testAndExerciseSvc,
		ActionPlanStatusSvc: ActionPlanStatusSvc,
	}
}

func (teraps *TestAndExerciseRecoveryActionPlanService) GetAll(testAndExerciseId uint, params *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListDTO, error) {
	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseId); err != nil {
		return nil, err
	}
	return teraps.repo.GetAll(testAndExerciseId, params)
}

func (teraps *TestAndExerciseRecoveryActionPlanService) GetByID(testAndExerciseID uint, id uint) (*models.TestAndExerciseRecoveryActionPlan, error) {
	actionPlan, err := teraps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return actionPlan, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) GetFullDetailsByID(testAndExerciseID uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanDetailDTO, error) {
	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseID); err != nil {
		return nil, err
	}
	var recoveryActionPlan businessContinuityDTOs.TestAndExerciseRecoveryActionPlanDetailDTO

	actionPlan, err := teraps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(actionPlan, &recoveryActionPlan)
	if err != nil {
		return nil, err
	}

	if actionPlan.TargetDate != nil {
		targetDate := actionPlan.TargetDate.Format("2006-01-02")
		recoveryActionPlan.TargetDateString = &targetDate
	}

	recoveryActionPlan.CreatedAt = actionPlan.CreatedAt
	recoveryActionPlan.UpdatedAt = actionPlan.UpdatedAt

	return &recoveryActionPlan, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) Create(testAndExerciseID uint, currentUserID uint, recoveryActionPlanFormDTO *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) (uint, error) {
	tx := teraps.repo.BeginTransaction()

	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseID); err != nil {
		return 0, err
	}

	responsibleUnitID, err := teraps.userSvc.GetUserPrimaryUnitID(recoveryActionPlanFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(errors.New("organization unit not set for responsible user"), "not_exist", "responsible unit")
	}

	actionPlanStatus, err := teraps.ActionPlanStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	dates, err := teraps.parseTestAndExerciseRecoveryActionPlanDates(recoveryActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	actionPlan := &models.TestAndExerciseRecoveryActionPlan{
		TestAndExerciseID:  testAndExerciseID,
		BcpPhaseID:         recoveryActionPlanFormDTO.BcpPhaseID,
		ProcessID:          recoveryActionPlanFormDTO.ProcessID,
		SeverityLevelID:    recoveryActionPlanFormDTO.SeverityLevelID,
		Title:              recoveryActionPlanFormDTO.Title,
		ResponsibleUserID:  recoveryActionPlanFormDTO.ResponsibleUserID,
		ResponsibleUnitID:  responsibleUnitID,
		TargetDate:         dates.TargetDate,
		ActionPlanStatusId: &actionPlanStatus.ID,
		IsRecovery:         recoveryActionPlanFormDTO.IsRecovery,
		Percentage:         recoveryActionPlanFormDTO.Percentage,
		CreatedBYID:        currentUserID,
	}

	id, err := teraps.repo.Create(actionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingTestAndExerciseRecoveryActionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, actionPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = teraps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TestAndExerciseRecoveryActionPlan{}.TableName(),
			id,
			existingTestAndExerciseRecoveryActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TestAndExerciseRecoveryActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) Update(testAndExerciseID uint, id uint, currentUserID uint, recoveryActionPlanFormDTO *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) error {
	tx := teraps.repo.BeginTransaction()

	existingRecoveryPlan, err := teraps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldRecoveryActionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingRecoveryPlan.ResponsibleUserID != recoveryActionPlanFormDTO.ResponsibleUserID {
		responsibleUnitID, err := teraps.userSvc.GetUserPrimaryUnitID(existingRecoveryPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingRecoveryPlan.ResponsibleUnitID = responsibleUnitID
	}

	dates, err := teraps.parseTestAndExerciseRecoveryActionPlanDates(recoveryActionPlanFormDTO)
	if err != nil {
		return err
	}

	existingRecoveryPlan.BcpPhaseID = recoveryActionPlanFormDTO.BcpPhaseID
	existingRecoveryPlan.ProcessID = recoveryActionPlanFormDTO.ProcessID
	existingRecoveryPlan.SeverityLevelID = recoveryActionPlanFormDTO.SeverityLevelID
	existingRecoveryPlan.Title = recoveryActionPlanFormDTO.Title
	existingRecoveryPlan.ResponsibleUserID = recoveryActionPlanFormDTO.ResponsibleUserID
	existingRecoveryPlan.TargetDate = dates.TargetDate
	existingRecoveryPlan.IsRecovery = recoveryActionPlanFormDTO.IsRecovery
	existingRecoveryPlan.Percentage = recoveryActionPlanFormDTO.Percentage
	existingRecoveryPlan.UpdatedBYID = &currentUserID

	err = teraps.repo.Update(existingRecoveryPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingActionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", oldRecoveryActionPlan.ResponsibleUser.FirstNameEN, oldRecoveryActionPlan.ResponsibleUser.LastNameEN),
			"bcp_phase":        oldRecoveryActionPlan.BcpPhase.TitleEn,
			"process":          oldRecoveryActionPlan.Process.TitleEn,
			"severity_level":   oldRecoveryActionPlan.SeverityLevel.TitleEn,
			"target_date":      oldRecoveryActionPlan.TargetDateString,
		}

		newData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"bcp_phase":        existingActionPlan.BcpPhase.TitleEn,
			"process":          existingActionPlan.Process.TitleEn,
			"severity_level":   existingActionPlan.SeverityLevel.TitleEn,
			"target_date":      existingActionPlan.TargetDateString,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := teraps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.TestAndExerciseRecoveryActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TestAndExerciseRecoveryActionPlan{}.TableName(),
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

func (teraps *TestAndExerciseRecoveryActionPlanService) Delete(testAndExerciseID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := teraps.repo.BeginTransaction()

	if _, err := teraps.testAndExerciseSvc.GetByID(testAndExerciseID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedActionPlans []map[string]interface{}

	for _, id := range ids {
		actionPlan, err := teraps.GetFullDetailsByID(testAndExerciseID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedActionPlans = append(deletedActionPlans, map[string]interface{}{
			"id":                id,
			"test_and_exercise": actionPlan.ID,
		})

		if err := teraps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	if len(successfulIDs) > 0 {
		go func() {
			err := teraps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.TestAndExerciseRecoveryActionPlan{}.TableName(),
				successfulIDs,
				deletedActionPlans,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}

func (teraps *TestAndExerciseRecoveryActionPlanService) parseTestAndExerciseRecoveryActionPlanDates(dto *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanFormDTO) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanParsedDates, error) {
	parsed := &businessContinuityDTOs.TestAndExerciseRecoveryActionPlanParsedDates{}

	var err error

	if &dto.TargetDate != nil {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
