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

type PenetrationTestEnvironmentPreparationServiceInterface interface {
	GetByID(penetrationTestID uint) (*models.PenetrationTestEnvironmentPreparation, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestEnvironmentPreparationDetailDTO, error)
	Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestEnvironmentPreparationFormDTO) (uint, error)
	Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestEnvironmentPreparationFormDTO) error
	Delete(penetrationTestID uint, currentUserID uint) (uint, error)
}

type PenetrationTestEnvironmentPreparationService struct {
	repo               penetrationTestRepositories.PenetrationTestEnvironmentPreparationRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	penetrationTestSvc PenetrationTestServiceInterface
	readinessSvc       PenetrationTestEnvironmentReadinessServiceInterface
}

func NewPenetrationTestEnvironmentPreparationService(
	repo penetrationTestRepositories.PenetrationTestEnvironmentPreparationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
	readinessSvc PenetrationTestEnvironmentReadinessServiceInterface,
) PenetrationTestEnvironmentPreparationServiceInterface {
	return &PenetrationTestEnvironmentPreparationService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		penetrationTestSvc: penetrationTestSvc,
		readinessSvc:       readinessSvc,
	}
}

func (ptps PenetrationTestEnvironmentPreparationService) GetByID(penetrationTestID uint) (*models.PenetrationTestEnvironmentPreparation, error) {
	if _, err := ptps.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	environmentPreparation, err := ptps.repo.GetByID(penetrationTestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}
	return environmentPreparation, nil
}

func (ptps PenetrationTestEnvironmentPreparationService) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestEnvironmentPreparationDetailDTO, error) {
	var environmentPreparationDTO penetrationTestDTOs.PenetrationTestEnvironmentPreparationDetailDTO

	if _, err := ptps.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	environmentPreparation, err := ptps.repo.GetFullDetailsByID(penetrationTestID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}

	err = mapstructure.Decode(environmentPreparation, &environmentPreparationDTO)
	if err != nil {
		return nil, err
	}

	environmentNotifyDate := environmentPreparation.EnvironmentNotifyDate.Format("2006-01-02")
	environmentPreparationDTO.EnvironmentNotifyDateString = &environmentNotifyDate

	expectedReadyDate := environmentPreparation.ExpectedReadyDate.Format("2006-01-02")
	environmentPreparationDTO.ExpectedReadyDateString = &expectedReadyDate

	environmentPreparationDTO.CreatedAt = environmentPreparation.CreatedAt
	environmentPreparationDTO.UpdatedAt = environmentPreparation.UpdatedAt

	return &environmentPreparationDTO, nil
}

func (ptps *PenetrationTestEnvironmentPreparationService) Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestEnvironmentPreparationFormDTO) (uint, error) {
	tx := ptps.repo.BeginTransaction()

	if _, err := ptps.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	if exists, err := ptps.repo.ExistsByPenetrationTestID(penetrationTestID, tx); err != nil {
		tx.Rollback()
		return 0, err
	} else if exists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("environment preparation already exists"), "unique", "environment_preparation")
	}
	if _, err := ptps.readinessSvc.GetByID(form.EnvironmentReadinessID); err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid environment readiness"), "invalid", "environment_readiness_id")
	}

	environmentNotifyDate, err := utility.ParseMandatoryDate(form.EnvironmentNotifyDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid environment notify date"), "invalid_date", "environment_notify_date")
	}

	expectedReadyDate, err := utility.ParseMandatoryDate(form.ExpectedReadyDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid expected ready date"), "invalid_date", "expected_ready_date")
	}

	if expectedReadyDate.Before(environmentNotifyDate) {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("expected ready date cannot be before notification"), "environment_ready_before_notify", "expected_ready_date")
	}

	environmentPreparation := &models.PenetrationTestEnvironmentPreparation{
		PenetrationTestID:      penetrationTestID,
		EnvironmentNotifyDate:  environmentNotifyDate,
		ExpectedReadyDate:      expectedReadyDate,
		EnvironmentReadinessID: form.EnvironmentReadinessID,
		CreatedBYID:            currentUserID,
	}

	id, err := ptps.repo.Create(environmentPreparation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingEnvironmentPreparation, err := ptps.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test environment preparation", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestEnvironmentPreparation{}.TableName(),
			id,
			existingEnvironmentPreparation,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "environment preparation", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.PenetrationTestEnvironmentPreparation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ptps *PenetrationTestEnvironmentPreparationService) Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestEnvironmentPreparationFormDTO) error {
	tx := ptps.repo.BeginTransaction()

	if _, err := ptps.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return err
	}

	existingEnvironmentPreparation, err := ptps.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := ptps.readinessSvc.GetByID(form.EnvironmentReadinessID); err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid environment readiness"), "invalid", "environment_readiness_id")
	}

	environmentNotifyDate, err := utility.ParseMandatoryDate(form.EnvironmentNotifyDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid environment notify date"), "invalid_date", "environment_notify_date")
	}

	expectedReadyDate, err := utility.ParseMandatoryDate(form.ExpectedReadyDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid expected ready date"), "invalid_date", "expected_ready_date")
	}

	if expectedReadyDate.Before(environmentNotifyDate) {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("expected ready date cannot be before notification"), "environment_ready_before_notify", "expected_ready_date")
	}

	existingEnvironmentPreparation.EnvironmentNotifyDate = environmentNotifyDate
	existingEnvironmentPreparation.ExpectedReadyDate = expectedReadyDate
	existingEnvironmentPreparation.EnvironmentReadinessID = form.EnvironmentReadinessID
	existingEnvironmentPreparation.UpdatedBYID = &currentUserID

	err = ptps.repo.Update(existingEnvironmentPreparation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingEntry, err := ptps.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test environment preparation", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestEnvironmentPreparation{}.TableName(),
			existingEnvironmentPreparation.ID,
			existingEntry,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "environment preparation", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.PenetrationTestEnvironmentPreparation{}.TableName(),
				"recordID":  existingEnvironmentPreparation.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ptps *PenetrationTestEnvironmentPreparationService) Delete(penetrationTestID uint, currentUserID uint) (uint, error) {
	tx := ptps.repo.BeginTransaction()
	var failedId uint
	var deletedEnvironmentPreparation []map[string]interface{}

	if _, err := ptps.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	environmentPreparation, err := ptps.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = ptps.repo.Delete(environmentPreparation.ID, tx); err != nil {
		tx.Rollback()
		failedId = environmentPreparation.ID
		return failedId, err
	}

	deletedEnvironmentPreparation = append(deletedEnvironmentPreparation, map[string]interface{}{
		"id":                  environmentPreparation.ID,
		"penetration_test_id": penetrationTestID,
	})

	tx.Commit()

	go func() {
		err := ptps.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PenetrationTestEnvironmentPreparation{}.TableName(),
			environmentPreparation.ID,
			deletedEnvironmentPreparation,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "environment preparation", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
