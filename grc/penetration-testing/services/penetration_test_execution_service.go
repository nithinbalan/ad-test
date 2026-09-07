package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	penetrationTestRepositories "jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestExecutionServiceInterface interface {
	GetByID(penetrationTestID uint) (*models.PenetrationTestExecution, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestExecutionDetailDTO, error)
	Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestExecutionFormDTO) (uint, error)
	Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestExecutionFormDTO) error
	Delete(penetrationTestID uint, currentUserID uint) (uint, error)
}

type PenetrationTestExecutionService struct {
	repo               penetrationTestRepositories.PenetrationTestExecutionRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	penetrationTestSvc PenetrationTestServiceInterface
	executionStatusSvc PenetrationTestExecutionStatusServiceInterface
	securityResultSvc  PenetrationTestSecurityResultServiceInterface
}

func NewPenetrationTestExecutionService(
	repo penetrationTestRepositories.PenetrationTestExecutionRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
	executionStatusSvc PenetrationTestExecutionStatusServiceInterface,
	securityResultSvc PenetrationTestSecurityResultServiceInterface,
) PenetrationTestExecutionServiceInterface {
	return &PenetrationTestExecutionService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		penetrationTestSvc: penetrationTestSvc,
		executionStatusSvc: executionStatusSvc,
		securityResultSvc:  securityResultSvc,
	}
}

func (ptes PenetrationTestExecutionService) GetByID(penetrationTestID uint) (*models.PenetrationTestExecution, error) {
	if _, err := ptes.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	execution, err := ptes.repo.GetByID(penetrationTestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}
	return execution, nil
}

func (ptes PenetrationTestExecutionService) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestExecutionDetailDTO, error) {
	var executionDTO penetrationTestDTOs.PenetrationTestExecutionDetailDTO

	if _, err := ptes.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	execution, err := ptes.repo.GetFullDetailsByID(penetrationTestID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}

	err = mapstructure.Decode(execution, &executionDTO)
	if err != nil {
		return nil, err
	}

	testStartDate := execution.TestStartDate.Format("2006-01-02")
	executionDTO.TestStartDateString = &testStartDate

	testEndDate := execution.TestEndDate.Format("2006-01-02")
	executionDTO.TestEndDateString = &testEndDate

	executionDTO.CreatedAt = execution.CreatedAt
	executionDTO.UpdatedAt = execution.UpdatedAt

	return &executionDTO, nil
}

func (ptes *PenetrationTestExecutionService) Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestExecutionFormDTO) (uint, error) {
	tx := ptes.repo.BeginTransaction()

	if _, err := ptes.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	if _, err := ptes.executionStatusSvc.GetByID(form.ExecutionStatusID); err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid execution status"), "invalid", "execution_status_id")
	}

	if _, err := ptes.securityResultSvc.GetByID(form.SecurityResultID); err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid security result"), "invalid", "security_result_id")
	}

	testStartDate, err := utility.ParseMandatoryDate(form.TestStartDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid test start date"), "invalid_date", "test_start_date")
	}

	testEndDate, err := utility.ParseMandatoryDate(form.TestEndDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid test end date"), "invalid_date", "test_end_date")
	}

	if testEndDate.Before(testStartDate) {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("test end date cannot be before test start date"), "test_end_before_start", "test_end_date")
	}

	execution := &models.PenetrationTestExecution{
		PenetrationTestID: penetrationTestID,
		TestStartDate:     testStartDate,
		TestEndDate:       testEndDate,
		ExecutionStatusID: form.ExecutionStatusID,
		SecurityResultID:  form.SecurityResultID,
		CreatedBYID:       currentUserID,
	}

	id, err := ptes.repo.Create(execution, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingExecution, err := ptes.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test execution", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptes.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestExecution{}.TableName(),
			id,
			existingExecution,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "execution", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.PenetrationTestExecution{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ptes *PenetrationTestExecutionService) Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestExecutionFormDTO) error {
	tx := ptes.repo.BeginTransaction()

	if _, err := ptes.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return err
	}

	existingExecution, err := ptes.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := ptes.executionStatusSvc.GetByID(form.ExecutionStatusID); err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid execution status"), "invalid", "execution_status_id")
	}

	if _, err := ptes.securityResultSvc.GetByID(form.SecurityResultID); err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid security result"), "invalid", "security_result_id")
	}

	testStartDate, err := utility.ParseMandatoryDate(form.TestStartDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid test start date"), "invalid_date", "test_start_date")
	}

	testEndDate, err := utility.ParseMandatoryDate(form.TestEndDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid test end date"), "invalid_date", "test_end_date")
	}

	if testEndDate.Before(testStartDate) {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("test end date cannot be before test start date"), "test_end_before_start", "test_end_date")
	}

	existingExecution.TestStartDate = testStartDate
	existingExecution.TestEndDate = testEndDate
	existingExecution.ExecutionStatusID = form.ExecutionStatusID
	existingExecution.SecurityResultID = form.SecurityResultID
	existingExecution.UpdatedBYID = &currentUserID

	err = ptes.repo.Update(existingExecution, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingEntry, err := ptes.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test execution", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptes.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestExecution{}.TableName(),
			existingExecution.ID,
			existingEntry,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "execution", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.PenetrationTestExecution{}.TableName(),
				"recordID":  existingExecution.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ptes *PenetrationTestExecutionService) Delete(penetrationTestID uint, currentUserID uint) (uint, error) {
	tx := ptes.repo.BeginTransaction()
	var failedId uint
	var deletedExecution []map[string]interface{}

	if _, err := ptes.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	execution, err := ptes.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = ptes.repo.Delete(execution.ID, tx); err != nil {
		tx.Rollback()
		failedId = execution.ID
		return failedId, err
	}

	deletedExecution = append(deletedExecution, map[string]interface{}{
		"id":                  execution.ID,
		"penetration_test_id": penetrationTestID,
	})

	tx.Commit()

	go func() {
		err := ptes.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PenetrationTestExecution{}.TableName(),
			execution.ID,
			deletedExecution,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "execution", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
