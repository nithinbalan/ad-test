package services

import (
	"errors"
	"fmt"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
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
	repo               repositories.PenetrationTestRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	typeSvc            PenetrationTestTypeServiceInterface
	planSvc            PenetrationTestPlanServiceInterface
	classificationSvc  PenetrationTestClassificationServiceInterface
	executorSvc        PenetrationTestExecutorServiceInterface
	userSvc            organizationServices.UserServiceInterface
}

func NewPenetrationTestService(
	repo repositories.PenetrationTestRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	typeSvc PenetrationTestTypeServiceInterface,
	planSvc PenetrationTestPlanServiceInterface,
	classificationSvc PenetrationTestClassificationServiceInterface,
	executorSvc PenetrationTestExecutorServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) PenetrationTestServiceInterface {
	return &PenetrationTestService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		typeSvc:            typeSvc,
		planSvc:            planSvc,
		classificationSvc:  classificationSvc,
		executorSvc:        executorSvc,
		userSvc:            userSvc,
	}
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
		detail.EnvironmentNotifyDate = &value.EnvironmentPreparation.EnvironmentNotifyDate
		detail.ExpectedReadyDate = &value.EnvironmentPreparation.ExpectedReadyDate
		if value.EnvironmentPreparation.EnvironmentReadiness != nil {
			detail.EnvironmentReadiness = &penetrationTestDTOs.PenetrationTestEnvironmentReadinessDTO{
				ID:      value.EnvironmentPreparation.EnvironmentReadiness.ID,
				TitleEn: value.EnvironmentPreparation.EnvironmentReadiness.TitleEn,
				TitleAr: value.EnvironmentPreparation.EnvironmentReadiness.TitleAr,
				Status:  value.EnvironmentPreparation.EnvironmentReadiness.Status,
			}
		}
	}
	if value.Execution != nil {
		detail.TestStartDate = &value.Execution.TestStartDate
		detail.TestEndDate = &value.Execution.TestEndDate
		if value.Execution.ExecutionStatus != nil {
			detail.ExecutionStatus = &penetrationTestDTOs.PenetrationTestExecutionStatusDTO{
				ID:      value.Execution.ExecutionStatus.ID,
				TitleEn: value.Execution.ExecutionStatus.TitleEn,
				TitleAr: value.Execution.ExecutionStatus.TitleAr,
				Type:    value.Execution.ExecutionStatus.Type,
				Status:  value.Execution.ExecutionStatus.Status,
			}
		}
		if value.Execution.SecurityResult != nil {
			detail.SecurityResult = &penetrationTestDTOs.PenetrationTestSecurityResultDTO{
				ID:      value.Execution.SecurityResult.ID,
				TitleEn: value.Execution.SecurityResult.TitleEn,
				TitleAr: value.Execution.SecurityResult.TitleAr,
				Type:    value.Execution.SecurityResult.Type,
				Status:  value.Execution.SecurityResult.Status,
			}
		}
	}
	if value.ManagementSharing != nil {
		detail.ResultsSharedWithManagement = value.ManagementSharing.ResultsSharedWithManagement
		detail.ManagementNotifyDate = value.ManagementSharing.ManagementNotifyDate
		detail.ManagementResponseDate = value.ManagementSharing.ManagementResponseDate
		detail.RemediationNotes = value.ManagementSharing.RemediationNotes
	}
	counts, err := ptls.repo.GetSeverityCounts(id)
	if err != nil {
		return nil, err
	}
	detail.CriticalCount = counts.CriticalCount
	detail.HighCount = counts.HighCount
	detail.MediumCount = counts.MediumCount
	detail.LowCount = counts.LowCount
	detail.TotalCount = counts.TotalCount
	return detail, nil
}

func (ptls *PenetrationTestService) Create(currentUserID uint, formDTO *penetrationTestDTOs.PenetrationTestFormDTO) (uint, error) {
	tx := ptls.repo.BeginTransaction()
	if err := ptls.validatePenetrationTestFormDTO(formDTO); err != nil {
		tx.Rollback()
		return 0, err
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
	value.CreatedBYID = currentUserID
	referenceNumber, err := ptls.repo.GenerateReferenceNumber(tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	value.ReferenceNumber = referenceNumber

	id, err := ptls.repo.Create(value, tx)
	if err != nil {
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
	if err := ptls.validatePenetrationTestFormDTO(form); err != nil {
		tx.Rollback()
		return err
	}
	value, err := ptls.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	value.SystemApplicationName = form.SystemApplicationName
	value.PenetrationTestTypeID = form.PenetrationTestTypeID
	value.PenetrationTestPlanID = form.PenetrationTestPlanID
	value.Year = form.Year
	value.Period = form.Period
	value.PenetrationTestClassificationID = form.PenetrationTestClassificationID
	value.PenetrationTestExecutorID = form.PenetrationTestExecutorID
	value.CSResponsibleEmployeeID = form.CSResponsibleEmployeeID
	value.Scope = form.Scope
	value.UpdatedBYID = &currentUserID
	now := time.Now()
	value.UpdatedAt = &now
	if err := ptls.repo.Update(value, tx); err != nil {
		tx.Rollback()
		return err
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
		if err := ptls.repo.DeletePenetrationTestDetails(id, tx); err != nil {
			tx.Rollback()
			failedIDs = append(failedIDs, id)
			continue
		}
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

func (ptls *PenetrationTestService) validatePenetrationTestFormDTO(formDTO *penetrationTestDTOs.PenetrationTestFormDTO) error {
	if _, err := ptls.typeSvc.GetByID(formDTO.PenetrationTestTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test type does not exist"), "not_exist", "penetration test type")
	}
	if _, err := ptls.planSvc.GetByID(formDTO.PenetrationTestPlanID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test plan does not exist"), "not_exist", "penetration test plan")
	}
	if _, err := ptls.classificationSvc.GetByID(formDTO.PenetrationTestClassificationID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test classification does not exist"), "not_exist", "penetration test classification")
	}
	if _, err := ptls.executorSvc.GetByID(formDTO.PenetrationTestExecutorID); err != nil {
		return utility.NewValidationError(fmt.Errorf("penetration test executor does not exist"), "not_exist", "penetration test executor")
	}
	if _, err := ptls.userSvc.GetByID(formDTO.CSResponsibleEmployeeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("CS responsible employee does not exist"), "not_exist", "CS responsible employee")
	}
	return nil
}
