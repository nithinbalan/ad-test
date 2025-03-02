package services

import (
	"errors"
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	TestAndExerciseCommunicationActionPlanRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TestAndExerciseCommunicationActionPlanServiceInterface interface {
	GetAll(testAndExerciseId uint, params *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListDTO, error)
	GetFullDetailsByID(testAndExerciseId uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanDetailDTO, error)
	GetByID(testAndExerciseId uint, id uint) (*models.TestAndExerciseCommunicationActionPlan, error)
	Create(testAndExerciseId uint, currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanFormDTO) (uint, error)
	Update(testAndExerciseId uint, id uint, currentUserID uint, TestAndExerciseCommunicationActionPlan *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanFormDTO) error
	Delete(testAndExerciseId uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type TestAndExerciseCommunicationActionPlanService struct {
	repo                             TestAndExerciseCommunicationActionPlanRepositories.TestAndExerciseCommunicationActionPlanRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	communicationActionPlanStatusSvc CommunicationActionPlanStatusServiceInterface
}

func NewTestAndExerciseCommunicationActionPlanService(
	repo TestAndExerciseCommunicationActionPlanRepositories.TestAndExerciseCommunicationActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	communicationActionPlanStatusSvc CommunicationActionPlanStatusServiceInterface,
) TestAndExerciseCommunicationActionPlanServiceInterface {
	return &TestAndExerciseCommunicationActionPlanService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		userSvc:                          userSvc,
		communicationActionPlanStatusSvc: communicationActionPlanStatusSvc,
	}
}

func (tecaps *TestAndExerciseCommunicationActionPlanService) GetAll(testAndExerciseId uint, params *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListDTO, error) {
	return tecaps.repo.GetAll(testAndExerciseId, params)
}

func (tecaps TestAndExerciseCommunicationActionPlanService) GetByID(testAndExerciseId uint, id uint) (*models.TestAndExerciseCommunicationActionPlan, error) {
	TestAndExerciseCommunicationActionPlan, err := tecaps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return TestAndExerciseCommunicationActionPlan, nil
}

func (tecaps TestAndExerciseCommunicationActionPlanService) GetFullDetailsByID(testAndExerciseId uint, id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanDetailDTO, error) {
	// Initialize your DTO
	var testAndExerciseDTO businessContinuityDTOs.TestAndExerciseCommunicationActionPlanDetailDTO

	testAndExercise, err := tecaps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(testAndExercise, &testAndExerciseDTO)
	if err != nil {
		return nil, err
	}

	testAndExerciseDTO.CreatedAt = testAndExercise.CreatedAt
	testAndExerciseDTO.UpdatedAt = testAndExercise.UpdatedAt

	return &testAndExerciseDTO, nil
}

func (tecaps *TestAndExerciseCommunicationActionPlanService) Create(testAndExerciseId uint, currentUserID uint, testAndExerciseFormDTO *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanFormDTO) (uint, error) {
	tx := tecaps.repo.BeginTransaction()

	responsibleUnitID, err := tecaps.userSvc.GetUserPrimaryUnitID(testAndExerciseFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "test and exercise responsible unit")
	}

	communicationActionPlanStatus, err := tecaps.communicationActionPlanStatusSvc.GetByType("pending")
	if err != nil {
		return 0, err
	}

	TestAndExerciseCommunicationActionPlan := &models.TestAndExerciseCommunicationActionPlan{
		TestAndExerciseID:               testAndExerciseId,
		BcpPhaseID:                      testAndExerciseFormDTO.BcpPhaseID,
		ProcessID:                       testAndExerciseFormDTO.ProcessID,
		Title:                           testAndExerciseFormDTO.Title,
		ResponsibleUserID:               testAndExerciseFormDTO.ResponsibleUserID,
		ResponsibleUnitID:               responsibleUnitID,
		Frequency:                       testAndExerciseFormDTO.Frequency,
		Recipient:                       testAndExerciseFormDTO.Recipient,
		Channel:                         testAndExerciseFormDTO.Channel,
		Description:                     testAndExerciseFormDTO.Description,
		CommunicationActionPlanStatusID: &communicationActionPlanStatus.ID,
		CreatedBYID:                     currentUserID,
	}

	id, err := tecaps.repo.Create(TestAndExerciseCommunicationActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingTestAndExerciseCommunicationActionPlan, err := tecaps.GetFullDetailsByID(testAndExerciseId, TestAndExerciseCommunicationActionPlan.ID)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = tecaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TestAndExerciseCommunicationActionPlan{}.TableName(),
			id,
			existingTestAndExerciseCommunicationActionPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TestAndExerciseCommunicationActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (tecaps *TestAndExerciseCommunicationActionPlanService) Update(testAndExerciseId uint, id uint, currentUserID uint, testAndExerciseFormDTO *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanFormDTO) error {
	tx := tecaps.repo.BeginTransaction()

	existingTestAndExerciseCommunicationActionPlan, err := tecaps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldTestAndExerciseCommunicationActionPlan, err := tecaps.GetFullDetailsByID(testAndExerciseId, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingTestAndExerciseCommunicationActionPlan.ResponsibleUserID != testAndExerciseFormDTO.ResponsibleUserID {
		responsibleUnitID, err := tecaps.userSvc.GetUserPrimaryUnitID(existingTestAndExerciseCommunicationActionPlan.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingTestAndExerciseCommunicationActionPlan.ResponsibleUnitID = responsibleUnitID
	}

	existingTestAndExerciseCommunicationActionPlan.TestAndExerciseID = testAndExerciseId
	existingTestAndExerciseCommunicationActionPlan.BcpPhaseID = testAndExerciseFormDTO.BcpPhaseID
	existingTestAndExerciseCommunicationActionPlan.ProcessID = testAndExerciseFormDTO.ProcessID
	existingTestAndExerciseCommunicationActionPlan.Title = testAndExerciseFormDTO.Title
	existingTestAndExerciseCommunicationActionPlan.ResponsibleUserID = testAndExerciseFormDTO.ResponsibleUserID
	existingTestAndExerciseCommunicationActionPlan.Frequency = testAndExerciseFormDTO.Frequency
	existingTestAndExerciseCommunicationActionPlan.Channel = testAndExerciseFormDTO.Channel
	existingTestAndExerciseCommunicationActionPlan.Recipient = testAndExerciseFormDTO.Recipient
	existingTestAndExerciseCommunicationActionPlan.Description = testAndExerciseFormDTO.Description
	existingTestAndExerciseCommunicationActionPlan.UpdatedBYID = &currentUserID

	err = tecaps.repo.Update(existingTestAndExerciseCommunicationActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingTestAndExerciseCommunicationActionPlan, err := tecaps.GetFullDetailsByID(testAndExerciseId, id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", oldTestAndExerciseCommunicationActionPlan.ResponsibleUser.FirstNameEN, oldTestAndExerciseCommunicationActionPlan.ResponsibleUser.LastNameEN),
			"bcp_phase":        oldTestAndExerciseCommunicationActionPlan.BcpPhase.TitleEn,
			"process":          oldTestAndExerciseCommunicationActionPlan.Process.TitleEn,
			"description":      oldTestAndExerciseCommunicationActionPlan.Description,
			"frequency":        oldTestAndExerciseCommunicationActionPlan.Frequency,
			"recipient":        oldTestAndExerciseCommunicationActionPlan.Recipient,
			"channel":          oldTestAndExerciseCommunicationActionPlan.Channel,
		}

		newData := map[string]interface{}{
			"responsible_user": fmt.Sprintf("%s %s", existingTestAndExerciseCommunicationActionPlan.ResponsibleUser.FirstNameEN, existingTestAndExerciseCommunicationActionPlan.ResponsibleUser.LastNameEN),
			"bcp_phase":        existingTestAndExerciseCommunicationActionPlan.BcpPhase.TitleEn,
			"process":          existingTestAndExerciseCommunicationActionPlan.Process.TitleEn,
			"description":      existingTestAndExerciseCommunicationActionPlan.Description,
			"frequency":        existingTestAndExerciseCommunicationActionPlan.Frequency,
			"recipient":        existingTestAndExerciseCommunicationActionPlan.Recipient,
			"channel":          existingTestAndExerciseCommunicationActionPlan.Channel,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := tecaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.TestAndExerciseCommunicationActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TestAndExerciseCommunicationActionPlan{}.TableName(),
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

func (tecaps *TestAndExerciseCommunicationActionPlanService) Delete(testAndExerciseId uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := tecaps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedTestAndExerciseCommunicationActionPlan []map[string]interface{}

	for _, id := range ids {
		testAndExercise, err := tecaps.GetFullDetailsByID(testAndExerciseId, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedTestAndExerciseCommunicationActionPlan = append(deletedTestAndExerciseCommunicationActionPlan, map[string]interface{}{
			"id":                id,
			"test_and_exercise": testAndExercise.ID,
		})

		// Delete the test and exercise
		if err := tecaps.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := tecaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.TestAndExerciseCommunicationActionPlan{}.TableName(),
				successfulIds,
				deletedTestAndExerciseCommunicationActionPlan,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

		}()
	}

	return failedIds, nil
}
