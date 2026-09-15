package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	penetrationTestRepositories "jethur/grc/penetration-testing/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PenetrationTestClosureServiceInterface interface {
	GetByID(penetrationTestID uint) (*models.PenetrationTestClosure, error)
	GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestClosureDetailDTO, error)
	Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestClosureFormDTO) (uint, error)
	Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestClosureFormDTO) error
	Delete(penetrationTestID uint, currentUserID uint) (uint, error)
}

type PenetrationTestClosureService struct {
	repo               penetrationTestRepositories.PenetrationTestClosureRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	penetrationTestSvc PenetrationTestServiceInterface
	userSvc            organizationServices.UserServiceInterface
}

func NewPenetrationTestClosureService(
	repo penetrationTestRepositories.PenetrationTestClosureRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	penetrationTestSvc PenetrationTestServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) PenetrationTestClosureServiceInterface {
	return &PenetrationTestClosureService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		penetrationTestSvc: penetrationTestSvc,
		userSvc:            userSvc,
	}
}

func (ptcs PenetrationTestClosureService) GetByID(penetrationTestID uint) (*models.PenetrationTestClosure, error) {
	if _, err := ptcs.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	closure, err := ptcs.repo.GetByID(penetrationTestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}
	return closure, nil
}

func (ptcs PenetrationTestClosureService) GetFullDetailsByID(penetrationTestID uint, tx ...utility.DatabaseTransaction) (*penetrationTestDTOs.PenetrationTestClosureDetailDTO, error) {
	var closureDTO penetrationTestDTOs.PenetrationTestClosureDetailDTO

	if _, err := ptcs.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		return nil, err
	}

	closure, err := ptcs.repo.GetFullDetailsByID(penetrationTestID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, penetrationTestID)
		}
		return nil, err
	}

	err = mapstructure.Decode(closure, &closureDTO)
	if err != nil {
		return nil, err
	}

	closureDate := closure.ClosureDate.Format("2006-01-02")
	closureDTO.ClosureDateString = &closureDate

	closureDTO.CreatedAt = closure.CreatedAt
	closureDTO.UpdatedAt = closure.UpdatedAt

	return &closureDTO, nil
}

func (ptcs *PenetrationTestClosureService) Create(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestClosureFormDTO) (uint, error) {
	tx := ptcs.repo.BeginTransaction()

	if _, err := ptcs.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	if exists, err := ptcs.repo.ExistsByPenetrationTestID(penetrationTestID, tx); err != nil {
		tx.Rollback()
		return 0, err
	} else if exists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("closure already exists"), "closure_already_exists", "closure")
	}

	if hasOpenFindings, err := ptcs.repo.HasOpenFindings(penetrationTestID, tx); err != nil {
		tx.Rollback()
		return 0, err
	} else if hasOpenFindings {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("all findings must be closed before closing the penetration test"), "findings_not_closed", "findings")
	}

	if _, err := ptcs.userSvc.GetByID(form.ClosureOwnerID); err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("closure owner does not exist"), "not_exist", "closure_owner")
	}

	closureDate, err := utility.ParseMandatoryDate(form.ClosureDate)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("invalid closure date"), "invalid_date", "closure_date")
	}

	closure := &models.PenetrationTestClosure{
		PenetrationTestID: penetrationTestID,
		ClosureOwnerID:    form.ClosureOwnerID,
		ClosureDate:       closureDate,
		ClosureRemarks:    form.ClosureRemarks,
		CreatedBYID:       currentUserID,
	}

	id, err := ptcs.repo.Create(closure, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := ptcs.penetrationTestSvc.UpdateStatus(penetrationTestID, "closed", tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingClosure, err := ptcs.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test closure", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.PenetrationTestClosure{}.TableName(),
			id,
			existingClosure,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "closure", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.PenetrationTestClosure{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ptcs *PenetrationTestClosureService) Update(penetrationTestID uint, currentUserID uint, form *penetrationTestDTOs.PenetrationTestClosureFormDTO) error {
	tx := ptcs.repo.BeginTransaction()

	if _, err := ptcs.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return err
	}

	existingClosure, err := ptcs.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if _, err := ptcs.userSvc.GetByID(form.ClosureOwnerID); err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("closure owner does not exist"), "not_exist", "closure_owner")
	}

	closureDate, err := utility.ParseMandatoryDate(form.ClosureDate)
	if err != nil {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("invalid closure date"), "invalid_date", "closure_date")
	}

	existingClosure.ClosureOwnerID = form.ClosureOwnerID
	existingClosure.ClosureDate = closureDate
	existingClosure.ClosureRemarks = form.ClosureRemarks
	existingClosure.UpdatedBYID = &currentUserID

	err = ptcs.repo.Update(existingClosure, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingEntry, err := ptcs.GetFullDetailsByID(penetrationTestID)
		if err != nil {
			utility.LogError("Failed to fetch penetration test closure", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ptcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.PenetrationTestClosure{}.TableName(),
			existingClosure.ID,
			existingEntry,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "closure", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.PenetrationTestClosure{}.TableName(),
				"recordID":  existingClosure.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ptcs *PenetrationTestClosureService) Delete(penetrationTestID uint, currentUserID uint) (uint, error) {
	tx := ptcs.repo.BeginTransaction()
	var failedId uint
	var deletedClosure []map[string]interface{}

	if _, err := ptcs.penetrationTestSvc.GetByID(penetrationTestID); err != nil {
		tx.Rollback()
		return 0, err
	}

	closure, err := ptcs.GetByID(penetrationTestID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = ptcs.repo.Delete(closure.ID, tx); err != nil {
		tx.Rollback()
		failedId = closure.ID
		return failedId, err
	}

	if err := ptcs.penetrationTestSvc.UpdateStatus(penetrationTestID, "draft", tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	deletedClosure = append(deletedClosure, map[string]interface{}{
		"id":                  closure.ID,
		"penetration_test_id": penetrationTestID,
	})

	tx.Commit()

	go func() {
		err := ptcs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.PenetrationTestClosure{}.TableName(),
			closure.ID,
			deletedClosure,
			currentUserID,
			"penetration_testing", "penetration_testing_pentest_log", "closure", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}
