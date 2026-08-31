package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestServiceInterface interface {
	GetAll(params *penetrationTestDTOs.PenetrationTestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestListResponseDTO, error)
	GetByID(id uint) (*models.PenetrationTest, error)
	GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestDetailDTO, error)
	Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type PenetrationTestService struct {
	repo           repositories.PenetrationTestRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestService(repo repositories.PenetrationTestRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestServiceInterface {
	return &PenetrationTestService{repo: repo, activityLogSvc: activityLogSvc}
}

func (ptls *PenetrationTestService) GetAll(params *penetrationTestDTOs.PenetrationTestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestListResponseDTO, error) {
	return ptls.repo.GetAll(params)
}

func (ptls *PenetrationTestService) GetByID(id uint) (*models.PenetrationTest, error) {
	value, err := ptls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (ptls *PenetrationTestService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestDetailDTO, error) {
	value, err := ptls.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &penetrationTestDTOs.PenetrationTestDetailDTO{}
	if err := mapstructure.Decode(value, detail); err != nil {
		return nil, err
	}
	detail.CreatedAt = value.CreatedAt
	detail.UpdatedAt = value.UpdatedAt
	if value.EnvironmentPreparation != nil {
		detail.EnvironmentNotifyDate = value.EnvironmentPreparation.EnvironmentNotifyDate
		detail.ExpectedReadyDate = value.EnvironmentPreparation.ExpectedReadyDate
		detail.EnvironmentReadiness = penetrationTestDTOs.PenetrationTestEnvironmentReadinessDTO{
			ID:      value.EnvironmentPreparation.EnvironmentReadiness.ID,
			TitleEn: value.EnvironmentPreparation.EnvironmentReadiness.TitleEn,
			TitleAr: value.EnvironmentPreparation.EnvironmentReadiness.TitleAr,
			Status:  value.EnvironmentPreparation.EnvironmentReadiness.Status,
		}
	}
	if value.Execution != nil {
		detail.TestStartDate = value.Execution.TestStartDate
		detail.TestEndDate = value.Execution.TestEndDate
		detail.ExecutionStatus = penetrationTestDTOs.PenetrationTestExecutionStatusDTO{
			ID:      value.Execution.ExecutionStatus.ID,
			TitleEn: value.Execution.ExecutionStatus.TitleEn,
			TitleAr: value.Execution.ExecutionStatus.TitleAr,
			Type:    value.Execution.ExecutionStatus.Type,
			Status:  value.Execution.ExecutionStatus.Status,
		}
		detail.SecurityResult = penetrationTestDTOs.PenetrationTestSecurityResultDTO{
			ID:      value.Execution.SecurityResult.ID,
			TitleEn: value.Execution.SecurityResult.TitleEn,
			TitleAr: value.Execution.SecurityResult.TitleAr,
			Type:    value.Execution.SecurityResult.Type,
			Status:  value.Execution.SecurityResult.Status,
		}
	}
	if value.ManagementSharing != nil {
		detail.ResultsSharedWithManagement = value.ManagementSharing.ResultsSharedWithManagement
		detail.ManagementNotifyDate = value.ManagementSharing.ManagementNotifyDate
		detail.ManagementResponseDate = value.ManagementSharing.ManagementResponseDate
		detail.RemediationNotes = value.ManagementSharing.RemediationNotes
	}
	return detail, nil
}

func (ptls *PenetrationTestService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFormDTO) (uint, error) {
	tx := ptls.repo.BeginTransaction()
	environmentNotifyDate, err := utility.ParseMandatoryDate(formDTO.EnvironmentNotifyDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(err, "invalid", "environment_notify_date")
	}
	expectedReadyDate, err := utility.ParseMandatoryDate(formDTO.ExpectedReadyDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(err, "invalid", "expected_ready_date")
	}
	testStartDate, err := utility.ParseMandatoryDate(formDTO.TestStartDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(err, "invalid", "test_start_date")
	}
	testEndDate, err := utility.ParseMandatoryDate(formDTO.TestEndDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(err, "invalid", "test_end_date")
	}
	if testEndDate.Before(testStartDate) {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("test end date cannot be before test start date"), "invalid", "test_end_date")
	}
	var managementNotifyDate *time.Time
	if formDTO.ManagementNotifyDate != nil && *formDTO.ManagementNotifyDate != "" {
		managementNotifyDate, err = utility.ParseDate(*formDTO.ManagementNotifyDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(err, "invalid", "management_notify_date")
		}
	}
	var managementResponseDate *time.Time
	if formDTO.ManagementResponseDate != nil && *formDTO.ManagementResponseDate != "" {
		managementResponseDate, err = utility.ParseDate(*formDTO.ManagementResponseDate)
		if err != nil {
			tx.Rollback()
			return 0, utility.NewValidationError(err, "invalid", "management_response_date")
		}
	}
	value := &models.PenetrationTest{
		SystemApplicationName:           formDTO.SystemApplicationName,
		PenetrationTestTypeID:           formDTO.PenetrationTestTypeID,
		PenetrationTestPlanID:           formDTO.PenetrationTestPlanID,
		Year:                            formDTO.Year,
		Period:                          formDTO.Period,
		PenetrationTestClassificationID: formDTO.PenetrationTestClassificationID,
		PenetrationTestExecutorID:       formDTO.PenetrationTestExecutorID,
		CSResponsibleEmployeeID:         formDTO.CSResponsibleEmployeeID,
		Scope:                           formDTO.Scope,
	}
	if err := ptls.repo.ValidateReferences(value); err != nil {
		tx.Rollback()
		return 0, err
	}
	value.CreatedBYID = currentUserID
	value.ReferenceNumber, err = ptls.repo.GenerateReferenceNumber(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := ptls.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.EnvironmentPreparation = &models.PenetrationTestEnvironmentPreparation{
		PenetrationTestID:      id,
		EnvironmentNotifyDate:  environmentNotifyDate,
		ExpectedReadyDate:      expectedReadyDate,
		EnvironmentReadinessID: formDTO.EnvironmentReadinessID,
		CreatedBYID:            currentUserID,
	}
	value.Execution = &models.PenetrationTestExecution{
		PenetrationTestID: id,
		TestStartDate:     testStartDate,
		TestEndDate:       testEndDate,
		ExecutionStatusID: formDTO.ExecutionStatusID,
		SecurityResultID:  formDTO.SecurityResultID,
		CreatedBYID:       currentUserID,
	}
	value.ManagementSharing = &models.PenetrationTestManagementSharing{
		PenetrationTestID:           id,
		ResultsSharedWithManagement: formDTO.ResultsSharedWithManagement,
		ManagementNotifyDate:        managementNotifyDate,
		ManagementResponseDate:      managementResponseDate,
		RemediationNotes:            formDTO.RemediationNotes,
		CreatedBYID:                 currentUserID,
	}
	if _, err := ptls.repo.CreatePenetrationTestEnvironmentPreparation(value.EnvironmentPreparation, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if _, err := ptls.repo.CreatePenetrationTestExecution(value.Execution, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if _, err := ptls.repo.CreatePenetrationTestManagementSharing(value.ManagementSharing, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}
	if ptls.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptls.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}
			if err := ptls.activityLogSvc.LogActivity(models.CreateAction, models.PenetrationTest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_setting", "create penetration test", nil); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}
	return id, nil
}

func (ptls *PenetrationTestService) Update(id uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestFormDTO) error {
	tx := ptls.repo.BeginTransaction()
	value, err := ptls.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	environmentNotifyDate, err := utility.ParseMandatoryDate(form.EnvironmentNotifyDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(err, "invalid", "environment_notify_date")
	}
	expectedReadyDate, err := utility.ParseMandatoryDate(form.ExpectedReadyDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(err, "invalid", "expected_ready_date")
	}
	testStartDate, err := utility.ParseMandatoryDate(form.TestStartDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(err, "invalid", "test_start_date")
	}
	testEndDate, err := utility.ParseMandatoryDate(form.TestEndDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(err, "invalid", "test_end_date")
	}
	if testEndDate.Before(testStartDate) {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("test end date cannot be before test start date"), "invalid", "test_end_date")
	}
	var managementNotifyDate *time.Time
	if form.ManagementNotifyDate != nil && *form.ManagementNotifyDate != "" {
		managementNotifyDate, err = utility.ParseDate(*form.ManagementNotifyDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid", "management_notify_date")
		}
	}
	var managementResponseDate *time.Time
	if form.ManagementResponseDate != nil && *form.ManagementResponseDate != "" {
		managementResponseDate, err = utility.ParseDate(*form.ManagementResponseDate)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(err, "invalid", "management_response_date")
		}
	}
	updated := &models.PenetrationTest{
		SystemApplicationName:           form.SystemApplicationName,
		PenetrationTestTypeID:           form.PenetrationTestTypeID,
		PenetrationTestPlanID:           form.PenetrationTestPlanID,
		Year:                            form.Year,
		Period:                          form.Period,
		PenetrationTestClassificationID: form.PenetrationTestClassificationID,
		PenetrationTestExecutorID:       form.PenetrationTestExecutorID,
		CSResponsibleEmployeeID:         form.CSResponsibleEmployeeID,
		Scope:                           form.Scope,
	}
	updated.EnvironmentPreparation = &models.PenetrationTestEnvironmentPreparation{
		EnvironmentNotifyDate: environmentNotifyDate, ExpectedReadyDate: expectedReadyDate,
		EnvironmentReadinessID: form.EnvironmentReadinessID,
	}
	updated.Execution = &models.PenetrationTestExecution{
		TestStartDate: testStartDate, TestEndDate: testEndDate,
		ExecutionStatusID: form.ExecutionStatusID,
		SecurityResultID:  form.SecurityResultID,
	}
	updated.ManagementSharing = &models.PenetrationTestManagementSharing{
		ResultsSharedWithManagement: form.ResultsSharedWithManagement,
		ManagementNotifyDate:        managementNotifyDate, ManagementResponseDate: managementResponseDate,
		RemediationNotes: form.RemediationNotes,
	}
	if err := ptls.repo.ValidateReferences(updated); err != nil {
		tx.Rollback()
		return err
	}
	value.SystemApplicationName = updated.SystemApplicationName
	value.PenetrationTestTypeID = updated.PenetrationTestTypeID
	value.PenetrationTestPlanID = updated.PenetrationTestPlanID
	value.Year = updated.Year
	value.Period = updated.Period
	value.PenetrationTestClassificationID = updated.PenetrationTestClassificationID
	value.PenetrationTestExecutorID = updated.PenetrationTestExecutorID
	value.CSResponsibleEmployeeID = updated.CSResponsibleEmployeeID
	value.Scope = updated.Scope
	value.UpdatedBYID = &currentUserID
	now := time.Now()
	value.UpdatedAt = &now
	if err := ptls.repo.Update(value, tx); err != nil {
		tx.Rollback()
		return err
	}
	updated.EnvironmentPreparation.PenetrationTestID = id
	updated.Execution.PenetrationTestID = id
	updated.ManagementSharing.PenetrationTestID = id
	updated.EnvironmentPreparation.UpdatedBYID = &currentUserID
	updated.EnvironmentPreparation.UpdatedAt = &now
	updated.Execution.UpdatedBYID = &currentUserID
	updated.Execution.UpdatedAt = &now
	updated.ManagementSharing.UpdatedBYID = &currentUserID
	updated.ManagementSharing.UpdatedAt = &now
	if value.EnvironmentPreparation != nil {
		updated.EnvironmentPreparation.ID = value.EnvironmentPreparation.ID
		updated.EnvironmentPreparation.CreatedBYID = value.EnvironmentPreparation.CreatedBYID
		updated.EnvironmentPreparation.CreatedAt = value.EnvironmentPreparation.CreatedAt
		if err := ptls.repo.UpdatePenetrationTestEnvironmentPreparation(updated.EnvironmentPreparation, tx); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		updated.EnvironmentPreparation.CreatedBYID = currentUserID
		updated.EnvironmentPreparation.CreatedAt = now
		if _, err := ptls.repo.CreatePenetrationTestEnvironmentPreparation(updated.EnvironmentPreparation, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if value.Execution != nil {
		updated.Execution.ID = value.Execution.ID
		updated.Execution.CreatedBYID = value.Execution.CreatedBYID
		updated.Execution.CreatedAt = value.Execution.CreatedAt
		if err := ptls.repo.UpdatePenetrationTestExecution(updated.Execution, tx); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		updated.Execution.CreatedBYID = currentUserID
		updated.Execution.CreatedAt = now
		if _, err := ptls.repo.CreatePenetrationTestExecution(updated.Execution, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if value.ManagementSharing != nil {
		updated.ManagementSharing.ID = value.ManagementSharing.ID
		updated.ManagementSharing.CreatedBYID = value.ManagementSharing.CreatedBYID
		updated.ManagementSharing.CreatedAt = value.ManagementSharing.CreatedAt
		if err := ptls.repo.UpdatePenetrationTestManagementSharing(updated.ManagementSharing, tx); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		updated.ManagementSharing.CreatedBYID = currentUserID
		updated.ManagementSharing.CreatedAt = now
		if _, err := ptls.repo.CreatePenetrationTestManagementSharing(updated.ManagementSharing, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if ptls.activityLogSvc != nil {
		go func() {
			existingEntry, err := ptls.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch penetration test", map[string]interface{}{"error": err.Error()})
				return
			}
			if err := ptls.activityLogSvc.LogActivity(models.UpdateAction, models.PenetrationTest{}.TableName(), id, existingEntry, currentUserID, "penetration_testing", "penetration_testing_setting", "update penetration test", nil); err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return nil
}

func (ptls *PenetrationTestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedRecords []map[string]interface{}
	for _, id := range ids {
		value, getErr := ptls.repo.GetByID(id)
		if getErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		tx := ptls.repo.BeginTransaction()
		if err := ptls.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		successfulIDs = append(successfulIDs, id)
		deletedRecords = append(deletedRecords, map[string]interface{}{"id": id, "reference_number": value.ReferenceNumber})
	}
	if len(successfulIDs) > 0 && ptls.activityLogSvc != nil {
		go func() {
			if err := ptls.activityLogSvc.LogBulkActivities(models.DeleteAction, models.PenetrationTest{}.TableName(), successfulIDs, deletedRecords, currentUserID, nil); err != nil {
				utility.LogError("Failed to log bulk penetration test deletion", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return failedIDs, nil
}
