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
	Update(id uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type PenetrationTestService struct {
	repo           repositories.PenetrationTestRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewPenetrationTestService(repo repositories.PenetrationTestRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) PenetrationTestServiceInterface {
	return &PenetrationTestService{repo: repo, activityLogSvc: activityLogSvc}
}

func (service *PenetrationTestService) GetAll(params *penetrationTestDTOs.PenetrationTestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestListResponseDTO, error) {
	return service.repo.GetAll(params)
}

func (service *PenetrationTestService) GetByID(id uint) (*models.PenetrationTest, error) {
	value, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return value, nil
}

func (service *PenetrationTestService) GetFullDetailsByID(id uint) (*penetrationTestDTOs.PenetrationTestDetailDTO, error) {
	value, err := service.GetByID(id)
	if err != nil {
		return nil, err
	}

	detail := &penetrationTestDTOs.PenetrationTestDetailDTO{}
	if err := mapstructure.Decode(value, detail); err != nil {
		return nil, err
	}
	detail.CreatedAt = value.CreatedAt
	detail.UpdatedAt = value.UpdatedAt
	return detail, nil
}

func (service *PenetrationTestService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFormDTO) (uint, error) {
	tx := service.repo.BeginTransaction()
	value, err := service.buildModel(formDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := service.repo.ValidateReferences(value); err != nil {
		tx.Rollback()
		return 0, err
	}
	value.CreatedBYID = currentUserID
	value.ReferenceNumber, err = service.repo.GenerateReferenceNumber(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	id, err := service.repo.Create(value, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}
	if service.activityLogSvc != nil {
		go service.logActivity(models.CreateAction, id, currentUserID, nil)
	}
	return id, nil
}

func (service *PenetrationTestService) Update(id uint, currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFormDTO) error {
	tx := service.repo.BeginTransaction()
	value, err := service.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	updated, err := service.buildModel(formDTO)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := service.repo.ValidateReferences(updated); err != nil {
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
	value.EnvironmentNotifyDate = updated.EnvironmentNotifyDate
	value.ExpectedReadyDate = updated.ExpectedReadyDate
	value.EnvironmentReadinessID = updated.EnvironmentReadinessID
	value.TestStartDate = updated.TestStartDate
	value.TestEndDate = updated.TestEndDate
	value.ExecutionStatusID = updated.ExecutionStatusID
	value.SecurityResultID = updated.SecurityResultID
	value.ResultsSharedWithManagement = updated.ResultsSharedWithManagement
	value.ManagementNotifyDate = updated.ManagementNotifyDate
	value.ManagementResponseDate = updated.ManagementResponseDate
	value.RemediationNotes = updated.RemediationNotes
	value.UpdatedBYID = &currentUserID
	now := time.Now()
	value.UpdatedAt = &now
	if err := service.repo.Update(value, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}
	if service.activityLogSvc != nil {
		go service.logActivity(models.UpdateAction, id, currentUserID, nil)
	}
	return nil
}

func (service *PenetrationTestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedRecords []map[string]interface{}
	for _, id := range ids {
		value, getErr := service.repo.GetByID(id)
		if getErr != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		tx := service.repo.BeginTransaction()
		if err := service.repo.Delete(id, tx); err != nil {
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
	if len(successfulIDs) > 0 && service.activityLogSvc != nil {
		go func() {
			if err := service.activityLogSvc.LogBulkActivities(models.DeleteAction, models.PenetrationTest{}.TableName(), successfulIDs, deletedRecords, currentUserID, nil); err != nil {
				utility.LogError("Failed to log bulk penetration test deletion", map[string]interface{}{"error": err.Error()})
			}
		}()
	}
	return failedIDs, nil
}

func (service *PenetrationTestService) logActivity(action models.Action, id uint, currentUserID uint, data interface{}) {
	value, err := service.GetFullDetailsByID(id)
	if err != nil {
		utility.LogError("Failed to fetch penetration test for activity log", map[string]interface{}{"error": err.Error()})
		return
	}
	if data == nil {
		data = value
	}
	if err := service.activityLogSvc.LogActivity(action, models.PenetrationTest{}.TableName(), id, data, currentUserID, "penetration_testing", "penetration_testing_setting", "", nil); err != nil {
		utility.LogError("Failed to log penetration test activity", map[string]interface{}{"error": err.Error()})
	}
}

func (service *PenetrationTestService) buildModel(formDTO *penetrationTestDTOs.PenetrationTestFormDTO) (*models.PenetrationTest, error) {
	environmentNotifyDate, err := utility.ParseMandatoryDate(formDTO.EnvironmentNotifyDate)
	if err != nil {
		return nil, utility.NewValidationError(err, "invalid", "environment_notify_date")
	}
	expectedReadyDate, err := utility.ParseMandatoryDate(formDTO.ExpectedReadyDate)
	if err != nil {
		return nil, utility.NewValidationError(err, "invalid", "expected_ready_date")
	}
	testStartDate, err := utility.ParseMandatoryDate(formDTO.TestStartDate)
	if err != nil {
		return nil, utility.NewValidationError(err, "invalid", "test_start_date")
	}
	testEndDate, err := utility.ParseMandatoryDate(formDTO.TestEndDate)
	if err != nil {
		return nil, utility.NewValidationError(err, "invalid", "test_end_date")
	}
	if testEndDate.Before(testStartDate) {
		return nil, utility.NewValidationError(fmt.Errorf("test end date cannot be before test start date"), "invalid", "test_end_date")
	}
	managementNotifyDate, err := parseOptionalDate(formDTO.ManagementNotifyDate, "management_notify_date")
	if err != nil {
		return nil, err
	}
	managementResponseDate, err := parseOptionalDate(formDTO.ManagementResponseDate, "management_response_date")
	if err != nil {
		return nil, err
	}
	return &models.PenetrationTest{
		SystemApplicationName: formDTO.SystemApplicationName, PenetrationTestTypeID: formDTO.PenetrationTestTypeID,
		PenetrationTestPlanID: formDTO.PenetrationTestPlanID, Year: formDTO.Year, Period: formDTO.Period,
		PenetrationTestClassificationID: formDTO.PenetrationTestClassificationID, PenetrationTestExecutorID: formDTO.PenetrationTestExecutorID,
		CSResponsibleEmployeeID: formDTO.CSResponsibleEmployeeID, Scope: formDTO.Scope, EnvironmentNotifyDate: environmentNotifyDate,
		ExpectedReadyDate: expectedReadyDate, EnvironmentReadinessID: formDTO.EnvironmentReadinessID, TestStartDate: testStartDate,
		TestEndDate: testEndDate, ExecutionStatusID: formDTO.ExecutionStatusID, SecurityResultID: formDTO.SecurityResultID,
		ResultsSharedWithManagement: formDTO.ResultsSharedWithManagement, ManagementNotifyDate: managementNotifyDate,
		ManagementResponseDate: managementResponseDate, RemediationNotes: formDTO.RemediationNotes,
	}, nil
}

func parseOptionalDate(value *string, field string) (*time.Time, error) {
	if value == nil || *value == "" {
		return nil, nil
	}
	parsed, err := utility.ParseDate(*value)
	if err != nil {
		return nil, utility.NewValidationError(err, "invalid", field)
	}
	return parsed, nil
}
